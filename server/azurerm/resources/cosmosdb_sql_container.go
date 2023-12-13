package resources

import (
	"github.com/mitchellh/mapstructure"
)

// cosmosdbSqlContainerValues is holds the values that we need to be able
// to calculate the price of the CosmosdbSqlContainer
type cosmosdbSqlContainerValues struct {
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

// decodeCosmosdbSqlContainerValues decodes and returns cosmosdbSqlContainerValues from a Terraform values map.
func decodeCosmosdbSqlContainerValues(tfVals map[string]interface{}) (cosmosdbSqlContainerValues, error) {
	var v cosmosdbSqlContainerValues
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

// newCosmosdbSqlDatabase initializes a new CosmosdbTable from the provider
func (p *Provider) newCosmosdbSqlContainer(vals cosmosdbSqlContainerValues) *Cosmosdb {
	if vals.CosmosdbAccountName == nil {
		return nil
	}
	inst := &Cosmosdb{
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
