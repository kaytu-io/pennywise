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
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
)

func TestElasticIP_Components(t *testing.T) {
	logger, err := zap.NewProduction()
	require.NoError(t, err)

	p, err := NewProvider("aws", "us-east-1", logger)
	require.NoError(t, err)

	t.Run("EIP", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_eip.test",
			Type:         "aws_eip",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"vpc": true,
			},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{
			{
				Name:           "Elastic IP",
				Details:        []string{"ElasticIP:IdleAddress"},
				HourlyQuantity: decimal.NewFromInt(1),
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AmazonEC2"),
					Family:   util.StringPtr("IP Address"),
					Location: util.StringPtr("us-east-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "Group", Value: util.StringPtr("ElasticIP:IdleAddress")},
					},
				},
				PriceFilter: &price.Filter{
					Unit: util.StringPtr("Hrs"),
					AttributeFilters: []*price.AttributeFilter{
						{Key: "TermType", Value: util.StringPtr("OnDemand")},
						{Key: "StartingRange", Value: util.StringPtr("1")},
					},
				},
			},
		}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})

	t.Run("EIPCustomerOwnedIpv4Pool", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_eip.test",
			Type:         "aws_eip",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"customer_owned_ipv4_pool": "customer-owned-ipv4-pool-id",
			},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})

	t.Run("EIPInstance", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_eip.test",
			Type:         "aws_eip",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"instance": "instance-id",
			},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})

	t.Run("EIPNetworkInterface", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_eip.test",
			Type:         "aws_eip",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"network_interface": "network-interface-id",
			},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})
}
