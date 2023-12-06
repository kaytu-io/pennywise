package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
	"strings"
)

type PrivateDNSZone struct {
	location string
}

type privateDNSZoneValue struct {
	Value struct {
		Location string `mapstructure:"location"`
	} `mapstructure:"value"`
}

func (p *Provider) newPrivateDNSZone(vals privateDNSZoneValue) *PrivateDNSZone {
	inst := &PrivateDNSZone{
		location: vals.Value.Location,
	}
	return inst
}

func decoderPrivateDnsZone(tfVals map[string]interface{}) (privateDNSZoneValue, error) {
	var v privateDNSZoneValue
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

func (inst *PrivateDNSZone) component() []query.Component {
	costComponent := make([]query.Component, 0)
	costComponent = append(costComponent, PrivateDNSZoneCostComponent(inst.location))

	return costComponent
}

func PrivateDNSZoneCostComponent(region string) query.Component {
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

	return hostedPublicZoneCostComponent(region)
}
