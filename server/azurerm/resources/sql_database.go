package resources

import (
	"fmt"
	"github.com/kaytu-io/infracost/external/schema"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

// SQLDatabase is the entity that holds the logic to calculate price
// of the azurerm_public_ip
type SQLDatabase struct {
	provider *Provider

	location          string
	sku               string
	isElasticPool     bool
	licenseType       string
	tier              string
	family            string
	cores             *int64
	maxSizeGB         *float64
	readReplicaCount  *int64
	zoneRedundant     bool
	backupStorageType string

	// Usage
	// receive override number of GBs used by extra data storage.
	extraDataStorageGB *float64
	// receive monthly number of used vCore-hours for serverless compute.
	monthlyVCoreHours *int64
	// receive number of GBs used by long-term retention backup storage.
	longTermRetentionStorageGB *int64
	// receive number of GBs used by Point-In-Time Restore (PITR) backup storage
	backupStorageGB *int64
}

// sqlDatabaseValues is holds the values that we need to be able
// to calculate the price of the SQLDatabase
type sqlDatabaseValues struct {
	Location                      string   `mapstructure:"location"`
	Sku                           string   `mapstructure:"sku"`
	LicenseType                   string   `mapstructure:"license_type"`
	Tier                          string   `mapstructure:"tier"`
	Family                        string   `mapstructure:"family"`
	Cores                         *int64   `mapstructure:"cores"`
	MaxSizeGB                     *float64 `mapstructure:"max_size_gb"`
	ReadReplicaCount              *int64   `mapstructure:"read_replica_count"`
	ZoneRedundant                 bool     `mapstructure:"zone_redundant"`
	BackupStorageType             string   `mapstructure:"backup_storage_type"`
	Edition                       string   `mapstructure:"edition"`
	RequestedServiceObjectiveName *string  `mapstructure:"requested_service_objective_name"`
	ReadScale                     *bool    `mapstructure:"read_scale"`

	Usage struct {
		ExtraDataStorageGB         *float64 `mapstructure:"extra_data_storage_gb"`
		MonthlyVCoreHours          *int64   `mapstructure:"monthly_vcore_hours"`
		LongTermRetentionStorageGB *int64   `mapstructure:"long_term_retention_storage_gb"`
		BackupStorageGB            *int64   `mapstructure:"backup_storage_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodeSqlDatabaseValues decodes and returns sqlDatabaseValues from a Terraform values map.
func decodeSqlDatabaseValues(tfVals map[string]interface{}) (sqlDatabaseValues, error) {
	var v sqlDatabaseValues
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

func parseSKU(sku string) (skuConfig, error) {
	if dtuMap.usesDTUUnits(sku) {
		return skuConfig{
			sku: sku,
		}, nil
	}

	return parseMSSQLSku(sku)
}

// newSQLDatabase initializes a new SQLDatabase from the provider
func (p *Provider) newSQLDatabase(vals sqlDatabaseValues) *SQLDatabase {
	config := skuConfig{
		sku:    "GP_Gen5_2",
		tier:   "General Purpose",
		family: "Compute Gen5",
		cores:  util.IntPtr(2),
	}

	edition := vals.Edition
	if edition != "" {
		config = skuConfig{
			sku: edition,
		}

		if val, ok := sqlEditionMapping[edition]; ok {
			config = val
		}
	}

	if vals.RequestedServiceObjectiveName != nil {
		var err error
		config, err = parseSKU(*vals.RequestedServiceObjectiveName)
		if err != nil {
			return nil
		}
	}

	var maxSizeGB *float64
	if vals.MaxSizeGB != nil {
		maxBytes := *vals.MaxSizeGB
		if maxBytes > 0 {
			val := maxBytes / 1073741824
			maxSizeGB = &val
		}
	}

	var readReplicas *int64
	if vals.ReadScale != nil {
		if *vals.ReadScale {
			var i int64 = 1
			readReplicas = &i
		}
	}

	inst := &SQLDatabase{
		provider: p,

		location:          vals.Location,
		sku:               config.sku,
		isElasticPool:     false,
		tier:              config.tier,
		family:            config.family,
		cores:             config.cores,
		maxSizeGB:         maxSizeGB,
		readReplicaCount:  readReplicas,
		zoneRedundant:     vals.ZoneRedundant,
		backupStorageType: "Geo",

		extraDataStorageGB:         vals.Usage.ExtraDataStorageGB,
		monthlyVCoreHours:          vals.Usage.MonthlyVCoreHours,
		longTermRetentionStorageGB: vals.Usage.LongTermRetentionStorageGB,
		backupStorageGB:            vals.Usage.BackupStorageGB,
	}
	return inst
}

func (inst *SQLDatabase) Components() []query.Component {
	var costComponents []query.Component
	if inst.isElasticPool {
		costComponents = inst.elasticPoolCostComponents()
		GetCostComponentNamesAndSetLogger(costComponents, inst.provider.logger)
		return costComponents
	}

	if inst.cores != nil {
		costComponents = inst.vCoreCostComponents()
		GetCostComponentNamesAndSetLogger(costComponents, inst.provider.logger)
		return costComponents
	}

	costComponents = inst.dtuCostComponents()
	GetCostComponentNamesAndSetLogger(costComponents, inst.provider.logger)
	return costComponents
}

const (
	sqlServerlessTier = "general purpose - serverless"
	sqlHyperscaleTier = "hyperscale"
)

type skuConfig struct {
	sku    string
	tier   string
	family string
	cores  *int64
}

var (
	mssqlTierMapping = map[string]string{
		"b": "Basic",
		"p": "Premium",
		"s": "Standard",
	}

	mssqlPremiumDTUIncludedStorage = map[string]float64{
		"p1":  500,
		"p2":  500,
		"p4":  500,
		"p6":  500,
		"p11": 4096,
		"p15": 4096,
	}

	mssqlStorageRedundancyTypeMapping = map[string]string{
		"geo":   "RA-GRS",
		"local": "LRS",
		"zone":  "ZRS",
	}

	sqlEditionMapping = map[string]skuConfig{
		"GeneralPurpose": {
			sku:    "GP_Gen5_2",
			tier:   "General Purpose",
			family: "Compute Gen5",
			cores:  util.IntPtr(2),
		},
		"BusinessCritical": {
			sku:    "BC_Gen5_2",
			tier:   "Business Critical",
			family: "Compute Gen5",
			cores:  util.IntPtr(2),
		},
		"Hyperscale": {
			sku:    "HS_Gen5_2",
			tier:   "Hyperscale",
			family: "Compute Gen5",
			cores:  util.IntPtr(2),
		},
		"Standard": {
			sku: "S0",
		},
		"Premium": {
			sku: "P1",
		},
		"DataWarehouse": {
			sku: "DW100c",
		},
		"Stretch": {
			sku: "DS100",
		},
	}
)

func (inst *SQLDatabase) dtuCostComponents() []query.Component {
	skuName := strings.ToLower(inst.sku)
	if skuName == "basic" {
		skuName = "b"
		inst.sku = "B"
	}

	daysInMonth := schema.HourToMonthUnitMultiplier.DivRound(decimal.NewFromInt(24), 24)

	components := []query.Component{
		{
			Name:            fmt.Sprintf("Compute (%s)", strings.ToTitle(inst.sku)),
			Unit:            "hours",
			MonthlyQuantity: daysInMonth,
			ProductFilter: &product.Filter{
				Provider: util.StringPtr("azurerm"),
				Location: util.StringPtr(inst.location),
				Service:  util.StringPtr("SQL Database"),
				Family:   util.StringPtr("Databases"),
				AttributeFilters: []*product.AttributeFilter{
					{Key: "product_name", ValueRegex: util.StringPtr("^SQL Database Single")},
					{Key: "sku_name", ValueRegex: util.StringPtr(fmt.Sprintf("^%s$", inst.sku))},
					{Key: "meter_name", ValueRegex: util.StringPtr("DTU(s)?$")},
				},
			},
			PriceFilter: &price.Filter{
				AttributeFilters: []*price.AttributeFilter{
					{Key: "type", Value: util.StringPtr("Consumption")},
				},
			},
		},
	}

	var extraStorageGB float64

	if !strings.HasPrefix(skuName, "b") && inst.extraDataStorageGB != nil {
		extraStorageGB = *inst.extraDataStorageGB
	} else if strings.HasPrefix(skuName, "s") && inst.maxSizeGB != nil {
		includedStorageGB := 250.0
		extraStorageGB = *inst.maxSizeGB - includedStorageGB
	} else if strings.HasPrefix(skuName, "p") && inst.maxSizeGB != nil {
		includedStorageGB, ok := mssqlPremiumDTUIncludedStorage[skuName]
		if ok {
			extraStorageGB = *inst.maxSizeGB - includedStorageGB
		}
	}

	if extraStorageGB > 0 {
		c := inst.extraDataStorageCostComponent(extraStorageGB)
		if c != nil {
			components = append(components, *c)
		}
	}

	components = append(components, inst.longTermRetentionCostComponent())
	components = append(components, inst.pitrBackupCostComponent())

	return components
}

func (inst *SQLDatabase) vCoreCostComponents() []query.Component {
	components := inst.computeHoursCostComponents()

	if strings.ToLower(inst.tier) == sqlHyperscaleTier {
		components = append(components, inst.readReplicaCostComponent())
	}

	if strings.ToLower(inst.tier) != sqlServerlessTier && strings.ToLower(inst.licenseType) == "licenseincluded" {
		components = append(components, inst.mssqlLicenseCostComponent())
	}

	components = append(components, inst.mssqlStorageCostComponent())

	if strings.ToLower(inst.tier) != sqlHyperscaleTier {
		components = append(components, inst.longTermRetentionCostComponent())
		components = append(components, inst.pitrBackupCostComponent())
	}

	return components
}

func (inst *SQLDatabase) elasticPoolCostComponents() []query.Component {
	return []query.Component{
		inst.longTermRetentionCostComponent(),
		inst.pitrBackupCostComponent(),
	}
}

func (inst *SQLDatabase) computeHoursCostComponents() []query.Component {
	if strings.ToLower(inst.tier) == sqlServerlessTier {
		return inst.serverlessComputeHoursCostComponents()
	}

	return inst.provisionedComputeCostComponents()
}

func (inst *SQLDatabase) serverlessComputeHoursCostComponents() []query.Component {
	productNameRegex := fmt.Sprintf(".*%s - %s.*", inst.tier, inst.family)

	var vCoreHours decimal.Decimal
	if inst.monthlyVCoreHours != nil {
		vCoreHours = decimal.NewFromInt(*inst.monthlyVCoreHours)
	}

	costComponents := []query.Component{
		{
			Name:            fmt.Sprintf("Compute (serverless, %s)", inst.sku),
			Unit:            "vCore-hours",
			MonthlyQuantity: vCoreHours,
			ProductFilter: &product.Filter{
				Provider: util.StringPtr("azurerm"),
				Location: util.StringPtr(inst.location),
				Service:  util.StringPtr("SQL Database"),
				Family:   util.StringPtr("Databases"),
				AttributeFilters: []*product.AttributeFilter{
					{Key: "product_name", ValueRegex: util.StringPtr(productNameRegex)},
					{Key: "sku_name", Value: util.StringPtr("1 vCore")},
					{Key: "meter_name", ValueRegex: util.StringPtr("^(?!.* - Free$).*$")},
				},
			},
			PriceFilter: &price.Filter{
				AttributeFilters: []*price.AttributeFilter{
					{Key: "type", Value: util.StringPtr("Consumption")},
				},
			},
		},
	}

	if inst.zoneRedundant {
		costComponents = append(costComponents, query.Component{
			Name:            fmt.Sprintf("Zone redundancy (serverless, %s)", inst.sku),
			Unit:            "vCore-hours",
			MonthlyQuantity: vCoreHours,
			ProductFilter: &product.Filter{
				Provider: util.StringPtr("azurerm"),
				Location: util.StringPtr(inst.location),
				Service:  util.StringPtr("SQL Database"),
				Family:   util.StringPtr("Databases"),
				AttributeFilters: []*product.AttributeFilter{
					{Key: "product_name", ValueRegex: util.StringPtr(productNameRegex)},
					{Key: "sku_name", Value: util.StringPtr("1 vCore Zone Redundancy")},
					{Key: "meter_name", ValueRegex: util.StringPtr("^(?!.* - Free$).*$")},
				},
			},
			PriceFilter: &price.Filter{
				AttributeFilters: []*price.AttributeFilter{
					{Key: "type", Value: util.StringPtr("Consumption")},
				},
			},
		})
	}

	return costComponents
}

func (inst *SQLDatabase) provisionedComputeCostComponents() []query.Component {
	var cores int64
	if inst.cores != nil {
		cores = *inst.cores
	}

	productNameRegex := fmt.Sprintf(".*%s - %s.*", inst.tier, inst.family)
	name := fmt.Sprintf("Compute (provisioned, %s)", inst.sku)

	components := []query.Component{
		{
			Name:           name,
			Unit:           "hours",
			HourlyQuantity: decimal.NewFromInt(1),
			ProductFilter: &product.Filter{
				Provider: util.StringPtr("azurerm"),
				Location: util.StringPtr(inst.location),
				Service:  util.StringPtr("SQL Database"),
				Family:   util.StringPtr("Databases"),
				AttributeFilters: []*product.AttributeFilter{
					{Key: "product_name", ValueRegex: util.StringPtr(productNameRegex)},
					{Key: "sku_name", Value: util.StringPtr(fmt.Sprintf("%d vCore", cores))},
				},
			},
			PriceFilter: &price.Filter{
				AttributeFilters: []*price.AttributeFilter{
					{Key: "type", Value: util.StringPtr("Consumption")},
				},
			},
		},
	}

	if inst.zoneRedundant {
		components = append(components, query.Component{
			Name:           fmt.Sprintf("Zone redundancy (provisioned, %s)", inst.sku),
			Unit:           "hours",
			HourlyQuantity: decimal.NewFromInt(1),
			ProductFilter: &product.Filter{
				Provider: util.StringPtr("azurerm"),
				Location: util.StringPtr(inst.location),
				Service:  util.StringPtr("SQL Database"),
				Family:   util.StringPtr("Databases"),
				AttributeFilters: []*product.AttributeFilter{
					{Key: "product_name", ValueRegex: util.StringPtr(productNameRegex)},
					{Key: "sku_name", Value: util.StringPtr(fmt.Sprintf("%d vCore Zone Redundancy", cores))},
				},
			},
			PriceFilter: &price.Filter{
				AttributeFilters: []*price.AttributeFilter{
					{Key: "type", Value: util.StringPtr("Consumption")},
				},
			},
		})
	}

	return components
}

func (inst *SQLDatabase) readReplicaCostComponent() query.Component {
	productNameRegex := fmt.Sprintf(".*%s - %s.*", inst.tier, inst.family)
	skuName := mssqlSkuName(*inst.cores, inst.zoneRedundant)

	var replicaCount decimal.Decimal
	if inst.readReplicaCount != nil {
		replicaCount = decimal.NewFromInt(*inst.readReplicaCount)
	}

	return query.Component{
		Name:           "Read replicas",
		Unit:           "hours",
		HourlyQuantity: replicaCount,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("SQL Database"),
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

func (inst *SQLDatabase) longTermRetentionCostComponent() query.Component {
	var retention decimal.Decimal
	if inst.longTermRetentionStorageGB != nil {
		retention = decimal.NewFromInt(*inst.longTermRetentionStorageGB)
	}

	redundancyType, ok := mssqlStorageRedundancyTypeMapping[strings.ToLower(inst.backupStorageType)]
	if !ok {
		redundancyType = "RA-GRS"
	}

	return query.Component{
		Name:            fmt.Sprintf("Long-term retention (%s)", redundancyType),
		Unit:            "GB",
		MonthlyQuantity: retention,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("SQL Database"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("SQL Database - LTR Backup Storage")},
				{Key: "sku_name", Value: util.StringPtr(fmt.Sprintf("Backup %s", redundancyType))},
				{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf(".*%s Data Stored.*", redundancyType))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *SQLDatabase) pitrBackupCostComponent() query.Component {
	var pitrGB decimal.Decimal
	if inst.backupStorageGB != nil {
		pitrGB = decimal.NewFromInt(*inst.backupStorageGB)
	}

	redundancyType, ok := mssqlStorageRedundancyTypeMapping[strings.ToLower(inst.backupStorageType)]
	if !ok {
		redundancyType = "RA-GRS"
	}

	return query.Component{
		Name:            fmt.Sprintf("PITR backup storage (%s)", redundancyType),
		Unit:            "GB",
		MonthlyQuantity: pitrGB,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("SQL Database"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(".*PITR Backup Storage.*")},
				{Key: "sku_name", Value: util.StringPtr(fmt.Sprintf("Backup %s", redundancyType))},
				{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf(".*%s Data Stored.*", redundancyType))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *SQLDatabase) extraDataStorageCostComponent(extraStorageGB float64) *query.Component {
	tier := inst.tier
	if tier == "" {
		var ok bool
		tier, ok = mssqlTierMapping[strings.ToLower(inst.sku)[:1]]

		if !ok {
			return nil
		}
	}
	component := mssqlExtraDataStorageCostComponent(inst.location, tier, extraStorageGB)
	return &component
}

func mssqlSkuName(cores int64, zoneRedundant bool) string {
	sku := fmt.Sprintf("%d vCore", cores)

	if zoneRedundant {
		sku += " Zone Redundancy"
	}
	return sku
}

func (inst *SQLDatabase) mssqlProductFilter(filters []*product.AttributeFilter) *product.Filter {
	return &product.Filter{
		Provider:         util.StringPtr("azurerm"),
		Location:         util.StringPtr(inst.location),
		Service:          util.StringPtr("SQL Database"),
		Family:           util.StringPtr("Databases"),
		AttributeFilters: filters,
	}
}

func mssqlExtraDataStorageCostComponent(region string, tier string, extraStorageGB float64) query.Component {
	return query.Component{
		Name:            "Extra data storage",
		Unit:            "GB",
		MonthlyQuantity: decimal.NewFromFloat(extraStorageGB),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(region),
			Service:  util.StringPtr("SQL Database"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", ValueRegex: util.StringPtr(fmt.Sprintf("^%s$", tier))},
				{Key: "meter_name", Value: util.StringPtr("Data Stored")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *SQLDatabase) mssqlLicenseCostComponent() query.Component {
	licenseRegion := "Global"
	if strings.Contains(inst.location, "usgov") {
		licenseRegion = "US Gov"
	}

	if strings.Contains(inst.location, "china") {
		licenseRegion = "China"
	}

	if strings.Contains(inst.location, "germany") {
		licenseRegion = "Germany"
	}

	coresVal := int64(1)
	if inst.cores != nil {
		coresVal = *inst.cores
	}

	return query.Component{
		Name:           "SQL license",
		Unit:           "vCore-hours",
		HourlyQuantity: decimal.NewFromInt(coresVal),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(licenseRegion),
			Service:  util.StringPtr("SQL Database"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(fmt.Sprintf(".*%s - %s.*", inst.tier, "SQL License"))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *SQLDatabase) mssqlStorageCostComponent() query.Component {
	storageGB := decimal.NewFromInt(5)
	if inst.maxSizeGB != nil {
		storageGB = decimal.NewFromFloat(*inst.maxSizeGB)
	}

	storageTier := inst.tier
	if strings.ToLower(storageTier) == "general purpose - serverless" {
		storageTier = "General Purpose"
	}

	skuName := storageTier
	if inst.zoneRedundant {
		skuName += " Zone Redundancy"
	}

	productNameRegex := fmt.Sprintf(".*%s - Storage.*", storageTier)

	return query.Component{
		Name:            "Storage",
		Unit:            "GB",
		MonthlyQuantity: storageGB,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("SQL Database"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(productNameRegex)},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
				{Key: "meter_name", ValueRegex: util.StringPtr(".*Data Stored")},
			},
		},
	}
}
