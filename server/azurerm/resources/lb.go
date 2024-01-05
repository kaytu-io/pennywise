package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/kaytu-io/pennywise/server/resource"
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
	// receives monthly inbound and outbound data processed in GB
	monthlyDataProceed decimal.Decimal
}

// loadBalancerValues is holds the values that we need to be able
// to calculate the price of the ComputeInstance
type loadBalancerValues struct {
	Sku         string  `mapstructure:"sku"`
	Location    string  `mapstructure:"location"`
	RulesNumber float64 `mapstructure:"rules_number"`
	SkuTier     string  `mapstructure:"sku_tier"`

	Usage struct {
		MonthlyDataProceed float64 `mapstructure:"monthly_data_processed_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodeLoadBalancerValues decodes and returns loadBalancerValues from a Terraform values map.
func decodeLoadBalancerValues(tfVals map[string]interface{}) (loadBalancerValues, error) {
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

// newLoadBalancer initializes a new LoadBalancer from the provider
func (p *Provider) newLoadBalancer(vals loadBalancerValues) *LoadBalancer {
	inst := &LoadBalancer{
		provider: p,

		location:           convertRegion(vals.Location),
		sku:                "Basic",
		rulesNumber:        decimal.NewFromFloat(vals.RulesNumber),
		skuTier:            "Regional",
		monthlyDataProceed: decimal.NewFromFloat(vals.Usage.MonthlyDataProceed),
	}
	if vals.Sku != "" {
		inst.sku = vals.Sku
	}
	if vals.SkuTier != "" {
		inst.skuTier = vals.SkuTier
	}

	return inst
}

func (inst *LoadBalancer) Components() []resource.Component {
	var components []resource.Component

	if inst.sku == "Basic" {
		return nil
	}

	if inst.skuTier == "Regional" {
		// TODO: Check include rules. (They were calculated in azure web page but not in infracost)
		//components = append(components, RegionalIncludedRulesComponent(inst.provider.key, inst.location))
		// NAT rules are free.

		if inst.rulesNumber.InexactFloat64() > 5 {
			components = append(components, RegionalOverageRulesComponent(inst.provider.key, inst.location, inst.rulesNumber.Sub(decimal.NewFromInt(5))))
		}

		components = append(components, regionalDataProceedComponent(inst.provider.key, inst.location, inst.monthlyDataProceed))
	} else if inst.skuTier == "Global" {
		// TODO: Check include rules. (They were calculated in azure web page but not in infracost)
		//components = append(components, GlobalIncludedRulesComponent(inst.provider.key, inst.location))
		// NAT rules are free.

		if inst.rulesNumber.InexactFloat64() > 5 {
			components = append(components, GlobalOverageRulesComponent(inst.provider.key, inst.location, inst.rulesNumber.Sub(decimal.NewFromInt(5))))
		}

		components = append(components, inst.globalDataProceedComponent(inst.provider.key, inst.location, inst.monthlyDataProceed))
	}
	GetCostComponentNamesAndSetLogger(components, inst.provider.logger)

	return components
}

func regionalDataProceedComponent(key, location string, dataProceed decimal.Decimal) resource.Component {
	return resource.Component{
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
				{Key: "meter_name", Value: util.StringPtr("Standard Data Processed")},
			},
		},
	}
}

func (inst *LoadBalancer) globalDataProceedComponent(key, location string, dataProceed decimal.Decimal) resource.Component {
	return resource.Component{
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
				{Key: "meter_name", ValueRegex: util.StringPtr("Global Data Processed")},
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
