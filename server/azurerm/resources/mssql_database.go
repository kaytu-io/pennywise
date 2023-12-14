package resources

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"strconv"
	"strings"
)

type ServerId struct {
	Values struct {
		Location string `mapstructure:"location"`
	} `mapstructure:"values"`
}

// mssqlDatabaseValues is holds the values that we need to be able
// to calculate the price of the MssqlDatabase
type mssqlDatabaseValues struct {
	SkuName            *string  `mapstructure:"sku_name"`
	MaxSizeGb          *float64 `mapstructure:"max_size_gb"`
	ServerId           ServerId `mapstructure:"server_id"`
	ReadReplicaCount   *int64   `mapstructure:"read_replica_count"`
	LicenseType        *string  `mapstructure:"license_type"`
	StorageAccountType *string  `mapstructure:"storage_account_type"`
	ZoneRedundant      *bool    `mapstructure:"zone_redundant"`
	ElasticPoolId      *string  `mapstructure:"elastic_pool_id"`

	Usage struct {
	} `mapstructure:"pennywise_usage"`
}

func parseMSSQLSku(sku string) (skuConfig, error) {
	s := strings.Split(sku, "_")
	if len(s) < 3 {
		return skuConfig{}, fmt.Errorf("unrecognized MSSQL SKU format for resource: %s", sku)
	}

	tierKey := strings.ToLower(strings.Join(s[0:len(s)-2], "_"))
	tier, ok := sqlTierMapping[tierKey]
	if !ok {
		return skuConfig{}, fmt.Errorf("invalid tier in MSSQL SKU for resourcs: %s", sku)
	}

	familyKey := strings.ToLower(s[len(s)-2])
	family, ok := sqlFamilyMapping[familyKey]
	if !ok {
		return skuConfig{}, fmt.Errorf("invalid family in MSSQL SKU for resource: %s", sku)
	}

	cores, err := strconv.ParseInt(s[len(s)-1], 10, 64)
	if err != nil {
		return skuConfig{}, fmt.Errorf("invalid core count in MSSQL SKU for resource: %s", sku)
	}

	return skuConfig{
		sku:    sku,
		tier:   tier,
		family: family,
		cores:  &cores,
	}, nil
}

// decodeMssqlDatabaseValues decodes and returns mssqlDatabaseValues from a Terraform values map.
func decodeMssqlDatabaseValues(tfVals map[string]interface{}) (mssqlDatabaseValues, error) {
	var v mssqlDatabaseValues
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

// newMssqlDatabase initializes a new MssqlDatabase from the provider
func (p *Provider) newMssqlDatabase(vals mssqlDatabaseValues) *SQLDatabase {

	region := vals.ServerId.Values.Location

	sku := "GP_S_Gen5_2"
	if vals.SkuName != nil {
		sku = *vals.SkuName
	}

	var maxSize *float64
	if vals.MaxSizeGb != nil {
		val := *vals.MaxSizeGb
		maxSize = &val
	}

	var replicaCount *int64
	if vals.ReadReplicaCount != nil {
		val := *vals.ReadReplicaCount
		replicaCount = &val
	}

	licenseType := "LicenseIncluded"
	if vals.LicenseType != nil {
		licenseType = *vals.LicenseType
	}
	storageAccountType := "Geo"
	if vals.StorageAccountType != nil {
		storageAccountType = *vals.StorageAccountType
	}
	zoneRedundant := false
	if vals.ZoneRedundant != nil {
		zoneRedundant = *vals.ZoneRedundant
	}

	inst := SQLDatabase{
		location:          region,
		sku:               sku,
		licenseType:       licenseType,
		maxSizeGB:         maxSize,
		readReplicaCount:  replicaCount,
		zoneRedundant:     zoneRedundant,
		backupStorageType: storageAccountType,
	}

	if strings.ToLower(sku) == "elasticpool" || vals.ElasticPoolId != nil {
		inst.isElasticPool = true
	} else if !dtuMap.usesDTUUnits(sku) {
		c, err := parseMSSQLSku(sku)
		if err != nil {
			return nil
		}

		inst.tier = c.tier
		inst.family = c.family
		inst.cores = c.cores
	}

	return &inst
}

var (
	sqlTierMapping = map[string]string{
		"gp":   "General Purpose",
		"gp_s": "General Purpose - Serverless",
		"hs":   "Hyperscale", // TODO: SingleDB or not
		"bc":   "Business Critical",
	}

	sqlFamilyMapping = map[string]string{
		"gen5": "Compute Gen5",
		"gen4": "Compute Gen4",
		"m":    "Compute M Series",
	}

	dtuMap = dtuMapping{
		"free":  true,
		"basic": true,

		"s": true,
		"d": true,
		"p": true,
	}
)

type dtuMapping map[string]bool

func (d dtuMapping) usesDTUUnits(sku string) bool {
	sanitized := strings.ToLower(sku)
	if d[sanitized] {
		return true
	}

	if sanitized == "" {
		return false
	}

	return d[sanitized[0:1]]
}
