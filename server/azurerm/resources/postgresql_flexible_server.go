package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"regexp"
	"strings"
)

// PostgresqlFlexibleServer is the entity that holds the logic to calculate price
// of the azurerm_postgresql_flexible_server
type PostgresqlFlexibleServer struct {
	provider *Provider

	location        string
	sku             string
	tier            string
	instanceType    string
	instanceVersion string
	storage         int64

	// Usage
	// receive additional backup storage in GB. If geo-redundancy is enabled, you should set this to twice the required storage capacity.
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

func (inst *PostgresqlFlexibleServer) Components() []resource.Component {
	var components []resource.Component

	components = append(components, inst.computeCostComponent(), inst.backupCostComponent(), inst.storageCostComponent())
	GetCostComponentNamesAndSetLogger(components, inst.provider.logger)

	return components
}

func (inst *PostgresqlFlexibleServer) computeCostComponent() resource.Component {
	attrs := getFlexibleServerFilterAttributes(inst.tier, inst.instanceType, inst.instanceVersion)

	return resource.Component{
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

func (inst *PostgresqlFlexibleServer) storageCostComponent() resource.Component {
	var quantity decimal.Decimal
	if inst.storage > 0 {
		quantity = decimal.NewFromInt(inst.storage / 1024)
	}

	return resource.Component{
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

func (inst *PostgresqlFlexibleServer) backupCostComponent() resource.Component {
	var quantity decimal.Decimal
	if inst.additionalBackupStorageGb != nil {
		quantity = decimal.NewFromFloat(*inst.additionalBackupStorageGb)
	}

	return resource.Component{
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

type flexibleServerFilterAttributes struct {
	SKUName   string
	TierName  string
	MeterName string
	Series    string
}

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
