package resources

import (
	"github.com/mitchellh/mapstructure"
)

// cosmosdbSqlDatabaseValues is holds the values that we need to be able
// to calculate the price of the CosmosdbSqlDatabase
type cosmosdbSqlDatabaseValues struct {
	CosmosdbAccountName *CosmosdbAccountName `mapstructure:"account_name"`
	ResourceGroupName   *ResourceGroupName   `mapstructure:"resource_group_name"`
	Throughput          *int64               `mapstructure:"throughput"`
	AutoscaleSetting    []AutoscaleSetting   `mapstructure:"autoscale_settings"`

	Usage struct {
		// receives monthly number of serverless request units
		MonthlyServerlessRequestUnits *int64 `mapstructure:"monthly_serverless_request_units"`
		// receives Average utilisation of the maximum RU/s, starting at 10%. Possible values from 10 to 100
		MaxRequestUnitsUtilizationPercentage *float64 `mapstructure:"max_request_units_utilization_percentage"`
		// receives monthly number of read analytical storage operations\
		MonthlyAnalyticalStorageReadOperations *int64 `mapstructure:"monthly_analytical_storage_read_operations"`
		// receives monthly number of write analytical storage operations.
		MonthlyAnalyticalStorageWriteOperations *int64 `mapstructure:"monthly_analytical_storage_write_operations"`
		// receives total size of storage in GB
		StorageGb *int64 `mapstructure:"storage_gb"`
		//receives monthly total amount of point-in-time restore data in GB
		MonthlyRestoredDataGb *int64 `mapstructure:"monthly_restored_data_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodeCosmosdbSqlDatabaseValues decodes and returns cosmosdbSqlDatabaseValues from a Terraform values map.
func decodeCosmosdbSqlDatabaseValues(tfVals map[string]interface{}) (cosmosdbSqlDatabaseValues, error) {
	var v cosmosdbSqlDatabaseValues
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
func (p *Provider) newCosmosdbSqlDatabase(vals cosmosdbSqlDatabaseValues) *Cosmosdb {
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
