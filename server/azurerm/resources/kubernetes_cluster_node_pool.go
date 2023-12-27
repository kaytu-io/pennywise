package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"regexp"
	"strings"
)

type kubernetesClusterNodePool struct {
	provider *Provider

	location     string
	minCount     *int64
	nodeCount    *int64
	vmSize       string
	osSku        *string
	osDiskType   *string
	osDiskSizeGB *int
	osType       *string
	// Usage
	// receives node count for the default node pool
	nodes *int64
	// receives monthly hours for the default node pool
	monthlyHrs *float64
}
type KubernetesClusterIdStruct struct {
	Values struct {
		Location string `mapstructure:"location"`
	} `mapstructure:"values"`
}
type kubernetesClusterNodePoolValues struct {
	KubernetesClusterIdStruct KubernetesClusterIdStruct `mapstructure:"kubernetes_cluster_id"`

	MinCount     *int64  `mapstructure:"min_count"`
	NodeCount    *int64  `mapstructure:"node_count"`
	VmSize       string  `mapstructure:"vm_size"`
	OSSku        *string `mapstructure:"os_sku"`
	OSType       *string `mapstructure:"os_type"`
	OSDiskType   *string `mapstructure:"os_disk_type"`
	OsDiskSizeGB *int    `mapstructure:"os_disk_size_gb"`

	Usage struct {
		Nodes      *int64   `mapstructure:"nodes"`
		MonthlyHrs *float64 `mapstructure:"monthly_hrs"`
	} `mapstructure:"pennywise_usage"`
}

// decodeKubernetesClusterNodePool decodes and returns kubernetesClusterValues from a Terraform values map.
func decodeKubernetesClusterNodePool(tfVals map[string]interface{}) (kubernetesClusterNodePoolValues, error) {
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
		provider: p,
		location: vals.KubernetesClusterIdStruct.Values.Location,

		minCount:     vals.MinCount,
		nodeCount:    vals.NodeCount,
		vmSize:       vals.VmSize,
		osSku:        vals.OSSku,
		osType:       vals.OSType,
		osDiskType:   vals.OSDiskType,
		osDiskSizeGB: vals.OsDiskSizeGB,
		//usage :
		nodes:      vals.Usage.Nodes,
		monthlyHrs: vals.Usage.MonthlyHrs,
	}
	return inst
}

func (inst *kubernetesClusterNodePool) Components() []query.Component {
	nodeCount := decimal.NewFromInt(1)
	region := getLocationName(inst.location)
	// if the node count is not set explicitly let's take the min_count.
	if inst.minCount != nil {
		nodeCount = decimal.NewFromInt(*inst.minCount)
	}

	if inst.nodeCount != nil {
		nodeCount = decimal.NewFromInt(*inst.nodeCount)
	}

	if inst.nodes != nil {
		nodeCount = decimal.NewFromInt(*inst.nodes)
	}

	costComponents := aksClusterNodePool(inst.provider.key, inst.osDiskSizeGB, inst.osDiskType, inst.osSku, inst.osType,
		inst.vmSize, region, nodeCount, inst.monthlyHrs)

	GetCostComponentNamesAndSetLogger(costComponents, inst.provider.logger)
	return costComponents
}

func aksClusterNodePool(key string, osDiskSizeGB *int, OSDiskType *string, osSku *string, osType *string,
	vmsSize string, region string, nodeCount decimal.Decimal, monthlyHrsUsage *float64) []query.Component {
	var costComponents []query.Component

	for i := int64(0); i < nodeCount.CoefficientInt64(); i++ {
		var monthlyHrs *decimal.Decimal
		if monthlyHrsUsage != nil {
			monthlyHrs = util.DecimalPtr(decimal.NewFromFloat(*monthlyHrsUsage))
		}

		instanceType := vmsSize

		os := "Linux"
		if osSku != nil {
			os = *osSku
		}

		if osType != nil {
			if strings.HasSuffix(strings.ToLower(*osType), "windows") {
				os = "Windows"
			}
		}

		if strings.Contains(strings.ToLower(os), "windows") {
			purchaseOption := "Consumption"
			costComponents = append(costComponents, windowsVirtualMachineComponent(key, region, instanceType, purchaseOption, monthlyHrs))
		} else {
			costComponents = append(costComponents, linuxVirtualMachineComponent(key, region, instanceType, monthlyHrs))
		}

		osDiskType := "Managed"
		if OSDiskType != nil {
			osDiskType = *OSDiskType
		}

		if strings.ToLower(osDiskType) == "managed" {
			diskSize := 128
			if osDiskSizeGB != nil {
				diskSize = *osDiskSizeGB
			}
			osDisk := aksOSDiskSubResource(key, region, instanceType, diskSize)
			if osDisk != nil {
				costComponents = append(costComponents, *osDisk)
			}
		}
	}

	return costComponents
}

func aksOSDiskSubResource(key string, region string, instanceType string, diskSize int) *query.Component {
	diskType := aksGetStorageType(instanceType)
	storageReplicationType := "LRS"

	diskName := mapDiskName(diskType, diskSize)
	if diskName == "" {
		fmt.Printf("Could not map disk type %s and size %d to disk name", diskType, diskSize)
		return nil
	}

	productName, ok := diskProductNameMap[diskType]
	if !ok {
		fmt.Printf("Could not map disk type %s to product name", diskType)
		return nil
	}

	managedStorageInst := ManagedDisk{}
	costComponent := managedStorageInst.managedStorageComponent(key, region, diskName, storageReplicationType, productName)
	return &costComponent
}

func aksGetStorageType(instanceType string) string {
	parts := strings.Split(instanceType, "_")

	subfamily := ""
	if len(parts) > 1 {
		subfamily = parts[1]
	}

	// Check if the subfamily is a known premium type
	premiumPrefixes := []string{"ds", "gs", "m"}
	for _, p := range premiumPrefixes {
		if strings.HasPrefix(strings.ToLower(subfamily), p) {
			return "Premium"
		}
	}

	// Otherwise check if it contains an s as an 'Additive Feature'
	// as per https://learn.microsoft.com/en-us/azure/virtual-machines/vm-naming-conventions
	re := regexp.MustCompile(`\d+[A-Za-z]*(s)`)
	matches := re.FindStringSubmatch(subfamily)

	if len(matches) > 0 {
		return "Premium"
	}

	return "Standard"
}
