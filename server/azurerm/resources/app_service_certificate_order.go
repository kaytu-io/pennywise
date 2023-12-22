package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

// AppServiceCertificateOrder is the entity that holds the logic to calculate price
// of the azurerm_app_service_certificate_order
type AppServiceCertificateOrder struct {
	provider *Provider

	location    string
	productType string
}

// appServiceCertificateOrderValues is holds the values that we need to be able
// to calculate the price of the AppServiceCertificateOrder
type appServiceCertificateOrderValues struct {
	Location    string  `mapstructure:"location"`
	ProductType *string `mapstructure:"product_type"`
}

// decodeAppServiceCertificateOrderValues decodes and returns appServiceCertificateOrderValues from a Terraform values map.
func decodeAppServiceCertificateOrderValues(tfVals map[string]interface{}) (appServiceCertificateOrderValues, error) {
	var v appServiceCertificateOrderValues
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

// newApplicationGateway initializes a new AppServiceCertificateOrder from the provider
func (p *Provider) newAppServiceCertificateOrder(vals appServiceCertificateOrderValues) *AppServiceCertificateOrder {
	productType := "Standard"
	if vals.ProductType != nil {
		productType = *vals.ProductType
	}

	inst := &AppServiceCertificateOrder{
		provider: p,

		location:    vals.Location,
		productType: productType,
	}
	if inst.location == "global" {
		inst.location = "Global"
	} else if strings.Contains(inst.location, "usgov") {
		inst.location = "US Gov"
	}
	return inst
}

func (inst *AppServiceCertificateOrder) Components() []query.Component {
	var components []query.Component

	components = append(components, inst.AppServiceSSLCertificateComponent())

	return components
}

func (inst *AppServiceCertificateOrder) AppServiceSSLCertificateComponent() query.Component {
	return query.Component{
		Name: fmt.Sprintf("SSL certificate (%s)", inst.productType),
		Unit: "years",

		MonthlyQuantity: decimal.NewFromInt(1).Div(decimal.NewFromInt(12)),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.provider.key),
			Location: util.StringPtr(inst.location),
			Service:  util.StringPtr("Azure App Service"),
			Family:   util.StringPtr("Compute"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", ValueRegex: util.StringPtr(fmt.Sprintf("%s SSL - 1 Year", inst.productType))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
