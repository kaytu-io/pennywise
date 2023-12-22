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

// ApiManagement is the entity that holds the logic to calculate price
// of the azurerm_api_management
type ApiManagement struct {
	provider *Provider

	location string
	skuName  string

	// Usage
	selfHostedGatewayCount *int64
	monthlyApiCalls        *int64
}

// apiManagementValues is holds the values that we need to be able
// to calculate the price of the ApiManagement
type apiManagementValues struct {
	Location string `mapstructure:"location"`
	SkuName  string `mapstructure:"sku_name"`

	Usage struct {
		SelfHostedGatewayCount *int64 `mapstructure:"self_hosted_gateway_count"`
		MonthlyApiCalls        *int64 `mapstructure:"monthly_api_calls"`
	} `mapstructure:"pennywise_usage"`
}

// decodeApiManagementValues decodes and returns apiManagementValues from a Terraform values map.
func decodeApiManagementValues(tfVals map[string]interface{}) (apiManagementValues, error) {
	var v apiManagementValues
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

// newAppServiceCertificateBinding initializes a new AppServiceCertificateBinding from the provider
func (p *Provider) newApiManagement(vals apiManagementValues) *ApiManagement {
	inst := &ApiManagement{
		provider: p,

		location: vals.Location,
		skuName:  vals.SkuName,

		selfHostedGatewayCount: vals.Usage.SelfHostedGatewayCount,
		monthlyApiCalls:        vals.Usage.MonthlyApiCalls,
	}
	return inst
}

func (inst *ApiManagement) Components() []query.Component {
	var components []query.Component

	var tier string
	var capacity decimal.Decimal
	if s := strings.Split(inst.skuName, "_"); len(s) == 2 {
		tier = strings.ToLower(s[0])
		capacity, _ = decimal.NewFromString(s[1])
	}

	if tier != "consumption" {
		components = append(components, inst.apiManagementCostComponent(
			fmt.Sprintf("API management (%s)", tier),
			"units",
			tier,
			capacity))

	} else {
		var apiCalls *decimal.Decimal
		if inst.monthlyApiCalls != nil {
			apiCalls = decimalPtr(decimal.NewFromInt(*inst.monthlyApiCalls))
		}

		if apiCalls != nil {
			apiCalls = decimalPtr(apiCalls.Div(decimal.NewFromInt(10000)))
			components = append(components, inst.consumptionAPICostComponent(tier, *apiCalls))
		} else {
			components = append(components, inst.consumptionAPICostComponent(tier, decimal.Zero))
		}
	}

	if tier == "premium" {
		var selfHostedGateways decimal.Decimal
		if inst.selfHostedGatewayCount != nil {
			selfHostedGateways = decimal.NewFromInt(*inst.selfHostedGatewayCount)
		}
		components = append(components, inst.apiManagementCostComponent(
			"Self hosted gateway",
			"gateways",
			"Gateway",
			selfHostedGateways,
		))
	}

	return components
}

func (inst *ApiManagement) apiManagementCostComponent(name, unit, tier string, quantity decimal.Decimal) query.Component {
	return query.Component{
		Name:           name,
		Unit:           unit,
		HourlyQuantity: quantity,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("API Management"),
			Family:   util.StringPtr("Developer Tools"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", ValueRegex: util.StringPtr(fmt.Sprintf("/^%s$/i", tier))},
				{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf("/^%s unit$/i", tier))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *ApiManagement) consumptionAPICostComponent(tier string, quantity decimal.Decimal) query.Component {
	return query.Component{
		Name:            "API management (consumption)",
		Unit:            "1M calls",
		MonthlyQuantity: quantity,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("API Management"),
			Family:   util.StringPtr("Developer Tools"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", ValueRegex: util.StringPtr(fmt.Sprintf("/^%s$/i", tier))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr("100")},
			},
		},
	}
}
