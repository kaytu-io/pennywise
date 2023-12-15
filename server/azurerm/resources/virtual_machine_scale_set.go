package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

type VirtualMachineScaleSetSku struct {
	Name     string `mapstructure:"name"`
	Tier     string `mapstructure:"tier"`
	Capacity int64  `mapstructure:"capacity"`
}

// VirtualMachineScaleSet is the entity that holds the logic to calculate price
// of the azurerm_virtual_machine_scale_set
type VirtualMachineScaleSet struct {
	provider *Provider

	location                     string
	sku                          []VirtualMachineScaleSetSku
	licenseType                  *string
	additionalCapabilities       []VirtualMachineScaleSetAdditionalCapability
	osDisk                       []OsDisk
	osProfileWindowsConfig       *interface{}
	storageProfileImageReference []StorageImageReference
	storageProfileOsDisk         []StorageDisk
	storageProfileDataDisk       []StorageDisk

	// Usage
	monthlyHours              *decimal.Decimal
	osDiskMonthlyOperations   *decimal.Decimal
	dataDiskMonthlyOperations *decimal.Decimal
	instances                 *int64
}

// virtualMachineScaleSetValues is holds the values that we need to be able
// to calculate the price of the ComputeInstance
type virtualMachineScaleSetValues struct {
	Size                         string                                       `mapstructure:"size"`
	Location                     string                                       `mapstructure:"location"`
	Sku                          []VirtualMachineScaleSetSku                  `mapstructure:"sku"`
	LicenseType                  *string                                      `mapstructure:"license_type"`
	AdditionalCapabilities       []VirtualMachineScaleSetAdditionalCapability `mapstructure:"additional_capabilities"`
	OsDisk                       []OsDisk                                     `mapstructure:"os_disk"`
	OsProfileWindowsConfig       *interface{}                                 `mapstructure:"os_profile_windows_config"`
	StorageProfileImageReference []StorageImageReference                      `mapstructure:"storage_profile_image_reference"`
	StorageProfileOsDisk         []StorageDisk                                `mapstructure:"storage_profile_os_disk"`
	StorageProfileDataDisk       []StorageDisk                                `mapstructure:"storage_profile_data_disk"`

	Usage struct {
		MonthlyHours              *float64 `mapstructure:"monthly_hours"`
		OsDiskMonthlyOperations   *float64 `mapstructure:"os_disk_monthly_operations"`
		DataDiskMonthlyOperations *float64 `mapstructure:"data_disk_monthly_operations"`
		Instances                 *int64   `mapstructure:"instances"`
	} `mapstructure:"pennywise_usage"`
}

// decodeVirtualMachineScaleSetValues decodes and returns computeInstanceValues from a Terraform values map.
func decodeVirtualMachineScaleSetValues(tfVals map[string]interface{}) (virtualMachineScaleSetValues, error) {
	var v virtualMachineScaleSetValues
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

// newVirtualMachineScaleSet initializes a new VirtualMachineScaleSet from the provider
func (p *Provider) newVirtualMachineScaleSet(vals virtualMachineScaleSetValues) *VirtualMachineScaleSet {
	var monthlyHours *decimal.Decimal
	var osDiskMonthlyOperations *decimal.Decimal
	var dataDiskMonthlyOperations *decimal.Decimal

	if vals.Usage.MonthlyHours != nil {
		tmp := decimal.NewFromFloat(*vals.Usage.MonthlyHours)
		monthlyHours = &tmp
	}

	if vals.Usage.OsDiskMonthlyOperations != nil {
		tmp := decimal.NewFromFloat(*vals.Usage.OsDiskMonthlyOperations)
		osDiskMonthlyOperations = &tmp
	}

	if vals.Usage.DataDiskMonthlyOperations != nil {
		tmp := decimal.NewFromFloat(*vals.Usage.DataDiskMonthlyOperations)
		dataDiskMonthlyOperations = &tmp
	}

	inst := &VirtualMachineScaleSet{
		provider: p,

		location:                     getLocationName(vals.Location),
		sku:                          vals.Sku,
		licenseType:                  vals.LicenseType,
		additionalCapabilities:       vals.AdditionalCapabilities,
		osDisk:                       vals.OsDisk,
		osProfileWindowsConfig:       vals.OsProfileWindowsConfig,
		storageProfileImageReference: vals.StorageProfileImageReference,
		storageProfileOsDisk:         vals.StorageProfileOsDisk,
		storageProfileDataDisk:       vals.StorageProfileDataDisk,

		monthlyHours:              monthlyHours,
		osDiskMonthlyOperations:   osDiskMonthlyOperations,
		dataDiskMonthlyOperations: dataDiskMonthlyOperations,
		instances:                 vals.Usage.Instances,
	}

	return inst
}

// Components returns the price component queries that make up this Instance.
func (inst *VirtualMachineScaleSet) Components() []query.Component {
	var components []query.Component

	instances := inst.sku[0].Capacity
	if inst.instances != nil {
		instances = *inst.instances
	}

	for i := int64(0); i < instances; i++ {
		os := "Linux"
		if inst.osProfileWindowsConfig != nil {
			os = "Windows"
		}
		if len(inst.storageProfileOsDisk) > 0 {
			if inst.storageProfileOsDisk[0].OsType == "windows" {
				os = "Windows"
			}
		}
		if len(inst.storageProfileImageReference) > 0 {
			if inst.storageProfileImageReference[0].Offer == "windowsserver" {
				os = "Windows"
			}
		}

		if os == "Linux" {
			components = append(components, linuxVirtualMachineComponent(inst.provider.key, inst.location, inst.sku[0].Name, inst.monthlyHours))
		}

		if os == "Windows" {
			licenseType := "Windows_Client"
			if inst.licenseType != nil {
				licenseType = *inst.licenseType
			}
			purchaseOption := "Consumption"
			if strings.ToLower(licenseType) == "windows_client" || strings.ToLower(licenseType) == "windows_server" {
				purchaseOption = "DevTestConsumption"
			}
			components = append(components, windowsVirtualMachineComponent(inst.provider.key, inst.location, inst.sku[0].Name, purchaseOption, inst.monthlyHours))
		}
	}
	var osDiskMonthlyOperations float64
	if inst.dataDiskMonthlyOperations != nil {
		osDiskMonthlyOperations = inst.osDiskMonthlyOperations.InexactFloat64()
	}
	if len(inst.storageProfileOsDisk) > 0 {
		managedStorage := inst.provider.newManagedStorage(managedDiskValues{
			StorageAccountType: inst.storageProfileOsDisk[0].ManagedDiskType,
			Location:           inst.location,
			DiskSizeGb:         0,
			DiskIopsReadWrite:  0,
			BurstingEnabled:    false,
			DiskMbpsReadWrite:  0,

			Usage: struct {
				MonthlyDiskOperations float64 `mapstructure:"monthly_disk_operations"`
			}{MonthlyDiskOperations: osDiskMonthlyOperations},
		})
		components = append(components, managedStorage.Components()...)
	}

	var dataDiskMonthlyOperations float64
	if inst.dataDiskMonthlyOperations != nil {
		dataDiskMonthlyOperations = inst.dataDiskMonthlyOperations.InexactFloat64()
	}
	if len(inst.storageProfileDataDisk) > 0 {
		for _, disk := range inst.storageProfileDataDisk {
			managedStorage := inst.provider.newManagedStorage(managedDiskValues{
				StorageAccountType: disk.ManagedDiskType,
				Location:           inst.location,
				DiskSizeGb:         0,
				DiskIopsReadWrite:  0,
				BurstingEnabled:    false,
				DiskMbpsReadWrite:  0,

				Usage: struct {
					MonthlyDiskOperations float64 `mapstructure:"monthly_disk_operations"`
				}{MonthlyDiskOperations: dataDiskMonthlyOperations},
			})
			components = append(components, managedStorage.Components()...)
		}
	}

	return components
}
