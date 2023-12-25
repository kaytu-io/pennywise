package resources

import (
	"fmt"
	"github.com/kaytu-io/infracost/external/usage"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

// StorageAccount is the entity that holds the logic to calculate price
// of the azurerm_storage_account
type StorageAccount struct {
	provider *Provider

	location               string
	accountReplicationType string
	accountKind            string
	accountTier            string
	accessTier             string
	nfsv3                  bool

	// Usage
	monthlyStorageGB                        *decimal.Decimal
	monthlyIterativeReadOperations          *decimal.Decimal
	monthlyReadOperations                   *decimal.Decimal
	monthlyIterativeWriteOperations         *decimal.Decimal
	monthlyWriteOperations                  *decimal.Decimal
	monthlyListAndCreateContainerOperations *decimal.Decimal
	monthlyOtherOperations                  *decimal.Decimal
	monthlyDataRetrievalGB                  *decimal.Decimal
	monthlyDataWriteGB                      *decimal.Decimal
	blobIndexTags                           *decimal.Decimal
	dataAtRestStorageGB                     *decimal.Decimal
	snapshotsStorageGB                      *decimal.Decimal
	metadataAtRestStorageGB                 *decimal.Decimal
	earlyDeletionGB                         *decimal.Decimal
}

// storageAccountValues is holds the values that we need to be able
// to calculate the price of the StorageAccount
type storageAccountValues struct {
	Location               string  `mapstructure:"location"`
	AccountKind            *string `mapstructure:"account_kind"`
	AccountTier            string  `mapstructure:"account_tier"`
	AccessTier             *string `mapstructure:"access_tier"`
	NFSv3                  *bool   `mapstructure:"nfsv3_enabled"`
	AccountReplicationType string  `mapstructure:"account_replication_type"`

	Usage struct {
		// receive total size of storage in GB.
		MonthlyStorageGB *float64 `mapstructure:"storage_gb"`
		// receive monthly number of Iterative read operations (GPv2).
		MonthlyIterativeReadOperations *float64 `mapstructure:"monthly_iterative_read_operations"`
		// receive monthly number of Read operations.
		MonthlyReadOperations *float64 `mapstructure:"monthly_read_operations"`
		// receive monthly number of Iterative write operations (GPv2).
		MonthlyIterativeWriteOperations *float64 `mapstructure:"monthly_iterative_write_operations"`
		// receive monthly number of Write operations.
		MonthlyWriteOperations *float64 `mapstructure:"monthly_write_operations"`
		// receive monthly number of List and Create Container operations
		MonthlyListAndCreateContainerOperations *float64 `mapstructure:"monthly_list_and_create_container_operations"`
		// receive monthly number of All other operations.
		MonthlyOtherOperations *float64 `mapstructure:"monthly_other_operations"`
		// receive monthly number of data retrieval in GB.
		MonthlyDataRetrievalGB *float64 `mapstructure:"monthly_data_retrieval_gb"`
		// receive monthly number of data write in GB.
		MonthlyDataWriteGB *float64 `mapstructure:"monthly_data_write_gb"`
		// receive total number of Blob indexes.
		BlobIndexTags *float64 `mapstructure:"blob_index_tags"`
		// receive total size of Data at Rest in GB (File storage).
		DataAtRestStorageGB *float64 `mapstructure:"data_at_rest_storage_gb"`
		// receive total size of Snapshots in GB (File storage).
		SnapshotsStorageGB *float64 `mapstructure:"snapshots_storage_gb"`
		// receive total size of Metadata in GB (File storage).
		MetadataAtRestStorageGB *float64 `mapstructure:"metadata_at_rest_storage_gb"`
		// receive total size of Early deletion data in GB.
		EarlyDeletionGB *float64 `mapstructure:"early_deletion_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodeStorageAccountValues decodes and returns storageAccountValues from a Terraform values map.
func decodeStorageAccountValues(tfVals map[string]interface{}) (storageAccountValues, error) {
	var v storageAccountValues
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

// newStorageAccount initializes a new StorageQueue from the provider
func (p *Provider) newStorageAccount(vals storageAccountValues) *StorageAccount {
	accountKind := "StorageV2"
	if vals.AccountKind != nil {
		accountKind = *vals.AccountKind
	}

	accountReplicationType := vals.AccountReplicationType
	switch strings.ToLower(accountReplicationType) {
	case "ragrs":
		accountReplicationType = "RA-GRS"
	case "ragzrs":
		accountReplicationType = "RA-GZRS"
	}

	accountTier := vals.AccountTier

	accessTier := "Hot"
	if vals.AccessTier != nil {
		accessTier = *vals.AccessTier
	}

	nfsv3 := false
	if vals.NFSv3 != nil {
		nfsv3 = *vals.NFSv3
	}

	inst := StorageAccount{
		provider: p,

		location:               vals.Location,
		accessTier:             accessTier,
		accountKind:            accountKind,
		accountReplicationType: accountReplicationType,
		accountTier:            accountTier,
		nfsv3:                  nfsv3,

		monthlyStorageGB:                        util.FloatToDecimal(vals.Usage.MonthlyStorageGB),
		monthlyIterativeReadOperations:          util.FloatToDecimal(vals.Usage.MonthlyIterativeReadOperations),
		monthlyReadOperations:                   util.FloatToDecimal(vals.Usage.MonthlyReadOperations),
		monthlyIterativeWriteOperations:         util.FloatToDecimal(vals.Usage.MonthlyIterativeWriteOperations),
		monthlyWriteOperations:                  util.FloatToDecimal(vals.Usage.MonthlyWriteOperations),
		monthlyListAndCreateContainerOperations: util.FloatToDecimal(vals.Usage.MonthlyListAndCreateContainerOperations),
		monthlyOtherOperations:                  util.FloatToDecimal(vals.Usage.MonthlyOtherOperations),
		monthlyDataRetrievalGB:                  util.FloatToDecimal(vals.Usage.MonthlyDataRetrievalGB),
		monthlyDataWriteGB:                      util.FloatToDecimal(vals.Usage.MonthlyDataWriteGB),
		blobIndexTags:                           util.FloatToDecimal(vals.Usage.BlobIndexTags),
		dataAtRestStorageGB:                     util.FloatToDecimal(vals.Usage.DataAtRestStorageGB),
		snapshotsStorageGB:                      util.FloatToDecimal(vals.Usage.SnapshotsStorageGB),
		metadataAtRestStorageGB:                 util.FloatToDecimal(vals.Usage.MetadataAtRestStorageGB),
		earlyDeletionGB:                         util.FloatToDecimal(vals.Usage.EarlyDeletionGB),
	}
	return &inst
}

func (inst *StorageAccount) Components() []query.Component {
	var components []query.Component

	if !inst.isReplicationTypeSupported() {
		return nil
	}

	if strings.ToLower(inst.accountTier) == "premium" {
		inst.accessTier = "Premium"
	}

	if strings.ToLower(inst.accountKind) == "storage" {
		inst.accessTier = "Standard"
	}

	components = append(components, inst.storageCostComponents()...)

	components = append(components, inst.dataAtRestCostComponents()...)
	components = append(components, inst.snapshotsCostComponents()...)
	components = append(components, inst.metadataAtRestCostComponents()...)

	components = append(components, inst.iterativeWriteOperationsCostComponents()...)
	components = append(components, inst.writeOperationsCostComponents()...)
	components = append(components, inst.listAndCreateContainerOperationsCostComponents()...)
	components = append(components, inst.iterativeReadOperationsCostComponents()...)
	components = append(components, inst.readOperationsCostComponents()...)
	components = append(components, inst.otherOperationsCostComponents()...)
	components = append(components, inst.dataRetrievalCostComponents()...)
	components = append(components, inst.dataWriteCostComponents()...)
	components = append(components, inst.blobIndexTagsCostComponents()...)

	components = append(components, inst.earlyDeletionCostComponents()...)

	return components
}

// buildProductFilter returns a product filter for the Storage Account's products.
func (inst *StorageAccount) buildProductFilter(meterName string) *product.Filter {
	var productName string

	switch {
	case strings.ToLower(inst.accountKind) == "blockblobstorage":
		productName = map[string]string{
			"Standard": "Blob Storage",
			"Premium":  "Premium Block Blob",
		}[inst.accountTier]
	case strings.ToLower(inst.accountKind) == "storage":
		productName = map[string]string{
			"Standard": "General Block Blob",
			"Premium":  "Premium Block Blob",
		}[inst.accountTier]
	case strings.ToLower(inst.accountKind) == "storagev2":
		if inst.nfsv3 {
			productName = map[string]string{
				"Standard": "General Block Blob v2 Hierarchical Namespace",
				"Premium":  "Premium Block Blob v2 Hierarchical Namespace",
			}[inst.accountTier]
		} else if strings.EqualFold(inst.accountReplicationType, "lrs") && strings.ToLower(inst.accessTier) == "hot" {
			// For some reason the Azure pricing doesn't contain all the LRS costs for all regions under "General Block Blob v2" product name.
			// But, the same pricing is available under "Blob Storage" product name.
			productName = map[string]string{
				"Standard": "Blob Storage",
				"Premium":  "Premium Block Blob",
			}[inst.accountTier]
		} else {
			productName = map[string]string{
				"Standard": "General Block Blob v2",
				"Premium":  "Premium Block Blob",
			}[inst.accountTier]
		}
	case strings.ToLower(inst.accountKind) == "blobstorage":
		productName = map[string]string{
			"Standard": "Blob Storage",
			"Premium":  "Premium Block Blob",
		}[inst.accountTier]
	case strings.ToLower(inst.accountKind) == "filestorage":
		productName = map[string]string{
			"Standard": "Files v2",
			"Premium":  "Premium Files",
		}[inst.accountTier]
	}

	skuName := fmt.Sprintf("%s %s", cases.Title(language.English).String(inst.accessTier), strings.ToUpper(inst.accountReplicationType))

	return &product.Filter{
		Provider: util.StringPtr(inst.provider.key),
		Location: util.StringPtr(inst.location),
		Service:  util.StringPtr("Storage"),
		Family:   util.StringPtr("Storage"),
		AttributeFilters: []*product.AttributeFilter{
			{Key: "product_name", Value: util.StringPtr(productName)},
			{Key: "sku_name", Value: util.StringPtr(skuName)},
			{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf(".*%s", meterName))},
		},
	}
}

func (inst *StorageAccount) storageCostComponents() []query.Component {
	var components []query.Component

	if strings.ToLower(inst.accountKind) == "filestorage" {
		return components
	}

	var quantity decimal.Decimal
	name := "Capacity"

	if inst.monthlyStorageGB == nil {
		components = append(components, inst.buildStorageCostComponent(
			name,
			"0",
			quantity,
		))
		return components
	}

	if inst.monthlyStorageGB != nil {
		quantity = *inst.monthlyStorageGB
	}

	// Only Hot storage has pricing tiers, others have a single price for any
	// amount.
	if !(strings.ToLower(inst.accessTier) == "hot") {
		components = append(components, inst.buildStorageCostComponent(
			name,
			"0",
			quantity,
		))
		return components
	}

	type dataTier struct {
		name       string
		startUsage string
	}

	data := []dataTier{
		{name: fmt.Sprintf("%s (first 50TB)", name), startUsage: "0"},
		{name: fmt.Sprintf("%s (next 450TB)", name), startUsage: "51200"},
		{name: fmt.Sprintf("%s (over 500TB)", name), startUsage: "512000"},
	}

	tierLimits := []int{51200, 512000}
	tiers := usage.CalculateTierBuckets(quantity, tierLimits)

	for i, d := range data {
		if i < len(tiers) && tiers[i].GreaterThan(decimal.Zero) {
			components = append(components, inst.buildStorageCostComponent(
				d.name,
				d.startUsage,
				tiers[i],
			))
		}
	}
	return components
}

func (inst *StorageAccount) iterativeWriteOperationsCostComponents() []query.Component {
	var components []query.Component

	if !(strings.ToLower(inst.accountKind) == "storagev2") || !inst.nfsv3 || strings.ToLower(inst.accountTier) == "premium" {
		return components
	}

	var quantity decimal.Decimal
	itemsPerCost := 100

	if inst.monthlyIterativeWriteOperations != nil {
		value := *inst.monthlyIterativeWriteOperations
		quantity = value.Div(decimal.NewFromInt(int64(itemsPerCost)))
	}

	meterName := "Iterative Write Operations"

	components = append(components, query.Component{
		Name:            "Iterative write operations",
		Unit:            "100 operations",
		MonthlyQuantity: quantity,
		ProductFilter:   inst.buildProductFilter(meterName),
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	return components
}

func (inst *StorageAccount) writeOperationsCostComponents() []query.Component {
	var components []query.Component

	if strings.ToLower(inst.accountKind) == "filestorage" && strings.ToLower(inst.accountTier) == "premium" {
		return components
	}

	var quantity decimal.Decimal
	itemsPerCost := 10000

	if inst.monthlyWriteOperations != nil {
		value := *inst.monthlyWriteOperations
		quantity = value.Div(decimal.NewFromInt(int64(itemsPerCost)))
	}

	meterName := "Write Operations"
	if strings.ToLower(inst.accountKind) == "storagev2" && inst.nfsv3 {
		meterName = "(?<!Iterative) Write Operations"
	}

	components = append(components, query.Component{
		Name:            "Write operations",
		Unit:            "10k operations",
		MonthlyQuantity: quantity,
		ProductFilter:   inst.buildProductFilter(meterName),
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	return components
}

func (inst *StorageAccount) listAndCreateContainerOperationsCostComponents() []query.Component {
	var components []query.Component

	if strings.ToLower(inst.accountKind) == "filestorage" && strings.ToLower(inst.accountTier) == "premium" {
		return components
	}

	var quantity decimal.Decimal
	itemsPerCost := 10000

	if inst.monthlyListAndCreateContainerOperations != nil {
		value := *inst.monthlyListAndCreateContainerOperations
		quantity = value.Div(decimal.NewFromInt(int64(itemsPerCost)))
	}

	name := "List and create container operations"
	meterName := "List and Create Container Operations"

	if strings.ToLower(inst.accountKind) == "filestorage" {
		name = "List operations"
		meterName = "List Operations"
	}

	components = append(components, query.Component{
		Name:            name,
		Unit:            "10k operations",
		MonthlyQuantity: quantity,
		ProductFilter:   inst.buildProductFilter(meterName),
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	return components
}

func (inst *StorageAccount) iterativeReadOperationsCostComponents() []query.Component {
	var components []query.Component

	if !(strings.ToLower(inst.accountKind) == "storagev2") || !inst.nfsv3 || strings.ToLower(inst.accountTier) == "premium" {
		return components
	}

	var quantity decimal.Decimal
	itemsPerCost := 10000

	if inst.monthlyIterativeReadOperations != nil {
		value := *inst.monthlyIterativeReadOperations
		quantity = value.Div(decimal.NewFromInt(int64(itemsPerCost)))
	}

	meterName := "Iterative Read Operations"

	components = append(components, query.Component{
		Name:            "Iterative read operations",
		Unit:            "10k operations",
		MonthlyQuantity: quantity,
		ProductFilter:   inst.buildProductFilter(meterName),
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	return components
}

func (inst *StorageAccount) readOperationsCostComponents() []query.Component {
	var components []query.Component

	if strings.ToLower(inst.accountKind) == "filestorage" && strings.ToLower(inst.accountTier) == "premium" {
		return components
	}

	var quantity decimal.Decimal
	itemsPerCost := 10000

	if inst.monthlyReadOperations != nil {
		value := *inst.monthlyReadOperations
		quantity = value.Div(decimal.NewFromInt(int64(itemsPerCost)))
	}

	meterName := "Read Operations"
	if strings.ToLower(inst.accountKind) == "storagev2" && inst.nfsv3 {
		meterName = "(?<!Iterative) Read Operations"
	}
	if strings.ToLower(inst.accountKind) == "storage" && contains([]string{"LRS", "GRS", "RA-GRS"}, strings.ToUpper(inst.accountReplicationType)) {
		// Storage V1 GRS/LRS/RA-GRS doesn't always have a Read Operations meter name, but we can use this regex
		// to match Read or Other Operations meter since they are the same price.
		meterName = "(Other|Read) Operations"
	}

	filter := inst.buildProductFilter(meterName)
	components = append(components, query.Component{
		Name:            "Read operations",
		Unit:            "10k operations",
		MonthlyQuantity: quantity,
		ProductFilter:   filter,
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	return components
}

func (inst *StorageAccount) otherOperationsCostComponents() []query.Component {
	var components []query.Component

	if strings.ToLower(inst.accountKind) == "filestorage" && strings.ToLower(inst.accountTier) == "premium" {
		return components
	}

	var quantity decimal.Decimal
	itemsPerCost := 10000

	if inst.monthlyOtherOperations != nil {
		value := *inst.monthlyOtherOperations
		quantity = value.Div(decimal.NewFromInt(int64(itemsPerCost)))
	}

	meterName := "Other Operations"
	if strings.ToLower(inst.accountKind) == "storage" {
		meterName = "Delete Operations"
	}

	components = append(components, query.Component{
		Name:            "All other operations",
		Unit:            "10k operations",
		MonthlyQuantity: quantity,
		ProductFilter:   inst.buildProductFilter(meterName),
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	return components
}

func (inst *StorageAccount) dataRetrievalCostComponents() []query.Component {
	var components []query.Component

	if !(strings.ToLower(inst.accessTier) == "cool") {
		return components
	}

	var quantity decimal.Decimal

	if inst.monthlyDataRetrievalGB != nil {
		quantity = *inst.monthlyDataRetrievalGB
	}

	meterName := "Data Retrieval"

	components = append(components, query.Component{
		Name:            "Data retrieval",
		Unit:            "GB",
		MonthlyQuantity: quantity,
		ProductFilter:   inst.buildProductFilter(meterName),
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	return components
}

func (inst *StorageAccount) dataWriteCostComponents() []query.Component {
	var components []query.Component

	if !(strings.ToLower(inst.accountKind) == "blockblobstorage" && !(strings.ToLower(inst.accountKind) == "blobstorage")) || !(strings.ToLower(inst.accessTier) == "cool") {
		return components
	}

	var quantity decimal.Decimal

	if inst.monthlyDataWriteGB != nil {
		quantity = *inst.monthlyDataWriteGB
	}

	meterName := "Data Write"

	components = append(components, query.Component{
		Name:            "Data write",
		Unit:            "GB",
		MonthlyQuantity: quantity,
		ProductFilter:   inst.buildProductFilter(meterName),
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	return components
}

func (inst *StorageAccount) blobIndexTagsCostComponents() []query.Component {
	var components []query.Component

	isBlockPremium := strings.ToLower(inst.accountKind) == "blockblobstorage" && strings.ToLower(inst.accountTier) == "premium"
	isBlobPremium := strings.ToLower(inst.accountKind) == "blobstorage" && strings.ToLower(inst.accountTier) == "premium"
	isV2NFSv3 := strings.ToLower(inst.accountKind) == "storagev2" && (inst.nfsv3 || strings.ToLower(inst.accountTier) == "premium")
	if strings.ToLower(inst.accountKind) == "filestorage" || strings.ToLower(inst.accountKind) == "storage" || isBlockPremium || isBlobPremium || isV2NFSv3 {
		return components
	}

	var quantity decimal.Decimal
	itemsPerCost := 10000

	if inst.blobIndexTags != nil {
		value := *inst.blobIndexTags
		quantity = value.Div(decimal.NewFromInt(int64(itemsPerCost)))
	}

	meterName := "Index Tags"

	components = append(components, query.Component{
		Name:            "Blob index",
		Unit:            "10k tags",
		MonthlyQuantity: quantity,
		ProductFilter:   inst.buildProductFilter(meterName),
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	return components
}

func (inst *StorageAccount) dataAtRestCostComponents() []query.Component {
	var components []query.Component

	if !(strings.ToLower(inst.accountKind) == "filestorage") {
		return components
	}

	var quantity decimal.Decimal

	if inst.dataAtRestStorageGB != nil {
		quantity = *inst.dataAtRestStorageGB
	}

	meterName := "Data Stored"
	if strings.ToLower(inst.accountTier) == "premium" {
		meterName = "Provisioned"
	}

	components = append(components, query.Component{
		Name:            "Data at rest",
		Unit:            "GB",
		MonthlyQuantity: quantity,
		ProductFilter:   inst.buildProductFilter(meterName),
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	return components
}

func (inst *StorageAccount) snapshotsCostComponents() []query.Component {
	var components []query.Component

	if !(strings.ToLower(inst.accountKind) == "filestorage") {
		return components
	}

	var quantity decimal.Decimal

	if inst.snapshotsStorageGB != nil {
		quantity = *inst.snapshotsStorageGB
	}

	meterName := "Data Stored"
	if strings.ToLower(inst.accountTier) == "premium" {
		meterName = "Snapshots"
	}

	components = append(components, query.Component{
		Name:            "Snapshots",
		Unit:            "GB",
		MonthlyQuantity: quantity,
		ProductFilter:   inst.buildProductFilter(meterName),
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	return components
}

func (inst *StorageAccount) metadataAtRestCostComponents() []query.Component {
	var components []query.Component

	if !(strings.ToLower(inst.accountKind) == "filestorage") || strings.ToLower(inst.accountTier) == "premium" {
		return components
	}

	var quantity decimal.Decimal

	if inst.metadataAtRestStorageGB != nil {
		quantity = *inst.metadataAtRestStorageGB
	}

	meterName := "Metadata"

	components = append(components, query.Component{
		Name:            "Metadata at rest",
		Unit:            "GB",
		MonthlyQuantity: quantity,
		ProductFilter:   inst.buildProductFilter(meterName),
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	return components
}

func (inst *StorageAccount) earlyDeletionCostComponents() []query.Component {
	var components []query.Component

	if strings.ToLower(inst.accountKind) == "storage" || strings.ToLower(inst.accountKind) == "blockblobstorage" || strings.ToLower(inst.accountKind) == "blobstorage" || !(strings.ToLower(inst.accessTier) == "cool") {
		return components
	}

	var quantity decimal.Decimal
	if inst.earlyDeletionGB != nil {
		quantity = *inst.earlyDeletionGB
	}

	meterName := "Early Delete"

	components = append(components, query.Component{
		Name:            "Early deletion",
		Unit:            "GB",
		MonthlyQuantity: quantity,
		ProductFilter:   inst.buildProductFilter(meterName),
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	return components
}

// buildStorageCostComponent builds one cost component for storage amount costs.
func (inst *StorageAccount) buildStorageCostComponent(name string, startUsage string, quantity decimal.Decimal) query.Component {
	meterName := "Data Stored"

	return query.Component{
		Name:            name,
		Unit:            "GB",
		MonthlyQuantity: quantity,
		ProductFilter:   inst.buildProductFilter(meterName),
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr(startUsage)},
			},
		},
	}
}

func (inst *StorageAccount) isFileStorage() bool {
	return strings.ToLower(inst.accountKind) == "filestorage"
}

func (inst *StorageAccount) isReplicationTypeSupported() bool {
	var validReplicationTypes []string

	switch {
	case strings.ToLower(inst.accountTier) == "premium":
		validReplicationTypes = []string{"LRS", "ZRS"}
	case strings.ToLower(inst.accountKind) == "blockblobstorage":
		validReplicationTypes = []string{"LRS", "GRS", "RA-GRS"}
	case strings.ToLower(inst.accountKind) == "storage":
		validReplicationTypes = []string{"LRS", "ZRS", "GRS", "RA-GRS"}
	case strings.ToLower(inst.accountKind) == "storagev2":
		validReplicationTypes = []string{"LRS", "ZRS", "GRS", "RA-GRS", "GZRS", "RA-GZRS"}
	case strings.ToLower(inst.accountKind) == "blobstorage":
		validReplicationTypes = []string{"LRS", "GRS", "RA-GRS"}
	case strings.ToLower(inst.accountKind) == "filestorage":
		validReplicationTypes = []string{"LRS", "GRS", "ZRS"}
	}

	if validReplicationTypes != nil {
		return contains(validReplicationTypes, strings.ToUpper(inst.accountReplicationType))
	}

	return true
}
