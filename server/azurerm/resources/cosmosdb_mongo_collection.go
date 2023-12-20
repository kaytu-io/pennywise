package resources

import (
	"github.com/mitchellh/mapstructure"
)

type MongoCollectionDatabaseName struct {
	Values struct {
		CosmosdbAccountName *CosmosdbAccountName `mapstructure:"account_name"`
		ResourceGroupName   *ResourceGroupName   `mapstructure:"resource_group_name"`
		Throughput          *int64               `mapstructure:"throughput"`
	} `mapstructure:"values"`
}

// cosmosdbCassandraTableValues is holds the values that we need to be able
// to calculate the price of the CosmosdbCassandraTable
type cosmosdbMongoCollectionValues struct {
	DatabaseName        *MongoCollectionDatabaseName `mapstructure:"database_name"`
	CosmosdbAccountName *CosmosdbAccountName         `mapstructure:"account_name"`
	ResourceGroupName   *ResourceGroupName           `mapstructure:"resource_group_name"`
	Throughput          *int64                       `mapstructure:"throughput"`
	AutoscaleSetting    []AutoscaleSetting           `mapstructure:"autoscale_settings"`

	Usage struct {
		// receives monthly number of serverless request units
		MonthlyServerlessRequestUnits *int64 `mapstructure:"monthly_serverless_request_units"`
		// receives Average utilisation of the maximum RU/s, starting at 10%. Possible values from 10 to 100
		MaxRequestUnitsUtilizationPercentage *float64 `mapstructure:"max_request_units_utilization_percentage"`
		// receives monthly number of read analytical storage operations
		MonthlyAnalyticalStorageReadOperations *int64 `mapstructure:"monthly_analytical_storage_read_operations"`
		// receives monthly number of write analytical storage operations.
		MonthlyAnalyticalStorageWriteOperations *int64 `mapstructure:"monthly_analytical_storage_write_operations"`
		// receives total size of storage in GB
		StorageGb *int64 `mapstructure:"storage_gb"`
		//receives monthly total amount of point-in-time restore data in GB
		MonthlyRestoredDataGb *int64 `mapstructure:"monthly_restored_data_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodeCosmosdbMongoCollectionValues decodes and returns cosmosdbMongoCollectionValues from a Terraform values map.
func decodeCosmosdbMongoCollectionValues(tfVals map[string]interface{}) (cosmosdbMongoCollectionValues, error) {
	var v cosmosdbMongoCollectionValues
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

// newCosmosdbCassandraKeyspace initializes a new CosmosdbTable from the provider
func (p *Provider) newCosmosdbMongoCollection(vals cosmosdbMongoCollectionValues) *Cosmosdb {
	if vals.DatabaseName == nil {
		return nil
	}

	var cosmosdbAccount *CosmosdbAccountName
	if vals.CosmosdbAccountName != nil {
		cosmosdbAccount = vals.CosmosdbAccountName
	}
	if vals.DatabaseName.Values.CosmosdbAccountName != nil {
		cosmosdbAccount = vals.DatabaseName.Values.CosmosdbAccountName
	}

	var resourceGroup *ResourceGroupName
	if vals.ResourceGroupName != nil {
		resourceGroup = vals.ResourceGroupName
	}
	if vals.DatabaseName.Values.ResourceGroupName != nil {
		resourceGroup = vals.DatabaseName.Values.ResourceGroupName
	}

	var throughput *int64
	if vals.DatabaseName.Values.Throughput != nil {
		throughput = vals.DatabaseName.Values.Throughput
	}
	if vals.Throughput != nil {
		throughput = vals.Throughput
	}

	inst := &Cosmosdb{
		provider: p,

		location:         cosmosdbAccount.Values.Location,
		cosmosdbAccount:  cosmosdbAccount,
		resourceGroup:    resourceGroup,
		autoscaleSetting: vals.AutoscaleSetting,
		throughput:       throughput,

		monthlyServerlessRequestUnits:           vals.Usage.MonthlyServerlessRequestUnits,
		maxRequestUnitsUtilizationPercentage:    vals.Usage.MaxRequestUnitsUtilizationPercentage,
		monthlyAnalyticalStorageReadOperations:  vals.Usage.MonthlyAnalyticalStorageReadOperations,
		monthlyAnalyticalStorageWriteOperations: vals.Usage.MonthlyAnalyticalStorageWriteOperations,
		storageGb:                               vals.Usage.StorageGb,
		monthlyRestoredDataGb:                   vals.Usage.MonthlyRestoredDataGb,
	}
	return inst
}
