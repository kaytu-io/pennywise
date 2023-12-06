package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

type RMDNSZone struct {
	location string
}

type RMDNSZoneValue struct {
	Location string `mapstructure:"resource_group_name"`
}

func (p *Provider) newRMDNSZone(vals RMDNSZoneValue) *RMDNSZone {
	inst := &RMDNSZone{
		location: vals.Location,
	}
	return inst
}

func decoderRMDNSZone(tfVals map[string]interface{}) (RMDNSZoneValue, error) {
	var v RMDNSZoneValue
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &v,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return v, err
	}

	if err := decoder.Decode(tfVals); err != nil {
		return v, err
	}
	return v, nil
}

func (inst *RMDNSZone) component() []query.Component {
	region := inst.location

	if strings.HasPrefix(strings.ToLower(region), "usgov") {
		region = "US Gov Zone 1"
	} else if strings.HasPrefix(strings.ToLower(region), "germany") {
		region = "DE Zone 1"
	} else if strings.HasPrefix(strings.ToLower(region), "china") {
		region = "Zone 1 (China)"
	} else {
		region = "Zone 1"
	}

	costComponents := make([]query.Component, 0)
	costComponents = append(costComponents, hostedPublicZoneCostComponent(region))

	return costComponents
}

func hostedPublicZoneCostComponent(region string) query.Component {
	return query.Component{
		Name:            "Hosted zone",
		Unit:            "months",
		MonthlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Location: util.StringPtr(region),
			Service:  util.StringPtr("Azure DNS"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", ValueRegex: util.StringPtr("Public Zone(s)?")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "Start_usage_amount", Value: util.StringPtr("0")},
			},
		},
	}
}
