package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// KeyVaultCertificate is the entity that holds the logic to calculate price
// of the azurerm_key_vault_certificate
type KeyVaultCertificate struct {
	provider *Provider

	location   string
	keyVaultId KeyVaultId

	// Usage
	monthlyCertificateRenewalRequests *int64
	monthlyCertificateOtherOperations *int64
}

// keyVaultCertificateValues is holds the values that we need to be able
// to calculate the price of the KeyVaultCertificate
type keyVaultCertificateValues struct {
	KeyVaultId KeyVaultId `mapstructure:"key_vault_id"`

	Usage struct {
		// receives monthly number of certificate renewal requests
		MonthlyCertificateRenewalRequests *int64 `mapstructure:"monthly_certificate_renewal_requests"`
		// receives monthly number of non-renewal certificate operations
		MonthlyCertificateOtherOperations *int64 `mapstructure:"monthly_certificate_other_operations"`
	} `mapstructure:"pennywise_usage"`
}

// decodeKeyVaultCertificateValues decodes and returns keyVaultCertificateValues from a Terraform values map.
func decodeKeyVaultCertificateValues(tfVals map[string]interface{}) (keyVaultCertificateValues, error) {
	var v keyVaultCertificateValues
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
func (p *Provider) newKeyVaultCertificate(vals keyVaultCertificateValues) *KeyVaultCertificate {
	inst := &KeyVaultCertificate{
		provider: p,

		location:   vals.KeyVaultId.Values.Location,
		keyVaultId: vals.KeyVaultId,

		monthlyCertificateOtherOperations: vals.Usage.MonthlyCertificateOtherOperations,
		monthlyCertificateRenewalRequests: vals.Usage.MonthlyCertificateRenewalRequests,
	}
	return inst
}

func (inst *KeyVaultCertificate) Components() []query.Component {
	var components []query.Component

	skuName := cases.Title(language.English).String(inst.keyVaultId.Values.SkuName)

	var certificateRenewals, certificateOperations *decimal.Decimal
	if inst.monthlyCertificateRenewalRequests != nil {
		certificateRenewals = decimalPtr(decimal.NewFromInt(*inst.monthlyCertificateRenewalRequests))
	}
	components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "Certificate renewals", "requests", skuName, "Certificate Renewal Request", "0", certificateRenewals, 1))

	if inst.monthlyCertificateRenewalRequests != nil {
		certificateOperations = decimalPtr(decimal.NewFromInt(*inst.monthlyCertificateRenewalRequests))
	}
	components = append(components, vaultKeysCostComponent(inst.provider.key, inst.location, "Certificate operations", "10K transactions", skuName, "Operations", "0", certificateOperations, 10000))
	GetCostComponentNamesAndSetLogger(components, inst.provider.logger)

	return components
}
