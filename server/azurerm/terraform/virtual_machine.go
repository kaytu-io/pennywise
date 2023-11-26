package terraform

import (
	"github.com/kaytu.io/pennywise/server/internal/price"
	"github.com/kaytu.io/pennywise/server/internal/product"
	"github.com/kaytu.io/pennywise/server/internal/query"
	"github.com/kaytu.io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

type OS string

const (
	WindowsOS OS = "Windows"
	LinuxOS   OS = "Linux"
)

type OsDisk struct {
	storageAccountType string
	diskSizeGb         decimal.Decimal
}

// VirtualMachine is the entity that holds the logic to calculate price
// of the google_compute_instance
type VirtualMachine struct {
	provider *Provider

	location        string
	vmSize          string
	operatingSystem OS
	licenseType     string
	storageOsDisk   bool
	storageDataDisk bool
	managedDiskType string

	// Usage
	monthlyOsDiskOperations   decimal.Decimal
	monthlyDataDiskOperations decimal.Decimal
	monthlyHours              decimal.Decimal
}

// virtualMachineValues is holds the values that we need to be able
// to calculate the price of the ComputeInstance
type virtualMachineValues struct {
	VMSize          string `mapstructure:"vm_size"`
	Location        string `mapstructure:"location"`
	OperatingSystem OS     `mapstructure:"operating_system"`
	LicenseType     string `mapstructure:"license_type"`
	StorageOsDisk   bool   `mapstructure:"storage_os_disk"`
	StorageDataDisk bool   `mapstructure:"storage_data_disk"`
	ManagedDiskType string `mapstructure:"managed_disk_type"`

	Usage struct {
		MonthlyOsDiskOperations   float64 `mapstructure:"monthly_os_disk_operations"`
		MonthlyDataDiskOperations float64 `mapstructure:"monthly_data_disk_operations"`
		MonthlyHours              float64 `mapstructure:"monthly_hours"`
	} `mapstructure:"tc_usage"`
}

// decodeVirtualMachineValues decodes and returns computeInstanceValues from a Terraform values map.
func decodeVirtualMachineValues(tfVals map[string]interface{}) (virtualMachineValues, error) {
	var v virtualMachineValues
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

// newVirtualMachine initializes a new VirtualMachine from the provider
func (p *Provider) newVirtualMachine(vals virtualMachineValues) *VirtualMachine {
	inst := &VirtualMachine{
		provider: p,

		location:                  getLocationName(vals.Location),
		vmSize:                    vals.VMSize,
		operatingSystem:           vals.OperatingSystem,
		licenseType:               vals.LicenseType,
		storageOsDisk:             vals.StorageOsDisk,
		storageDataDisk:           vals.StorageDataDisk,
		monthlyOsDiskOperations:   decimal.NewFromFloat(vals.Usage.MonthlyOsDiskOperations),
		monthlyDataDiskOperations: decimal.NewFromFloat(vals.Usage.MonthlyDataDiskOperations),
		monthlyHours:              decimal.NewFromFloat(vals.Usage.MonthlyHours),
	}

	return inst
}

// Components returns the price component queries that make up this Instance.
func (inst *VirtualMachine) Components() []query.Component {
	var components []query.Component
	if inst.operatingSystem == WindowsOS {
		windowsInst := inst.provider.newWindowsVirtualMachine(windowsVirtualMachineValues{Size: inst.vmSize, Location: inst.location, LicenseType: inst.licenseType, Usage: struct {
			MonthlyHours float64 `mapstructure:"monthly_hours"`
		}{MonthlyHours: inst.monthlyHours.InexactFloat64()}})
		components = []query.Component{windowsInst.windowsVirtualMachineComponent()}
	} else if inst.operatingSystem == LinuxOS {
		linuxInst := inst.provider.newLinuxVirtualMachine(linuxVirtualMachineValues{Size: inst.vmSize, Location: inst.location, Usage: struct {
			MonthlyHours float64 `mapstructure:"monthly_hours"`
		}{MonthlyHours: inst.monthlyHours.InexactFloat64()}})
		components = []query.Component{linuxInst.linuxVirtualMachineComponent()}
	}

	if inst.storageOsDisk {
		managedStorage := inst.provider.newManagedStorage(managedDiskValues{
			StorageAccountType: inst.managedDiskType,
			Location:           inst.location,
			DiskSizeGb:         1024,
			DiskIopsReadWrite:  2048,
			BurstingEnabled:    false,
			DiskMbpsReadWrite:  8,

			Usage: struct {
				MonthlyDiskOperations float64 `mapstructure:"monthly_disk_operations"`
			}{MonthlyDiskOperations: inst.monthlyOsDiskOperations.InexactFloat64()},
		})
		components = append(components, managedStorage.Components()...)
	}

	if inst.storageDataDisk {
		managedStorage := inst.provider.newManagedStorage(managedDiskValues{
			StorageAccountType: inst.managedDiskType,
			Location:           inst.location,
			DiskSizeGb:         1024,
			DiskIopsReadWrite:  2048,
			BurstingEnabled:    false,
			DiskMbpsReadWrite:  8,

			Usage: struct {
				MonthlyDiskOperations float64 `mapstructure:"monthly_disk_operations"`
			}{MonthlyDiskOperations: inst.monthlyOsDiskOperations.InexactFloat64()},
		})
		components = append(components, managedStorage.Components()...)
	}

	return components
}

func ultraSSDReservationCostComponent(key, location string) *query.Component {
	return &query.Component{
		Name:           "Ultra disk reservation (if unattached)",
		Unit:           "vCPU",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Location: util.StringPtr(location),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Ultra Disks")},
				{Key: "sku_name", Value: util.StringPtr("Ultra LRS")},
				{Key: "meter_name", Value: util.StringPtr("Ultra LRS Reservation per vCPU Provisioned")},
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

func osDiskSubResource(provider *Provider, location string, osDisk []OsDisk, monthlyDiskOperations *decimal.Decimal) []query.Component {
	managedStorage := provider.newManagedStorage(managedDiskValues{
		StorageAccountType: osDisk[0].storageAccountType,
		Location:           location,
		DiskSizeGb:         osDisk[0].diskSizeGb.InexactFloat64(),
	})
	return managedStorage.Components()
}
