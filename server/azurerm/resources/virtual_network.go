package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

// VirtualNetwork is the entity that holds the logic to calculate price
// of the azure_network_virtualnetwork
type VirtualNetwork struct {
	provider *Provider

	location         string
	peeringLocations []string

	// Usage
	monthlyDataTransferGB decimal.Decimal
}

// virtualNetworkValues is holds the values that we need to be able
// to calculate the price of the Virtual Network Values
type virtualNetworkValues struct {
	Location         string   `mapstructure:"location"`
	PeeringLocations []string `mapstructure:"peering_locations"`

	Usage struct {
		MonthlyDataTransferGB float64 `mapstructure:"monthly_data_transfer_gb"`
	} `mapstructure:"tc_usage"`
}

// decodeVirtualNetworkValues decodes and returns computeInstanceValues from a Terraform values map.
func decodeVirtualNetworkValues(tfVals map[string]interface{}) (virtualNetworkValues, error) {
	var v virtualNetworkValues
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

	if v.Usage.MonthlyDataTransferGB == 0 {
		v.Usage.MonthlyDataTransferGB = 100
	}
	return v, nil
}

func (p *Provider) newVirtualNetwork(vals virtualNetworkValues) *VirtualNetwork {
	inst := &VirtualNetwork{
		provider: p,

		location:              vals.Location,
		peeringLocations:      vals.PeeringLocations,
		monthlyDataTransferGB: decimal.NewFromFloat(vals.Usage.MonthlyDataTransferGB),
	}
	return inst
}

func (inst *VirtualNetwork) Components() ([]query.Component, error) {
	var components []query.Component

	for _, loc := range inst.peeringLocations {
		virtualNetworkP, err := decodeVirtualNetworkPeeringValues(virtualNetworkPeeringValues{
			SourceLocation:      inst.location,
			DestinationLocation: loc,
			Usage: struct {
				MonthlyDataTransferGB float64 `mapstructure:"monthly_data_transfer_gb"`
			}{MonthlyDataTransferGB: inst.monthlyDataTransferGB.InexactFloat64()},
		})
		if err != nil {
			return nil, err
		}
		peering := inst.provider.newVirtualNetworkPeering(virtualNetworkP)
		components = append(components, peering.Components()...)
	}

	return components, nil
}
