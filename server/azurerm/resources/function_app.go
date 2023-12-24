package resources

import (
	"fmt"
	"github.com/kaytu-io/infracost/external/resources/azure"
	"github.com/kaytu-io/infracost/external/schema"
	"github.com/kaytu-io/pennywise/server/internal/query"
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
	} `mapstructure:"values"`
}

// functionAppValues is holds the values that we need to be able
// to calculate the price of the FunctionApp
type functionAppValues struct {
	Location         string     `mapstructure:"location"`
	AppServicePlanId AppService `mapstructure:"app_service_plan_id"`

	Usage struct {
		MonthlyExecutions   *int64 `mapstructure:"monthly_executions"`
		ExecutionDurationMs *int64 `mapstructure:"execution_duration_ms"`
		MemoryMb            *int64 `mapstructure:"memory_mb"`
		Instances           *int64 `mapstructure:"instances"`
	} `mapstructure:"pennywise_usage"`
}

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

// newAppServiceCertificateBinding initializes a new AppServiceCertificateBinding from the provider
func (p *Provider) newFunctionApp(vals functionAppValues) *FunctionApp {
	var appService AppService
	tier := "standard"
	// support for the legacy azurerm_app_service_plan resource. This is only applicable for the legacy azurerm_function_app resource.
	if len(appService.Values.Sku) > 0 {
		skuTier := strings.ToLower(appService.Values.Sku[0].Tier)
		skuSize := strings.ToLower(appService.Values.Sku[0].Size)
		kind := strings.ToLower(appService.Values.Kind)

		if strings.ToLower(skuSize) != "y1" && (strings.ToLower(kind) == "elastic" || strings.ToLower(skuTier) == "elasticpremium") {
			tier = "premium"
		}

		return &FunctionApp{
			Address: d.Address,
			Region:  region,
			SKUName: skuSize,
			Tier:    tier,
			OSType:  kind,
		}
	}

	skuName := data.Get("sku_name").String()
	if strings.HasPrefix(strings.ToLower(skuName), "ep") {
		tier = "premium"
	}

	instt := FunctionApp{
		location: region,
		skuName:  strings.ToLower(skuName),
		tier:     tier,
		osType:   strings.ToLower(vals.),

		monthlyExecutions:   vals.Usage.MonthlyExecutions,
		executionDurationMs: vals.Usage.ExecutionDurationMs,
		memoryMb:            vals.Usage.MemoryMb,
		instances:           vals.Usage.Instances,
	}
	return inst
}

func (inst *FunctionApp) Components() []query.Component {
	var components []query.Component

	return components
}

func (r *FunctionApp) appFunctionPremiumCPUCostComponent() *schema.CostComponent {
	var skuCPU *int64

	if val, ok := functionAppSkuMapCPU[r.SKUName]; ok {
		skuCPU = &val
	}

	if skuCPU == nil {
		return nil
	}

	instances := decimal.NewFromInt(1)
	if r.Instances != nil {
		instances = decimal.NewFromInt(*r.Instances)
	}

	return &schema.CostComponent{
		Name:           fmt.Sprintf("vCPU (%s)", strings.ToUpper(r.SKUName)),
		Unit:           "vCPU",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: decimalPtr(instances.Mul(decimal.NewFromInt(*skuCPU))),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Functions"),
			ProductFamily: strPtr("Compute"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "meterName", ValueRegex: regexPtr("vCPU Duration$")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
	}
}

func (r *FunctionApp) appFunctionPremiumMemoryCostComponent() *schema.CostComponent {
	var skuMemory *float64

	if val, ok := functionAppSkuMapMem[r.SKUName]; ok {
		skuMemory = &val
	}

	if skuMemory == nil {
		return nil
	}

	instances := decimal.NewFromInt(1)
	if r.Instances != nil {
		instances = decimal.NewFromInt(*r.Instances)
	}

	return &schema.CostComponent{
		Name:           fmt.Sprintf("Memory (%s)", strings.ToUpper(r.SKUName)),
		Unit:           "GB",
		UnitMultiplier: schema.HourToMonthUnitMultiplier,
		HourlyQuantity: decimalPtr(instances.Mul(decimal.NewFromFloat(*skuMemory))),
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Functions"),
			ProductFamily: strPtr("Compute"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "meterName", ValueRegex: regexPtr("Memory Duration$")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption: strPtr("Consumption"),
		},
	}
}

func (r *FunctionApp) appFunctionConsumptionExecutionTimeCostComponent() *schema.CostComponent {
	gbSeconds := r.calculateFunctionAppGBSeconds()
	return &schema.CostComponent{
		Name:            "Execution time",
		Unit:            "GB-seconds",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: gbSeconds,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Functions"),
			ProductFamily: strPtr("Compute"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "meterName", ValueRegex: regexPtr("Execution Time$")},
				{Key: "skuName", Value: strPtr("Standard")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption:   strPtr("Consumption"),
			StartUsageAmount: strPtr("400000"),
		},
	}
}

func (r *FunctionApp) appFunctionConsumptionExecutionsCostComponent() *schema.CostComponent {
	// Azure's pricing API returns prices per 10 executions so if the user has provided
	// the number of executions, we should divide it by 10
	var executions *decimal.Decimal
	if r.MonthlyExecutions != nil {
		executions = decimalPtr(decimal.NewFromInt(*r.MonthlyExecutions).Div(decimal.NewFromInt(10)))
	}

	return &schema.CostComponent{
		Name:            "Executions",
		Unit:            "1M requests",
		UnitMultiplier:  decimal.NewFromInt(100000),
		MonthlyQuantity: executions,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("azure"),
			Region:        strPtr(r.Region),
			Service:       strPtr("Functions"),
			ProductFamily: strPtr("Compute"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "meterName", ValueRegex: regexPtr("Total Executions$")},
				{Key: "skuName", Value: strPtr("Standard")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			PurchaseOption:   strPtr("Consumption"),
			StartUsageAmount: strPtr("100000"),
		},
	}
}

func (r *FunctionApp) calculateFunctionAppGBSeconds() *decimal.Decimal {
	if r.MemoryMb == nil || r.ExecutionDurationMs == nil || r.MonthlyExecutions == nil {
		return nil
	}

	memorySize := decimal.NewFromInt(*r.MemoryMb)
	averageRequestDuration := decimal.NewFromInt(*r.ExecutionDurationMs)
	monthlyRequests := decimal.NewFromInt(*r.MonthlyExecutions)

	// Use a min of 128MB, and round-up to nearest 128MB
	if memorySize.LessThan(decimal.NewFromInt(128)) {
		memorySize = decimal.NewFromInt(128)
	}
	roundedMemory := memorySize.Div(decimal.NewFromInt(128)).Ceil().Mul(decimal.NewFromInt(128))
	// Apply the minimum request duration
	if averageRequestDuration.LessThan(decimal.NewFromInt(100)) {
		averageRequestDuration = decimal.NewFromInt(100)
	}
	durationSeconds := monthlyRequests.Mul(averageRequestDuration).Mul(decimal.NewFromFloat(0.001))
	gbSeconds := durationSeconds.Mul(roundedMemory).Div(decimal.NewFromInt(1024))

	return &gbSeconds
}
