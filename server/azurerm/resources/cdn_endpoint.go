package resources

import (
	"fmt"
	"github.com/kaytu-io/infracost/external/usage"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

type ProfileNameStruct struct {
	Values struct {
		SKU *string `mapstructure:"sku"`
	} `json:"mapstructure"`
}

type AzureRMCDNEndpoint struct {
	provider *Provider

	location           string
	globalDeliveryRule []map[string]interface{}
	deliveryRule       []map[string]interface{}
	sku                *string
	optimizationType   *string

	//usage
	monthlyOutboundGB          *int64
	monthlyRulesEngineRequests *int64
}

type AzureRMCDNEndpointValue struct {
	Location           string                   `mapstructure:"location"`
	GlobalDeliveryRule []map[string]interface{} `mapstructure:"global_delivery_rule"`
	DeliveryRule       []map[string]interface{} `mapstructure:"delivery_rule"`
	ProfileName        []ProfileNameStruct      `mapstructure:"profile_name"`
	OptimizationType   *string                  `mapstructure:"optimization_type"`

	Usage struct {
		MonthlyOutboundGB          *int64 `mapstructure:"monthly_outbound_gb"`
		MonthlyRulesEngineRequests *int64 `mapstructure:"monthly_rules_engine_requests"`
	} `mapstructure:"pennywise_usage"`
}

func (p *Provider) newCDNEndpoint(vals AzureRMCDNEndpointValue) *AzureRMCDNEndpoint {
	var sku *string
	if len(vals.ProfileName) > 0 {
		sku = vals.ProfileName[0].Values.SKU
	}
	inst := &AzureRMCDNEndpoint{
		provider:                   p,
		location:                   vals.Location,
		globalDeliveryRule:         vals.GlobalDeliveryRule,
		deliveryRule:               vals.DeliveryRule,
		sku:                        sku,
		optimizationType:           vals.OptimizationType,
		monthlyOutboundGB:          vals.Usage.MonthlyOutboundGB,
		monthlyRulesEngineRequests: vals.Usage.MonthlyRulesEngineRequests,
	}
	return inst
}

func decodeCDNEndpoint(tfVals map[string]interface{}) (AzureRMCDNEndpointValue, error) {
	var v AzureRMCDNEndpointValue
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

func (inst AzureRMCDNEndpoint) Component() []query.Component {
	region := getLocationName(inst.location)
	region = regionToZone(region)

	var costComponents []query.Component
	sku := ""
	if inst.sku != nil {
		fmt.Printf("sku : %v \n ", *inst.sku)
		sku = *inst.sku
	}

	if len(strings.Split(sku, "_")) != 2 || strings.ToLower(sku) == "standard_chinacdn" {
		fmt.Printf("Unrecognized/unsupported CDN sku format for resource %s: %s", "CDN", sku)
		return nil
	}

	costComponents = append(costComponents, cdnOutboundDataCostComponents(inst.provider.key, region, sku, inst.monthlyOutboundGB)...)

	if strings.ToLower(sku) == "standard_microsoft" {
		numberOfRules := 0
		if inst.globalDeliveryRule != nil {
			numberOfRules += len(inst.globalDeliveryRule)
		}
		if inst.deliveryRule != nil {
			numberOfRules += len(inst.deliveryRule)
		}

		if numberOfRules > 5 {
			numberOfRules -= 5

			costComponents = append(costComponents, cdnCostComponent(
				inst.provider.key,
				"Rules engine rules (over 5)",
				"rules",
				region,
				"Azure CDN from Microsoft",
				"Standard",
				"Rule",
				"5",
				decimal.NewFromInt(int64(numberOfRules)),
			))
		}

		if numberOfRules > 0 {
			var rulesRequests decimal.Decimal
			if inst.monthlyRulesEngineRequests != nil {
				rulesRequests = decimal.NewFromInt(*inst.monthlyRulesEngineRequests / 1000000)
			}
			costComponents = append(costComponents, cdnCostComponent(
				inst.provider.key,
				"Rules engine requests",
				"1M requests",
				region,
				"Azure CDN from Microsoft",
				"Standard",
				"Requests",
				"0",
				rulesRequests,
			))
		}
	}

	if strings.ToLower(sku) == "standard_akamai" || strings.ToLower(sku) == "standard_verizon" {
		if inst.optimizationType != nil {
			if strings.ToLower(*inst.optimizationType) == "dynamicsiteacceleration" {
				costComponents = append(costComponents, cdnAccelerationDataTransfersCostComponents(inst.provider.key, region, sku, inst.monthlyOutboundGB)...)
			}
		}
	}

	return costComponents
}

func cdnOutboundDataCostComponents(key string, region, sku string, monthlyOutboundGBInput *int64) []query.Component {
	var costComponents []query.Component

	type dataTier struct {
		name       string
		startUsage string
	}
	var name, productName, skuName, meterName string
	if s := strings.Split(sku, "_"); len(s) == 2 {
		productName = fmt.Sprintf("Azure CDN from %s", s[1])
		skuName = s[0]
		if strings.ToLower(s[1]) == "verizon" {
			name = fmt.Sprintf("Outbound data transfer (%s, ", s[0]+" "+s[1])
		} else {
			name = fmt.Sprintf("Outbound data transfer (%s, ", s[1])
		}
	}

	data := []dataTier{
		{name: fmt.Sprintf("%s%s", name, "first 10TB)"), startUsage: "0"},
		{name: fmt.Sprintf("%s%s", name, "next 40TB)"), startUsage: "10000"},
		{name: fmt.Sprintf("%s%s", name, "next 100TB)"), startUsage: "50000"},
		{name: fmt.Sprintf("%s%s", name, "next 350TB)"), startUsage: "150000"},
		{name: fmt.Sprintf("%s%s", name, "next 500TB)"), startUsage: "500000"},
		{name: fmt.Sprintf("%s%s", name, "next 4000TB)"), startUsage: "1e+06"},
		{name: fmt.Sprintf("%s%s", name, "over 5000TB)"), startUsage: "5e+06"},
	}
	meterName = fmt.Sprintf("%s Data Transfer", skuName)

	var monthlyOutboundGb *decimal.Decimal
	if monthlyOutboundGBInput != nil {
		monthlyOutboundGb = decimalPtr(decimal.NewFromInt(*monthlyOutboundGBInput))
		tierLimits := []int{10000, 40000, 100000, 350000, 500000, 4000000}
		tiers := usage.CalculateTierBuckets(*monthlyOutboundGb, tierLimits)

		for i, d := range data {
			if tiers[i].GreaterThan(decimal.Zero) {
				costComponents = append(costComponents, cdnCostComponent(
					key,
					d.name,
					"GB",
					region,
					productName,
					skuName,
					meterName,
					d.startUsage,
					tiers[i]))
			}
		}
	} else {
		costComponents = append(costComponents, cdnCostComponent(
			key,
			data[0].name,
			"GB",
			region,
			productName,
			skuName,
			meterName,
			data[0].startUsage,
			decimal.Zero))
	}
	//for _, v := range costComponents {
	//	fmt.Printf("monthly : %v \n", v.MonthlyQuantity)
	//
	//	fmt.Printf("location : %v \n", *v.ProductFilter.Location)
	//	for _, v1 := range v.PriceFilter.AttributeFilters {
	//		fmt.Printf("start usage : %v \n ", *v1.Value)
	//	}
	//
	//	for k1, v1 := range v.ProductFilter.AttributeFilters {
	//		fmt.Printf("key : %v , value : key : %v /value : %v \n", k1, v1.Key, *v1.ValueRegex)
	//	}
	//	fmt.Printf("\n\n")
	//}
	return costComponents
}

func cdnAccelerationDataTransfersCostComponents(key, region, sku string, monthlyOutboundGBInput *int64) []query.Component {
	var costComponents []query.Component

	type dataTier struct {
		name       string
		startUsage string
	}

	name := "Acceleration outbound data transfer "

	data := []dataTier{
		{name: fmt.Sprintf("%s%s", name, "(first 50TB)"), startUsage: "0"},
		{name: fmt.Sprintf("%s%s", name, "(next 100TB)"), startUsage: "50000"},
		{name: fmt.Sprintf("%s%s", name, "(next 350TB)"), startUsage: "150000"},
		{name: fmt.Sprintf("%s%s", name, "(next 500TB)"), startUsage: "500000"},
		{name: fmt.Sprintf("%s%s", name, "(over 1000TB)"), startUsage: "1e+06"},
	}

	var productName, skuName, meterName string
	if s := strings.Split(sku, "_"); len(s) == 2 {
		productName = fmt.Sprintf("Azure CDN from %s", s[1])
		skuName = s[0]
	}
	meterName = "Standard Acceleration Data Transfer"

	var monthlyOutboundGb *decimal.Decimal
	if monthlyOutboundGBInput != nil {
		monthlyOutboundGb = decimalPtr(decimal.NewFromInt(*monthlyOutboundGBInput))
		tierLimits := []int{50000, 100000, 350000, 500000, 1000000}
		tiers := usage.CalculateTierBuckets(*monthlyOutboundGb, tierLimits)

		for i, d := range data {
			if tiers[i].GreaterThan(decimal.Zero) {
				costComponents = append(costComponents, cdnCostComponent(
					key,
					d.name,
					"GB",
					region,
					productName,
					skuName,
					meterName,
					d.startUsage,
					tiers[i]))
			}
		}
	} else {
		costComponents = append(costComponents, cdnCostComponent(
			key,
			data[0].name,
			"GB",
			region,
			productName,
			skuName,
			meterName,
			data[0].startUsage,
			decimal.Zero))
	}

	return costComponents
}

func cdnCostComponent(key, name, unit, region, productName, skuName, meterName, startUsage string, quantity decimal.Decimal) query.Component {
	return query.Component{
		Name:            name,
		Unit:            unit,
		MonthlyQuantity: quantity,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Location: util.StringPtr(region),
			Service:  util.StringPtr("Content Delivery Network"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(fmt.Sprintf("^%s$", productName))},
				{Key: "sku_name", ValueRegex: util.StringPtr(fmt.Sprintf("^%s$", skuName))},
				{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf("%s$", meterName))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr(startUsage)},
			},
		},
	}
}
