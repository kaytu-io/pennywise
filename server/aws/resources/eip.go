package resources

import (
	"github.com/kaytu-io/pennywise/server/azurerm/resources"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/kaytu-io/pennywise/server/aws/region"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/util"
)

// ElasticIP represents an ElasticIP instance definition that can be cost-estimated.
type ElasticIP struct {
	providerKey           string
	logger                *zap.Logger
	region                region.Code
	customerOwnedIpv4Pool string
	instance              string
	networkInterface      string
}

type elasticIPValues struct {
	CustomerOwnedIpv4Pool string `mapstructure:"customer_owned_ipv4_pool"`
	Instance              string `mapstructure:"instance"`
	NetworkInterface      string `mapstructure:"network_interface"`
}

func decodeElasticIPValues(tfVals map[string]interface{}) (elasticIPValues, error) {
	var v elasticIPValues
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
func (p *Provider) newElasticIP(vals elasticIPValues) *ElasticIP {

	inst := &ElasticIP{
		providerKey:           p.key,
		logger:                p.logger,
		region:                p.region,
		customerOwnedIpv4Pool: vals.CustomerOwnedIpv4Pool,
		instance:              vals.Instance,
		networkInterface:      vals.NetworkInterface,
	}

	return inst
}

// Components returns the price component queries that make up this Instance.
func (inst *ElasticIP) Components() []resource.Component {
	// An Elastic IP address doesn’t incur charges as long as all the following conditions are true:
	// * The Elastic IP address is associated with an EC2 instance.
	// * The instance associated with the Elastic IP address is running.
	// * The instance has only one Elastic IP address attached to it.
	// * The Elastic IP address is associated with an attached network interface
	if len(inst.customerOwnedIpv4Pool) > 0 || len(inst.instance) > 0 || len(inst.networkInterface) > 0 {
		return []resource.Component{}
	}

	components := []resource.Component{inst.elasticIPInstanceComponent()}
	resources.GetCostComponentNamesAndSetLogger(components, inst.logger)
	return components
}

func (inst *ElasticIP) elasticIPInstanceComponent() resource.Component {

	attrFilters := []*product.AttributeFilter{
		{Key: "Group", Value: util.StringPtr("ElasticIP:IdleAddress")},
	}

	return resource.Component{
		Name:           "Elastic IP",
		Details:        []string{"ElasticIP:IdleAddress"},
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider:         util.StringPtr(inst.providerKey),
			Service:          util.StringPtr("AmazonEC2"),
			Family:           util.StringPtr("IP Address"),
			Location:         util.StringPtr(inst.region.String()),
			AttributeFilters: attrFilters,
		},
		PriceFilter: &price.Filter{
			Unit: util.StringPtr("Hrs"),
			AttributeFilters: []*price.AttributeFilter{
				{Key: "TermType", Value: util.StringPtr("OnDemand")},
				{Key: "StartingRange", Value: util.StringPtr("1")},
			},
		},
	}
}
