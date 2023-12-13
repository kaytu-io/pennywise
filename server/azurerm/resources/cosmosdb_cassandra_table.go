package resources

import (
	"github.com/mitchellh/mapstructure"
)

type CassandraKeyspaceId struct {
	Values struct {
		CosmosdbAccountName *CosmosdbAccountName `mapstructure:"account_name"`
		ResourceGroupName   *ResourceGroupName   `mapstructure:"resource_group_name"`
		Throughput          *int64               `mapstructure:"throughput"`
		AutoscaleSetting    []AutoscaleSetting   `mapstructure:"autoscale_settings"`
	} `mapstructure:"values"`
}

// cosmosdbCassandraTableValues is holds the values that we need to be able
// to calculate the price of the CosmosdbCassandraTable
type cosmosdbCassandraTableValues struct {
	CassandraKeyspaceId *CassandraKeyspaceId `mapstructure:"cassandra_keyspace_id"`
	Throughput          *int64               `mapstructure:"throughput"`

	Usage struct {
		MonthlyServerlessRequestUnits           *int64   `mapstructure:"monthly_serverless_request_units"`
		MaxRequestUnitsUtilizationPercentage    *float64 `mapstructure:"max_request_units_utilization_percentage"`
		MonthlyAnalyticalStorageReadOperations  *int64   `mapstructure:"monthly_analytical_storage_read_operations"`
		MonthlyAnalyticalStorageWriteOperations *int64   `mapstructure:"monthly_analytical_storage_write_operations"`
		StorageGb                               *int64   `mapstructure:"storage_gb"`
		MonthlyRestoredDataGb                   *int64   `mapstructure:"monthly_restored_data_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodeCosmosdbCassandraTableValues decodes and returns cosmosdbCassandraKeyspaceValues from a Terraform values map.
func decodeCosmosdbCassandraTableValues(tfVals map[string]interface{}) (cosmosdbCassandraTableValues, error) {
	var v cosmosdbCassandraTableValues
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
func (p *Provider) newCosmosdbCassandraTable(vals cosmosdbCassandraTableValues) *Cosmosdb {
	if vals.CassandraKeyspaceId == nil {
		return nil
	}
	if vals.CassandraKeyspaceId.Values.CosmosdbAccountName == nil {
		return nil
	}

	var throughput *int64
	if vals.CassandraKeyspaceId.Values.Throughput != nil {
		throughput = vals.CassandraKeyspaceId.Values.Throughput
	}
	if vals.Throughput != nil {
		throughput = vals.Throughput
	}

	inst := &Cosmosdb{
		provider: p,

		location:         vals.CassandraKeyspaceId.Values.CosmosdbAccountName.Values.Location,
		cosmosdbAccount:  vals.CassandraKeyspaceId.Values.CosmosdbAccountName,
		resourceGroup:    vals.CassandraKeyspaceId.Values.ResourceGroupName,
		autoscaleSetting: vals.CassandraKeyspaceId.Values.AutoscaleSetting,
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
