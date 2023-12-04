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

// ManagedDisk is the entity that holds the logic to calculate price
// of the google_compute_instance
type ManagedDisk struct {
	provider *Provider

	location           string
	storageAccountType string
	diskSizeGb         float64
	burstingEnabled    bool
	diskMbpsReadWrite  float64
	diskIopsReadWrite  float64

	// Usage
	monthlyDiskOperations decimal.Decimal
}

// managedDiskValues is holds the values that we need to be able
// to calculate the price of the ComputeInstance
type managedDiskValues struct {
	StorageAccountType string  `mapstructure:"storage_account_type"`
	Location           string  `mapstructure:"location"`
	DiskSizeGb         float64 `mapstructure:"disk_size_gb"`
	BurstingEnabled    bool    `mapstructure:"on_demand_bursting_enabled"`
	DiskMbpsReadWrite  float64 `mapstructure:"disk_mbps_read_write"`
	DiskIopsReadWrite  float64 `mapstructure:"disk_iops_read_write"`

	Usage struct {
		MonthlyDiskOperations float64 `mapstructure:"monthly_disk_operations"`
	} `mapstructure:"tc_usage"`
}

// decodeManagedStorageValues decodes and returns computeInstanceValues from a Terraform values map.
func decodeManagedStorageValues(tfVals map[string]interface{}) (managedDiskValues, error) {
	var v managedDiskValues
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

// newManagedStorage initializes a new VirtualMachine from the provider
func (p *Provider) newManagedStorage(vals managedDiskValues) *ManagedDisk {
	inst := &ManagedDisk{
		provider: p,

		location:              getLocationName(vals.Location),
		storageAccountType:    vals.StorageAccountType,
		diskSizeGb:            vals.DiskSizeGb,
		burstingEnabled:       vals.BurstingEnabled,
		diskMbpsReadWrite:     vals.DiskMbpsReadWrite,
		diskIopsReadWrite:     vals.DiskIopsReadWrite,
		monthlyDiskOperations: decimal.NewFromFloat(vals.Usage.MonthlyDiskOperations),
	}

	return inst
}

// Components returns the price component queries that make up this Instance.
func (inst *ManagedDisk) Components() []query.Component {
	var components []query.Component

	sku := strings.Split(inst.storageAccountType, "_")
	if sku[0] == "PremiumV2" {
		return nil // Not Supported
	} else if sku[0] == "UltraSSD" {
		requestedSize := float64(1024)
		iops := float64(2048)
		throughput := float64(8)

		if inst.diskSizeGb != 0 {
			requestedSize = inst.diskSizeGb
		}

		if inst.diskIopsReadWrite != 0 {
			iops = inst.diskIopsReadWrite
		}

		if inst.diskMbpsReadWrite != 0 {
			throughput = inst.diskMbpsReadWrite
		}
		components = append(components, inst.ultraLRSThroughputComponent(inst.provider.key, inst.location, throughput))
		components = append(components, inst.ultraLRSCapacityComponent(inst.provider.key, inst.location, requestedSize))
		components = append(components, inst.ultraLRSIOPsComponent(inst.provider.key, inst.location, iops))
	} else {
		requestedSize := 30
		if int(inst.diskSizeGb) != 0 {
			requestedSize = int(inst.diskSizeGb)
		}
		skuName := mapDiskName(sku[0], requestedSize)
		productName, ok := diskProductNameMap[sku[0]]
		if !ok {
			return nil
		}
		components = []query.Component{inst.managedStorageComponent(inst.provider.key, inst.location, skuName, sku[1], productName)}

		if (sku[0] == "Premium") && (inst.diskSizeGb >= 1000) && inst.burstingEnabled {
			components = append(components, inst.enableBurstingComponent(inst.provider.key, inst.location))
		}
	}

	if strings.ToLower(sku[0]) == "standard" || strings.ToLower(sku[0]) == "standardssd" {
		var opsQty decimal.Decimal

		opsQty = inst.monthlyDiskOperations.Div(decimal.NewFromInt(10000))
		inst.diskOperationsComponent(inst.provider.key, inst.location, inst.storageAccountType, opsQty)
	}
	return components
}

// ultraLRSThroughputComponent Throughput of Ultra LRS
func (inst *ManagedDisk) ultraLRSThroughputComponent(key, location string, throughput float64) query.Component {
	return query.Component{
		Name:           "Ultra LRS Throughput",
		HourlyQuantity: decimal.NewFromFloat(throughput),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", Value: util.StringPtr("Ultra LRS")},
				{Key: "meter_name", Value: util.StringPtr("Ultra LRS Provisioned Throughput (MBps)")},
			},
		},
		PriceFilter: &price.Filter{
			Unit: util.StringPtr("1/Hour"),
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

// ultraLRSCapacityComponent Capacity of Ultra LRS
func (inst *ManagedDisk) ultraLRSCapacityComponent(key, location string, diskSize float64) query.Component {
	return query.Component{
		Name:           "Ultra LRS Capacity",
		HourlyQuantity: decimal.NewFromFloat(diskSize),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", Value: util.StringPtr("Ultra LRS")},
				{Key: "meter_name", Value: util.StringPtr("Ultra LRS Provisioned Capacity")},
			},
		},
		PriceFilter: &price.Filter{
			Unit: util.StringPtr("1 GiB/Hour"),
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

// ultraLRSIOPsComponent IOPs for Ultra LRS
func (inst *ManagedDisk) ultraLRSIOPsComponent(key, location string, iops float64) query.Component {
	return query.Component{
		Name:           "Ultra LRS IOPs",
		HourlyQuantity: decimal.NewFromFloat(iops),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", Value: util.StringPtr("Ultra LRS")},
				{Key: "meter_name", Value: util.StringPtr("Ultra LRS Provisioned IOPS")},
			},
		},
		PriceFilter: &price.Filter{
			Unit: util.StringPtr("1/Hour"),
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

// managedStorageComponent is the component for Premium and Standard Managed Storages
func (inst *ManagedDisk) managedStorageComponent(key, location, diskName, storageReplicationType, productName string) query.Component {
	return query.Component{
		Name:            "Managed Storage",
		MonthlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr(productName)},
				{Key: "sku_name", Value: util.StringPtr(fmt.Sprintf("%s %s", diskName, storageReplicationType))},
				{Key: "meter_name", Value: util.StringPtr(fmt.Sprintf("%s %s Disk", diskName, storageReplicationType))},
			},
		},
		PriceFilter: &price.Filter{
			Unit: util.StringPtr("1/Month"),
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

// diskOperationsComponent is the component for Standard Managed Storages disk operations
func (inst *ManagedDisk) diskOperationsComponent(key, location, skuName string, quantity decimal.Decimal) query.Component {
	return query.Component{
		Name:            "Disk operations",
		MonthlyQuantity: quantity,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(".*Managed Disks.*")},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
				{Key: "meter_name", Value: util.StringPtr(fmt.Sprintf("%s Disk Operations", skuName))},
			},
		},
		PriceFilter: &price.Filter{
			Unit: util.StringPtr("1/Month"),
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

// enableBurstingComponent component for when the Bursting is enabled for the managed storage
func (inst *ManagedDisk) enableBurstingComponent(key, location string) query.Component {
	return query.Component{
		Name:            "Enable Bursting",
		MonthlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", Value: util.StringPtr("Burst Enablement LRS")},
			},
		},
		PriceFilter: &price.Filter{
			Unit: util.StringPtr("1/Month"),
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

type skuDetails struct {
	DiskOption string
	DiskSize   float64
}

const Standard = "Standard"
const StandardSSD = "StandardSSD"
const Premium = "Premium"

var diskSizeMap = map[string][]struct {
	Name string
	Size int
}{
	// The mapping is from https://docs.microsoft.com/en-us/azure/virtual-machines/disks-types
	// sizes of disks don't depend on replication types. meaning, Standard_LRS disk sizes are the same as Standard_ZRS
	Standard: {
		{"S4", 32},
		{"S6", 64},
		{"S10", 128},
		{"S15", 256},
		{"S20", 512},
		{"S30", 1024},
		{"S40", 2048},
		{"S50", 4096},
		{"S60", 8192},
		{"S70", 16384},
		{"S80", 32767},
	},
	StandardSSD: {
		{"E1", 4},
		{"E2", 8},
		{"E3", 16},
		{"E4", 32},
		{"E6", 64},
		{"E10", 128},
		{"E15", 256},
		{"E20", 512},
		{"E30", 1024},
		{"E40", 2048},
		{"E50", 4096},
		{"E60", 8192},
		{"E70", 16384},
		{"E80", 32767},
	},
	Premium: {
		{"P1", 4},
		{"P2", 8},
		{"P3", 16},
		{"P4", 32},
		{"P6", 64},
		{"P10", 128},
		{"P15", 256},
		{"P20", 512},
		{"P30", 1024},
		{"P40", 2048},
		{"P50", 4096},
		{"P60", 8192},
		{"P70", 16384},
		{"P80", 32767},
	},
}

func mapDiskName(diskType string, requestedSize int) string {
	diskTypeMap, ok := diskSizeMap[diskType]
	if !ok {
		return ""
	}

	name := ""
	for _, v := range diskTypeMap {
		name = v.Name
		if v.Size >= requestedSize {
			break
		}
	}

	if requestedSize > diskTypeMap[len(diskTypeMap)-1].Size {
		return ""
	}

	return name
}

var diskProductNameMap = map[string]string{
	Standard:    "Standard HDD Managed Disks",
	StandardSSD: "Standard SSD Managed Disks",
	Premium:     "Premium SSD Managed Disks",
}
