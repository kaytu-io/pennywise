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

// MysqlFlexibleServer is the entity that holds the logic to calculate price
// of the azurerm_mysql_flexible_server
type MysqlFlexibleServer struct {
	provider *Provider

	location        string
	sku             string
	tier            string
	instanceType    string
	instanceVersion string
	storage         int64
	iops            int64

	// Usage
	// receive additional backup storage in GB. If geo-redundancy is enabled, you should set this to twice the required storage capacity.
	additionalBackupStorageGb *float64
}

// mysqlFlexibleServerValues is holds the values that we need to be able
// to calculate the price of the MysqlFlexibleServer
type mysqlFlexibleServerValues struct {
	Location string `mapstructure:"location"`
	SkuName  string `mapstructure:"sku_name"`
	Storage  []struct {
		Iops   *int64 `mapstructure:"iops"`
		SizeGb *int64 `mapstructure:"size_gb"`
	} `mapstructure:"storage"`

	Usage struct {
		AdditionalBackupStorageGb *float64 `mapstructure:"additional_backup_storage_gb"`
	} `mapstructure:"pennywise_usage"`
}

func decodeMysqlFlexibleServerValues(tfVals map[string]interface{}) (mysqlFlexibleServerValues, error) {
	var v mysqlFlexibleServerValues
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

func (p *Provider) newMysqlFlexibleServer(vals mysqlFlexibleServerValues) *MysqlFlexibleServer {
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

	storage := int64(0)
	iops := int64(0)
	if len(vals.Storage) > 0 {
		if vals.Storage[0].SizeGb != nil {
			storage = *vals.Storage[0].SizeGb
		}
		if vals.Storage[0].Iops != nil {
			iops = *vals.Storage[0].Iops
		}
	}

	inst := &MysqlFlexibleServer{
		provider: p,

		location:        getLocationName(vals.Location),
		sku:             vals.SkuName,
		tier:            tier,
		instanceType:    size,
		instanceVersion: version,
		storage:         storage,
		iops:            iops,

		additionalBackupStorageGb: vals.Usage.AdditionalBackupStorageGb,
	}
	return inst
}

func (inst *MysqlFlexibleServer) Components() []query.Component {
	var components []query.Component

	components = append(components, inst.computeCostComponent(), inst.backupCostComponent(), inst.storageCostComponent(), inst.iopsCostComponent())

	return components
}

func (inst *MysqlFlexibleServer) computeCostComponent() query.Component {
	attrs := getFlexibleServerFilterAttributes(inst.tier, inst.instanceType, inst.instanceVersion)

	tierName := attrs.TierName
	if tierName == "Memory Optimized" {
		tierName = "Business Critical"
	}

	if tierName == "Business Critical" && attrs.Series == "Edsv4" {
		attrs.Series = ""
	}

	return query.Component{
		Name:           fmt.Sprintf("Compute (%s)", inst.sku),
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Azure Database for MySQL"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(fmt.Sprintf(".*%s.*", tierName))},
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

func (inst *MysqlFlexibleServer) storageCostComponent() query.Component {
	var quantity decimal.Decimal
	if inst.storage == 0 {
		quantity = decimal.NewFromInt(20)
	} else {
		quantity = decimal.NewFromInt(inst.storage)
	}

	return query.Component{
		Name:            "Storage",
		Unit:            "GB",
		MonthlyQuantity: quantity,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Azure Database for MySQL"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Az DB for MySQL Flexible Server Storage")},
				{Key: "meter_name", Value: util.StringPtr("Storage Data Stored")},
			},
		},
	}
}

func (inst *MysqlFlexibleServer) iopsCostComponent() query.Component {
	var freeIOPS int64 = 360

	iops := inst.iops
	if iops == 0 {
		iops = freeIOPS
	}

	additionalIOPS := iops - freeIOPS

	if additionalIOPS < 0 {
		additionalIOPS = 0
	}

	return query.Component{
		Name:            "Additional IOPS",
		Unit:            "IOPS",
		MonthlyQuantity: decimal.NewFromInt(additionalIOPS),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Azure Database for MySQL"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Az DB for MySQL Flexible Server Storage")},
				{Key: "sku_name", Value: util.StringPtr("Additional IOPS")},
			},
		},
	}
}

func (inst *MysqlFlexibleServer) backupCostComponent() query.Component {
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
			Service:  util.StringPtr("Azure Database for MySQL"),
			Family:   util.StringPtr("Databases"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Az DB for MySQL Flex Svr Backup Storage")},
				{Key: "meter_name", Value: util.StringPtr("Backup Storage LRS Data Stored")},
			},
		},
	}
}
