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

// StorageShare is the entity that holds the logic to calculate price
// of the azurerm_storage_share
type StorageShare struct {
	provider *Provider

	location               string
	accountReplicationType string
	accessTierString       string
	accountKind            *string
	quota                  int64

	// Usage
	storageGB               *int64
	snapshotsStorageGB      *int64
	monthlyReadOperations   *int64
	monthlyWriteOperations  *int64
	monthlyListOperations   *int64
	monthlyOtherOperations  *int64
	monthlyDataRetrievalGB  *int64
	metadataAtRestStorageGB *int64
}

// storageShareValues is holds the values that we need to be able
// to calculate the price of the StorageShare
type storageShareValues struct {
	StorageAccountName StorageAccountName `mapstructure:"storage_account_name"`
	Quota              int64              `mapstructure:"quota"`
	AccessTier         string             `mapstructure:"access_tier"`

	Usage struct {
		// receive Total size of storage in GB. Overrides any provided `quota`
		StorageGb *int64 `mapstructure:"storage_gb"`
		// receive total size of Snapshots in GB
		SnapshotsStorageGB *int64 `mapstructure:"snapshots_storage_gb"`
		// receive monthly number of Read operations.
		MonthlyReadOperations *int64 `mapstructure:"monthly_read_operations"`
		// receive monthly number of Write operations.
		MonthlyWriteOperations *int64 `mapstructure:"monthly_write_operations"`
		// receive monthly number of List and Create Container operations.
		MonthlyListOperations *int64 `mapstructure:"monthly_list_operations"`
		// receive monthly number of All other operations.
		MonthlyOtherOperations *int64 `mapstructure:"monthly_other_operations"`
		// receive monthly number of data retrieval in GB.
		MonthlyDataRetrievalGB *int64 `mapstructure:"monthly_data_retrieval_gb"`
		// receive total size of Metadata in GB
		MetadataAtRestStorageGB *int64 `mapstructure:"metadata_at_rest_storage_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodeStorageShareValues decodes and returns publicIPValues from a Terraform values map.
func decodeStorageShareValues(tfVals map[string]interface{}) (storageShareValues, error) {
	var v storageShareValues
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

// newStorageShare initializes a new PublicIP from the provider
func (p *Provider) newStorageShare(vals storageShareValues) *StorageShare {
	accountReplicationType := "LRS"

	accessTier := vals.AccessTier
	if accessTier == "" {
		accessTier = "TransactionOptimized"
	}

	if vals.StorageAccountName.Values.AccountReplicationType != nil {
		accountReplicationType = *vals.StorageAccountName.Values.AccountReplicationType
	}

	inst := &StorageShare{
		provider:               p,
		location:               vals.StorageAccountName.Values.Location,
		accountReplicationType: accountReplicationType,
		accessTierString:       accessTier,
		quota:                  vals.Quota,
		accountKind:            vals.StorageAccountName.Values.AccountKind,

		storageGB:               vals.Usage.StorageGb,
		snapshotsStorageGB:      vals.Usage.SnapshotsStorageGB,
		monthlyReadOperations:   vals.Usage.MonthlyReadOperations,
		monthlyWriteOperations:  vals.Usage.MonthlyWriteOperations,
		monthlyListOperations:   vals.Usage.MonthlyListOperations,
		monthlyOtherOperations:  vals.Usage.MonthlyOtherOperations,
		monthlyDataRetrievalGB:  vals.Usage.MonthlyDataRetrievalGB,
		metadataAtRestStorageGB: vals.Usage.MetadataAtRestStorageGB,
	}
	return inst
}

func (inst *StorageShare) Components() []query.Component {
	var components []query.Component

	if inst.accountKind != nil {
		accountKind := *inst.accountKind
		if strings.EqualFold(inst.accessTierString, "premium") && !strings.EqualFold(accountKind, "filestorage") {
			return components
		}
	}

	components = append(components, inst.dataStorageCostComponent())
	components = append(components, inst.snapshotCostComponents())
	components = append(components, inst.metadataCostComponents()...)
	components = append(components, inst.readOperationsCostComponents()...)
	components = append(components, inst.writeOperationsCostComponents()...)
	components = append(components, inst.listOperationsCostComponents()...)
	components = append(components, inst.otherOperationsCostComponents()...)
	components = append(components, inst.dataRetrievalCostComponents()...)

	return components
}

func (inst *StorageShare) productName() string {
	if inst.accessTier() == "Premium" {
		return "Premium Files"
	}

	return "Files v2"
}

func (inst *StorageShare) accessTier() string {
	return map[string]string{
		"hot":                  "Hot",
		"cool":                 "Cool",
		"transactionoptimized": "Standard",
		"premium":              "Premium",
	}[strings.ToLower(inst.accessTierString)]
}

func (inst *StorageShare) dataStorageCostComponent() query.Component {
	var qty decimal.Decimal

	if inst.accessTier() == "Premium" {
		qty = decimal.NewFromInt(inst.quota)
	}

	if inst.storageGB != nil {
		qty = decimal.NewFromInt(*inst.storageGB)
	}

	skuName := fmt.Sprintf("%s %s", inst.accessTier(), strings.ToUpper(inst.accountReplicationType))
	meterName := "Data Stored"
	if inst.accessTier() == "Premium" {
		meterName = "Provisioned"
	}

	return query.Component{
		Name:            "Data at rest",
		Unit:            "GB",
		MonthlyQuantity: qty,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr(inst.productName())},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
				{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf(".*%s", meterName))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr("0")},
			},
		},
	}
}

func (inst *StorageShare) snapshotCostComponents() query.Component {
	var qty decimal.Decimal
	if inst.snapshotsStorageGB != nil {
		qty = decimal.NewFromInt(*inst.snapshotsStorageGB)
	}

	skuName := fmt.Sprintf("%s %s", inst.accessTier(), strings.ToUpper(inst.accountReplicationType))
	meterName := "Data Stored"
	if inst.accessTier() == "Premium" {
		meterName = "Snapshots"
	}

	return query.Component{
		Name:            "Snapshots",
		Unit:            "GB",
		MonthlyQuantity: qty,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr(inst.productName())},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
				{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf(".*%s", meterName))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr("0")},
			},
		},
	}
}

func (inst *StorageShare) metadataCostComponents() []query.Component {
	if contains([]string{"Premium", "Standard"}, inst.accessTier()) {
		return []query.Component{}
	}

	var qty decimal.Decimal
	if inst.metadataAtRestStorageGB != nil {
		qty = decimal.NewFromInt(*inst.metadataAtRestStorageGB)
	}

	skuName := fmt.Sprintf("%s %s", inst.accessTier(), strings.ToUpper(inst.accountReplicationType))

	return []query.Component{{
		Name:            "Metadata at rest",
		Unit:            "GB",
		MonthlyQuantity: qty,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr(inst.productName())},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
				{Key: "meter_name", ValueRegex: util.StringPtr(".*Metadata")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr("0")},
			},
		},
	}}
}

func (inst *StorageShare) readOperationsCostComponents() []query.Component {
	if inst.accessTier() == "Premium" {
		return []query.Component{}
	}

	var qty decimal.Decimal
	if inst.monthlyReadOperations != nil {
		qty = decimal.NewFromInt(*inst.monthlyReadOperations).Div(decimal.NewFromInt(10000))
	}

	skuName := fmt.Sprintf("%s %s", inst.accessTier(), strings.ToUpper(inst.accountReplicationType))

	return []query.Component{{
		Name:            "Read operations",
		Unit:            "10k operations",
		MonthlyQuantity: qty,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr(inst.productName())},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
				{Key: "meter_name", ValueRegex: util.StringPtr(".*Read Operations")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr("0")},
			},
		},
	}}
}

func (inst *StorageShare) writeOperationsCostComponents() []query.Component {
	if inst.accessTier() == "Premium" {
		return []query.Component{}
	}

	var qty decimal.Decimal
	if inst.monthlyWriteOperations != nil {
		qty = decimal.NewFromInt(*inst.monthlyWriteOperations).Div(decimal.NewFromInt(10000))
	}

	skuName := fmt.Sprintf("%s %s", inst.accessTier(), strings.ToUpper(inst.accountReplicationType))

	return []query.Component{{
		Name:            "Write operations",
		Unit:            "10k operations",
		MonthlyQuantity: qty,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr(inst.productName())},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
				{Key: "meter_name", ValueRegex: util.StringPtr(".*Write Operations")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr("0")},
			},
		},
	}}
}

func (inst *StorageShare) listOperationsCostComponents() []query.Component {
	if inst.accessTier() == "Premium" {
		return []query.Component{}
	}

	var qty decimal.Decimal
	if inst.monthlyListOperations != nil {
		qty = decimal.NewFromInt(*inst.monthlyListOperations).Div(decimal.NewFromInt(10000))
	}

	skuName := fmt.Sprintf("%s %s", inst.accessTier(), strings.ToUpper(inst.accountReplicationType))

	return []query.Component{{
		Name:            "List operations",
		Unit:            "10k operations",
		MonthlyQuantity: qty,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr(inst.productName())},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
				{Key: "meter_name", ValueRegex: util.StringPtr(".*List Operations")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr("0")},
			},
		},
	}}
}

func (inst *StorageShare) otherOperationsCostComponents() []query.Component {
	if inst.accessTier() == "Premium" {
		return []query.Component{}
	}

	var qty decimal.Decimal
	if inst.monthlyOtherOperations != nil {
		qty = decimal.NewFromInt(*inst.monthlyOtherOperations).Div(decimal.NewFromInt(10000))
	}

	skuName := fmt.Sprintf("%s %s", inst.accessTier(), strings.ToUpper(inst.accountReplicationType))
	meterName := "Other Operations"
	if inst.accessTier() == "Standard" {
		meterName = "Protocol Operations"
	}

	return []query.Component{
		{
			Name:            "Other operations",
			Unit:            "10k operations",
			MonthlyQuantity: qty,
			ProductFilter: &product.Filter{
				Provider: util.StringPtr(inst.provider.key),
				Location: util.StringPtr(inst.location),
				Service:  util.StringPtr("Storage"),
				Family:   util.StringPtr("Storage"),
				AttributeFilters: []*product.AttributeFilter{
					{Key: "product_name", Value: util.StringPtr(inst.productName())},
					{Key: "sku_name", Value: util.StringPtr(skuName)},
					{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf(".*%s", meterName))},
				},
			},
			PriceFilter: &price.Filter{
				AttributeFilters: []*price.AttributeFilter{
					{Key: "type", Value: util.StringPtr("Consumption")},
					{Key: "tier_minimum_units", Value: util.StringPtr("0")},
				},
			},
		}}
}

func (inst *StorageShare) dataRetrievalCostComponents() []query.Component {
	if contains([]string{"Premium", "Standard", "Hot"}, inst.accessTier()) || strings.ToUpper(inst.accountReplicationType) == "GZRS" {
		return []query.Component{}
	}

	var qty decimal.Decimal
	if inst.monthlyDataRetrievalGB != nil {
		qty = decimal.NewFromInt(*inst.monthlyDataRetrievalGB)
	}

	skuName := fmt.Sprintf("%s %s", inst.accessTier(), strings.ToUpper(inst.accountReplicationType))

	return []query.Component{
		{
			Name:            "Data retrieval",
			Unit:            "GB",
			MonthlyQuantity: qty,
			ProductFilter: &product.Filter{
				Provider: util.StringPtr(inst.provider.key),
				Location: util.StringPtr(inst.location),
				Service:  util.StringPtr("Storage"),
				Family:   util.StringPtr("Storage"),
				AttributeFilters: []*product.AttributeFilter{
					{Key: "product_name", Value: util.StringPtr(inst.productName())},
					{Key: "sku_name", Value: util.StringPtr(skuName)},
					{Key: "meter_name", ValueRegex: util.StringPtr(".*Data Retrieval")},
				},
			},
			PriceFilter: &price.Filter{
				AttributeFilters: []*price.AttributeFilter{
					{Key: "type", Value: util.StringPtr("Consumption")},
					{Key: "tier_minimum_units", Value: util.StringPtr("0")},
				},
			},
		}}
}
