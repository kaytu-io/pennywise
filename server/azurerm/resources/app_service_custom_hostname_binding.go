package resources

import (
	"github.com/mitchellh/mapstructure"
)

// appServiceCustomHostnameBindingValues is holds the values that we need to be able
// to calculate the price of the AppServiceCustomHostnameBinding
type appServiceCustomHostnameBindingValues struct {
	ResourceGroupName *ResourceGroupName `mapstructure:"resource_group_name"`
	SSLState          string             `mapstructure:"product_type"`
}

// decodeAppServiceCustomHostnameBindingValues decodes and returns appServiceCustomHostnameBindingValues from a Terraform values map.
func decodeAppServiceCustomHostnameBindingValues(tfVals map[string]interface{}) (appServiceCustomHostnameBindingValues, error) {
	var v appServiceCustomHostnameBindingValues
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

// newAppServiceCustomHostnameBinding initializes a new AppServiceCertificateBinding from the provider
func (p *Provider) newAppServiceCustomHostnameBinding(vals appServiceCustomHostnameBindingValues) *AppServiceCertificateBinding {
	location := "Global"
	if vals.ResourceGroupName != nil {
		location = vals.ResourceGroupName.Values.Location
	}

	inst := &AppServiceCertificateBinding{
		provider: p,

		location: location,
		sslState: vals.SSLState,
	}
	return inst
}
