package resources

import "github.com/mitchellh/mapstructure"

// postgresqlServerValues is holds the values that we need to be able
// to calculate the price of the MariadbServer
type postgresqlServerValues struct {
	Location                  string `mapstructure:"location"`
	SkuName                   string `mapstructure:"sku_name"`
	StorageMb                 int64  `mapstructure:"storage_mb"`
	GeoRedundantBackupEnabled *bool  `mapstructure:"geo_redundant_backup_enabled"`

	Usage struct {
		AdditionalBackupStorageGb *int64 `mapstructure:"additional_backup_storage_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodePostgresqlServerValues decodes and returns postgresqlServerValues from a Terraform values map.
func decodePostgresqlServerValues(tfVals map[string]interface{}) (postgresqlServerValues, error) {
	var v postgresqlServerValues
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

// newPostgresqlServer initializes a new MariadbServer from the provider
func (p *Provider) newPostgresqlServer(vals postgresqlServerValues) *MariadbServer {
	inst := &MariadbServer{
		provider: p,

		serviceName:               "Azure Database for PostgreSQL",
		location:                  vals.Location,
		skuName:                   vals.SkuName,
		storageMb:                 vals.StorageMb,
		geoRedundantBackupEnabled: vals.GeoRedundantBackupEnabled,

		additionalBackupStorageGb: vals.Usage.AdditionalBackupStorageGb,
	}
	return inst
}
