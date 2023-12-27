package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/azurerm/resources"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"go.uber.org/zap"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"

	"github.com/kaytu-io/pennywise/server/aws/region"
)

// EKSCluster represents an EKSCluster instance definition that can be cost-estimated.
type EKSCluster struct {
	providerKey string
	logger      *zap.Logger
	region      region.Code

	// tenancy describes the tenancy of an instance.
	// Valid values include: Shared, Dedicated, Host.
	// Note: only "Shared" and "Dedicated" are supported at the moment.
	// Seems only used with fargate
	// tenancy string
}

type eKSClusterValues struct {
	// VpcConfig 						string `mapstructure:"vpc_config "`
}

func decodeEKSClusterValues(tfVals map[string]interface{}) (eKSClusterValues, error) {
	var v eKSClusterValues
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
func (p *Provider) newEKSCluster(vals eKSClusterValues) *EKSCluster {
	inst := &EKSCluster{
		providerKey: p.key,
		logger:      p.logger,
		region:      p.region,
		// tenancy:     "Shared",
	}

	return inst
}

// Components returns the price component queries that make up this Instance.
func (inst *EKSCluster) Components() []query.Component {
	components := []query.Component{inst.eKSClusterInstanceComponent()}

	resources.GetCostComponentNamesAndSetLogger(components, inst.logger)
	return components
}

func (inst *EKSCluster) eKSClusterInstanceComponent() query.Component {

	// EU-AmazonEKS-Hours:perCluster
	// Get us-east-1
	// Convert to USE1

	// eu-west-1 is the only exception where it need to be EU only
	region := ""
	if inst.region.String() == "eu-west-1" {
		region = "EU"
	} else {
		splitedRegion := strings.Split(inst.region.String(), "-")
		region = fmt.Sprintf("%s%s%s", strings.ToUpper(splitedRegion[0]), strings.ToUpper(splitedRegion[1][0:1]), splitedRegion[2])
	}

	return query.Component{
		Name:           "EKS Cluster",
		Details:        []string{"EKSCluster:Compute"},
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(inst.providerKey),
			Service:  util.StringPtr("AmazonEKS"),
			Family:   util.StringPtr("Compute"),
			Location: util.StringPtr(inst.region.String()),
			AttributeFilters: []*product.AttributeFilter{
				// {Key: "Tenancy", Value: util.StringPtr(inst.tenancy)},
				{Key: "UsageType", Value: util.StringPtr(fmt.Sprintf("%s-AmazonEKS-Hours:perCluster", region))},
			},
		},
		PriceFilter: &price.Filter{
			// Unit: util.StringPtr("Hours"),
			AttributeFilters: []*price.AttributeFilter{
				{Key: "TermType", Value: util.StringPtr("OnDemand")},
			},
		},
	}
}
