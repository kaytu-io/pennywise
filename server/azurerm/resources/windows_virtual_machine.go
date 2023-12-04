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

// WindowsVirtualMachine is the entity that holds the logic to calculate price
// of the google_compute_instance
type WindowsVirtualMachine struct {
	provider *Provider

	location    string
	size        string
	licenseType string
	osDisk      []OsDisk

	// Usage
	monthlyHours decimal.Decimal
}

// windowsVirtualMachineValues is holds the values that we need to be able
// to calculate the price of the ComputeInstance
type windowsVirtualMachineValues struct {
	Size        string `mapstructure:"size"`
	Location    string `mapstructure:"location"`
	LicenseType string `mapstructure:"license_type"`

	OsDisk []struct {
		StorageAccountType string  `mapstructure:"storage_account_type"`
		DiskSizeGb         float64 `mapstructure:"disk_size_gb"`
	} `mapstructure:"os_disk"`

	Usage struct {
		MonthlyHours float64 `mapstructure:"monthly_hours"`
	} `mapstructure:"tc_usage"`
}

// decodeWindowsVirtualMachineValues decodes and returns computeInstanceValues from a Terraform values map.
func decodeWindowsVirtualMachineValues(tfVals map[string]interface{}) (windowsVirtualMachineValues, error) {
	var v windowsVirtualMachineValues
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

// newWindowsVirtualMachine initializes a new WindowsVirtualMachine from the provider
func (p *Provider) newWindowsVirtualMachine(vals windowsVirtualMachineValues) *WindowsVirtualMachine {
	var osDisks []OsDisk
	for _, disk := range vals.OsDisk {
		osDisks = append(osDisks, OsDisk{storageAccountType: disk.StorageAccountType, diskSizeGb: decimal.NewFromFloat(disk.DiskSizeGb)})
	}
	inst := &WindowsVirtualMachine{
		provider: p,

		location:     getLocationName(vals.Location),
		size:         vals.Size,
		licenseType:  vals.LicenseType,
		osDisk:       osDisks,
		monthlyHours: decimal.NewFromFloat(vals.Usage.MonthlyHours),
	}

	return inst
}

// Components returns the price component queries that make up this Instance.
func (inst *WindowsVirtualMachine) Components() []query.Component {
	// TODO: check if we have ultra ssd or not
	components := []query.Component{inst.windowsVirtualMachineComponent()}
	components = append(components, osDiskSubResource(inst.provider, inst.location, inst.osDisk, nil)...)
	return components
}

// linuxVirtualMachineComponent returns the query needed to be able to calculate the price
func (inst *WindowsVirtualMachine) windowsVirtualMachineComponent() query.Component {
	purchaseOption := "Consumption"
	if strings.ToLower(inst.licenseType) == "windows_client" || strings.ToLower(inst.licenseType) == "windows_server" {
		purchaseOption = "DevTestConsumption"
	}
	return windowsVirtualMachineComponent(inst.provider.key, inst.location, inst.size, purchaseOption, inst.monthlyHours)
}

// linuxVirtualMachineComponent is the abstraction of the same LinuxVirtualMachine.linuxVirtualMachineComponent
// so it can be reused
func windowsVirtualMachineComponent(key, location, size, purchaseOption string, qty decimal.Decimal) query.Component {
	return query.Component{
		Name:            "Compute",
		MonthlyQuantity: qty,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Virtual Machines"),
			Family:   util.StringPtr("Compute"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "arm_sku_name", Value: util.StringPtr(size)},
				{Key: "priority", Value: util.StringPtr("regular")},
				{Key: "product_name", ValueRegex: util.StringPtr(".*Windows.*")},
			},
		},
		PriceFilter: &price.Filter{
			Unit: util.StringPtr("1 Hour"),
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr(purchaseOption)},
			},
		},
	}
}
