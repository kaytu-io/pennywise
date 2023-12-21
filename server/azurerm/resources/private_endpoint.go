package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/tier_request"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

// PrivateEndpoint is the entity that holds the logic to calculate price
// of the azurerm_private_endpoint
type PrivateEndpoint struct {
	provider *Provider

	location string

	// Usage
	// receive monthly inbound data processed in GB.
	monthlyInboundDataProcessedGb *int64
	// receive monthly outbound data processed in GB.
	monthlyOutboundDataProcessedGb *int64
}

// privateEndpointValues is holds the values that we need to be able
// to calculate the price of the PrivateEndpoint
type privateEndpointValues struct {
	Location string `mapstructure:"location"`

	Usage struct {
		MonthlyInboundDataProcessedGb  *int64 `mapstructure:"monthly_inbound_data_processed_gb"`
		MonthlyOutboundDataProcessedGb *int64 `mapstructure:"monthly_outbound_data_processed_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodePrivateEndpointValues decodes and returns privateEndpointValues from a Terraform values map.
func decodePrivateEndpointValues(tfVals map[string]interface{}) (privateEndpointValues, error) {
	var v privateEndpointValues
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

// newPrivateEndpoint initializes a new PrivateEndpoint from the provider
func (p *Provider) newPrivateEndpoint(vals privateEndpointValues) *PrivateEndpoint {
	inst := &PrivateEndpoint{
		provider: p,

		location:                       convertRegion(vals.Location),
		monthlyInboundDataProcessedGb:  vals.Usage.MonthlyInboundDataProcessedGb,
		monthlyOutboundDataProcessedGb: vals.Usage.MonthlyOutboundDataProcessedGb,
	}
	return inst
}

func (inst *PrivateEndpoint) Components() []query.Component {
	var components []query.Component

	components = append(components, inst.privateEndpointCostComponent("Private endpoint", "Standard Private Endpoint"))

	if inst.monthlyInboundDataProcessedGb != nil {
		inboundTiers := []int{1_000_000, 4_000_000}
		inboundQuantities := tier_request.CalculateTierBuckets(decimal.NewFromInt(*inst.monthlyInboundDataProcessedGb), inboundTiers)

		if len(inboundQuantities) > 0 {
			components = append(components, inst.privateEndpointDataCostComponent("Inbound data processed (first 1PB)", "Standard Data Processed - Ingress", "0", inboundQuantities[0]))
		}
		if len(inboundQuantities) > 1 && inboundQuantities[1].GreaterThan(decimal.Zero) {
			components = append(components, inst.privateEndpointDataCostComponent("Inbound data processed (next 4PB)", "Standard Data Processed - Ingress", "1000000", inboundQuantities[1]))
		}
		if len(inboundQuantities) > 2 && inboundQuantities[2].GreaterThan(decimal.Zero) {
			components = append(components, inst.privateEndpointDataCostComponent("Inbound data processed (over 5PB)", "Standard Data Processed - Ingress", "5000000", inboundQuantities[2]))
		}
	} else {
		components = append(components, inst.privateEndpointDataCostComponent("Inbound data processed (first 1PB)", "Standard Data Processed - Ingress", "0", decimal.Zero))
	}

	if inst.monthlyOutboundDataProcessedGb != nil {
		outboundTiers := []int{1_000_000, 4_000_000}
		outboundQuantities := tier_request.CalculateTierBuckets(decimal.NewFromInt(*inst.monthlyOutboundDataProcessedGb), outboundTiers)

		if len(outboundQuantities) > 0 {
			components = append(components, inst.privateEndpointDataCostComponent("Outbound data processed (first 1PB)", "Standard Data Processed - Egress", "0", outboundQuantities[0]))
		}

		if len(outboundQuantities) > 1 && outboundQuantities[1].GreaterThan(decimal.Zero) {
			components = append(components, inst.privateEndpointDataCostComponent("Outbound data processed (next 4PB)", "Standard Data Processed - Egress", "1000000", outboundQuantities[1]))
		}

		if len(outboundQuantities) > 2 && outboundQuantities[2].GreaterThan(decimal.Zero) {
			components = append(components, inst.privateEndpointDataCostComponent("Outbound data processed (over 5PB)", "Standard Data Processed - Egress", "5000000", outboundQuantities[2]))
		}
	} else {
		components = append(components, inst.privateEndpointDataCostComponent("Outbound data processed (first 1PB)", "Standard Data Processed - Egress", "0", decimal.Zero))
	}

	return components
}

func (inst *PrivateEndpoint) privateEndpointCostComponent(name, meterName string) query.Component {
	return query.Component{
		Name:            name,
		Unit:            "hour",
		MonthlyQuantity: decimal.NewFromInt(730),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Virtual Network"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Virtual Network Private Link")},
				{Key: "meter_name", Value: util.StringPtr(meterName)},
			},
		},
	}
}

func (inst *PrivateEndpoint) privateEndpointDataCostComponent(name, meterName, startUsage string, quantity decimal.Decimal) query.Component {
	return query.Component{
		Name:            name,
		Unit:            "GB",
		MonthlyQuantity: quantity,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Virtual Network"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Virtual Network Private Link")},
				{Key: "meter_name", Value: util.StringPtr(meterName)},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "tier_minimum_units", Value: util.StringPtr(startUsage)},
			},
		},
	}
}
