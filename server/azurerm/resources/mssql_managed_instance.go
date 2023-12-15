package resources

import "github.com/mitchellh/mapstructure"

// mssqlManagedInstanceValues is holds the values that we need to be able
// to calculate the price of the MssqlManagedInstance
type mssqlManagedInstanceValues struct {
	Location           string `mapstructure:"location"`
	SkuName            string `mapstructure:"sku_name"`
	Cores              int64  `mapstructure:"vcores"`
	LicenseType        string `mapstructure:"license_type"`
	StorageAccountType string `mapstructure:"storage_account_type"`
	StorageSizeInGb    int64  `mapstructure:"storage_size_in_gb"`

	Usage struct {
		LongTermRetentionStorageGB *int64 `mapstructure:"long_term_retention_storage_gb"`
		BackupStorageGB            *int64 `mapstructure:"backup_storage_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodeMssqlManagedInstanceValues decodes and returns mssqlManagedInstanceValues from a Terraform values map.
func decodeMssqlManagedInstanceValues(tfVals map[string]interface{}) (mssqlManagedInstanceValues, error) {
	var v mssqlManagedInstanceValues
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

// newMssqlManagedInstance initializes a new SqlManagedInstance from the provider
func (p *Provider) newMssqlManagedInstance(vals mssqlManagedInstanceValues) *SqlManagedInstance {
	storageAccountType := vals.StorageAccountType
	if storageAccountType == "" {
		storageAccountType = "LRS"
	} else if storageAccountType == "GRS" {
		storageAccountType = "RA-GRS"
	}

	inst := &SqlManagedInstance{
		provider: p,

		location:           getLocationName(vals.Location),
		sku:                vals.SkuName,
		cores:              vals.Cores,
		licenseType:        vals.LicenseType,
		storageAccountType: vals.StorageAccountType,
		storageSizeInGb:    vals.StorageSizeInGb,
		mssql:              false,
	}
	return inst
}
