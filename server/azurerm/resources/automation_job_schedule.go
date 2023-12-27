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

type AutomationJobSchedule struct {
	provider          *Provider
	location          string
	monthlyJobRunMins *int64
}

type AutomationJobScheduleValue struct {
	ResourceGroupName ResourceGroupName `mapstructure:"resource_group_name"`
	Usage             struct {
		MonthlyJobRunMins *int64 `mapstructure:"monthly_job_run_mins"`
	} `mapstructure:"pennywise_usage"`
}

func (p *Provider) newAutomationJOBSchedule(vals AutomationJobScheduleValue) *AutomationJobSchedule {
	inst := &AutomationJobSchedule{
		provider:          p,
		location:          vals.ResourceGroupName.Values.Location,
		monthlyJobRunMins: vals.Usage.MonthlyJobRunMins,
	}
	return inst
}

func decodeAutomationJOBSchedule(tfVals map[string]interface{}) (AutomationJobScheduleValue, error) {
	var v AutomationJobScheduleValue
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

func (inst *AutomationJobSchedule) Component() []query.Component {
	location := getLocationName(inst.location)

	monthlyJobRunMins := decimal.Zero
	if inst.monthlyJobRunMins != nil {
		monthlyJobRunMins = decimal.NewFromInt(*inst.monthlyJobRunMins)
	}

	costComponents := make([]query.Component, 0)
	costComponents = append(costComponents, automationRunTimeCostComponent(location, "500", "Basic Runtime", "Basic", monthlyJobRunMins))

	GetCostComponentNamesAndSetLogger(costComponents, inst.provider.logger)
	return costComponents
}

func automationRunTimeCostComponent(location, startUsage, meterName, skuName string, monthlyQuantity decimal.Decimal) query.Component {
	return query.Component{
		Name:            "Job run time",
		Unit:            "minutes",
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
				{Key: "tier_minimum_units", Value: util.StringPtr(startUsage)}},
		},
	}
}
