package terraform

import (
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

type Image struct {
	provider *Provider

	location  string
	imageType string
	storageGB decimal.Decimal
}

type Disk struct {
	ManagedDiskId []StorageDisk `mapstructure:"managed_disk_id"`
	SizeGb        *float64      `mapstructure:"size_gb"`
}

type sourceVirtualMachineValues struct {
	Index   string `mapstructure:"index"`
	Address string `mapstructure:"address"`
	Values  struct {
		StorageOsDisk   []StorageDisk `mapstructure:"storage_os_disk"`
		StorageDataDisk []StorageDisk `mapstructure:"storage_data_disk"`
	} `mapstructure:"values"`
}

type imageValues struct {
	Location string `mapstructure:"location"`

	OsDisk               []Disk                      `mapstructure:"os_disk"`
	DataDisk             []Disk                      `mapstructure:"data_disk"`
	SourceVirtualMachine *sourceVirtualMachineValues `mapstructure:"source_virtual_machine_id"`
	StorageGB            *float64                    `mapstructure:"storage_gb"`
}

// decodeImageValues decodes and returns computeInstanceValues from a Terraform values map.
func decodeImageValues(tfVals map[string]interface{}) (imageValues, error) {
	var v imageValues
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

// newImage initializes a new ManagedImage from the provider
func (p *Provider) newImage(vals imageValues) *Image {
	sSize := float64(0)
	if imageStorageSize(vals) != nil {
		sSize = *imageStorageSize(vals)
	}

	return &Image{
		provider:  p,
		location:  vals.Location,
		storageGB: decimal.NewFromFloat(sSize),
	}
}

func imageStorageSize(vals imageValues) *float64 {
	diskSize := getImageDiskStorage(vals)

	source := vals.SourceVirtualMachine
	if diskSize == 0 && source != nil {
		diskSize += getVMStorageSize(*source)
	}

	if diskSize == 0 {
		return nil
	}

	return &diskSize
}

func getImageDiskStorage(vals imageValues) float64 {
	var diskSize float64

	if vals.OsDisk != nil && len(vals.OsDisk) > 0 {
		managedDiskId := vals.OsDisk[0].ManagedDiskId

		diskSize += getDiskSizeGB(vals.OsDisk[0], managedDiskId, 0)
	}

	var refsDiskSize []StorageDisk
	for _, v := range vals.DataDisk {
		refsDiskSize = append(refsDiskSize, v.ManagedDiskId...)
	}

	for i, disk := range vals.DataDisk {
		diskSize += getDiskSizeGB(disk, refsDiskSize, i)
	}

	return diskSize
}

func getVMStorageSize(source sourceVirtualMachineValues) float64 {
	var size float64 = 128
	for _, disk := range source.Values.StorageOsDisk {
		if disk.DiskSizeGb != nil {
			size = *disk.DiskSizeGb
		}
	}

	for _, disk := range source.Values.StorageDataDisk {
		if disk.DiskSizeGb != nil {
			size += *disk.DiskSizeGb
		}
	}

	return size
}

func getDiskSizeGB(disk Disk, refs []StorageDisk, i int) float64 {
	if disk.SizeGb != nil {
		return *disk.SizeGb
	}

	if disk.ManagedDiskId != nil && len(refs) > i {
		ref := refs[i]
		return *ref.DiskSizeGb
	}

	return 0
}

func (inst *Image) Components() []query.Component {
	return []query.Component{{
		Name:            "Storage",
		Unit:            "1 GB/Month",
		MonthlyQuantity: inst.storageGB,
		ProductFilter: &product.Filter{
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Storage"),
			Family:   util.StringPtr("Storage"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", Value: util.StringPtr("Snapshots LRS")},
				{Key: "product_name", Value: util.StringPtr("Standard HDD Managed Disks")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}}
}
