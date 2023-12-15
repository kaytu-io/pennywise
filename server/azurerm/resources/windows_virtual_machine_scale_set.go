package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

// WindowsVirtualMachineScaleSet is the entity that holds the logic to calculate price
// of the azurerm_windows_virtual_machine_scale_set
type WindowsVirtualMachineScaleSet struct {
	provider *Provider

	location               string
	sku                    string
	licenseType            string
	instances              int64
	additionalCapabilities []VirtualMachineScaleSetAdditionalCapability
	osDisk                 []OsDisk

	// Usage
	monthlyHours            *decimal.Decimal
	osDiskMonthlyOperations *decimal.Decimal
}

// windowsVirtualMachineScaleSetValues is holds the values that we need to be able
// to calculate the price of the ComputeInstance
type windowsVirtualMachineScaleSetValues struct {
	Size                   string                                       `mapstructure:"size"`
	Location               string                                       `mapstructure:"location"`
	Sku                    string                                       `mapstructure:"sku"`
	LicenseType            string                                       `mapstructure:"license_type"`
	Instances              int64                                        `mapstructure:"instances"`
	AdditionalCapabilities []VirtualMachineScaleSetAdditionalCapability `mapstructure:"additional_capabilities"`
	OsDisk                 []OsDisk                                     `mapstructure:"os_disk"`

	Usage struct {
		MonthlyHours            *float64 `mapstructure:"monthly_hours"`
		OsDiskMonthlyOperations *float64 `mapstructure:"os_disk_monthly_operations"`
	} `mapstructure:"pennywise_usage"`
}

// decodeWindowsVirtualMachineScaleSetValues decodes and returns computeInstanceValues from a Terraform values map.
func decodeWindowsVirtualMachineScaleSetValues(tfVals map[string]interface{}) (windowsVirtualMachineScaleSetValues, error) {
	var v windowsVirtualMachineScaleSetValues
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

// newWindowsVirtualMachineScaleSet initializes a new WindowsVirtualMachineScaleSet from the provider
func (p *Provider) newWindowsVirtualMachineScaleSet(vals windowsVirtualMachineScaleSetValues) *WindowsVirtualMachineScaleSet {
	var monthlyHours *decimal.Decimal
	var osDiskMonthlyOperations *decimal.Decimal

	if vals.Usage.MonthlyHours != nil {
		tmp := decimal.NewFromFloat(*vals.Usage.MonthlyHours)
		monthlyHours = &tmp
	}

	if vals.Usage.OsDiskMonthlyOperations != nil {
		tmp := decimal.NewFromFloat(*vals.Usage.OsDiskMonthlyOperations)
		osDiskMonthlyOperations = &tmp
	}

	inst := &WindowsVirtualMachineScaleSet{
		provider: p,

		location:               getLocationName(vals.Location),
		sku:                    vals.Sku,
		licenseType:            vals.LicenseType,
		instances:              vals.Instances,
		additionalCapabilities: vals.AdditionalCapabilities,
		osDisk:                 vals.OsDisk,

		monthlyHours:            monthlyHours,
		osDiskMonthlyOperations: osDiskMonthlyOperations,
	}

	return inst
}

// Components returns the price component queries that make up this Instance.
func (inst *WindowsVirtualMachineScaleSet) Components() []query.Component {
	purchaseOption := "Consumption"
	if strings.ToLower(inst.licenseType) == "windows_client" || strings.ToLower(inst.licenseType) == "windows_server" {
		purchaseOption = "DevTestConsumption"
	}
	var components []query.Component

	for i := int64(0); i < inst.instances; i++ {
		components = append(components, windowsVirtualMachineComponent(inst.provider.key, inst.location, inst.sku, purchaseOption, inst.monthlyHours))
		if len(inst.additionalCapabilities) > 0 {
			if inst.additionalCapabilities[0].UltraSsdEnabled {
				components = append(components, ultraSSDReservationCostComponent(inst.provider.key, inst.location))
			}
		}
		if len(inst.osDisk) > 0 {
			components = append(components, osDiskSubResource(inst.provider, inst.location, inst.osDisk, inst.osDiskMonthlyOperations)...)
		}
	}

	return components
}
