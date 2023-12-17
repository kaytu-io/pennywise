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

func TestEKSNodeGroup_Components(t *testing.T) {
	p, err := NewProvider("aws", "eu-west-1")
	require.NoError(t, err)

	t.Run("EKSNodeGroupDefault", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "module.test.aws_eks_node_group.test",
			Type:         "aws_eks_node_group",
			Name:         "test",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"scaling_config": []interface{}{map[string]interface{}{
					"desired_size": 3,
					"min_size":     2,
					"max_size":     5,
				}},
				"disk_size":      30,
				"instance_types": []string{"t3.large"},
			},
		}

		rss := map[string]resource.Resource{
			"aws_launch_template.test": resource.Resource{
				Address:      "aws_launch_template.test",
				Type:         "aws_launch_template",
				Name:         "test",
				ProviderName: "aws",
				Values: map[string]interface{}{
					"instance_type": "m5.xlarge",
					"placement":     []interface{}{map[string]interface{}{"availability_zone": "eu-west-1a"}},
					"block_device_mappings": []interface{}{
						map[string]interface{}{
							"device_name": "/dev/sda1",
							"ebs": []interface{}{
								map[string]interface{}{"volume_size": float64(30)},
							},
						},
					},
				},
			},
		}

		expected := []query.Component{
			{
				Name:           "Compute",
				HourlyQuantity: decimal.NewFromInt(3),
				Details:        []string{"Linux", "on-demand", "t3.large"},
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AmazonEC2"),
					Family:   util.StringPtr("Compute Instance"),
					Location: util.StringPtr("eu-west-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "CapacityStatus", Value: util.StringPtr("Used")},
						{Key: "InstanceType", Value: util.StringPtr("t3.large")},
						{Key: "Tenancy", Value: util.StringPtr("Shared")},
						{Key: "OperatingSystem", Value: util.StringPtr("Linux")},
						{Key: "PreInstalledSW", Value: util.StringPtr("NA")},
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
				Name:            "Root volume: Storage",
				MonthlyQuantity: decimal.NewFromFloat(30),
				Unit:            "GB",
				Details:         []string{"gp3"},
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AmazonEC2"),
					Family:   util.StringPtr("Storage"),
					Location: util.StringPtr("eu-west-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "VolumeAPIName", Value: util.StringPtr("gp3")},
					},
				},
			},
		}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})

	t.Run("EKSNodeGroupLaunchTemplate", func(t *testing.T) {
		tfres := resource.Resource{
			Address:      "module.test.aws_eks_node_group.lt",
			Type:         "aws_eks_node_group",
			Name:         "lt",
			ProviderName: "aws",
			Values: map[string]interface{}{
				"scaling_config": []interface{}{map[string]interface{}{
					"desired_size": 2,
					"min_size":     1,
					"max_size":     5,
				}},
				"launch_template": []interface{}{map[string]interface{}{"id": "aws_launch_template.test"}},
			},
		}

		rss := map[string]resource.Resource{
			"aws_launch_template.test": resource.Resource{
				Address:      "aws_launch_template.test",
				Type:         "aws_launch_template",
				Name:         "test",
				ProviderName: "aws",
				Values: map[string]interface{}{
					"instance_type": "m5.xlarge",
					"placement":     []interface{}{map[string]interface{}{"availability_zone": "eu-west-1a"}},
					"block_device_mappings": []interface{}{
						map[string]interface{}{
							"device_name": "/dev/sda1",
							"ebs": []interface{}{
								map[string]interface{}{"volume_size": float64(50)},
							},
						},
					},
				},
			},
		}

		expected := []query.Component{
			{
				Name:           "Compute",
				HourlyQuantity: decimal.NewFromInt(2),
				Details:        []string{"Linux", "on-demand", "m5.xlarge"},
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AmazonEC2"),
					Family:   util.StringPtr("Compute Instance"),
					Location: util.StringPtr("eu-west-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "CapacityStatus", Value: util.StringPtr("Used")},
						{Key: "InstanceType", Value: util.StringPtr("m5.xlarge")},
						{Key: "Tenancy", Value: util.StringPtr("Shared")},
						{Key: "OperatingSystem", Value: util.StringPtr("Linux")},
						{Key: "PreInstalledSW", Value: util.StringPtr("NA")},
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
				Name:            "Root volume: Storage",
				MonthlyQuantity: decimal.NewFromFloat(50),
				Unit:            "GB",
				Details:         []string{"gp3"},
				ProductFilter: &product.Filter{
					Provider: util.StringPtr("aws"),
					Service:  util.StringPtr("AmazonEC2"),
					Family:   util.StringPtr("Storage"),
					Location: util.StringPtr("eu-west-1"),
					AttributeFilters: []*product.AttributeFilter{
						{Key: "VolumeAPIName", Value: util.StringPtr("gp3")},
					},
				},
			},
		}

		actual := p.ResourceComponents(rss, tfres)
		assert.Equal(t, expected, actual)
	})
}
