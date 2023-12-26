package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

type SourceLoadBalancer struct {
	Values struct {
		Sku      string `mapstructure:"sku"`
		Location string `mapstructure:"location"`
		SkuTier  string `mapstructure:"sku_tier"`
	} `mapstructure:"values"`
}

type LoadBalancerRule struct {
	provider *Provider

	loadBalancer SourceLoadBalancer
}

type loadBalancerRuleValues struct {
	LoadBalancer SourceLoadBalancer `mapstructure:"loadbalancer_id"`
}

// decodeLoadBalancerRuleValues decodes and returns computeInstanceValues from a Terraform values map.
func decodeLoadBalancerRuleValues(tfVals map[string]interface{}) (loadBalancerRuleValues, error) {
	var v loadBalancerRuleValues
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

func (p *Provider) newLoadBalancerRule(vals loadBalancerRuleValues) *LoadBalancerRule {
	lbRule := &LoadBalancerRule{
		provider: p,

		loadBalancer: SourceLoadBalancer{
			Values: struct {
				Sku      string `mapstructure:"sku"`
				Location string `mapstructure:"location"`
				SkuTier  string `mapstructure:"sku_tier"`
			}{
				Sku:      "Basic",
				Location: convertRegion(vals.LoadBalancer.Values.Location),
				SkuTier:  "Regional",
			},
		},
	}
	if vals.LoadBalancer.Values.Sku != "" {
		lbRule.loadBalancer.Values.Sku = vals.LoadBalancer.Values.Sku
	}
	if vals.LoadBalancer.Values.SkuTier != "" {
		lbRule.loadBalancer.Values.SkuTier = vals.LoadBalancer.Values.SkuTier
	}
	return lbRule
}

func (inst LoadBalancerRule) Components() []query.Component {
	if inst.loadBalancer.Values.Sku == "Basic" {
		return nil
	}
	costComponent := []query.Component{RegionalOverageRulesComponent(inst.provider.key, inst.loadBalancer.Values.Location, decimal.NewFromInt(1))}

	GetCostComponentNamesAndSetLogger(costComponent, inst.provider.logger)
	return costComponent
}

func RegionalIncludedRulesComponent(key, location string) query.Component {
	return query.Component{
		Name:            "Regional Included Rules",
		MonthlyQuantity: decimal.NewFromInt(1),
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

func GlobalIncludedRulesComponent(key, location string) query.Component {
	return query.Component{
		Name:            "Global Included Rules",
		MonthlyQuantity: decimal.NewFromInt(1),
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

func RegionalOverageRulesComponent(key, location string, overageRules decimal.Decimal) query.Component {
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

func GlobalOverageRulesComponent(key, location string, overageRules decimal.Decimal) query.Component {
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
