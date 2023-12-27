package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

type AutomationDSCConfigurationValue struct {
	Location string `mapstructure:"location"`
	Usage    struct {
		NonAzureConfigNodeCount *int64 `mapstructure:"non_azure_config_node_count"`
	} `mapstructure:"pennywise_usage"`
}

type AutomationDSCConfiguration struct {
	provider                *Provider
	location                string
	nonAzureConfigNodeCount *int64
}

func (p *Provider) newAutomationDSCConfiguration(vals AutomationDSCConfigurationValue) *AutomationDSCConfiguration {
	inst := &AutomationDSCConfiguration{
		provider:                p,
		location:                vals.Location,
		nonAzureConfigNodeCount: vals.Usage.NonAzureConfigNodeCount,
	}
	return inst
}

func decodeAutomationDNSConfiguration(tfVals map[string]interface{}) (AutomationDSCConfigurationValue, error) {
	var v AutomationDSCConfigurationValue
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

func (inst *AutomationDSCConfiguration) Component() []query.Component {
	costComponent := automationDSCNodesCostComponent(inst.location, inst.nonAzureConfigNodeCount)

	GetCostComponentNamesAndSetLogger(costComponent, inst.provider.logger)
	return costComponent
}

func automationDSCNodesCostComponent(location string, nonAzureConfigNodeCount *int64) []query.Component {
	region := getLocationName(location)
	var nonAzureConfigNodeCountDec decimal.Decimal

	if nonAzureConfigNodeCount != nil {
		nonAzureConfigNodeCountDec = decimal.NewFromInt(*nonAzureConfigNodeCount)
	}
	costComponents := make([]query.Component, 0)
	costComponents = append(costComponents, nonautomationDSCNodesCostComponent(region, "5", "Non-Azure Node", "Non-Azure", nonAzureConfigNodeCountDec))
	return costComponents
}

func nonautomationDSCNodesCostComponent(location, startUsage, meterName, skuName string, monthlyQuantity decimal.Decimal) query.Component {
	return query.Component{
		Name:            "Non-azure config nodes",
		Unit:            "nodes",
		MonthlyQuantity: monthlyQuantity,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(location),
			Service:  util.StringPtr("Automation"),
			Family:   util.StringPtr("Management and Governance"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", Value: util.StringPtr(fmt.Sprintf("%s", meterName))},
				{Key: "sku_name", Value: util.StringPtr(fmt.Sprintf("%s", skuName))},
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
