package resources

import "github.com/mitchellh/mapstructure"

// mysqlServerValues is holds the values that we need to be able
// to calculate the price of the MariadbServer
type mysqlServerValues struct {
	Location                  string `mapstructure:"location"`
	SkuName                   string `mapstructure:"sku_name"`
	StorageMb                 int64  `mapstructure:"storage_mb"`
	GeoRedundantBackupEnabled *bool  `mapstructure:"geo_redundant_backup_enabled"`

	Usage struct {
		AdditionalBackupStorageGb *int64 `mapstructure:"additional_backup_storage_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodeMysqlServerValues decodes and returns mysqlServerValues from a Terraform values map.
func decodeMysqlServerValues(tfVals map[string]interface{}) (mysqlServerValues, error) {
	var v mysqlServerValues
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

// newMysqlServer initializes a new MariadbServer from the provider
func (p *Provider) newMysqlServer(vals mysqlServerValues) *MariadbServer {
	inst := &MariadbServer{
		provider: p,

		serviceName:               "Azure Database for MySQL",
		location:                  vals.Location,
		skuName:                   vals.SkuName,
		storageMb:                 vals.StorageMb,
		geoRedundantBackupEnabled: vals.GeoRedundantBackupEnabled,

		additionalBackupStorageGb: vals.Usage.AdditionalBackupStorageGb,
	}
	return inst
}
