package resources

import (
	"github.com/kaytu-io/pennywise/server/azurerm/resources"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/kaytu-io/pennywise/server/aws/region"
)

// NatGateway represents an NatGateway instance definition that can be cost-estimated.
type NatGateway struct {
	providerKey string
	logger      *zap.Logger
	region      region.Code

	// Usage
	monthlyDataProcessedGB decimal.Decimal
}

type natGatewayValues struct {
	Usage struct {
		MonthlyDataProcessedGB float64 `mapstructure:"monthly_data_processed_gb"`
	} `mapstructure:"pennywise_usage"`
}

func decodeNatGatewayValues(tfVals map[string]interface{}) (natGatewayValues, error) {
	var v natGatewayValues
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

// NewInstance creates a new Instance from Terraform values.
func (p *Provider) newNatGateway(vals natGatewayValues) *NatGateway {

	inst := &NatGateway{
		providerKey:            p.key,
		logger:                 p.logger,
		region:                 p.region,
		monthlyDataProcessedGB: decimal.NewFromFloat(vals.Usage.MonthlyDataProcessedGB),
	}

	return inst
}

// Components returns the price component queries that make up this Instance.
func (inst *NatGateway) Components() []query.Component {
	components := []query.Component{inst.natGatewayInstanceComponent()}
	components = append(components, inst.natGatewayDataProcessedComponent())

	resources.GetCostComponentNamesAndSetLogger(components, inst.logger)
	return components
}

func (inst *NatGateway) natGatewayInstanceComponent() query.Component {
	return query.Component{
		Name:           "NAT gateway",
		Details:        []string{"NatGateway"},
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.providerKey),
			Service:  util.StringPtr("AmazonEC2"),
			Family:   util.StringPtr("NAT Gateway"),
			Location: util.StringPtr(inst.region.String()),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "UsageType", ValueRegex: util.StringPtr(".*NatGateway-Hours")},
			},
		},
		PriceFilter: &price.Filter{
			Unit: util.StringPtr("Hrs"),
			AttributeFilters: []*price.AttributeFilter{
				{Key: "TermType", Value: util.StringPtr("OnDemand")},
			},
		},
	}
}

func (inst *NatGateway) natGatewayDataProcessedComponent() query.Component {
	return query.Component{
		Name:            "NAT Data processed",
		Details:         []string{"NatGateway Data processed"},
		Usage:           true,
		MonthlyQuantity: inst.monthlyDataProcessedGB,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.providerKey),
			Service:  util.StringPtr("AmazonEC2"),
			Family:   util.StringPtr("NAT Gateway"),
			Location: util.StringPtr(inst.region.String()),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "UsageType", ValueRegex: util.StringPtr(".*NatGateway-Bytes")},
			},
		},
		PriceFilter: &price.Filter{
			Unit: util.StringPtr("GB"),
			AttributeFilters: []*price.AttributeFilter{
				{Key: "TermType", Value: util.StringPtr("OnDemand")},
			},
		},
	}
}
