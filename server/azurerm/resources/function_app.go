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

// FunctionApp is the entity that holds the logic to calculate price
// of the azurerm_function_app
type FunctionApp struct {
	provider *Provider

	location string
	skuName  string
	tier     string
	osType   string

	// Usage
	monthlyExecutions   *int64
	executionDurationMs *int64
	memoryMb            *int64
	instances           *int64
}

type AppService struct {
	Values struct {
		Kind string `mapstructure:"kind"`
		Sku  []struct {
			Tier string `mapstructure:"tier"`
			Size string `mapstructure:"size"`
		} `mapstructure:"sku"`
		SkuName string `mapstructure:"sku_name"`
		OsType  string `mapstructure:"os_type"`
	} `mapstructure:"values"`
}

// functionAppValues is holds the values that we need to be able
// to calculate the price of the FunctionApp
type functionAppValues struct {
	Location         string       `mapstructure:"location"`
	AppServicePlanId []AppService `mapstructure:"app_service_plan_id"`
	ServicePlanId    []AppService `mapstructure:"service_plan_id"`

	Usage struct {
		MonthlyExecutions   *int64 `mapstructure:"monthly_executions"`
		ExecutionDurationMs *int64 `mapstructure:"execution_duration_ms"`
		MemoryMb            *int64 `mapstructure:"memory_mb"`
		Instances           *int64 `mapstructure:"instances"`
	} `mapstructure:"pennywise_usage"`
}

var (
	functionAppSkuMapCPU = map[string]int64{
		"ep1": 1,
		"ep2": 2,
		"ep3": 4,
	}

	functionAppSkuMapMem = map[string]float64{
		"ep1": 3.5,
		"ep2": 7.0,
		"ep3": 14.0,
	}
)

// decodeFunctionAppValues decodes and returns functionAppValues from a Terraform values map.
func decodeFunctionAppValues(tfVals map[string]interface{}) (functionAppValues, error) {
	var v functionAppValues
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

// newFunctionApp initializes a new FunctionApp from the provider
func (p *Provider) newFunctionApp(vals functionAppValues) *FunctionApp {
	var appService AppService
	if len(vals.AppServicePlanId) == 0 && len(vals.ServicePlanId) == 0 {
		return &FunctionApp{
			provider: p,

			location: vals.Location,
			tier:     "standard",

			monthlyExecutions:   vals.Usage.MonthlyExecutions,
			executionDurationMs: vals.Usage.ExecutionDurationMs,
			memoryMb:            vals.Usage.MemoryMb,
			instances:           vals.Usage.Instances,
		}
	}
	if len(vals.AppServicePlanId) > 0 {
		appService = vals.AppServicePlanId[0]
	} else {
		appService = vals.ServicePlanId[0]
	}
	tier := "standard"
	if len(appService.Values.Sku) > 0 {
		skuTier := strings.ToLower(appService.Values.Sku[0].Tier)
		skuSize := strings.ToLower(appService.Values.Sku[0].Size)
		kind := strings.ToLower(appService.Values.Kind)

		if strings.ToLower(skuSize) != "y1" && (strings.ToLower(kind) == "elastic" || strings.ToLower(skuTier) == "elasticpremium") {
			tier = "premium"
		}

		return &FunctionApp{
			provider: p,

			location: vals.Location,
			skuName:  skuSize,
			tier:     tier,
			osType:   kind,

			monthlyExecutions:   vals.Usage.MonthlyExecutions,
			executionDurationMs: vals.Usage.ExecutionDurationMs,
			memoryMb:            vals.Usage.MemoryMb,
			instances:           vals.Usage.Instances,
		}
	}

	skuSize := appService.Values.SkuName
	if strings.HasPrefix(strings.ToLower(skuSize), "ep") {
		tier = "premium"
	}

	inst := FunctionApp{
		provider: p,

		location: vals.Location,
		skuName:  strings.ToLower(skuSize),
		tier:     tier,
		osType:   strings.ToLower(appService.Values.OsType),

		monthlyExecutions:   vals.Usage.MonthlyExecutions,
		executionDurationMs: vals.Usage.ExecutionDurationMs,
		memoryMb:            vals.Usage.MemoryMb,
		instances:           vals.Usage.Instances,
	}
	return &inst
}

func (inst *FunctionApp) Components() []query.Component {
	var components []query.Component

	if inst.tier == "premium" {
		cpu := inst.appFunctionPremiumCPUCostComponent()
		if cpu != nil {
			components = append(components, *cpu)
		}

		mem := inst.appFunctionPremiumMemoryCostComponent()
		if mem != nil {
			components = append(components, *mem)
		}

		return components
	}

	components = append(
		components,
		inst.appFunctionConsumptionExecutionTimeCostComponent(),
		inst.appFunctionConsumptionExecutionsCostComponent(),
	)

	return components
}

func (inst *FunctionApp) appFunctionPremiumCPUCostComponent() *query.Component {
	var skuCPU *int64

	if val, ok := functionAppSkuMapCPU[inst.skuName]; ok {
		skuCPU = &val
	}

	if skuCPU == nil {
		return nil
	}

	instances := decimal.NewFromInt(1)
	if inst.instances != nil {
		instances = decimal.NewFromInt(*inst.instances)
	}

	return &query.Component{
		Name:           fmt.Sprintf("vCPU (%s)", strings.ToUpper(inst.skuName)),
		Unit:           "vCPU",
		HourlyQuantity: instances.Mul(decimal.NewFromInt(*skuCPU)),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Functions"),
			Family:   util.StringPtr("Compute"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", ValueRegex: util.StringPtr("vCPU Duration$")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *FunctionApp) appFunctionPremiumMemoryCostComponent() *query.Component {
	var skuMemory *float64

	if val, ok := functionAppSkuMapMem[inst.skuName]; ok {
		skuMemory = &val
	}

	if skuMemory == nil {
		return nil
	}

	instances := decimal.NewFromInt(1)
	if inst.instances != nil {
		instances = decimal.NewFromInt(*inst.instances)
	}

	return &query.Component{
		Name:           fmt.Sprintf("Memory (%s)", strings.ToUpper(inst.skuName)),
		Unit:           "GB",
		HourlyQuantity: instances.Mul(decimal.NewFromFloat(*skuMemory)),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Functions"),
			Family:   util.StringPtr("Compute"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", ValueRegex: util.StringPtr("Memory Duration$")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *FunctionApp) appFunctionConsumptionExecutionTimeCostComponent() query.Component {
	var quantity decimal.Decimal
	gbSeconds := inst.calculateFunctionAppGBSeconds()
	if gbSeconds != nil {
		quantity = *gbSeconds
	}
	return query.Component{
		Name:            "Execution time",
		Unit:            "GB-seconds",
		MonthlyQuantity: quantity,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Functions"),
			Family:   util.StringPtr("Compute"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", ValueRegex: util.StringPtr("Execution Time$")},
				{Key: "sku_name", Value: util.StringPtr("Standard")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr("400000")},
			},
		},
	}
}

func (inst *FunctionApp) appFunctionConsumptionExecutionsCostComponent() query.Component {
	var executions decimal.Decimal
	if inst.monthlyExecutions != nil {
		executions = decimal.NewFromInt(*inst.monthlyExecutions).Div(decimal.NewFromInt(10))
	}

	return query.Component{
		Name:            "Executions",
		Unit:            "1M requests",
		MonthlyQuantity: executions,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Functions"),
			Family:   util.StringPtr("Compute"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", ValueRegex: util.StringPtr("Total Executions$")},
				{Key: "sku_name", Value: util.StringPtr("Standard")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr("100000")},
			},
		},
	}
}

func (inst *FunctionApp) calculateFunctionAppGBSeconds() *decimal.Decimal {
	if inst.memoryMb == nil || inst.executionDurationMs == nil || inst.monthlyExecutions == nil {
		return nil
	}

	memorySize := decimal.NewFromInt(*inst.memoryMb)
	averageRequestDuration := decimal.NewFromInt(*inst.executionDurationMs)
	monthlyRequests := decimal.NewFromInt(*inst.monthlyExecutions)

	if memorySize.LessThan(decimal.NewFromInt(128)) {
		memorySize = decimal.NewFromInt(128)
	}
	roundedMemory := memorySize.Div(decimal.NewFromInt(128)).Ceil().Mul(decimal.NewFromInt(128))
	if averageRequestDuration.LessThan(decimal.NewFromInt(100)) {
		averageRequestDuration = decimal.NewFromInt(100)
	}
	durationSeconds := monthlyRequests.Mul(averageRequestDuration).Mul(decimal.NewFromFloat(0.001))
	gbSeconds := durationSeconds.Mul(roundedMemory).Div(decimal.NewFromInt(1024))

	return &gbSeconds
}
