package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/shopspring/decimal"
	"strings"
)

// LoadBalancer is the entity that holds the logic to calculate price
// of the azurerm_load_balancer
type LoadBalancer struct {
	provider *Provider

	location    string
	skuName     string
	rulesNumber decimal.Decimal
	skuTier     string

	// Usage
	dailyDataProceed decimal.Decimal
}

// loadBalancerValues is holds the values that we need to be able
// to calculate the price of the ComputeInstance
type loadBalancerValues struct {
	SkuName     string  `mapstructure:"sku_name"`
	Location    string  `mapstructure:"location"`
	RulesNumber float64 `mapstructure:"rules_number"`
	SkuTier     string  `mapstructure:"sku_tier"`

	Usage struct {
		DailyDataProceed float64 `mapstructure:"daily_data_proceed"`
	} `mapstructure:"tc_usage"`
}

//// decodeLoadBalancerValues decodes and returns computeInstanceValues from a Terraform values map.
//func decodeLoadBalancerValues(request api.GetAzureLoadBalancerRequest) loadBalancerValues {
//	regionCode := convertRegion(request.RegionCode)
//	rulesNumber := int32(len(request.LoadBalancer.LoadBalancer.Properties.LoadBalancingRules) +
//		len(request.LoadBalancer.LoadBalancer.Properties.OutboundRules))
//	dailyDataProceed := int64(1000)
//	if request.DailyDataProceed != nil {
//		dailyDataProceed = *request.DailyDataProceed
//	}
//	return loadBalancerValues{
//		SkuName:          string(*request.LoadBalancer.LoadBalancer.SKU.Name),
//		Location:         regionCode,
//		RulesNumber:      rulesNumber,
//		SkuTier:          string(*request.LoadBalancer.LoadBalancer.SKU.Tier),
//		DailyDataProceed: dailyDataProceed,
//	}
//}

// newManagedStorage initializes a new LoadBalancer from the provider
func (p *Provider) newLoadBalancer(vals loadBalancerValues) *LoadBalancer {
	inst := &LoadBalancer{
		provider: p,

		location:         vals.Location,
		skuName:          vals.SkuName,
		rulesNumber:      decimal.NewFromFloat(vals.RulesNumber),
		skuTier:          vals.SkuTier,
		dailyDataProceed: decimal.NewFromFloat(vals.Usage.DailyDataProceed),
	}

	return inst
}

func (inst *LoadBalancer) Components() []query.Component {
	var components []query.Component

	if inst.skuTier == "Regional" {
		components = append(components, inst.regionalIncludedRulesComponent(inst.provider.key, inst.location))
		// NAT rules are free.

		if inst.rulesNumber.InexactFloat64() > 5 {
			components = append(components, inst.regionalOverageRulesComponent(inst.provider.key, inst.location, inst.rulesNumber.Sub(decimal.NewFromInt(5))))
		}

		components = append(components, regionalDataProceedComponent(inst.provider.key, inst.location, inst.dailyDataProceed))
	} else if inst.skuTier == "Global" {
		components = append(components, inst.globalIncludedRulesComponent(inst.provider.key, inst.location))
		// NAT rules are free.

		if inst.rulesNumber.InexactFloat64() > 5 {
			components = append(components, inst.globalOverageRulesComponent(inst.provider.key, inst.location, inst.rulesNumber.Sub(decimal.NewFromInt(5))))
		}

		components = append(components, inst.globalDataProceedComponent(inst.provider.key, inst.location, inst.dailyDataProceed))
	}

	return components
}

func (inst *LoadBalancer) regionalIncludedRulesComponent(key, location string) query.Component {
	return query.Component{
		Name:           "Regional Included Rules",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Load Balancer"),
			Family:   util.StringPtr("Networking"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", Value: util.StringPtr("Standard Included LB Rules and Outbound Rules")},
			},
		},
	}
}

func (inst *LoadBalancer) globalIncludedRulesComponent(key, location string) query.Component {
	return query.Component{
		Name:           "Global Included Rules",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Load Balancer"),
			Family:   util.StringPtr("Networking"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", Value: util.StringPtr("Global Included LB Rules and Outbound Rules")},
			},
		},
	}
}

func (inst *LoadBalancer) regionalOverageRulesComponent(key, location string, overageRules decimal.Decimal) query.Component {
	return query.Component{
		Name:           "Regional Overage Rules",
		HourlyQuantity: overageRules,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Load Balancer"),
			Family:   util.StringPtr("Networking"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", Value: util.StringPtr("Standard Overage LB Rules and Outbound Rules")},
			},
		},
	}
}

func (inst *LoadBalancer) globalOverageRulesComponent(key, location string, overageRules decimal.Decimal) query.Component {
	return query.Component{
		Name:           "Global Overage Rules",
		HourlyQuantity: overageRules,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Load Balancer"),
			Family:   util.StringPtr("Networking"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", Value: util.StringPtr("Global Overage LB Rules and Outbound Rules")},
			},
		},
	}
}

func regionalDataProceedComponent(key, location string, dataProceed decimal.Decimal) query.Component {
	return query.Component{
		Name:           "Regional Data Proceed",
		HourlyQuantity: dataProceed,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Load Balancer"),
			Family:   util.StringPtr("Networking"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", Value: util.StringPtr("Standard")},
				{Key: "meter_name", ValueRegex: util.StringPtr("Data Processed")},
			},
		},
	}
}

func (inst *LoadBalancer) globalDataProceedComponent(key, location string, dataProceed decimal.Decimal) query.Component {
	return query.Component{
		Name:           "Global Data Proceed",
		HourlyQuantity: dataProceed,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Load Balancer"),
			Family:   util.StringPtr("Networking"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", Value: util.StringPtr("Global")},
				{Key: "meter_name", ValueRegex: util.StringPtr("Data Processed")},
			},
		},
	}
}

func convertRegion(region string) string {
	if strings.Contains(strings.ToLower(region), "usgov") {
		return "US Gov"
	} else if strings.Contains(strings.ToLower(region), "china") {
		return "Сhina"
	} else {
		return "Global"
	}
}
