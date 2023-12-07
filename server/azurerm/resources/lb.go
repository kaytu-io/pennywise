package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

// LoadBalancer is the entity that holds the logic to calculate price
// of the azurerm_load_balancer
type LoadBalancer struct {
	provider *Provider

	location    string
	sku         string
	rulesNumber decimal.Decimal
	skuTier     string

	// Usage
	dailyDataProceed decimal.Decimal
}

// loadBalancerValues is holds the values that we need to be able
// to calculate the price of the ComputeInstance
type loadBalancerValues struct {
	Sku         string  `mapstructure:"sku"`
	Location    string  `mapstructure:"location"`
	RulesNumber float64 `mapstructure:"rules_number"`
	SkuTier     string  `mapstructure:"sku_tier"`

	Usage struct {
		DailyDataProceed float64 `mapstructure:"daily_data_proceed"`
	} `mapstructure:"tc_usage"`
}

// decodeLoadBalancerValues decodes and returns computeInstanceValues from a Terraform values map.
func decodeLoadBalancerValues(tfVals map[string]interface{}) (loadBalancerValues, error) {
	fmt.Println("TF VALS", tfVals)
	var v loadBalancerValues
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

// newManagedStorage initializes a new LoadBalancer from the provider
func (p *Provider) newLoadBalancer(vals loadBalancerValues) *LoadBalancer {
	inst := &LoadBalancer{
		provider: p,

		location:         convertRegion(vals.Location),
		sku:              "Basic",
		rulesNumber:      decimal.NewFromFloat(vals.RulesNumber),
		skuTier:          "Regional",
		dailyDataProceed: decimal.NewFromFloat(vals.Usage.DailyDataProceed),
	}
	if vals.Sku != "" {
		inst.sku = vals.Sku
	}
	if vals.SkuTier != "" {
		inst.skuTier = vals.SkuTier
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

func regionalDataProceedComponent(key, location string, dataProceed decimal.Decimal) query.Component {
	return query.Component{
		Name:            "Regional Data Proceed",
		Unit:            "GB",
		MonthlyQuantity: dataProceed,
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
		Name:            "Global Data Proceed",
		Unit:            "GB",
		MonthlyQuantity: dataProceed,
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
