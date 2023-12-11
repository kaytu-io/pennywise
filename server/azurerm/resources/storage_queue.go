package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

// StorageQueue is the entity that holds the logic to calculate price
// of the azurerm_public_ip
type StorageQueue struct {
	provider *Provider

	location         string
	allocationMethod string
	sku              *string
}

// storageQueueValues is holds the values that we need to be able
// to calculate the price of the StorageQueue
type storageQueueValues struct {
	Location         string  `mapstructure:"location"`
	AllocationMethod string  `mapstructure:"allocation_method"`
	Sku              *string `mapstructure:"sku"`
}

// decodeStorageQueueValues decodes and returns publicIPValues from a Terraform values map.
func decodeStorageQueueValues(tfVals map[string]interface{}) (storageQueueValues, error) {
	fmt.Println("TFVALS", tfVals)
	var v storageQueueValues
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

// newPublicIP initializes a new PublicIP from the provider
func (p *Provider) newStorageQueue(vals storageQueueValues) *StorageQueue {
	inst := &StorageQueue{
		provider: p,

		location:         vals.Location,
		allocationMethod: vals.AllocationMethod,
		sku:              vals.Sku,
	}
	return inst
}

func (inst *StorageQueue) Components() []query.Component {
	var components []query.Component

	return components
}
