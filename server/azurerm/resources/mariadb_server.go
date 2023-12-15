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

// MariadbServer is the entity that holds the logic to calculate price
// of the azurerm_mariadb_server
type MariadbServer struct {
	provider *Provider

	serviceName               string
	location                  string
	skuName                   string
	storageMb                 int64
	geoRedundantBackupEnabled *bool

	// Usage
	additionalBackupStorageGb *int64
}

// mariadbServerValues is holds the values that we need to be able
// to calculate the price of the MariadbServer
type mariadbServerValues struct {
	Location                  string `mapstructure:"location"`
	SkuName                   string `mapstructure:"sku_name"`
	StorageMb                 int64  `mapstructure:"storage_mb"`
	GeoRedundantBackupEnabled *bool  `mapstructure:"geo_redundant_backup_enabled"`

	Usage struct {
		AdditionalBackupStorageGb *int64 `mapstructure:"additional_backup_storage_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodeMariadbServerValues decodes and returns mariadbServerValues from a Terraform values map.
func decodeMariadbServerValues(tfVals map[string]interface{}) (mariadbServerValues, error) {
	var v mariadbServerValues
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

// newMariadbServer initializes a new MariadbServer from the provider
func (p *Provider) newMariadbServer(vals mariadbServerValues) *MariadbServer {
	inst := &MariadbServer{
		provider: p,

		serviceName:               "Azure Database for MariaDB",
		location:                  vals.Location,
		skuName:                   vals.SkuName,
		storageMb:                 vals.StorageMb,
		geoRedundantBackupEnabled: vals.GeoRedundantBackupEnabled,

		additionalBackupStorageGb: vals.Usage.AdditionalBackupStorageGb,
	}
	return inst
}

func (inst *MariadbServer) Components() []query.Component {
	var components []query.Component

	var tier, family, cores string
	if sku := strings.Split(inst.skuName, "_"); len(sku) == 3 {
		tier = sku[0]
		family = sku[1]
		cores = sku[2]
	} else {
		return components
	}

	tierName := map[string]string{
		"B":  "Basic",
		"GP": "General Purpose",
		"MO": "Memory Optimized",
	}[tier]

	productNameRegex := fmt.Sprintf(".*%s - Compute %s.*", tierName, family)
	skuName := fmt.Sprintf("%s vCore", cores)

	components = append(components, inst.databaseComputeInstance(fmt.Sprintf("Compute (%s)", inst.skuName), productNameRegex, skuName))

	storageGB := inst.storageMb / 1024

	// MO and GP storage cost are the same, and we don't have cost component for MO Storage now
	if strings.ToLower(tier) == "mo" {
		tierName = "General Purpose"
	}
	productNameRegex = fmt.Sprintf(".*%s - Storage.*", tierName)

	components = append(components, inst.databaseStorageComponent(productNameRegex, storageGB))

	var backupStorageGB decimal.Decimal
	if inst.additionalBackupStorageGb != nil {
		backupStorageGB = decimal.NewFromInt(*inst.additionalBackupStorageGb)
	}

	skuName = "Backup LRS"
	if inst.geoRedundantBackupEnabled != nil {
		if *inst.geoRedundantBackupEnabled {
			skuName = "Backup GRS"
		}
	}

	components = append(components, inst.databaseBackupStorageComponent(skuName, backupStorageGB))

	return components
}

func (inst *MariadbServer) databaseComputeInstance(name, productNameRegex, skuName string) query.Component {
	return query.Component{
		Name:           name,
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr(inst.serviceName),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(productNameRegex)},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *MariadbServer) databaseStorageComponent(productNameRegex string, storageGB int64) query.Component {
	return query.Component{
		Name:            "Storage",
		Unit:            "GB",
		MonthlyQuantity: decimal.NewFromInt(storageGB),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr(inst.serviceName),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(productNameRegex)},
			},
		},
	}
}

func (inst *MariadbServer) databaseBackupStorageComponent(skuName string, backupStorageGB decimal.Decimal) query.Component {
	return query.Component{
		Name:            "Additional backup storage",
		Unit:            "GB",
		MonthlyQuantity: backupStorageGB,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr(inst.serviceName),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(".*Single Server - Backup Storage.*")},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
			},
		},
	}
}
