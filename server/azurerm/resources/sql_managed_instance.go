package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

// SqlManagedInstance is the entity that holds the logic to calculate price
// of the azurerm_sql_managed_instance
type SqlManagedInstance struct {
	provider *Provider

	location           string
	sku                string
	cores              int64
	licenseType        string
	storageAccountType string
	storageSizeInGb    int64
	mssql              bool

	// Usage
	// receive number of GBs used by long-term retention backup storage.
	longTermRetentionStorageGB *int64
	// receive number of GBs used by Point-In-Time Restore (PITR) backup storage.
	backupStorageGB *int64
}

// sqlManagedInstanceValues is holds the values that we need to be able
// to calculate the price of the SqlManagedInstance
type sqlManagedInstanceValues struct {
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

// decodeSqlManagedInstanceValues decodes and returns sqlManagedInstanceValues from a Terraform values map.
func decodeSqlManagedInstanceValues(tfVals map[string]interface{}) (sqlManagedInstanceValues, error) {
	var v sqlManagedInstanceValues
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

// newSqlManagedInstance initializes a new SqlManagedInstance from the provider
func (p *Provider) newSqlManagedInstance(vals sqlManagedInstanceValues) *SqlManagedInstance {
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

func (inst *SqlManagedInstance) Components() []resource.Component {
	var components []resource.Component

	components = append(components, inst.managedInstanceComponent())

	if !inst.mssql || ((inst.storageSizeInGb - 32) > 0) {
		components = append(components, inst.sqlMIStorageCostComponent(), inst.sqlMIBackupCostComponent())
	}

	if inst.licenseType == "LicenseIncluded" {
		components = append(components, inst.sqlMILicenseCostComponent())
	}

	components = append(components, inst.sqlMILongTermRetentionStorageGBCostComponent())
	GetCostComponentNamesAndSetLogger(components, inst.provider.logger)

	return components
}

func (inst *SqlManagedInstance) managedInstanceComponent() resource.Component {
	return resource.Component{
		Name:           fmt.Sprintf("Compute (%s %d Cores)", strings.ToTitle(inst.sku), inst.cores),
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("SQL Managed Instance"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: inst.productDescription()},
				{Key: "sku_name", Value: inst.meteredName()},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *SqlManagedInstance) productDescription() *string {
	productDescription := ""

	if strings.Contains(inst.sku, "GP") {
		productDescription = "SQL Managed Instance General Purpose"
	} else if strings.Contains(inst.sku, "BC") {
		productDescription = "SQL Managed Instance Business Critical"
	}

	if strings.Contains(inst.sku, "Gen5") {
		productDescription = fmt.Sprintf("%s - %s", productDescription, "Compute Gen5")
	}

	return util.StringPtr(productDescription)
}

func (inst *SqlManagedInstance) meteredName() *string {
	meterName := fmt.Sprintf("%d %s", inst.cores, "vCore")

	return util.StringPtr(meterName)
}

func (inst *SqlManagedInstance) sqlMIStorageCostComponent() resource.Component {
	return resource.Component{
		Name:            "Additional Storage",
		Unit:            "GB",
		MonthlyQuantity: decimal.NewFromInt(inst.storageSizeInGb - 32),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("SQL Managed Instance"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("SQL Managed Instance General Purpose - Storage")},
				{Key: "meter_name", ValueRegex: util.StringPtr("Data Stored$")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *SqlManagedInstance) sqlMIBackupCostComponent() resource.Component {
	var backup decimal.Decimal

	if inst.backupStorageGB != nil {
		backup = decimal.NewFromInt(*inst.backupStorageGB)
	}

	return resource.Component{
		Name:            fmt.Sprintf("PITR backup storage (%s)", inst.storageAccountType),
		Unit:            "GB",
		MonthlyQuantity: backup,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("SQL Managed Instance"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("SQL Managed Instance PITR Backup Storage")},
				{Key: "meter_name", Value: util.StringPtr(fmt.Sprintf("%s Data Stored", inst.storageAccountType))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *SqlManagedInstance) sqlMILicenseCostComponent() resource.Component {
	return resource.Component{
		Name:           "SQL license",
		Unit:           "vCore-hours",
		HourlyQuantity: decimal.NewFromInt(inst.cores),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr("Global"),
			Service:  util.StringPtr("SQL Managed Instance"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("SQL Managed Instance General Purpose - SQL License")},
				{Key: "meter_name", Value: util.StringPtr("vCore")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *SqlManagedInstance) sqlMILongTermRetentionStorageGBCostComponent() resource.Component {
	var retention decimal.Decimal

	if inst.longTermRetentionStorageGB != nil {
		retention = decimal.NewFromInt(*inst.longTermRetentionStorageGB)
	}

	return resource.Component{
		Name:            fmt.Sprintf("LTR backup storage (%s)", inst.storageAccountType),
		Unit:            "GB",
		MonthlyQuantity: retention,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("SQL Managed Instance"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("SQL Managed Instance - LTR Backup Storage")},
				{Key: "meter_name", Value: util.StringPtr(fmt.Sprintf("LTR Backup %s Data Stored", inst.storageAccountType))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
