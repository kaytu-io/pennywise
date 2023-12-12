package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/tier_request"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

type DNSARecord struct {
	provider *Provider

	location string
	// Usage
	monthlyQueries *int64
}

type ResourceGroupNameStruct struct {
	Values struct {
		Location string `mapstructure:"location"`
	} `mapstructure:"values"`
}

type dnsARecordValues struct {
	ResourceGroupName ResourceGroupNameStruct `mapstructure:"resource_group_name"`

	Usage struct {
		MonthlyQueries *int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"tc_usage"`
}

func decodeDNSARecord(tfVals map[string]interface{}) (dnsARecordValues, error) {
	var v dnsARecordValues
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

func (p *Provider) newDNSARecord(vals dnsARecordValues) *DNSARecord {
	inst := &DNSARecord{
		location:       vals.ResourceGroupName.Values.Location,
		monthlyQueries: vals.Usage.MonthlyQueries,
	}
	return inst
}

func (inst *DNSARecord) component() []query.Component {

	return DNSQueriesCostComponent(inst.location, inst.monthlyQueries)
}

func DNSQueriesCostComponent(region string, monthlyQueries *int64) []query.Component {
	region = getLocationName(region)

	var monthlyQueriesDec decimal.Decimal
	var requestQuantities []decimal.Decimal
	var costComponents []query.Component
	requests := []int{1000000000}
	if strings.HasPrefix(strings.ToLower(region), "usgov") {
		region = "US Gov Zone 1"
	}
	if strings.HasPrefix(strings.ToLower(region), "germany") {
		region = "DE Zone 1"
	}
	if strings.HasPrefix(strings.ToLower(region), "china") {
		region = "Zone 1 (China)"
	}
	if region != "US Gov Zone 1" && region != "DE Zone 1" && region != "Zone 1 (China)" {
		region = "Zone 1"
	}

	if monthlyQueries != nil {
		monthlyQueriesDec = decimal.NewFromInt(*monthlyQueries)
		requestQuantities = tier_request.CalculateTierBuckets(monthlyQueriesDec, requests)

		firstBqueries := requestQuantities[0].Div(decimal.NewFromInt(1000000))
		overBqueries := requestQuantities[1].Div(decimal.NewFromInt(1000000))
		costComponents = append(costComponents, dnsQueriesFirstCostComponent(region, "DNS queries (first 1B)", "0", &firstBqueries))

		if requestQuantities[1].GreaterThan(decimal.NewFromInt(0)) {
			costComponents = append(costComponents, dnsQueriesFirstCostComponent(region, "DNS queries (over 1B)", "1000", &overBqueries))
		}
	} else {
		var unknown decimal.Decimal
		costComponents = append(costComponents, dnsQueriesFirstCostComponent(region, "DNS queries (first 1B)", "0", &unknown))
	}
	return costComponents
}

func dnsQueriesFirstCostComponent(region, name, startUsage string, monthlyQueries *decimal.Decimal) query.Component {
	return query.Component{
		Name:            name,
		Unit:            "1M queries",
		MonthlyQuantity: *monthlyQueries,
		ProductFilter: &product.Filter{
			Location: util.StringPtr(region),
			Service:  util.StringPtr("Azure DNS"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", Value: util.StringPtr("Public Queries")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr(startUsage)}},
		},
	}
}
