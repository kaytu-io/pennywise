package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

// PublicIP is the entity that holds the logic to calculate price
// of the azurerm_public_ip
type PublicIP struct {
	provider *Provider

	location         string
	allocationMethod string
	sku              *string
}

// publicIPValues is holds the values that we need to be able
// to calculate the price of the PublicIP
type publicIPValues struct {
	Location         string  `mapstructure:"location"`
	AllocationMethod string  `mapstructure:"allocation_method"`
	Sku              *string `mapstructure:"sku"`
}

// decodePublicIPValues decodes and returns publicIPValues from a Terraform values map.
func decodePublicIPValues(tfVals map[string]interface{}) (publicIPValues, error) {
	var v publicIPValues
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

// newPublicIP initializes a new PublicIP from the provider
func (p *Provider) newPublicIP(vals publicIPValues) *PublicIP {
	inst := &PublicIP{
		provider: p,

		location:         getLocationName(vals.Location),
		allocationMethod: vals.AllocationMethod,
		sku:              vals.Sku,
	}
	return inst
}

func (inst *PublicIP) Components() []query.Component {
	var components []query.Component

	var meterName string
	sku := "Basic"

	if inst.sku != nil {
		sku = *inst.sku
	}

	switch sku {
	case "Basic":
		meterName = "Basic IPv4 " + inst.allocationMethod + " Public IP"
	case "Standard":
		meterName = "Standard IPv4 " + inst.allocationMethod + " Public IP"
	}

	components = append(components, inst.publicIPCostComponent(fmt.Sprintf("IP address (%s)", strings.ToLower(inst.allocationMethod)), sku, meterName))

	return components
}

func (inst *PublicIP) publicIPCostComponent(name, sku, meterName string) query.Component {
	return query.Component{
		Name:           name,
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Virtual Network"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("IP Addresses")},
				{Key: "sku_name", Value: util.StringPtr(sku)},
				{Key: "meter_name", Value: util.StringPtr(meterName)},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
