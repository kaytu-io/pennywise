package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

type ContainerRegistry struct {
	location                string
	geoReplicationLocations int
	sKU                     string
	storageGB               *float64
	monthlyBuildVCPUHrs     *float64
}

type ContainerRegistryValue struct {
	Location       string                 `mapstructure:"location"`
	SKU            string                 `mapstructure:"sku"`
	GeoReplication map[string]interface{} `mapstructure:"georeplications"`

	// we should get this fields from user
	StorageGB           *float64 `mapstructure:"storage_gb"`
	MonthlyBuildVCPUHrs *float64 `mapstructure:"monthly_build_vcpu_hrs"`
}

func (p *Provider) newContainerRegistry(vals ContainerRegistryValue) *ContainerRegistry {
	inst := &ContainerRegistry{
		location:                vals.Location,
		sKU:                     vals.SKU,
		geoReplicationLocations: len(vals.GeoReplication),
		storageGB:               vals.StorageGB,
		monthlyBuildVCPUHrs:     vals.MonthlyBuildVCPUHrs,
	}
	return inst
}

func decoderContainerRegistry(tfVals map[string]interface{}) (ContainerRegistryValue, error) {
	var v ContainerRegistryValue
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

func decimalPtr(de decimal.Decimal) *decimal.Decimal {
	return &de
}

func (inst *ContainerRegistry) component() []query.Component {
	var locationsCount int
	var storageGB, includedStorage, monthlyBuildVCPU *decimal.Decimal
	var overStorage decimal.Decimal

	sku := "Classic"

	if inst.sKU != "" {
		sku = inst.sKU
	}

	switch sku {
	case "Classic":
		includedStorage = decimalPtr(decimal.NewFromFloat(10))
	case "Basic":
		includedStorage = decimalPtr(decimal.NewFromFloat(10))
	case "Standard":
		includedStorage = decimalPtr(decimal.NewFromFloat(100))
	case "Premium":
		includedStorage = decimalPtr(decimal.NewFromFloat(500))
	}

	locationsCount = inst.geoReplicationLocations
	fmt.Printf("sku : %v \n", sku)
	costComponents := make([]query.Component, 0)

	if locationsCount > 0 {
		suffix := fmt.Sprintf("%d locations", locationsCount)
		if locationsCount == 1 {
			suffix = fmt.Sprintf("%d location", locationsCount)
		}
		costComponents = append(costComponents, containerRegistryGeolocationCostComponent(fmt.Sprintf("Geo replication (%s)", suffix), sku, "westeurope", inst.geoReplicationLocations))
	}
	costComponents = append(costComponents, containerRegistryCostComponent(fmt.Sprintf("Registry usage (%s)", sku), sku, "westeurope"))

	if inst.storageGB != nil {
		storageGB = decimalPtr(decimal.NewFromFloat(*inst.storageGB))
		if storageGB.GreaterThan(*includedStorage) {
			overStorage = storageGB.Sub(*includedStorage)
			storageGB = &overStorage

			costComponents = append(costComponents, containerRegistryStorageCostComponent(fmt.Sprintf("Storage (over %sGB)", includedStorage), sku, "westeurope", storageGB))
		}
	} else {
		costComponents = append(costComponents, containerRegistryStorageCostComponent(fmt.Sprintf("Storage (over %sGB)", includedStorage), sku, "westeurope", storageGB))
	}

	if inst.monthlyBuildVCPUHrs != nil {
		monthlyBuildVCPU = decimalPtr(decimal.NewFromFloat(*inst.monthlyBuildVCPUHrs * 3600))
	}
	costComponents = append(costComponents, containerRegistryCPUCostComponent("Build vCPU", sku, "westeurope", monthlyBuildVCPU))
	return costComponents
}

func containerRegistryCostComponent(name, sku, location string) query.Component {
	return query.Component{
		Name:            name,
		Unit:            "Day",
		MonthlyQuantity: decimal.NewFromInt(30),
		ProductFilter: &product.Filter{
			Location: util.StringPtr(location),
			Service:  util.StringPtr("Container Registry"),
			Family:   util.StringPtr("Containers"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Container Registry")},
				{Key: "sku_name", Value: util.StringPtr(sku)},
				{Key: "meter_name", Value: util.StringPtr(fmt.Sprintf("%s Registry Unit", sku))},
			},
		},
	}
}

func containerRegistryGeolocationCostComponent(name, sku, location string, geoReplicationLocations int) query.Component {
	fmt.Printf("location ; %v \n ", location)
	return query.Component{
		Name:            name,
		Unit:            "Day",
		MonthlyQuantity: decimal.NewFromInt(30 * int64(geoReplicationLocations)),
		ProductFilter: &product.Filter{
			Location: util.StringPtr(location),
			Service:  util.StringPtr("Container Registry"),
			Family:   util.StringPtr("Containers"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Container Registry")},
				{Key: "sku_name", Value: util.StringPtr(sku)},
				{Key: "meter_name", Value: util.StringPtr(fmt.Sprintf("%s Registry Unit", sku))},
			},
		},
	}
}

func containerRegistryStorageCostComponent(name, sku, location string, storage *decimal.Decimal) query.Component {
	return query.Component{
		Name:            name,
		Unit:            "GB",
		MonthlyQuantity: *storage,
		ProductFilter: &product.Filter{
			Location: util.StringPtr(location),
			Service:  util.StringPtr("Container Registry"),
			Family:   util.StringPtr("Containers"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Container Registry")},
				{Key: "sku_name", Value: util.StringPtr(sku)},
				{Key: "meter_name", Value: util.StringPtr("Data Stored")},
			},
		},
	}
}

func containerRegistryCPUCostComponent(name, sku, location string, monthlyBuildVCPU *decimal.Decimal) query.Component {
	return query.Component{
		Name:            name,
		Unit:            "second",
		MonthlyQuantity: *monthlyBuildVCPU,
		ProductFilter: &product.Filter{
			Location: util.StringPtr(location),
			Service:  util.StringPtr("Container Registry"),
			Family:   util.StringPtr("Containers"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Container Registry")},
				{Key: "sku_name", Value: util.StringPtr(sku)},
				{Key: "meter_name", Value: util.StringPtr("Task vCPU Duration")},
			},
		},
		//PriceFilter: &price.Filter{
		//	AttributeFilters: []*price.AttributeFilter{
		//		{Key: "startUsageAmount", Value: util.StringPtr("6000")},
		//	},
		//},
	}
}
