package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

type AutomationAccount struct {
	provider                *Provider
	location                string
	monthlyJobRunMins       *int64
	nonAzureConfigNodeCount *int64
	monthlyWatcherHrs       *int64
}

type AutomationAccountValue struct {
	Location string `mapstructure:"location"`
	Usage    struct {
		MonthlyJobRunMins       *int64 `mapstructure:"monthly_job_run_mins"`
		NonAzureConfigNodeCount *int64 `mapstructure:"non_azure_config_node_count"`
		MonthlyWatcherHrs       *int64 `mapstructure:"monthly_watcher_hrs"`
	} `mapstructure:"pennywise_usage"`
}

func (p *Provider) newAutomationAccount(vals AutomationAccountValue) *AutomationAccount {
	inst := &AutomationAccount{
		provider:                p,
		location:                vals.Location,
		nonAzureConfigNodeCount: vals.Usage.NonAzureConfigNodeCount,
		monthlyJobRunMins:       vals.Usage.MonthlyJobRunMins,
		monthlyWatcherHrs:       vals.Usage.MonthlyWatcherHrs,
	}
	return inst
}

func decodeAutomationAccount(tfVals map[string]interface{}) (AutomationAccountValue, error) {
	var v AutomationAccountValue
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

func (inst *AutomationAccount) Component() []query.Component {
	var monthlyJobRunMins, nonAzureConfigNodeCount decimal.Decimal
	location := inst.location

	costComponents := make([]query.Component, 0)

	if inst.monthlyJobRunMins != nil {
		monthlyJobRunMins = decimal.NewFromInt(*inst.monthlyWatcherHrs)
		if monthlyJobRunMins.IsPositive() {
			costComponents = append(costComponents, automationRunTimeCostComponent(location, "500", "Basic Runtime", "Basic", monthlyJobRunMins))
		}
	} else {
		costComponents = append(costComponents, automationRunTimeCostComponent(location, "500", "Basic Runtime", "Basic", monthlyJobRunMins))
	}

	if inst.nonAzureConfigNodeCount != nil {
		nonAzureConfigNodeCount = decimal.NewFromInt(*inst.nonAzureConfigNodeCount)
		if nonAzureConfigNodeCount.IsPositive() {
			costComponents = append(costComponents, nonautomationDSCNodesCostComponent(location, "5", "Non-Azure Node", "Non-Azure", nonAzureConfigNodeCount))
		}
	} else {
		costComponents = append(costComponents, nonautomationDSCNodesCostComponent(location, "5", "Non-Azure Node", "Non-Azure", nonAzureConfigNodeCount))
	}

	costComponents = append(costComponents, watchersCostComponent(inst.provider.key, inst.monthlyWatcherHrs, location, "744", "Watcher", "Basic"))

	GetCostComponentNamesAndSetLogger(costComponents, inst.provider.logger)
	return costComponents
}

func watchersCostComponent(key string, monthlyWatcherHrs *int64, location string, startUsage, meterName, skuName string) query.Component {
	var monthlyQuantity decimal.Decimal
	if monthlyWatcherHrs != nil {
		monthlyQuantity = decimal.NewFromInt(*monthlyWatcherHrs)
	}
	return query.Component{
		Name:            "Watchers",
		Unit:            "hours",
		MonthlyQuantity: monthlyQuantity,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(location),
			Service:  util.StringPtr("Automation"),
			Family:   util.StringPtr("Management and Governance"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", Value: util.StringPtr(meterName)},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
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
