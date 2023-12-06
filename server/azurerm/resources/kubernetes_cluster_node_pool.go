package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

type kubernetesClusterNodePool struct {
	provider *Provider

	monthlyHrs   decimal.Decimal
	location     string
	minCount     *int64
	nodeCount    *int64
	nodes        *int64
	vmSize       string
	skuTier      string
	osSku        *string
	osDiskType   *string
	osDiskSizeGB *int
}

type kubernetesClusterNodePoolValues struct {
	Location        string `mapstructure:"location"`
	SkuTier         string `mapstructure:"sku_tier"`
	DefaultNodePool struct {
		MinCount     int64  `mapstructure:"min_count"`
		NodeCount    int64  `mapstructure:"node_count"`
		VmSize       string `mapstructure:"vm_size"`
		OSSku        string `mapstructure:"os_sku"`
		OSDiskType   string `json:"os_disk_type"`
		OsDiskSizeGB int    `json:"os_disk_size_gb"`
	} `mapstructure:"default_node_pool"`

	// TODO: we should get this fields from user
	Nodes      int64   `json:"nodes"`
	MonthlyHrs float64 `json:"monthly_hrs"`
}

// decoderKubernetesCluster decodes and returns kubernetesClusterValues from a Terraform values map.
func decoderKubernetesClusterNodePool(tfVals map[string]interface{}) (kubernetesClusterNodePoolValues, error) {
	var v kubernetesClusterNodePoolValues
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

func (p *Provider) newAzureRMKubernetesClusterNodePool(vals kubernetesClusterNodePoolValues) *kubernetesClusterNodePool {
	inst := &kubernetesClusterNodePool{
		location:   vals.Location,
		skuTier:    vals.SkuTier,
		minCount:   &vals.DefaultNodePool.MinCount,
		nodeCount:  &vals.DefaultNodePool.NodeCount,
		vmSize:     vals.DefaultNodePool.VmSize,
		osSku:      &vals.DefaultNodePool.OSSku,
		osDiskType: &vals.DefaultNodePool.OSDiskType,

		nodes:      &vals.Nodes,
		monthlyHrs: decimal.NewFromFloat(vals.MonthlyHrs),
	}
	return inst
}

func (inst *kubernetesClusterNodePool) Components() []query.Component {
	return aksClusterNodePool(inst, inst.location)
}

func aksClusterNodePool(vals *kubernetesClusterNodePool, region string) []query.Component {
	var costComponents []query.Component

	instanceType := vals.vmSize

	os := "Linux"
	if vals.osSku != nil {
		os = *vals.osSku
	}

	if vals.osSku != nil {
		if strings.HasSuffix(strings.ToLower(*vals.osSku), "windows") {
			os = "Windows"
		}
	}

	if strings.EqualFold(os, "windows") {
		purchaseOption := "Consumption"
		costComponents = append(costComponents, windowsVirtualMachineComponent(vals.provider.key, region, instanceType, purchaseOption, vals.monthlyHrs))
	} else {
		costComponents = append(costComponents, linuxVirtualMachineComponent(vals.provider.key, region, instanceType, vals.monthlyHrs))
	}

	osDiskType := "Managed"
	if vals.osDiskType != nil {
		osDiskType = *vals.osDiskType
	}

	if strings.ToLower(osDiskType) == "managed" {
		osDisk := aksOSDiskSubResource(vals)
		if osDisk != nil {
			costComponents = append(costComponents, *osDisk)
		}
	}

	return costComponents
}

func aksOSDiskSubResource(inst *kubernetesClusterNodePool) *query.Component {
	//managedStorageInst := ManagedDisk{}
	//costComponent := managedStorageInst.managedStorageComponent(inst.provider.key, inst.location, inst.skuTier)
	//return &costComponent
	return nil
}
