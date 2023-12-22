package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

type IntegrationServiceEnvironment struct {
	location string
	skuName  string
}

type IntegrationServiceEnvironmentValue struct {
	ResourceGroupName ResourceGroupNameStruct `mapstructure:"resource_group_name"`
	SkuName           string                  `mapstructure:"sku_name"`
}

func (p *Provider) newIntegrationServiceEnvironment(vals IntegrationServiceEnvironmentValue) *IntegrationServiceEnvironment {
	inst := &IntegrationServiceEnvironment{
		location: vals.ResourceGroupName.Values.Location,
		skuName:  vals.SkuName,
	}
	return inst
}

func decodeIntegrationServiceEnvironment(tfVals map[string]interface{}) (IntegrationServiceEnvironmentValue, error) {
	var v IntegrationServiceEnvironmentValue
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

func (inst *IntegrationServiceEnvironment) Component() []query.Component {
	region := getLocationName(inst.location)

	productName := "Logic Apps Integration Service Environment"
	skuName := inst.skuName
	sku := strings.ToLower(skuName[:strings.IndexByte(skuName, '_')])
	scaleNumber, _ := strconv.Atoi(skuName[strings.IndexByte(skuName, '_')+1:])

	costComponents := make([]query.Component, 0)

	if sku == "developer" {
		productName += " - Developer"
	}

	costComponents = append(costComponents, IntegrationBaseServiceEnvironmentCostComponent("Base units", region, productName))

	if sku == "premium" && scaleNumber > 0 {
		costComponents = append(costComponents, IntegrationScaleServiceEnvironmentCostComponent("Scale units", region, productName, scaleNumber))

	}
	return costComponents
}

func IntegrationBaseServiceEnvironmentCostComponent(name, region, productName string) query.Component {
	return query.Component{
		Name:           name,
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(region),
			Service:  util.StringPtr("Logic Apps"),
			Family:   util.StringPtr("Integration"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr(productName)},
				{Key: "sku_name", Value: util.StringPtr("Base")},
				{Key: "meter_name", Value: util.StringPtr("Base Unit")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
func IntegrationScaleServiceEnvironmentCostComponent(name, region, productName string, scaleNumber int) query.Component {
	return query.Component{
		Name:           name,
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(int64(scaleNumber)),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Service:  util.StringPtr("Logic Apps"),
			Location: util.StringPtr(region),
			Family:   util.StringPtr("Integration"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr(productName)},
				{Key: "sku_name", Value: util.StringPtr("Scale")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
