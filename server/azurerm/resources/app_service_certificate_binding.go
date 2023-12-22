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

// AppServiceCertificateBinding is the entity that holds the logic to calculate price
// of the azurerm_app_service_certificate_order
type AppServiceCertificateBinding struct {
	provider *Provider

	location string
	sslState string
}

type BindingCertificateId struct {
	Values struct {
		Location string `mapstructure:"location"`
	} `mapstructure:"values"`
}

// appServiceCertificateBindingValues is holds the values that we need to be able
// to calculate the price of the AppServiceCertificateBinding
type appServiceCertificateBindingValues struct {
	CertificateId BindingCertificateId `mapstructure:"certificate_id"`
	SSLState      string               `mapstructure:"product_type"`
}

// decodeAppServiceCertificateBindingValues decodes and returns appServiceCertificateBindingValues from a Terraform values map.
func decodeAppServiceCertificateBindingValues(tfVals map[string]interface{}) (appServiceCertificateBindingValues, error) {
	var v appServiceCertificateBindingValues
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

// newAppServiceCertificateBinding initializes a new AppServiceCertificateBinding from the provider
func (p *Provider) newAppServiceCertificateBinding(vals appServiceCertificateBindingValues) *AppServiceCertificateBinding {
	inst := &AppServiceCertificateBinding{
		provider: p,

		location: vals.CertificateId.Values.Location,
		sslState: vals.SSLState,
	}
	return inst
}

func (inst *AppServiceCertificateBinding) Components() []query.Component {
	var components []query.Component

	sslState := strings.ToUpper(inst.sslState)

	if !strings.HasPrefix(sslState, "IP") {
		return components
	}

	components = append(components, inst.AppServiceIpSslCertificateComponent())

	return components
}

func (inst *AppServiceCertificateBinding) AppServiceIpSslCertificateComponent() query.Component {
	return query.Component{
		Name:            "IP SSL certificate",
		Unit:            "months",
		MonthlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Azure App Service"),
			Family:   util.StringPtr("Compute"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", Value: util.StringPtr("IP SSL")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
