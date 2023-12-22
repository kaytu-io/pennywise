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

// AppServiceEnvironment is the entity that holds the logic to calculate price
// of the azurerm_app_service_environment
type AppServiceEnvironment struct {
	provider *Provider

	location    string
	pricingTier *string

	// Usage
	operatingSystem *string
}

// appServiceEnvironmentValues is holds the values that we need to be able
// to calculate the price of the AppServiceEnvironment
type appServiceEnvironmentValues struct {
	ResourceGroupName ResourceGroupName `mapstructure:"resource_group_name"`
	PricingTier       *string           `mapstructure:"pricing_tier"`

	Usage struct {
		OperatingSystem *string `mapstructure:"operating_system"`
	} `mapstructure:"usage"`
}

// decodeAppServiceEnvironmentValues decodes and returns appServiceEnvironmentValues from a Terraform values map.
func decodeAppServiceEnvironmentValues(tfVals map[string]interface{}) (appServiceEnvironmentValues, error) {
	var v appServiceEnvironmentValues
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

// newAppServiceEnvironment initializes a new AppServiceEnvironment from the provider
func (p *Provider) newAppServiceEnvironment(vals appServiceEnvironmentValues) *AppServiceEnvironment {
	inst := &AppServiceEnvironment{
		provider: p,

		location:    vals.ResourceGroupName.Values.Location,
		pricingTier: vals.PricingTier,

		operatingSystem: vals.Usage.OperatingSystem,
	}
	return inst
}

func (inst *AppServiceEnvironment) Components() []query.Component {
	var components []query.Component

	tier := "I1"
	if inst.pricingTier != nil {
		tier = *inst.pricingTier
	}

	stampFeeTiers := []string{"I1", "I2", "I3"}
	productName := "Isolated Plan"
	os := "linux"
	if inst.operatingSystem != nil {
		os = strings.ToLower(*inst.operatingSystem)
	}
	if os == "linux" {
		productName += " - Linux"
	}
	if contains(stampFeeTiers, tier) == true {
		components = append(components, inst.appIsolatedServicePlanCostComponentStampFee(productName))
	}
	components = append(components, inst.appIsolatedServicePlanCostComponent(fmt.Sprintf("Instance usage (%s)", tier), productName, tier))

	return components
}

func (inst *AppServiceEnvironment) appIsolatedServicePlanCostComponentStampFee(productName string) query.Component {
	return query.Component{
		Name:           "Stamp fee",
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Azure App Service"),
			Family:   util.StringPtr("Compute"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Azure App Service " + productName)},
				{Key: "sku_name", Value: util.StringPtr("Stamp")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
func (inst *AppServiceEnvironment) appIsolatedServicePlanCostComponent(name, productName, tier string) query.Component {
	return query.Component{
		Name:           name,
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Azure App Service"),
			Family:   util.StringPtr("Compute"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Azure App Service " + productName)},
				{Key: "sku_name", Value: util.StringPtr(tier)},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
