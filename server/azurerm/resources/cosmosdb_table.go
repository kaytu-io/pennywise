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

type GeoLocation struct {
	Location         string `mapstructure:"Location"`
	FailoverPriority int64  `mapstructure:"failover_priority"`
	ZoneRedundant    *bool  `mapstructure:"zone_redundant"`
}

type CosmosdbAccountBackup struct {
	Type              string   `mapstructure:"type"`
	IntervalInMinutes *float64 `mapstructure:"interval_in_minutes"`
	RetentionInHours  *float64 `mapstructure:"retention_in_hours"`
}

type CosmosdbAccountName struct {
	Values struct {
		Location                     string                  `mapstructure:"location"`
		OfferType                    string                  `mapstructure:"offer_type"`
		GeoLocation                  []GeoLocation           `mapstructure:"geo_location"`
		Backup                       []CosmosdbAccountBackup `mapstructure:"backup"`
		EnableMultipleWriteLocations *bool                   `mapstructure:"enable_multiple_write_locations"`
		AnalyticStorageEnabled       *bool                   `mapstructure:"analytical_storage_enabled"`
	}
}

type ResourceGroupName struct {
	Values struct {
		Location string `mapstructure:"location"`
	}
}

type AutoscaleSetting struct {
	MaxThroughput int64 `mapstructure:"max_throughput"`
}

// CosmosdbTable is the entity that holds the logic to calculate price
// of the azurerm_key_vault_key
type CosmosdbTable struct {
	provider *Provider

	cosmosdbAccount  *CosmosdbAccountName
	resourceGroup    *ResourceGroupName
	throughput       *int64
	autoscaleSetting []AutoscaleSetting
	location         string

	// Usage
	monthlyServerlessRequestUnits           *int64
	maxRequestUnitsUtilizationPercentage    *float64
	monthlyAnalyticalStorageReadOperations  *int64
	monthlyAnalyticalStorageWriteOperations *int64
	storageGb                               *int64
	monthlyRestoredDataGb                   *int64
}

// cosmosdbTableValues is holds the values that we need to be able
// to calculate the price of the CosmosdbTable
type cosmosdbTableValues struct {
	CosmosdbAccountName *CosmosdbAccountName `mapstructure:"account_name"`
	ResourceGroupName   *ResourceGroupName   `mapstructure:"resource_group_name"`
	Throughput          *int64               `mapstructure:"throughput"`
	AutoscaleSetting    []AutoscaleSetting   `mapstructure:"autoscale_settings"`

	Usage struct {
		MonthlyServerlessRequestUnits           *int64   `mapstructure:"monthly_serverless_request_units"`
		MaxRequestUnitsUtilizationPercentage    *float64 `mapstructure:"max_request_units_utilization_percentage"`
		MonthlyAnalyticalStorageReadOperations  *int64   `mapstructure:"monthly_analytical_storage_read_operations"`
		MonthlyAnalyticalStorageWriteOperations *int64   `mapstructure:"monthly_analytical_storage_write_operations"`
		StorageGb                               *int64   `mapstructure:"storage_gb"`
		MonthlyRestoredDataGb                   *int64   `mapstructure:"monthly_restored_data_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodeCosmosdbTableValues decodes and returns cosmosdbTableValues from a Terraform values map.
func decodeCosmosdbTableValues(tfVals map[string]interface{}) (cosmosdbTableValues, error) {
	var v cosmosdbTableValues
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

// newCosmosdbTable initializes a new CosmosdbTable from the provider
func (p *Provider) newCosmosdbTable(vals cosmosdbTableValues) *CosmosdbTable {
	if vals.CosmosdbAccountName == nil {
		return nil
	}
	inst := &CosmosdbTable{
		provider: p,

		location:         vals.CosmosdbAccountName.Values.Location,
		cosmosdbAccount:  vals.CosmosdbAccountName,
		resourceGroup:    vals.ResourceGroupName,
		autoscaleSetting: vals.AutoscaleSetting,
		throughput:       vals.Throughput,

		monthlyServerlessRequestUnits:           vals.Usage.MonthlyServerlessRequestUnits,
		maxRequestUnitsUtilizationPercentage:    vals.Usage.MaxRequestUnitsUtilizationPercentage,
		monthlyAnalyticalStorageReadOperations:  vals.Usage.MonthlyAnalyticalStorageReadOperations,
		monthlyAnalyticalStorageWriteOperations: vals.Usage.MonthlyAnalyticalStorageWriteOperations,
		storageGb:                               vals.Usage.StorageGb,
		monthlyRestoredDataGb:                   vals.Usage.MonthlyRestoredDataGb,
	}
	return inst
}

func (inst *CosmosdbTable) Components() []query.Component {
	var components []query.Component

	if inst == nil {
		return nil
	}
	if inst.cosmosdbAccount == nil {
		return nil
	}

	if inst.cosmosdbAccount == nil {
		return components
	}
	components = append(components, inst.cosmosDBCostComponents()...)

	return components
}

func (inst *CosmosdbTable) cosmosDBCostComponents() []query.Component {
	var components []query.Component

	// The geo_location attribute is a required attribute however it can be an empty list because of
	// expressions evaluating as nil, e.g. using a data block. If the geoLocations variable is empty
	// we set it as a sane default which is using the location from the parent region.
	if len(inst.cosmosdbAccount.Values.GeoLocation) == 0 {
		inst.cosmosdbAccount.Values.GeoLocation = append(inst.cosmosdbAccount.Values.GeoLocation, GeoLocation{
			Location:         inst.location,
			FailoverPriority: 1,
		})
	}

	model := "Provisioned"
	skuName := "RUs"
	if inst.cosmosdbAccount.Values.EnableMultipleWriteLocations != nil {
		if *inst.cosmosdbAccount.Values.EnableMultipleWriteLocations {
			skuName = "mRUs"
		}
	}

	var throughputs *decimal.Decimal
	if inst.throughput != nil {
		throughputs = decimalPtr(decimal.NewFromInt(*inst.throughput))
	} else if len(inst.autoscaleSetting) > 0 {
		throughputs = decimalPtr(decimal.NewFromInt(inst.autoscaleSetting[0].MaxThroughput))
		model = "Autoscale"
	} else {
		model = "Serverless"
		availabilityZone := false
		for _, geo := range inst.cosmosdbAccount.Values.GeoLocation {
			if geo.ZoneRedundant != nil {
				availabilityZone = *geo.ZoneRedundant
			}
		}
		location := inst.cosmosdbAccount.Values.GeoLocation[0].Location
		components = append(components, inst.serverlessCosmosCostComponent(location, availabilityZone))
	}
	if model == "Provisioned" || model == "Autoscale" {
		components = inst.provisionedCosmosCostComponents(model, throughputs, skuName)
	}

	components = append(components, inst.storageCosmosCostComponents(skuName)...)
	components = append(components, inst.backupStorageCosmosCostComponents()...)

	return components
}

func (inst *CosmosdbTable) provisionedCosmosCostComponents(model string, throughputs *decimal.Decimal, skuName string) []query.Component {
	var components []query.Component

	var meterName string
	if strings.ToLower(skuName) == "rus" {
		meterName = "100 RU/s"
	} else {
		meterName = "100 Multi-master RU/s"
	}

	name := "Provisioned throughput"
	if model == "Autoscale" {
		name = fmt.Sprintf("%s (autoscale", name)

		if inst.maxRequestUnitsUtilizationPercentage != nil {
			throughputs = decimalPtr(throughputs.Mul(decimal.NewFromFloat(*inst.maxRequestUnitsUtilizationPercentage / 100)))
		} else {
			throughputs = nil
		}
	} else {
		name = fmt.Sprintf("%s (provisioned", name)
	}

	if throughputs != nil {
		throughputs = decimalPtr(throughputs.Div(decimal.NewFromInt(100)))
	} else {
		return components
	}

	for _, g := range inst.cosmosdbAccount.Values.GeoLocation {
		quantity := *throughputs

		if model == "Autoscale" {
			if strings.ToLower(skuName) == "rus" {
				quantity = quantity.Mul(decimal.NewFromFloat(1.5))
			}
		} else {
			if strings.ToLower(skuName) == "rus" {
				if g.ZoneRedundant != nil {
					if *g.ZoneRedundant {
						quantity = quantity.Mul(decimal.NewFromFloat(1.25))
					}
				}
			}
		}

		if l := locationNameMapping(g.Location); l != "" {
			components = append(components, query.Component{
				Name:           fmt.Sprintf("%s, %s)", name, l),
				Unit:           "RU/s x 100",
				HourlyQuantity: quantity,
				ProductFilter: &product.Filter{
					Provider: util.StringPtr(inst.provider.key),
					Location: util.StringPtr(g.Location),
					Service:  util.StringPtr("Azure Cosmos DB"),
					Family:   util.StringPtr("Databases"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "meter_name", Value: util.StringPtr(meterName)},
						{Key: "sku_name", Value: util.StringPtr(skuName)},
					},
				},
				PriceFilter: &price.Filter{
					AttributeFilters: []*price.AttributeFilter{
						{Key: "type", Value: util.StringPtr("Consumption")},
					},
				},
			})
		}
	}

	return components
}

func (inst *CosmosdbTable) serverlessCosmosCostComponent(location string, availabilityZone bool) query.Component {
	var requestUnits decimal.Decimal
	if inst.monthlyServerlessRequestUnits != nil {
		requestUnits = decimal.NewFromInt(*inst.monthlyServerlessRequestUnits)
		requestUnits = requestUnits.Div(decimal.NewFromInt(1000000))
	}

	if availabilityZone {
		requestUnits = requestUnits.Mul(decimal.NewFromFloat(1.25))
	}

	return query.Component{
		Name:            "Provisioned throughput (serverless)",
		Unit:            "1M RU",
		MonthlyQuantity: requestUnits,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(location),
			Service:  util.StringPtr("Azure Cosmos DB"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Azure Cosmos DB serverless")},
				{Key: "sku_name", Value: util.StringPtr("RUs")},
				{Key: "meter_name", Value: util.StringPtr("1M RUs")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *CosmosdbTable) storageCosmosCostComponents(skuName string) []query.Component {
	var components []query.Component
	var storageGB decimal.Decimal
	if inst.storageGb != nil {
		storageGB = decimal.NewFromInt(*inst.storageGb)
	}

	for _, g := range inst.cosmosdbAccount.Values.GeoLocation {
		if l := locationNameMapping(g.Location); l != "" {
			components = append(components, storageCosmosCostComponent(
				fmt.Sprintf("Transactional storage (%s)", l),
				g.Location,
				skuName,
				"Azure Cosmos DB",
				storageGB))

			if inst.cosmosdbAccount.Values.AnalyticStorageEnabled != nil {
				if *inst.cosmosdbAccount.Values.AnalyticStorageEnabled {
					components = append(components, storageCosmosCostComponent(
						fmt.Sprintf("Analytical storage (%s)", l),
						g.Location,
						"Standard",
						"Azure Cosmos DB Analytics Storage",
						storageGB))

					var writeOperations, readOperations decimal.Decimal
					if inst.monthlyAnalyticalStorageWriteOperations != nil {
						writeOperations = decimal.NewFromInt(*inst.monthlyAnalyticalStorageWriteOperations)
						writeOperations = writeOperations.Div(decimal.NewFromInt(10000))
					}
					if inst.monthlyAnalyticalStorageReadOperations != nil {
						readOperations = decimal.NewFromInt(*inst.monthlyAnalyticalStorageReadOperations)
						readOperations = readOperations.Div(decimal.NewFromInt(10000))
					}
					components = append(components, inst.operationsCosmosCostComponent(
						fmt.Sprintf("Analytical write operations (%s)", l),
						g.Location,
						"Write Operations",
						writeOperations,
					))

					components = append(components, inst.operationsCosmosCostComponent(
						fmt.Sprintf("Analytical read operations (%s)", l),
						g.Location,
						"Read Operations",
						readOperations,
					))
				}
			}
		}
	}

	return components
}

func (inst *CosmosdbTable) backupStorageCosmosCostComponents() []query.Component {
	var components []query.Component
	var backupStorageGB decimal.Decimal
	if inst.storageGb != nil {
		backupStorageGB = decimal.NewFromInt(*inst.storageGb)
	}

	var name, meterName, skuName, productName string
	numberOfCopies := decimalPtr(decimal.NewFromInt(1))
	backupType := "Pereodic"
	if len(inst.cosmosdbAccount.Values.Backup) > 0 {
		backupType = inst.cosmosdbAccount.Values.Backup[0].Type
	}

	if strings.ToLower(backupType) == "periodic" {
		name = "Periodic backup"
		meterName = "Data Stored"
		skuName = "Standard"
		productName = "Azure Cosmos DB Snapshot"

		if !backupStorageGB.Equal(decimal.Zero) {
			intervalHrs := 4.0
			retentionHrs := 8.0

			if len(inst.cosmosdbAccount.Values.Backup) > 0 {
				if inst.cosmosdbAccount.Values.Backup[0].IntervalInMinutes != nil {
					intervalHrs = *inst.cosmosdbAccount.Values.Backup[0].IntervalInMinutes / 60
				}
			}
			if len(inst.cosmosdbAccount.Values.Backup) > 0 {
				if inst.cosmosdbAccount.Values.Backup[0].RetentionInHours != nil {

					retentionHrs = *inst.cosmosdbAccount.Values.Backup[0].RetentionInHours
				}
			}

			if retentionHrs > intervalHrs {
				numberOfCopies = decimalPtr(decimal.NewFromFloat(retentionHrs / intervalHrs).Floor().Sub(decimal.NewFromInt(2)))
			}
			backupStorageGB = backupStorageGB.Mul(*numberOfCopies)
		}
	} else {
		name = "Continuous backup"
		meterName = "Continuous Backup"
		skuName = "Backup"
		productName = "Azure Cosmos DB - PITR"
	}

	for _, g := range inst.cosmosdbAccount.Values.GeoLocation {
		if backupStorageGB.Equal(decimal.Zero) {
			break
		}
		if l := locationNameMapping(g.Location); l != "" {
			components = append(components, backupCosmosCostComponent(
				fmt.Sprintf("%s (%s)", name, l),
				g.Location,
				skuName,
				productName,
				meterName,
				backupStorageGB,
			))
		}
	}

	var pitr decimal.Decimal
	if inst.monthlyRestoredDataGb != nil {
		pitr = decimal.NewFromInt(*inst.monthlyRestoredDataGb)
	}

	components = append(components, backupCosmosCostComponent(
		"Restored data",
		inst.location,
		"Backup",
		"Azure Cosmos DB - PITR",
		".*Data Restore",
		pitr,
	))

	return components
}

func storageCosmosCostComponent(name, location, skuName, productName string, quantities decimal.Decimal) query.Component {
	return query.Component{
		Name:            name,
		Unit:            "GB",
		MonthlyQuantity: quantities,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(location),
			Service:  util.StringPtr("Azure Cosmos DB"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", ValueRegex: util.StringPtr(".*Data Stored")},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
				{Key: "product_name", Value: util.StringPtr(productName)},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func backupCosmosCostComponent(name, location, skuName, productName, meterName string, quantities decimal.Decimal) query.Component {
	return query.Component{
		Name:            name,
		Unit:            "GB",
		MonthlyQuantity: quantities,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(location),
			Service:  util.StringPtr("Azure Cosmos DB"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", ValueRegex: util.StringPtr(meterName)},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
				{Key: "product_name", Value: util.StringPtr(productName)},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *CosmosdbTable) operationsCosmosCostComponent(name, location, meterName string, quantities decimal.Decimal) query.Component {
	return query.Component{
		Name:            name,
		Unit:            "10K operations",
		MonthlyQuantity: quantities,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Azure Cosmos DB"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf(".*%s", meterName))},
				{Key: "sku_name", Value: util.StringPtr("Standard")},
				{Key: "product_name", Value: util.StringPtr("Azure Cosmos DB Analytics Storage")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
