package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

// PublicIPPrefix is the entity that holds the logic to calculate price
// of the azurerm_public_ip
type PublicIPPrefix struct {
	provider *Provider

	location string
}

// publicIPPrefixValues is holds the values that we need to be able
// to calculate the price of the PublicIPPrefix
type publicIPPrefixValues struct {
	Location string `mapstructure:"location"`
}

// decodePublicIPPrefixValues decodes and returns publicIPPrefixValues from a Terraform values map.
func decodePublicIPPrefixValues(tfVals map[string]interface{}) (publicIPPrefixValues, error) {
	var v publicIPPrefixValues
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

// newPublicIPPrefix initializes a new PublicIPPrefix from the provider
func (p *Provider) newPublicIPPrefix(vals publicIPPrefixValues) *PublicIPPrefix {
	inst := &PublicIPPrefix{
		provider: p,

		location: convertRegion(vals.Location),
	}
	return inst
}

func (inst *PublicIPPrefix) Components() []query.Component {
	var components []query.Component

	components = append(components, inst.publicIPPrefixCostComponent("IP prefix"))

	return components
}

func (inst *PublicIPPrefix) publicIPPrefixCostComponent(name string) query.Component {
	return query.Component{
		Name:           name,
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Virtual Network"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Public IP Prefix")},
				{Key: "meter_name", Value: util.StringPtr("Standard Static IP Addresses")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
