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

func TestElastiCache_Components(t *testing.T) {
	logger, err := zap.NewProduction()
	require.NoError(t, err)

	p, err := NewProvider("aws", "eu-west-1", logger)
	require.NoError(t, err)

	t.Run("RedisEngine", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_elasticache_cluster.test",
			Type:         "aws_elasticache_cluster",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"node_type":       "cache.m4.large",
				"engine":          "redis",
				"num_cache_nodes": 1,
			},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{
			{
				Name:           "Cache instance",
				HourlyQuantity: decimal.NewFromInt(1),
				Details:        []string{"Redis"},
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AmazonElastiCache"),
					Family:   util.StringPtr("Cache Instance"),
					Location: util.StringPtr("eu-west-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "InstanceType", Value: util.StringPtr("cache.m4.large")},
						{Key: "CacheEngine", Value: util.StringPtr("Redis")},
					},
				},
				PriceFilter: &price.Filter{
					Unit: util.StringPtr("Hrs"),
					AttributeFilters: []*price.AttributeFilter{
						{Key: "TermType", Value: util.StringPtr("OnDemand")},
					},
				},
			},
		}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})

	t.Run("RedisSnapShotRetentionLimit", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_elasticache_cluster.test",
			Type:         "aws_elasticache_cluster",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"node_type":                "cache.m4.large",
				"engine":                   "redis",
				"num_cache_nodes":          1,
				"snapshot_retention_limit": 5,
			},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{
			{
				Name:           "Cache instance",
				HourlyQuantity: decimal.NewFromInt(1),
				Details:        []string{"Redis"},
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AmazonElastiCache"),
					Family:   util.StringPtr("Cache Instance"),
					Location: util.StringPtr("eu-west-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "InstanceType", Value: util.StringPtr("cache.m4.large")},
						{Key: "CacheEngine", Value: util.StringPtr("Redis")},
					},
				},
				PriceFilter: &price.Filter{
					Unit: util.StringPtr("Hrs"),
					AttributeFilters: []*price.AttributeFilter{
						{Key: "TermType", Value: util.StringPtr("OnDemand")},
					},
				},
			},
			{
				Name:            "Backup storage",
				Details:         []string{"0"},
				MonthlyQuantity: decimal.NewFromInt(0),
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AmazonElastiCache"),
					Family:   util.StringPtr("Storage Snapshot"),
					Location: util.StringPtr("eu-west-1"),
				},
				PriceFilter: &price.Filter{
					Unit: util.StringPtr("GB-Mo"),
					AttributeFilters: []*price.AttributeFilter{
						{Key: "TermType", Value: util.StringPtr("OnDemand")},
					},
				},
			},
		}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})

	t.Run("RedisReplicationGroupID", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_elasticache_cluster.test",
			Type:         "aws_elasticache_cluster",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"node_type":            "cache.m4.large",
				"engine":               "redis",
				"replication_group_id": "replication-group-1",
			},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})

	t.Run("MemcacheEngine", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_elasticache_cluster.test",
			Type:         "aws_elasticache_cluster",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"node_type":       "cache.m4.large",
				"engine":          "memcached",
				"num_cache_nodes": 1,
			},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{
			{
				Name:           "Cache instance",
				HourlyQuantity: decimal.NewFromInt(1),
				Details:        []string{"Memcached"},
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AmazonElastiCache"),
					Family:   util.StringPtr("Cache Instance"),
					Location: util.StringPtr("eu-west-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "InstanceType", Value: util.StringPtr("cache.m4.large")},
						{Key: "CacheEngine", Value: util.StringPtr("Memcached")},
					},
				},
				PriceFilter: &price.Filter{
					Unit: util.StringPtr("Hrs"),
					AttributeFilters: []*price.AttributeFilter{
						{Key: "TermType", Value: util.StringPtr("OnDemand")},
					},
				},
			},
		}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})

	t.Run("MemcacheNumCacheNodes", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_elasticache_cluster.test",
			Type:         "aws_elasticache_cluster",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"node_type":       "cache.m4.large",
				"engine":          "memcached",
				"num_cache_nodes": 2,
			},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{
			{
				Name:           "Cache instance",
				HourlyQuantity: decimal.NewFromInt(2),
				Details:        []string{"Memcached"},
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AmazonElastiCache"),
					Family:   util.StringPtr("Cache Instance"),
					Location: util.StringPtr("eu-west-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "InstanceType", Value: util.StringPtr("cache.m4.large")},
						{Key: "CacheEngine", Value: util.StringPtr("Memcached")},
					},
				},
				PriceFilter: &price.Filter{
					Unit: util.StringPtr("Hrs"),
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
