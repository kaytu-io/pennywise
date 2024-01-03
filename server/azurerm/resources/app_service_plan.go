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

// AppServicePlan is the entity that holds the logic to calculate price
// of the azurerm_app_service_plan
type AppServicePlan struct {
	provider *Provider

	location    string
	skuSize     string
	skuCapacity int64
	kind        string
}

// appServicePlanValues is holds the values that we need to be able
// to calculate the price of the AppServicePlan
type appServicePlanValues struct {
	Location string `mapstructure:"location"`
	Kind     string `mapstructure:"kind"`
	Sku      []struct {
		Size     string `mapstructure:"size"`
		Capacity int64  `mapstructure:"capacity"`
	} `mapstructure:"sku"`
}

// decodeAppServicePlanValues decodes and returns appServicePlanValues from a Terraform values map.
func decodeAppServicePlanValues(tfVals map[string]interface{}) (appServicePlanValues, error) {
	var v appServicePlanValues
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
func (p *Provider) newAppServicePlan(vals appServicePlanValues) *AppServicePlan {
	inst := &AppServicePlan{
		provider: p,

		location:    vals.Location,
		skuSize:     vals.Sku[0].Size,
		skuCapacity: vals.Sku[0].Capacity,
		kind:        vals.Kind,
	}
	return inst
}

func (inst *AppServicePlan) Components() []query.Component {
	var components []query.Component

	sku := ""
	os := "windows"
	var capacity int64 = 1
	if inst.skuCapacity > 0 {
		capacity = inst.skuCapacity
	}
	productName := "Standard Plan"

	if len(inst.skuSize) < 2 || strings.ToLower(inst.skuSize[:2]) == "ep" || strings.ToLower(inst.skuSize[:2]) == "y1" || strings.ToLower(inst.skuSize[:2]) == "ws" {
		return components
	}

	var additionalAttributeFilters []*product.AttributeFilter

	switch strings.ToLower(inst.skuSize[:1]) {
	case "s":
		sku = "S" + inst.skuSize[1:]
	case "b":
		sku = "B" + inst.skuSize[1:]
		productName = "Basic Plan"
	case "p", "i":
		sku, productName, additionalAttributeFilters = getVersionedAppServicePlanSKU(inst.skuSize, os)
	}

	switch strings.ToLower(inst.skuSize[:2]) {
	case "pc":
		sku = "PC" + inst.skuSize[2:]
		productName = "Premium Windows Container Plan"
	case "y1":
		sku = "Shared"
		productName = "Shared Plan"
	}

	if inst.kind != "" {
		os = strings.ToLower(inst.kind)
	}
	if os == "app" {
		os = "windows"
	}
	if os != "windows" && productName != "Premium Plan" && productName != "Isolated Plan" {
		productName += " - Linux"
	}

	components = append(components, servicePlanCostComponent(
		inst.location,
		fmt.Sprintf("Instance usage (%s)", inst.skuSize),
		productName,
		sku,
		capacity,
		additionalAttributeFilters...,
	))

	return components
}

func servicePlanCostComponent(region, name, productName, skuRefactor string, capacity int64, additionalAttributeFilters ...*product.AttributeFilter) query.Component {
	return query.Component{
		Name:           name,
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(capacity),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(region),
			Service:  util.StringPtr("Azure App Service"),
			Family:   util.StringPtr("Compute"),
			AttributeFilters: append([]*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Azure App Service " + productName)},
				{Key: "sku_name", Value: util.StringPtr(skuRefactor)},
			}, additionalAttributeFilters...),
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func getVersionedAppServicePlanSKU(skuName, os string) (string, string, []*product.AttributeFilter) {
	tier := "Premium"
	if strings.ToLower(skuName[:1]) == "i" {
		tier = "Isolated"
	}

	version := strings.ToLower(skuName[2:])
	if version == "v1" {
		version = ""
	}
	var formattedSku string
	if strings.Contains(strings.ToLower(version), "m") || strings.Contains(strings.ToLower(skuName[:2]), "0") {
		formattedSku = strings.TrimSpace(skuName[:2] + version)
	} else {
		formattedSku = strings.TrimSpace(skuName[:2] + " " + version)
	}

	productVersion := version
	if len(version) > 0 && version[0] == 'm' {
		productVersion = version[1:]
	}
	productName := strings.ReplaceAll(tier+" "+productVersion+" Plan", "  ", " ")

	//if productVersion == "v3" && os == "linux" {
	//	return formattedSku, productName, []*product.AttributeFilter{
	//		{
	//			Key:        "arm_sku_name",
	//			ValueRegex: util.StringPtr(fmt.Sprintf("%s$", strings.ToLower(strings.ReplaceAll(formattedSku, " ", "_")))),
	//		},
	//	}
	//}

	return formattedSku, productName, nil
}
