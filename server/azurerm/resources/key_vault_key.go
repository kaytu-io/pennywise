package resources

import (
	"github.com/kaytu-io/infracost/external/usage"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

type KeyVaultId struct {
	Values struct {
		Location string `mapstructure:"location"`
		SkuName  string `mapstructure:"sku_name"`
	} `mapstructure:"values"`
}

// KeyVaultKey is the entity that holds the logic to calculate price
// of the azurerm_key_vault_key
type KeyVaultKey struct {
	provider *Provider

	location   string
	keyType    string
	keySize    *string
	keyVaultId KeyVaultId

	// Usage
	monthlySecretOperations        *int64
	monthlyKeyRotationRenewals     *int64
	monthlyProtectedKeysOperations *int64
	hsmProtectedKeys               *int64
}

// keyVaultKeyValues is holds the values that we need to be able
// to calculate the price of the KeyVaultKey
type keyVaultKeyValues struct {
	KeyType    string     `mapstructure:"key_type"`
	KeySize    *string    `mapstructure:"key_size"`
	KeyVaultId KeyVaultId `mapstructure:"key_vault_id"`

	Usage struct {
		MonthlySecretOperations        *int64 `mapstructure:"monthly_secrets_operations"`
		MonthlyKeyRotationRenewals     *int64 `mapstructure:"monthly_key_rotation_renewals"`
		MonthlyProtectedKeysOperations *int64 `mapstructure:"monthly_protected_keys_operations"`
		HsmProtectedKeys               *int64 `mapstructure:"hsm_protected_keys"`
	} `mapstructure:"pennywise_usage"`
}

// decodeKeyVaultKeyValues decodes and returns keyVaultKeyValues from a Terraform values map.
func decodeKeyVaultKeyValues(tfVals map[string]interface{}) (keyVaultKeyValues, error) {
	var v keyVaultKeyValues
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

// newKeyVaultKey initializes a new KeyVaultKey from the provider
func (p *Provider) newKeyVaultKey(vals keyVaultKeyValues) *KeyVaultKey {
	inst := &KeyVaultKey{
		provider: p,

		location:   vals.KeyVaultId.Values.Location,
		keyType:    vals.KeyType,
		keySize:    vals.KeySize,
		keyVaultId: vals.KeyVaultId,

		monthlySecretOperations:        vals.Usage.MonthlySecretOperations,
		monthlyKeyRotationRenewals:     vals.Usage.MonthlyKeyRotationRenewals,
		monthlyProtectedKeysOperations: vals.Usage.MonthlyProtectedKeysOperations,
		hsmProtectedKeys:               vals.Usage.HsmProtectedKeys,
	}
	return inst
}

func (inst *KeyVaultKey) Components() []query.Component {
	var components []query.Component

	var keySize string
	if inst.keySize != nil {
		keySize = *inst.keySize
	}

	unit := "10K transactions"

	skuName := cases.Title(language.English).String(inst.keyVaultId.Values.SkuName)

	var secretsTransactions *decimal.Decimal
	if inst.monthlySecretOperations != nil {
		secretsTransactions = decimalPtr(decimal.NewFromInt(*inst.monthlySecretOperations))
	}
	components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "Secrets operations", unit, skuName, "Operations", "0", secretsTransactions, 10000))

	var keyRotationRenewals *decimal.Decimal
	if inst.monthlyKeyRotationRenewals != nil {
		keyRotationRenewals = decimalPtr(decimal.NewFromInt(*inst.monthlyKeyRotationRenewals))
	}
	components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "Storage key rotations", "renewals", skuName, "Secret Renewal", "0", keyRotationRenewals, 1))

	if !strings.HasSuffix(inst.keyType, "HSM") {
		var softwareProtectedTransactions *decimal.Decimal
		if inst.monthlyProtectedKeysOperations != nil {
			softwareProtectedTransactions = decimalPtr(decimal.NewFromInt(*inst.monthlyProtectedKeysOperations))
		}

		if inst.keyType == "RSA" && keySize == "2048" {
			components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "Software-protected keys", unit, skuName, "Operations", "0", softwareProtectedTransactions, 10000))
		} else {
			components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "Software-protected keys", unit, skuName, "Advanced Key Operations", "0", softwareProtectedTransactions, 10000))
		}
	}

	if strings.HasSuffix(inst.keyType, "HSM") && strings.ToLower(skuName) == "premium" {
		var protectedKeys, hsmProtectedTransactions *decimal.Decimal

		keyUnit := "months"

		if inst.hsmProtectedKeys != nil {
			protectedKeys = decimalPtr(decimal.NewFromInt(*inst.hsmProtectedKeys))

			if inst.keyType == "RSA-HSM" && keySize == "2048" {
				components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "HSM-protected keys", keyUnit, skuName, "Premium HSM-protected RSA 2048-bit key", "0", protectedKeys, 1))
			} else {

				tierLimits := []int{250, 1250, 2500}
				keysQuantities := usage.CalculateTierBuckets(*protectedKeys, tierLimits)

				components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "HSM-protected keys (first 250)", keyUnit, skuName, "Premium HSM-protected Advanced Key", "0", &keysQuantities[0], 1))
				if keysQuantities[1].GreaterThan(decimal.Zero) {
					components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "HSM-protected keys (next 1250)", keyUnit, skuName, "Premium HSM-protected Advanced Key", "250", &keysQuantities[1], 1))
				}
				if keysQuantities[2].GreaterThan(decimal.Zero) {
					components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "HSM-protected keys (next 2500)", keyUnit, skuName, "Premium HSM-protected Advanced Key", "1500", &keysQuantities[2], 1))
				}
				if keysQuantities[3].GreaterThan(decimal.Zero) {
					components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "HSM-protected keys (over 4000)", keyUnit, skuName, "Premium HSM-protected Advanced Key", "4000", &keysQuantities[3], 1))
				}
			}
		} else {
			var unknown *decimal.Decimal
			components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "HSM-protected keys", keyUnit, skuName, "Premium HSM-protected Advanced Key", "0", unknown, 1))
		}

		if inst.monthlyProtectedKeysOperations != nil {
			hsmProtectedTransactions = decimalPtr(decimal.NewFromInt(*inst.monthlyProtectedKeysOperations))

			if inst.keyType == "RSA" && keySize == "2048" {
				components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "HSM-protected keys", unit, skuName, "Operations", "0", hsmProtectedTransactions, 10000))
			} else {
				components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "HSM-protected keys", unit, skuName, "Advanced Key Operations", "0", hsmProtectedTransactions, 10000))
			}
		}
	}

	return components
}

func vaultKeysCostComponent(key, region, name, unit, skuName, meterName, startUsage string, quantity *decimal.Decimal, multi int) query.Component {
	if quantity != nil {
		quantity = decimalPtr(quantity.Div(decimal.NewFromInt(int64(multi))))
	} else {
		quantity = decimalPtr(decimal.Zero)
	}

	return query.Component{
		Name:            name,
		Unit:            unit,
		MonthlyQuantity: *quantity,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Location: util.StringPtr(region),
			Service:  util.StringPtr("Key Vault"),
			Family:   util.StringPtr("Security"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", Value: util.StringPtr("Key Vault")},
				{Key: "sku_name", Value: util.StringPtr(skuName)},
				{Key: "meter_name", Value: util.StringPtr(meterName)},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr(startUsage)},
			},
		},
	}
}
