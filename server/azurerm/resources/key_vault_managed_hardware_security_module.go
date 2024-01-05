package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

// KeyVaultManagedHardwareSecurityModule is the entity that holds the logic to calculate price
// of the azurerm_key_vault_managed_hardware_security_module
type KeyVaultManagedHardwareSecurityModule struct {
	provider *Provider

	location string
}

// keyVaultManagedHardwareSecurityModuleValues is holds the values that we need to be able
// to calculate the price of the KeyVaultManagedHardwareSecurityModule
type keyVaultManagedHardwareSecurityModuleValues struct {
	Location string `mapstructure:"location"`
}

// decodeKeyVaultManagedHardwareSecurityModuleValues decodes and returns keyVaultCertificateValues from a Terraform values map.
func decodeKeyVaultManagedHardwareSecurityModuleValues(tfVals map[string]interface{}) (keyVaultManagedHardwareSecurityModuleValues, error) {
	var v keyVaultManagedHardwareSecurityModuleValues
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

// newKeyVaultKey initializes a new KeyVaultCertificate from the provider
func (p *Provider) newKeyVaultManagedHardwareSecurityModule(vals keyVaultManagedHardwareSecurityModuleValues) *KeyVaultManagedHardwareSecurityModule {
	inst := &KeyVaultManagedHardwareSecurityModule{
		provider: p,

		location: vals.Location,
	}
	return inst
}

func (inst *KeyVaultManagedHardwareSecurityModule) Components() []resource.Component {
	var components []resource.Component

	components = append(components, inst.hsmPoolComponent())
	GetCostComponentNamesAndSetLogger(components, inst.provider.logger)

	return components
}

func (inst *KeyVaultManagedHardwareSecurityModule) hsmPoolComponent() resource.Component {
	return resource.Component{
		Name:           "HSM pools",
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Key Vault"),
			Family:   util.StringPtr("Security"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Azure Dedicated HSM")},
				{Key: "sku_name", Value: util.StringPtr("Standard")},
				{Key: "meter_name", Value: util.StringPtr("Standard Instance")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr("0")},
			},
		},
	}
}
