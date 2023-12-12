package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

type VirtualMachineScaleSetAdditionalCapability struct {
	UltraSsdEnabled bool `mapstructure:"ultra_ssd_enabled"`
}

// LinuxVirtualMachineScaleSet is the entity that holds the logic to calculate price
// of the azurerm_linux_virtual_machine_scale_set
type LinuxVirtualMachineScaleSet struct {
	provider *Provider

	location               string
	sku                    string
	instances              int64
	additionalCapabilities []VirtualMachineScaleSetAdditionalCapability
	osDisk                 []OsDisk

	// Usage
	monthlyHours            *decimal.Decimal
	osDiskMonthlyOperations *decimal.Decimal
}

// linuxVirtualMachineScaleSetValues is holds the values that we need to be able
// to calculate the price of the ComputeInstance
type linuxVirtualMachineScaleSetValues struct {
	Size                   string                                       `mapstructure:"size"`
	Location               string                                       `mapstructure:"location"`
	Sku                    string                                       `mapstructure:"sku"`
	Instances              int64                                        `mapstructure:"instances"`
	AdditionalCapabilities []VirtualMachineScaleSetAdditionalCapability `mapstructure:"additional_capabilities"`
	OsDisk                 []OsDisk                                     `mapstructure:"os_disk"`

	Usage struct {
		MonthlyHours            *float64 `mapstructure:"monthly_hours"`
		OsDiskMonthlyOperations *float64 `mapstructure:"os_disk_monthly_operations"`
	} `mapstructure:"pennywise_usage"`
}

// decodeLinuxVirtualMachineScaleSetValues decodes and returns computeInstanceValues from a Terraform values map.
func decodeLinuxVirtualMachineScaleSetValues(tfVals map[string]interface{}) (linuxVirtualMachineScaleSetValues, error) {
	var v linuxVirtualMachineScaleSetValues
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

// newLinuxVirtualMachineScaleSet initializes a new LinuxVirtualMachineScaleSet from the provider
func (p *Provider) newLinuxVirtualMachineScaleSet(vals linuxVirtualMachineScaleSetValues) *LinuxVirtualMachineScaleSet {
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

	inst := &LinuxVirtualMachineScaleSet{
		provider: p,

		location:               getLocationName(vals.Location),
		instances:              vals.Instances,
		sku:                    vals.Sku,
		additionalCapabilities: vals.AdditionalCapabilities,
		osDisk:                 vals.OsDisk,

		monthlyHours:            monthlyHours,
		osDiskMonthlyOperations: osDiskMonthlyOperations,
	}

	return inst
}

// Components returns the price component queries that make up this Instance.
func (inst *LinuxVirtualMachineScaleSet) Components() []query.Component {
	var components []query.Component

	for i := int64(0); i < inst.instances; i++ {
		components = append(components, linuxVirtualMachineComponent(inst.provider.key, inst.location, inst.sku, decimal.NewFromFloat(730)))

		if len(inst.additionalCapabilities) > 0 {
			if inst.additionalCapabilities[0].UltraSsdEnabled {
				components = append(components, ultraSSDReservationCostComponent(inst.provider.key, inst.location))
			}
		}
		if len(inst.osDisk) > 0 {
			fmt.Println("OSDISK", inst.osDisk)
			components = append(components, osDiskSubResource(inst.provider, inst.location, inst.osDisk, inst.osDiskMonthlyOperations)...)
		}
	}

	return components
}
