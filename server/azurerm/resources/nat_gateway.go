package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

// NatGateway is the entity that holds the logic to calculate price
// of the azurerm_application_gateway
type NatGateway struct {
	provider *Provider

	location string

	// Usage
	monthlyDataProcessedGb *int64
}

// natGatewayValues is holds the values that we need to be able
// to calculate the price of the ComputeInstance
type natGatewayValues struct {
	Location string `mapstructure:"location"`

	Usage struct {
		MonthlyDataProcessedGb *int64 `mapstructure:"monthly_data_processed_gb"`
	} `mapstructure:"tc_usage"`
}

// decodeNatGatewayValues decodes and returns natGatewayValues from a Terraform values map.
func decodeNatGatewayValues(tfVals map[string]interface{}) (natGatewayValues, error) {
	var v natGatewayValues
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

// newNatGateway initializes a new NatGateway from the provider
func (p *Provider) newNatGateway(vals natGatewayValues) *NatGateway {
	inst := &NatGateway{
		provider: p,

		location:               convertRegion(vals.Location),
		monthlyDataProcessedGb: vals.Usage.MonthlyDataProcessedGb,
	}
	return inst
}

func (inst *NatGateway) Components() []query.Component {
	var components []query.Component

	var monthlyDataProcessedGb decimal.Decimal
	if inst.monthlyDataProcessedGb != nil {
		monthlyDataProcessedGb = decimal.NewFromInt(*inst.monthlyDataProcessedGb)
	}

	components = append(components, inst.natGatewayCostComponent("NAT gateway"))
	components = append(components, inst.dataProcessedCostComponent("Data processed", monthlyDataProcessedGb))

	return components
}

func (inst *NatGateway) natGatewayCostComponent(name string) query.Component {
	return query.Component{

		Name:           name,
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("NAT Gateway"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", Value: util.StringPtr("Standard Gateway")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
func (inst *NatGateway) dataProcessedCostComponent(name string, monthlyDataProcessedGb decimal.Decimal) query.Component {
	return query.Component{

		Name:            name,
		Unit:            "GB",
		MonthlyQuantity: monthlyDataProcessedGb,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("NAT Gateway"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", Value: util.StringPtr("Standard Data Processed")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
