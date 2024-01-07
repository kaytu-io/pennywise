package resources

import (
	"github.com/kaytu-io/pennywise/server/resource"
	"go.uber.org/zap"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/util"
)

func TestEKSCluster_Components(t *testing.T) {
	logger, err := zap.NewProduction()
	require.NoError(t, err)

	p, err := NewProvider("aws", "us-east-1", logger)
	require.NoError(t, err)

	t.Run("EKSCluster", func(t *testing.T) {
		tfres := resource.ResourceDef{
			Address:      "aws_eks_cluster.test",
			Type:         "aws_eks_cluster",
			Name:         "test",
			ProviderName: "aws",
			Values:       map[string]interface{}{},
		}
		rss := map[string]resource.ResourceDef{}

		expected := []resource.Component{
			{
				Name:           "EKS Cluster",
				Details:        []string{"EKSCluster:Compute"},
				HourlyQuantity: decimal.NewFromInt(1),
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AmazonEKS"),
					Family:   util.StringPtr("Compute"),
					Location: util.StringPtr("us-east-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "UsageType", Value: util.StringPtr("USE1-AmazonEKS-Hours:perCluster")},
					},
				},
				PriceFilter: &price.Filter{
					AttributeFilters: []*price.AttributeFilter{
						{Key: "TermType", Value: util.StringPtr("OnDemand")},
					},
				},
			},
		}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})
}
