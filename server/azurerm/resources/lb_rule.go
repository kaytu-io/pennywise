package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/shopspring/decimal"
)

func (inst *LoadBalancer) regionalIncludedRulesComponent(key, location string) query.Component {
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

func (inst *LoadBalancer) globalIncludedRulesComponent(key, location string) query.Component {
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
