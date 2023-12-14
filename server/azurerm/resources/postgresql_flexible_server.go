package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"regexp"
	"strings"
)

// PostgresqlFlexibleServer is the entity that holds the logic to calculate price
// of the azurerm_mariadb_server
type PostgresqlFlexibleServer struct {
	provider *Provider

	location        string
	sku             string
	tier            string
	instanceType    string
	instanceVersion string
	storage         int64

	// Usage
	additionalBackupStorageGb *float64
}

// postgresqlFlexibleServerValues is holds the values that we need to be able
// to calculate the price of the PostgresqlFlexibleServer
type postgresqlFlexibleServerValues struct {
	Location  string `mapstructure:"location"`
	SkuName   string `mapstructure:"sku_name"`
	StorageMb int64  `mapstructure:"storage_mb"`

	Usage struct {
		AdditionalBackupStorageGb *float64 `mapstructure:"additional_backup_storage_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodePostgresqlFlexibleServerValues decodes and returns postgresqlFlexibleServerValues from a Terraform values map.
func decodePostgresqlFlexibleServerValues(tfVals map[string]interface{}) (postgresqlFlexibleServerValues, error) {
	var v postgresqlFlexibleServerValues
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

// newPostgresqlFlexibleServer initializes a new PostgresqlFlexibleServer from the provider
func (p *Provider) newPostgresqlFlexibleServer(vals postgresqlFlexibleServerValues) *PostgresqlFlexibleServer {
	var tier, size, version string

	s := strings.Split(vals.SkuName, "_")
	if len(s) < 3 || len(s) > 4 {
		return nil
	}

	if len(s) > 2 {
		tier = strings.ToLower(s[0])
		size = s[2]
	}

	if len(s) > 3 {
		version = s[3]
	}

	supportedTiers := []string{"b", "gp", "mo"}
	if !contains(supportedTiers, tier) {
		return nil
	}

	if tier != "b" {
		coreRegex := regexp.MustCompile(`(\d+)`)
		match := coreRegex.FindStringSubmatch(size)
		if len(match) < 1 {
			return nil
		}
	}

	inst := &PostgresqlFlexibleServer{
		provider: p,

		location:        getLocationName(vals.Location),
		sku:             vals.SkuName,
		tier:            tier,
		instanceType:    size,
		instanceVersion: version,
		storage:         vals.StorageMb,

		additionalBackupStorageGb: vals.Usage.AdditionalBackupStorageGb,
	}
	return inst
}

func (inst *PostgresqlFlexibleServer) Components() []query.Component {
	var components []query.Component

	components = append(components, inst.computeCostComponent(), inst.backupCostComponent(), inst.storageCostComponent())

	return components
}

// computeCostComponent returns a cost component for server compute requirements.
func (inst *PostgresqlFlexibleServer) computeCostComponent() query.Component {
	attrs := getFlexibleServerFilterAttributes(inst.tier, inst.instanceType, inst.instanceVersion)

	return query.Component{
		Name:           fmt.Sprintf("Compute (%s)", inst.sku),
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Azure Database for PostgreSQL"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(fmt.Sprintf(".*%s.*", attrs.TierName))},
				{Key: "product_name", ValueRegex: util.StringPtr(fmt.Sprintf(".*%s.*", attrs.Series))},
				{Key: "sku_name", ValueRegex: util.StringPtr(fmt.Sprintf("^(?i)%s$", attrs.SKUName))},
				{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf("^(?i)%s$", attrs.MeterName))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

// storageCostComponent returns a cost component for server's storage.
func (inst *PostgresqlFlexibleServer) storageCostComponent() query.Component {
	var quantity decimal.Decimal
	if inst.storage > 0 {
		// Storage is in MB
		quantity = decimal.NewFromInt(inst.storage / 1024)
	}

	return query.Component{
		Name:            "Storage",
		Unit:            "GB",
		MonthlyQuantity: quantity,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Azure Database for PostgreSQL"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Az DB for PostgreSQL Flexible Server Storage")},
				{Key: "meter_name", Value: util.StringPtr("Storage Data Stored")},
			},
		},
	}
}

// backupCostComponent returns a cost component for additional backup storage.
func (inst *PostgresqlFlexibleServer) backupCostComponent() query.Component {
	var quantity decimal.Decimal
	if inst.additionalBackupStorageGb != nil {
		quantity = decimal.NewFromFloat(*inst.additionalBackupStorageGb)
	}

	return query.Component{
		Name:            "Additional backup storage",
		Unit:            "GB",
		MonthlyQuantity: quantity,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Azure Database for PostgreSQL"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Azure Database for PostgreSQL Flexible Server Backup Storage")},
				{Key: "meter_name", Value: util.StringPtr("Backup Storage LRS Data Stored")},
			},
		},
	}
}

// flexibleServerFilterAttributes defines CPAPI filter attributes for compute
// cost component derived from IaC provider's SKU.
type flexibleServerFilterAttributes struct {
	SKUName   string
	TierName  string
	MeterName string
	Series    string
}

// getFlexibleServerFilterAttributes returns a struct with CPAPI filter
// attributes based on values extracted from IaC provider's SKU.
func getFlexibleServerFilterAttributes(tier, instanceType, instanceVersion string) flexibleServerFilterAttributes {
	var skuName, meterName, series string

	tierName := map[string]string{
		"b":  "Burstable",
		"gp": "General Purpose",
		"mo": "Memory Optimized",
	}[tier]

	if tier == "b" {
		meterName = instanceType
		skuName = instanceType
		series = "BS"
	} else {
		meterName = "vCore"

		coreRegex := regexp.MustCompile(`(\d+)`)
		match := coreRegex.FindStringSubmatch(instanceType)
		cores := match[1]
		skuName = fmt.Sprintf("%s vCore", cores)

		series = coreRegex.ReplaceAllString(instanceType, "") + instanceVersion
	}

	return flexibleServerFilterAttributes{
		SKUName:   skuName,
		TierName:  tierName,
		MeterName: meterName,
		Series:    series,
	}
}
