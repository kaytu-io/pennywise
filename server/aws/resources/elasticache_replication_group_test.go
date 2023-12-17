package resources

import (
	"github.com/kaytu-io/pennywise/server/resource"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
)

func TestElastiCacheReplication_Components(t *testing.T) {
	p, err := NewProvider("aws", "eu-west-1")
	require.NoError(t, err)

	//1 group 1 node
	t.Run("RedisEngineDefault", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_elasticache_replication_group.test",
			Type:         "aws_elasticache_replication_group",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"node_type":          "cache.m4.large",
				"num_cache_clusters": 1,
				"availability_zones": []string{"eu-west-1a", "eu-west-1b"},
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

	t.Run("RedisEngineGlobalReplicationGroupID", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_elasticache_replication_group.test",
			Type:         "aws_elasticache_replication_group",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"num_cache_clusters":          1,
				"availability_zones":          []string{"eu-west-1a", "eu-west-1b"},
				"global_replication_group_id": "global-replication-group-1",
			},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})

	t.Run("RedisSnapShotRetentionLimit", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_elasticache_replication_group.test",
			Type:         "aws_elasticache_replication_group",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"node_type":                "cache.m4.large",
				"engine":                   "redis",
				"num_cache_clusters":       1,
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

	t.Run("RedisEngineNumCacheNodes", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_elasticache_replication_group.test",
			Type:         "aws_elasticache_replication_group",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"node_type":          "cache.m4.large",
				"engine":             "redis",
				"num_cache_clusters": 2,
			},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{
			{
				Name:           "Cache instance",
				HourlyQuantity: decimal.NewFromInt(2),
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

	t.Run("RedisEngineClusterMode", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "aws_elasticache_replication_group.test",
			Type:         "aws_elasticache_replication_group",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"node_type": "cache.m4.large",
				"engine":    "redis",
				"cluster_mode": []map[string]int{
					map[string]int{
						"replicas_per_node_group": 3,
						"num_node_groups":         2,
					},
				},
			},
		}
		rss := map[string]resource.Resource{}

		expected := []query.Component{
			{
				Name:           "Cache instance",
				HourlyQuantity: decimal.NewFromInt(8),
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

}
