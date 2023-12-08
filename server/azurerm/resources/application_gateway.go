package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/tier_request"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

type ApplicationGatewaySku struct {
	Name     string `mapstructure:"name"`
	Capacity *int64 `mapstructure:"capacity"`
	Tier     string `mapstructure:"tier"`
}

type ApplicationGatewayAutoscaleConfiguration struct {
	MaxCapacity int64 `mapstructure:"max_capacity"`
	MinCapacity int64 `mapstructure:"min_capacity"`
}

// ApplicationGateway is the entity that holds the logic to calculate price
// of the azurerm_application_gateway
type ApplicationGateway struct {
	provider *Provider

	location               string
	sku                    []ApplicationGatewaySku
	autoscaleConfiguration []ApplicationGatewayAutoscaleConfiguration

	// Usage
	capacityUnits          *int64
	monthlyDataProcessedGb *int64
}

// applicationGatewayValues is holds the values that we need to be able
// to calculate the price of the ApplicationGateway
type applicationGatewayValues struct {
	Location               string                                     `mapstructure:"location"`
	Sku                    []ApplicationGatewaySku                    `mapstructure:"sku"`
	AutoscaleConfiguration []ApplicationGatewayAutoscaleConfiguration `mapstructure:"autoscale_configuration"`

	Usage struct {
		CapacityUnits          *int64 `mapstructure:"capacity_units"`
		MonthlyDataProcessedGb *int64 `mapstructure:"monthly_data_processed_gb"`
	} `mapstructure:"tc_usage"`
}

// decodeApplicationGatewayValues decodes and returns applicationGatewayValues from a Terraform values map.
func decodeApplicationGatewayValues(tfVals map[string]interface{}) (applicationGatewayValues, error) {
	var v applicationGatewayValues
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

// newApplicationGateway initializes a new ApplicationGateway from the provider
func (p *Provider) newApplicationGateway(vals applicationGatewayValues) *ApplicationGateway {
	inst := &ApplicationGateway{
		provider: p,

		location:               vals.Location,
		sku:                    vals.Sku,
		autoscaleConfiguration: vals.AutoscaleConfiguration,
		capacityUnits:          vals.Usage.CapacityUnits,
		monthlyDataProcessedGb: vals.Usage.MonthlyDataProcessedGb,
	}
	return inst
}

func (inst *ApplicationGateway) Components() []query.Component {
	var components []query.Component

	var sku, tier string
	tierLimits := []int{10240, 30720}
	var capacityUnits int64 = 1

	if len(inst.autoscaleConfiguration) > 0 {
		capacityUnits = inst.autoscaleConfiguration[0].MinCapacity
	}

	if inst.capacityUnits != nil {
		capacityUnits = *inst.capacityUnits
	}

	if inst.sku[0].Capacity != nil {
		capacityUnits = *inst.sku[0].Capacity
	}

	skuNameParts := strings.Split(inst.sku[0].Name, "_")
	if len(skuNameParts) > 1 {
		sku = skuNameParts[1]
	}

	if sku == "v2" {
		if skuNameParts[0] == "Standard" {
			tier = "Basic v2"
			components = append(components, inst.fixedForV2CostComponent(fmt.Sprintf("Gateway usage (%s)", tier), "Standard v2"))
			components = append(components, inst.capacityUnitsCostComponent("basic", "Standard v2", capacityUnits))
		} else {
			tier = "WAF v2"
			components = append(components, inst.fixedForV2CostComponent(fmt.Sprintf("Gateway usage (%s)", tier), tier))
			components = append(components, inst.capacityUnitsCostComponent("WAF", tier, capacityUnits))
		}
	} else {
		if skuNameParts[0] == "Standard" {
			tier = "Basic"
		} else {
			tier = "WAF"
		}
		components = append(components, inst.gatewayCostComponent(fmt.Sprintf("Gateway usage (%s, %s)", tier, sku), tier, sku, capacityUnits))

		if inst.monthlyDataProcessedGb != nil {
			result := tier_request.CalculateTierBuckets(decimal.NewFromInt(*inst.monthlyDataProcessedGb), tierLimits)

			if sku == "Small" {
				if result[0].GreaterThan(decimal.Zero) {
					components = append(components, inst.dataProcessingCostComponent("Data processing (0-10TB)", sku, "0", result[0]))
				}
				if result[1].GreaterThan(decimal.Zero) {
					components = append(components, inst.dataProcessingCostComponent("Data processing (10-40TB)", sku, "0", result[1]))
				}
				if result[2].GreaterThan(decimal.Zero) {
					components = append(components, inst.dataProcessingCostComponent("Data processing (over 40TB)", sku, "0", result[2]))
				}
			}

			if sku == "Medium" {
				if result[1].GreaterThan(decimal.Zero) {
					components = append(components, inst.dataProcessingCostComponent("Data processing (10-40TB)", sku, "10240", result[1]))
				}
				if result[2].GreaterThan(decimal.Zero) {
					components = append(components, inst.dataProcessingCostComponent("Data processing (over 40TB)", sku, "10240", result[2]))
				}
			}

			if sku == "Large" {
				if result[2].GreaterThan(decimal.Zero) {
					components = append(components, inst.dataProcessingCostComponent("Data processing (over 40TB)", sku, "40960", result[2]))
				}
			}

		} else {
			components = append(components, inst.dataProcessingCostComponent("Data processing (0-10TB)", sku, "0", decimal.Zero))
		}
	}

	return components
}

func (inst *ApplicationGateway) gatewayCostComponent(name, tier, sku string, capacityUnits int64) query.Component {
	return query.Component{
		Name:           name,
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(capacityUnits),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Application Gateway"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(fmt.Sprintf("%s Application Gateway", tier))},
				{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf("%s Gateway", sku))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *ApplicationGateway) dataProcessingCostComponent(name, sku, startUsage string, qty decimal.Decimal) query.Component {
	return query.Component{
		Name:            name,
		Unit:            "GB",
		MonthlyQuantity: qty,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Application Gateway"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf("%s Data Processed", sku))},
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

func (inst *ApplicationGateway) capacityUnitsCostComponent(name, tier string, capacityUnits int64) query.Component {
	return query.Component{
		Name:           fmt.Sprintf("V2 capacity units (%s)", name),
		Unit:           "CU",
		HourlyQuantity: decimal.NewFromInt(capacityUnits),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Application Gateway"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(fmt.Sprintf("Application Gateway %s", tier))},
				{Key: "meter_name", Value: util.StringPtr("Standard Capacity Units")},
			},
		},

		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *ApplicationGateway) fixedForV2CostComponent(name, tier string) query.Component {
	return query.Component{
		Name:           name,
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Application Gateway"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(fmt.Sprintf("Application Gateway %s", tier))},
				{Key: "meter_name", ValueRegex: util.StringPtr("Fixed Cost")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
