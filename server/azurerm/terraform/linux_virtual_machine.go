package terraform

import (
	"fmt"
	"github.com/cycloidio/terracost/price"
	"github.com/cycloidio/terracost/product"
	"github.com/cycloidio/terracost/query"
	"github.com/cycloidio/terracost/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

// LinuxVirtualMachine is the entity that holds the logic to calculate price
// of the google_compute_instance
type LinuxVirtualMachine struct {
	provider *Provider

	location string
	size     string
	osDisk   []OsDisk

	// Usage
	monthlyHours decimal.Decimal
}

// linuxVirtualMachineValues is holds the values that we need to be able
// to calculate the price of the ComputeInstance
type linuxVirtualMachineValues struct {
	Size     string `mapstructure:"size"`
	Location string `mapstructure:"location"`

	OsDisk []struct {
		StorageAccountType string  `mapstructure:"storage_account_type"`
		DiskSizeGb         float64 `mapstructure:"disk_size_gb"`
	} `mapstructure:"os_disk"`

	Usage struct {
		MonthlyHours float64 `mapstructure:"monthly_hours"`
	} `mapstructure:"tc_usage"`
}

// decodeLinuxVirtualMachineValues decodes and returns computeInstanceValues from a Terraform values map.
func decodeLinuxVirtualMachineValues(tfVals map[string]interface{}) (linuxVirtualMachineValues, error) {
	var v linuxVirtualMachineValues
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

// newLinuxVirtualMachine initializes a new LinuxVirtualMachine from the provider
func (p *Provider) newLinuxVirtualMachine(vals linuxVirtualMachineValues) *LinuxVirtualMachine {
	var osDisks []OsDisk
	for _, disk := range vals.OsDisk {
		osDisks = append(osDisks, OsDisk{storageAccountType: disk.StorageAccountType, diskSizeGb: decimal.NewFromFloat(disk.DiskSizeGb)})
	}

	inst := &LinuxVirtualMachine{
		provider: p,

		location:     getLocationName(vals.Location),
		size:         vals.Size,
		osDisk:       osDisks,
		monthlyHours: decimal.NewFromFloat(vals.Usage.MonthlyHours),
	}

	return inst
}

// Components returns the price component queries that make up this Instance.
func (inst *LinuxVirtualMachine) Components() []query.Component {
	// TODO: check if we have ultra ssd or not
	components := []query.Component{inst.linuxVirtualMachineComponent()}
	components = append(components, osDiskSubResource(inst.provider, inst.location, inst.osDisk, nil)...)

	return components
}

// linuxVirtualMachineComponent returns the query needed to be able to calculate the price
func (inst *LinuxVirtualMachine) linuxVirtualMachineComponent() query.Component {
	return linuxVirtualMachineComponent(inst.provider.key, inst.location, inst.size, inst.monthlyHours)
}

// linuxVirtualMachineComponent is the abstraction of the same LinuxVirtualMachine.linuxVirtualMachineComponent
// so it can be reused
func linuxVirtualMachineComponent(key, location, size string, qty decimal.Decimal) query.Component {
	return query.Component{
		Name:            "Compute",
		MonthlyQuantity: qty,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Service:  util.StringPtr("Virtual Machines"),
			Family:   util.StringPtr("Compute"),
			Location: util.StringPtr(location),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "arm_sku_name", ValueRegex: util.StringPtr(fmt.Sprintf("%s", size))},
				{Key: "priority", Value: util.StringPtr("regular")},
				{Key: "product_name", ValueRegex: util.StringPtr("^(?!.*Windows).*")},
			},
		},
		PriceFilter: &price.Filter{
			Unit: util.StringPtr("1 Hour"),
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
