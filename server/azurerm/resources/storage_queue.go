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

type StorageAccountName struct {
	Values struct {
		Location               string  `mapstructure:"location"`
		AccountReplicationType *string `mapstructure:"account_replication_type"`
		AccountKind            *string `mapstructure:"account_kind"`
		AccountTier            *string `mapstructure:"account_tier"`
	} `mapstructure:"values"`
}

// StorageQueue is the entity that holds the logic to calculate price
// of the azurerm_storage_queue
type StorageQueue struct {
	provider *Provider

	location               string
	accountReplicationType string
	accountKind            string
	accountTier            *string

	// Usage
	monthlyGeoReplicationDataTransferGB *int64
	monthlyStorageGB                    *int64
	monthlyClass1Operations             *int64
	monthlyClass2Operations             *int64
}

// storageQueueValues is holds the values that we need to be able
// to calculate the price of the StorageQueue
type storageQueueValues struct {
	StorageAccountName StorageAccountName `mapstructure:"storage_account_name"`

	Usage struct {
		MonthlyStorageGb                    *int64 `mapstructure:"monthly_storage_gb"`
		MonthlyClass1Operations             *int64 `mapstructure:"monthly_class_1_operations"`
		MonthlyClass2Operations             *int64 `mapstructure:"monthly_class_2_operations"`
		MonthlyGeoReplicationDataTransferGB *int64 `mapstructure:"monthly_geo_replication_data_transfer_gb"`
	} `mapstructure:"tc_usage"`
}

// decodeStorageQueueValues decodes and returns storageQueueValues from a Terraform values map.
func decodeStorageQueueValues(tfVals map[string]interface{}) (storageQueueValues, error) {
	var v storageQueueValues
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

// newStorageQueue initializes a new StorageQueue from the provider
func (p *Provider) newStorageQueue(vals storageQueueValues) *StorageQueue {
	accountReplicationType := "LRS"
	accountKind := "StorageV2"

	if vals.StorageAccountName.Values.AccountKind != nil {
		accountKind = *vals.StorageAccountName.Values.AccountKind
	}
	if vals.StorageAccountName.Values.AccountReplicationType != nil {
		accountReplicationType = *vals.StorageAccountName.Values.AccountReplicationType
	}
	switch strings.ToLower(accountReplicationType) {
	case "ragrs":
		accountReplicationType = "RA-GRS"
	case "ragzrs":
		accountReplicationType = "RA-GZRS"
	}

	inst := &StorageQueue{
		provider: p,

		location:               vals.StorageAccountName.Values.Location,
		accountReplicationType: accountReplicationType,
		accountKind:            accountKind,
		accountTier:            vals.StorageAccountName.Values.AccountTier,

		monthlyStorageGB:                    vals.Usage.MonthlyStorageGb,
		monthlyClass1Operations:             vals.Usage.MonthlyClass1Operations,
		monthlyClass2Operations:             vals.Usage.MonthlyClass2Operations,
		monthlyGeoReplicationDataTransferGB: vals.Usage.MonthlyGeoReplicationDataTransferGB,
	}
	return inst
}

func (inst *StorageQueue) Components() []query.Component {
	var components []query.Component

	if !inst.isAccountKindSupported() {
		fmt.Println("!inst.isAccountKindSupported()")
		return nil
	}

	if !inst.isReplicationTypeSupported() {
		fmt.Println("!inst.isReplicationTypeSupported()")
		return nil
	}

	components = append(components, inst.dataStorageCostComponent())
	components = append(components, inst.operationsCostComponents()...)
	components = append(components, inst.geoReplicationDataTransferCostComponents()...)

	return components
}

func (inst *StorageQueue) isAccountKindSupported() bool {
	return inst.isStorageV1() || inst.isStorageV2()
}

func (inst *StorageQueue) isReplicationTypeSupported() bool {
	var validReplicationTypes []string

	switch {
	case inst.isStorageV1():
		validReplicationTypes = []string{"LRS", "GRS", "RA-GRS"}
	case inst.isStorageV2():
		validReplicationTypes = []string{"LRS", "ZRS", "GRS", "RA-GRS", "GZRS", "RA-GZRS"}
	}

	if validReplicationTypes != nil {
		return contains(validReplicationTypes, strings.ToUpper(inst.accountReplicationType))
	}

	return true
}

func (inst *StorageQueue) isStorageV1() bool {
	return strings.EqualFold(inst.accountKind, "storage")
}

func (inst *StorageQueue) isStorageV2() bool {
	return strings.EqualFold(inst.accountKind, "storagev2")
}

func (inst *StorageQueue) productName() string {
	if inst.isStorageV1() {
		return "Queues"
	}

	return "Queues v2"
}

func (inst *StorageQueue) dataStorageCostComponent() query.Component {
	var qty decimal.Decimal
	if inst.monthlyStorageGB != nil {
		qty = decimal.NewFromInt(*inst.monthlyStorageGB)
	}

	return query.Component{
		Name:            "Capacity",
		Unit:            "GB",
		MonthlyQuantity: qty,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr(inst.productName())},
				{Key: "sku_name", Value: util.StringPtr(fmt.Sprintf("Standard %s", strings.ToUpper(inst.accountReplicationType)))},
				{Key: "meter_name", Value: util.StringPtr(fmt.Sprintf("%s Data Stored", strings.ToUpper(inst.accountReplicationType)))},
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

func (inst *StorageQueue) operationsCostComponents() []query.Component {
	var components []query.Component

	if !contains([]string{"GZRS", "RA-GZRS"}, strings.ToUpper(inst.accountReplicationType)) {
		var class1Qty decimal.Decimal
		if inst.monthlyClass1Operations != nil {
			class1Qty = decimal.NewFromInt(*inst.monthlyClass1Operations).Div(decimal.NewFromInt(10000))
		}

		components = append(components, query.Component{
			Name:            "Class 1 operations",
			Unit:            "10k operations",
			MonthlyQuantity: class1Qty,
			ProductFilter: &product.Filter{
				Provider: util.StringPtr(inst.provider.key),
				Location: util.StringPtr(inst.location),
				Service:  util.StringPtr("Storage"),
				Family:   util.StringPtr("Storage"),
				AttributeFilters: []*product.AttributeFilter{
					{Key: "product_name", Value: util.StringPtr(inst.productName())},
					{Key: "sku_name", Value: util.StringPtr(fmt.Sprintf("Standard %s", strings.ToUpper(inst.accountReplicationType)))},
					{Key: "meter_name", ValueRegex: util.StringPtr(".*Class 1 Operations")},
				},
			},
			PriceFilter: &price.Filter{
				AttributeFilters: []*price.AttributeFilter{
					{Key: "type", Value: util.StringPtr("Consumption")},
					{Key: "tier_minimum_units", Value: util.StringPtr("0")},
				},
			},
		})
	}

	var class2Qty decimal.Decimal
	if inst.monthlyClass1Operations != nil {
		class2Qty = decimal.NewFromInt(*inst.monthlyClass2Operations).Div(decimal.NewFromInt(10000))
	}

	components = append(components, query.Component{
		Name:            "Class 2 operations",
		Unit:            "10k operations",
		MonthlyQuantity: class2Qty,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr(inst.productName())},
				{Key: "sku_name", Value: util.StringPtr(fmt.Sprintf("Standard %s", strings.ToUpper(inst.accountReplicationType)))},
				{Key: "meter_name", ValueRegex: util.StringPtr(".*Class 2 Operations")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr("0")},
			},
		},
	})

	return components
}

func (inst *StorageQueue) geoReplicationDataTransferCostComponents() []query.Component {
	if contains([]string{"LRS", "ZRS"}, strings.ToUpper(inst.accountReplicationType)) {
		return []query.Component{}
	}

	var qty decimal.Decimal
	if inst.monthlyGeoReplicationDataTransferGB != nil {
		qty = decimal.NewFromInt(*inst.monthlyGeoReplicationDataTransferGB)
	}

	return []query.Component{
		{
			Name:            "Geo-replication data transfer",
			Unit:            "GB",
			MonthlyQuantity: qty,
			ProductFilter: &product.Filter{
				Provider: util.StringPtr(inst.provider.key),
				Location: util.StringPtr(inst.location),
				Service:  util.StringPtr("Storage"),
				Family:   util.StringPtr("Storage"),
				AttributeFilters: []*product.AttributeFilter{
					{Key: "product_name", Value: util.StringPtr("Storage - Bandwidth")},
					{Key: "sku_name", Value: util.StringPtr("Geo-Replication v2")},
					{Key: "meter_name", Value: util.StringPtr("Geo-Replication v2 Data Transfer")},
				},
			},
			PriceFilter: &price.Filter{
				AttributeFilters: []*price.AttributeFilter{
					{Key: "type", Value: util.StringPtr("Consumption")},
					{Key: "tier_minimum_units", Value: util.StringPtr("0")},
				},
			},
		},
	}
}

func contains(arr []string, e string) bool {
	for _, a := range arr {
		if a == e {
			return true
		}
	}
	return false
}
