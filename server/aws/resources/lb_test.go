package resources

import (
	"github.com/kaytu-io/pennywise/server/resource"
	"testing"

	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLB_Components(t *testing.T) {
	p, err := NewProvider("aws", "eu-west-1")
	require.NoError(t, err)

	t.Run("DefaultValues", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_lb.test",
			Type:         "aws_lb",
			Name:         "test",
			ProviderName: "aws",
			Values:       map[string]interface{}{},
		}
		rss := map[string]resource.Resource{}
		expected := []query.Component{
			{
				Name:           "Application Load Balancer",
				HourlyQuantity: decimal.NewFromInt(1),
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AWSELB"),
					Family:   util.StringPtr("Load Balancer-Application"),
					Location: util.StringPtr("eu-west-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "UsageType", ValueRegex: util.StringPtr("LoadBalancerUsage")},
					},
				},
				PriceFilter: &price.Filter{
					Unit: util.StringPtr("Hrs"),
				},
			},
		}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})

	t.Run("NetworkLoadBalancer", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_lb.test",
			Type:         "aws_lb",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"load_balancer_type": "network",
			},
		}
		rss := map[string]resource.Resource{}
		expected := []query.Component{
			{
				Name:           "Network Load Balancer",
				HourlyQuantity: decimal.NewFromInt(1),
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AWSELB"),
					Family:   util.StringPtr("Load Balancer-Network"),
					Location: util.StringPtr("eu-west-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "UsageType", ValueRegex: util.StringPtr("LoadBalancerUsage")},
					},
				},
				PriceFilter: &price.Filter{
					Unit: util.StringPtr("Hrs"),
				},
			},
		}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})

	t.Run("GatewayLoadBalancer", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_lb.test",
			Type:         "aws_lb",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"load_balancer_type": "gateway",
			},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{
			{
				Name:           "Gateway Load Balancer",
				HourlyQuantity: decimal.NewFromInt(1),
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AWSELB"),
					Family:   util.StringPtr("Load Balancer-Gateway"),
					Location: util.StringPtr("eu-west-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "UsageType", ValueRegex: util.StringPtr("LoadBalancerUsage")},
					},
				},
				PriceFilter: &price.Filter{
					Unit: util.StringPtr("Hrs"),
				},
			},
		}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})

	t.Run("ClassicLoadBalancer", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_elb.test",
			Type:         "aws_elb",
			Name:         "test",
			ProviderName: "aws",
			Values:       map[string]interface{}{},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{
			{
				Name:           "Classic Load Balancer",
				HourlyQuantity: decimal.NewFromInt(1),
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AWSELB"),
					Family:   util.StringPtr("Load Balancer"),
					Location: util.StringPtr("eu-west-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "UsageType", ValueRegex: util.StringPtr("LoadBalancerUsage")},
					},
				},
				PriceFilter: &price.Filter{
					Unit: util.StringPtr("Hrs"),
				},
			},
		}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})
}
