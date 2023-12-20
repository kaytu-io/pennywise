package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

type KubernetesCluster struct {
	provider *Provider

	skuTier                       *string
	location                      string
	httpApplicationRoutingEnabled bool

	nodeCount    *int64
	minCount     *int64
	osSku        *string
	osType       *string
	vmSize       string
	osDiskType   *string
	osDiskSizeGB *int

	loadBalancerSku                          *string
	addonProfileHttpApplicationRoutingEnable *bool

	// Usage
	loadBalancerMonthlyDataProcessedGB *int64
	defaultNodePoolNodes               *int64
	monthlyHrs                         *float64
}

type DefaultNodePoolStruct struct {
	MinCount     *int64  `mapstructure:"min_count"`
	NodeCount    *int64  `mapstructure:"node_count"`
	OSSku        *string `mapstructure:"os_sku"`
	OSType       *string `mapstructure:"os_type"`
	VmSize       string  `mapstructure:"vm_size"`
	OSDiskType   *string `mapstructure:"os_disk_type"`
	OsDiskSizeGB *int    `mapstructure:"os_disk_size_gb"`
}

type AddonProfileStruct struct {
	HttpApplicationRouting struct {
		Enable *bool `mapstructure:"enable"`
	} `mapstructure:"http_application_routing"`
}

type kubernetesClusterValues struct {
	SkuTier                       *string `mapstructure:"sku_tier"`
	Location                      string  `mapstructure:"location"`
	HttpApplicationRoutingEnabled bool    `mapstructure:"http_application_routing_enabled"`

	DefaultNodePool []DefaultNodePoolStruct `mapstructure:"default_node_pool"`
	AddonProfile    []AddonProfileStruct    `mapstructure:"addon_profile"`

	NetworkProfile []struct {
		LoadBalancerSku *string `mapstructure:"load_balancer_sku"`
	} `mapstructure:"network_profile"`

	Usage struct {
		// receives node count for the default node pool
		Nodes                  *int64   `mapstructure:"nodes"`
		MonthlyHrs             *float64 `mapstructure:"monthly_hrs"`
		MonthlyDataProcessedGB *int64   `mapstructure:"monthly_data_processed_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decoderKubernetesCluster decodes and returns kubernetesClusterValues from a Terraform values map.
func decoderKubernetesCluster(tfVals map[string]interface{}) (kubernetesClusterValues, error) {
	var v kubernetesClusterValues
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

func (p *Provider) NewAzureRMKubernetesCluster(vals kubernetesClusterValues) *KubernetesCluster {
	var DefaultNodePool DefaultNodePoolStruct
	var AddonProfile AddonProfileStruct

	if len(vals.DefaultNodePool) > 0 {
		DefaultNodePool = vals.DefaultNodePool[0]
	}
	if len(vals.AddonProfile) > 0 {
		AddonProfile = vals.AddonProfile[0]
	}

	var loadBalancerSku *string
	if len(vals.NetworkProfile) > 0 {
		loadBalancerSku = vals.NetworkProfile[0].LoadBalancerSku
	}
	inst := &KubernetesCluster{
		provider:                                 p,
		skuTier:                                  vals.SkuTier,
		location:                                 vals.Location,
		httpApplicationRoutingEnabled:            vals.HttpApplicationRoutingEnabled,
		nodeCount:                                DefaultNodePool.NodeCount,
		minCount:                                 DefaultNodePool.MinCount,
		osSku:                                    DefaultNodePool.OSSku,
		osType:                                   DefaultNodePool.OSType,
		vmSize:                                   DefaultNodePool.VmSize,
		osDiskType:                               DefaultNodePool.OSDiskType,
		osDiskSizeGB:                             DefaultNodePool.OsDiskSizeGB,
		addonProfileHttpApplicationRoutingEnable: AddonProfile.HttpApplicationRouting.Enable,
		loadBalancerSku:                          loadBalancerSku,

		//usage
		defaultNodePoolNodes:               vals.Usage.Nodes,
		monthlyHrs:                         vals.Usage.MonthlyHrs,
		loadBalancerMonthlyDataProcessedGB: vals.Usage.MonthlyDataProcessedGB,
	}
	return inst
}

func (inst KubernetesCluster) Components() []query.Component {
	region := inst.location
	region = getLocationName(region)

	var costComponents []query.Component

	skuTier := "Free"
	if inst.skuTier != nil {
		skuTier = *inst.skuTier
	}

	// Azure switched from "Paid" to "Standard" in API version 2023-02-01
	// (Terraform Azure provider version v3.51.0)
	if contains([]string{"paid", "standard"}, strings.ToLower(skuTier)) {
		costComponents = append(costComponents, query.Component{
			Name:           "Uptime SLA",
			Unit:           "hours",
			HourlyQuantity: decimal.NewFromInt(1),
			ProductFilter: &product.Filter{
				Provider: util.StringPtr(inst.provider.key),
				Location: util.StringPtr(region),
				Service:  util.StringPtr("Azure Kubernetes Service"),
				Family:   util.StringPtr("Compute"),
				AttributeFilters: []*product.AttributeFilter{
					{Key: "meter_name", Value: util.StringPtr("Standard Uptime SLA")},
				},
			},
			PriceFilter: &price.Filter{
				AttributeFilters: []*price.AttributeFilter{
					{Key: "type", Value: util.StringPtr("Consumption")},
				},
			},
		})
	}

	nodeCount := decimal.NewFromInt(1)
	if inst.nodeCount != nil {
		nodeCount = decimal.NewFromInt(*inst.nodeCount)
	}

	// if the node count is not set explicitly let's take the min_count.
	if inst.minCount != nil && nodeCount.Equal(decimal.NewFromInt(1)) {
		nodeCount = decimal.NewFromInt(*inst.minCount)
	}
	if inst.defaultNodePoolNodes != nil {
		nodeCount = decimal.NewFromInt(*inst.defaultNodePoolNodes)
	}

	// TODO: check the input values if we had os_disk_type_db put its value as NodePool input
	costComponents = append(costComponents, aksClusterNodePool(inst.provider.key, inst.osDiskSizeGB, inst.osDiskType,
		inst.osSku, inst.osType, inst.vmSize, region, nodeCount, inst.monthlyHrs)...)

	if inst.loadBalancerSku != nil {
		if strings.ToLower(*inst.loadBalancerSku) == "standard" {
			region = convertRegion(region)
			var monthlyDataProcessedGb decimal.Decimal
			if inst.loadBalancerMonthlyDataProcessedGB != nil {
				monthlyDataProcessedGb = decimal.NewFromInt(*inst.loadBalancerMonthlyDataProcessedGB)
			}
			costComponents = append(costComponents, regionalDataProceedComponent(inst.provider.key, region, monthlyDataProcessedGb))
		}
	}

	routingEnabled := inst.httpApplicationRoutingEnabled
	// Deprecated and removed in v3
	if inst.addonProfileHttpApplicationRoutingEnable != nil {
		routingEnabled = *inst.addonProfileHttpApplicationRoutingEnable
	}

	if routingEnabled {
		if strings.HasPrefix(strings.ToLower(region), "usgov") {
			region = "US Gov Zone 1"
		} else if strings.HasPrefix(strings.ToLower(region), "germany") {
			region = "DE Zone 1"
		} else if strings.HasPrefix(strings.ToLower(region), "china") {
			region = "Zone 1 (China)"
		} else {
			region = "Zone 1"
		}

		costComponents = append(costComponents, hostedPublicZoneCostComponent(inst.provider.key, region))
	}

	return costComponents
}
