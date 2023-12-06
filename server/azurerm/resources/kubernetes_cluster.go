package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

type kubernetesCluster struct {
	provider *Provider

	skuTier                       string
	location                      string
	httpApplicationRoutingEnabled bool

	nodeCount                 *int64
	minCount                  *int64
	osSku                     *string
	vmSize                    string
	osDiskType                *string
	defaultNodePoolNodes      *int64
	defaultNodePoolMonthlyHrs int64

	loadBalancerSku *string

	LoadBalancerMonthlyDataProcessedGB *int64

	addonProfileHttpApplicationRoutingEnable *bool
}

type kubernetesClusterValues struct {
	SkuTier                       string `mapstructure:"sku_tier"`
	Location                      string `mapstructure:"location"`
	HttpApplicationRoutingEnabled bool   `mapstructure:"http_application_routing_enabled"`

	DefaultNodePool struct {
		MinCount   int64  `mapstructure:"min_count"`
		NodeCount  int64  `mapstructure:"node_count"`
		OSSku      string `mapstructure:"os_sku"`
		VmSize     string `mapstructure:"vm_size"`
		OSDiskType string `mapstructure:"os_disk_type"`

		//TODO:we should get Nodes , MonthlyHrs fields from user
		Nodes      int64 `mapstructure:"nodes"`
		MonthlyHrs int64 `mapstructure:"monthly_hrs"`
	} `mapstructure:"default_node_pool"`

	NetworkProfile struct {
		LoadBalancerSku string `mapstructure:"load_balancer_sku"`
	} `mapstructure:"network_profile"`

	AddonProfile struct {
		HttpApplicationRouting struct {
			Enable bool `mapstructure:"enable"`
		} `mapstructure:"http_application_routing"`
	} `mapstructure:"addon_profile"`

	//TODO:we should get this field from user
	LoadBalancer struct {
		MonthlyDataProcessedGB int64 `mapstructure:"monthly_data_processed_gb"`
	} `mapstructure:"load_balancer"`
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

func (p *Provider) NewAzureRMKubernetesCluster(vals kubernetesClusterValues) *kubernetesCluster {
	inst := &kubernetesCluster{
		skuTier:                       vals.SkuTier,
		location:                      vals.Location,
		httpApplicationRoutingEnabled: vals.HttpApplicationRoutingEnabled,
		nodeCount:                     &vals.DefaultNodePool.NodeCount,
		minCount:                      &vals.DefaultNodePool.MinCount,
		osSku:                         &vals.DefaultNodePool.OSSku,
		vmSize:                        vals.DefaultNodePool.VmSize,
		osDiskType:                    &vals.DefaultNodePool.OSDiskType,
		defaultNodePoolNodes:          &vals.DefaultNodePool.Nodes,
		defaultNodePoolMonthlyHrs:     vals.DefaultNodePool.MonthlyHrs,

		loadBalancerSku:                          &vals.NetworkProfile.LoadBalancerSku,
		LoadBalancerMonthlyDataProcessedGB:       &vals.LoadBalancer.MonthlyDataProcessedGB,
		addonProfileHttpApplicationRoutingEnable: &vals.AddonProfile.HttpApplicationRouting.Enable,
	}
	return inst
}

func contains(wordsThatWantToCheck []string, text string) bool {
	for _, a := range wordsThatWantToCheck {
		if strings.Contains(a, text) {
			return true
		}
	}
	return false
}

func (inst kubernetesCluster) Components() []query.Component {
	var costComponents []query.Component
	region := inst.location
	skuTier := "Free"
	if inst.skuTier != "" {
		skuTier = inst.skuTier
	}

	// Azure switched from "Paid" to "Standard" in API version 2023-02-01
	// (Terraform Azure provider version v3.51.0)
	if contains([]string{"paid", "standard"}, strings.ToLower(skuTier)) {
		costComponents = append(costComponents, query.Component{
			Name:           "Uptime SLA",
			Unit:           "hours",
			HourlyQuantity: decimal.NewFromInt(1),
			ProductFilter: &product.Filter{
				Location: util.StringPtr(region),
				Service:  util.StringPtr("Azure Kubernetes Service"),
				Family:   util.StringPtr("Compute"),
				AttributeFilters: []*product.AttributeFilter{
					{Key: "meter_name", Value: util.StringPtr("Standard Uptime SLA")},
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
	aksClusterNodePoolValue := &kubernetesClusterNodePool{
		provider:   inst.provider,
		osSku:      inst.osSku,
		vmSize:     inst.vmSize,
		osDiskType: inst.osDiskType,
		skuTier:    inst.skuTier,
	}
	kubernetesClusterNodePoolQueries := aksClusterNodePool(aksClusterNodePoolValue, region)
	for _, qu := range kubernetesClusterNodePoolQueries {
		costComponents = append(costComponents, qu)
	}

	if inst.loadBalancerSku != nil {
		if strings.ToLower(*inst.loadBalancerSku) != "standard" {
			region = convertRegion(region)
			var monthlyDataProcessedGb decimal.Decimal
			if inst.LoadBalancerMonthlyDataProcessedGB != nil {
				monthlyDataProcessedGb = decimal.NewFromInt(*inst.LoadBalancerMonthlyDataProcessedGB)
			}

			costComponents = append(costComponents, regionalDataProceedComponent(inst.provider.key, region, monthlyDataProcessedGb))
		}
	}

	routingEnabled := inst.httpApplicationRoutingEnabled
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
		costComponents = append(costComponents, hostedPublicZoneCostComponent(region))
	}
	return costComponents
}
