package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

type OS string

type OsDisk struct {
	Caching            string  `mapstructure:"caching"`
	StorageAccountType string  `mapstructure:"storage_account_type"`
	DiskSizeGb         float64 `mapstructure:"disk_size_gb"`
}

// VirtualMachine is the entity that holds the logic to calculate price
// of the google_compute_instance
type VirtualMachine struct {
	provider *Provider

	location              string
	vmSize                string
	licenseType           string
	storageOsDisk         []StorageDisk
	storageDataDisk       []StorageDisk
	managedDiskType       string
	storageImageReference []StorageImageReference `mapstructure:"storage_image_reference"`

	// Usage
	monthlyOsDiskOperations   *decimal.Decimal
	monthlyDataDiskOperations *decimal.Decimal
	monthlyHours              *decimal.Decimal
}

type StorageImageReference struct {
	Offer string `mapstructure:"offer"`
}

type StorageDisk struct {
	DiskSizeGb      *float64 `mapstructure:"disk_size_gb"`
	ManagedDiskType string   `mapstructure:"managed_disk_type"`
	OsType          string   `mapstructure:"os_type"`
}

// virtualMachineValues is holds the values that we need to be able
// to calculate the price of the ComputeInstance
type virtualMachineValues struct {
	VMSize                string                  `mapstructure:"vm_size"`
	Location              string                  `mapstructure:"location"`
	LicenseType           string                  `mapstructure:"license_type"`
	StorageOsDisk         []StorageDisk           `mapstructure:"storage_os_disk"`
	StorageDataDisk       []StorageDisk           `mapstructure:"storage_data_disk"`
	ManagedDiskType       string                  `mapstructure:"managed_disk_type"`
	StorageImageReference []StorageImageReference `mapstructure:"storage_image_reference"`

	Usage struct {
		// receive monthly number of main disk operations (writes, reads, deletes) using a unit size of 256KiB.
		MonthlyOsDiskOperations *float64 `mapstructure:"monthly_os_disk_operations"`
		// receive monthly number of disk operations (writes, reads, deletes) using a unit size of 256KiB per additional disk.
		MonthlyDataDiskOperations *float64 `mapstructure:"monthly_data_disk_operations"`
		// receive monthly number of hours the instance ran for.
		MonthlyHours *float64 `mapstructure:"monthly_hours"`
	} `mapstructure:"pennywise_usage"`
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
		licenseType:               vals.LicenseType,
		storageOsDisk:             vals.StorageOsDisk,
		storageDataDisk:           vals.StorageDataDisk,
		monthlyOsDiskOperations:   util.FloatToDecimal(vals.Usage.MonthlyOsDiskOperations),
		monthlyDataDiskOperations: util.FloatToDecimal(vals.Usage.MonthlyDataDiskOperations),
		monthlyHours:              util.FloatToDecimal(vals.Usage.MonthlyHours),
		storageImageReference:     vals.StorageImageReference,
	}
	return inst
}

// Components returns the price component queries that make up this Instance.
func (inst *VirtualMachine) Components() []query.Component {
	var components []query.Component

	os := "Linux"
	if len(inst.storageImageReference) > 0 {
		if strings.ToLower(inst.storageImageReference[0].Offer) == "windowsserver" {
			os = "Windows"
		}
	}
	if len(inst.storageOsDisk) > 0 {
		if strings.ToLower(inst.storageOsDisk[0].OsType) == "windows" {
			os = "Windows"
		}
	}

	if os == "Windows" {
		windowsInst := inst.provider.newWindowsVirtualMachine(windowsVirtualMachineValues{Size: inst.vmSize, Location: inst.location, LicenseType: inst.licenseType, Usage: struct {
			MonthlyHours *float64 `mapstructure:"monthly_hours"`
		}{MonthlyHours: util.DecimalToFloat(inst.monthlyHours)}})
		components = []query.Component{windowsInst.windowsVirtualMachineComponent()}
	} else if os == "Linux" {
		linuxInst := inst.provider.newLinuxVirtualMachine(linuxVirtualMachineValues{Size: inst.vmSize, Location: inst.location, Usage: struct {
			MonthlyHours *float64 `mapstructure:"monthly_hours"`
		}{MonthlyHours: util.DecimalToFloat(inst.monthlyHours)}})
		components = []query.Component{linuxInst.linuxVirtualMachineComponent()}
	}
	components = append(components, ultraSSDReservationCostComponent(inst.provider.key, inst.location))
	if len(inst.storageOsDisk) > 0 {
		var osDiskOperations *decimal.Decimal
		if inst.monthlyOsDiskOperations != nil {
			osDiskOperations = inst.monthlyOsDiskOperations
		}
		managedStorage := inst.provider.newManagedStorage(managedDiskValues{
			StorageAccountType: inst.storageOsDisk[0].ManagedDiskType,
			Location:           inst.location,
			DiskSizeGb:         0,
			DiskIopsReadWrite:  0,
			BurstingEnabled:    false,
			DiskMbpsReadWrite:  0,

			Usage: struct {
				MonthlyDiskOperations *float64 `mapstructure:"monthly_disk_operations"`
			}{MonthlyDiskOperations: util.DecimalToFloat(osDiskOperations)},
		})
		components = append(components, managedStorage.Components()...)
	}

	if len(inst.storageDataDisk) > 0 {
		for _, disk := range inst.storageDataDisk {
			managedStorage := inst.provider.newManagedStorage(managedDiskValues{
				StorageAccountType: disk.ManagedDiskType,
				Location:           inst.location,
				DiskSizeGb:         0,
				DiskIopsReadWrite:  0,
				BurstingEnabled:    false,
				DiskMbpsReadWrite:  0,

				Usage: struct {
					MonthlyDiskOperations *float64 `mapstructure:"monthly_disk_operations"`
				}{MonthlyDiskOperations: util.DecimalToFloat(inst.monthlyDataDiskOperations)},
			})
			components = append(components, managedStorage.Components()...)
		}
	}
	GetCostComponentNamesAndSetLogger(components, inst.provider.logger)

	return components
}

func ultraSSDReservationCostComponent(key, location string) query.Component {
	return query.Component{
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
	var diskOperations *float64
	if monthlyDiskOperations != nil {
		diskOperations = util.FloatPtr(monthlyDiskOperations.InexactFloat64())
	}
	managedStorage := provider.newManagedStorage(managedDiskValues{
		StorageAccountType: osDisk[0].StorageAccountType,
		Location:           location,
		DiskSizeGb:         osDisk[0].DiskSizeGb,
		Usage: struct {
			MonthlyDiskOperations *float64 `mapstructure:"monthly_disk_operations"`
		}{MonthlyDiskOperations: diskOperations},
	})
	return managedStorage.Components()
}
