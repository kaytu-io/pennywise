package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/resource"
)

const ProviderName = "azurerm"

var (
	locationDisplayToName = map[string]string{
		"West US":              "westus",
		"West US 2":            "westus2",
		"East US":              "eastus",
		"Central US":           "centralus",
		"Central US EUAP":      "centraluseuap",
		"South Central US":     "southcentralus",
		"North Central US":     "northcentralus",
		"West Central US":      "westcentralus",
		"East US 2":            "eastus2",
		"East US 2 EUAP":       "eastus2euap",
		"Brazil South":         "brazilsouth",
		"Brazil US":            "brazilus",
		"North Europe":         "northeurope",
		"West Europe":          "westeurope",
		"East Asia":            "eastasia",
		"Southeast Asia":       "southeastasia",
		"Japan West":           "japanwest",
		"Japan East":           "japaneast",
		"Korea Central":        "koreacentral",
		"Korea South":          "koreasouth",
		"South India":          "southindia",
		"West India":           "westindia",
		"Central India":        "centralindia",
		"Australia East":       "australiaeast",
		"Australia Southeast":  "australiasoutheast",
		"Canada Central":       "canadacentral",
		"Canada East":          "canadaeast",
		"UK South":             "uksouth",
		"UK West":              "ukwest",
		"France Central":       "francecentral",
		"France South":         "francesouth",
		"Australia Central":    "australiacentral",
		"Australia Central 2":  "australiacentral2",
		"UAE Central":          "uaecentral",
		"UAE North":            "uaenorth",
		"South Africa North":   "southafricanorth",
		"South Africa West":    "southafricawest",
		"Switzerland North":    "switzerlandnorth",
		"Switzerland West":     "switzerlandwest",
		"Germany North":        "germanynorth",
		"Germany West Central": "germanywestcentral",
		"Norway East":          "norwayeast",
		"Norway West":          "norwaywest",
		"Brazil Southeast":     "brazilsoutheast",
		"West US 3":            "westus3",
		"East US SLV":          "eastusslv",
		"Sweden Central":       "swedencentral",
		"Sweden South":         "swedensouth",
	}
)

// Provider is an implementation of the resources.Provider, used to extract component queries from
// terraform resources.
type Provider struct {
	key string
}

// NewProvider initializes a new Azure provider with key and region
func NewProvider(key string) (*Provider, error) {
	return &Provider{
		key: key,
	}, nil
}

// Name returns the Provider's common name.
func (p *Provider) Name() string { return p.key }

// ResourceComponents returns Component queries for a given terraform.Resource.
func (p *Provider) ResourceComponents(rss map[string]resource.Resource, tfRes resource.Resource) []query.Component {
	switch tfRes.Type {
	case "azurerm_linux_virtual_machine":
		vals, err := decodeLinuxVirtualMachineValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newLinuxVirtualMachine(vals).Components()
	case "azurerm_virtual_machine":
		vals, err := decodeVirtualMachineValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newVirtualMachine(vals).Components()
	case "azurerm_windows_virtual_machine":
		vals, err := decodeWindowsVirtualMachineValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newWindowsVirtualMachine(vals).Components()
	case "azurerm_managed_disk":
		vals, err := decodeManagedStorageValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newManagedStorage(vals).Components()
	case "azurerm_image":
		vals, err := decodeImageValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newImage(vals).Components()
	case "azurerm_snapshot":
		vals, err := decodeSnapshotValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newSnapshot(vals).Components()
	case "azurerm_virtual_network":
		_, err := decodeVirtualNetworkValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return nil
	//return p.newVirtualNetwork(vals).Components()
	case "azurerm_virtual_network_gateway":
		vals, err := decoderVirtualNetworkGateway(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newVirtualNetworkGateway(vals).Components()
	case "azurerm_virtual_network_gateway_connection":
		vals, err := decoderVirtualNetworkGatewayConnection(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newVirtualNetworkGatewayConnection(vals).Component()
	case "azurerm_load_balancer":
		return p.newLoadBalancer(loadBalancerValues{}).Components()
	case "azurerm_dns_a_record":
		vals, err := decoderDNSARecord(tfRes.Values)
		if err != nil {
			return nil
		}
		monthlyQueries := int64(1500000000)
		vals.MonthlyQueries = &monthlyQueries
		return p.newDNSARecord(vals).component()
	case "azurerm_dns_aaaa_record":
		vals, err := decoderDNSAAAARecord(tfRes.Values)
		if err != nil {
			return nil
		}
		monthlyQueries := int64(1500000000)
		vals.MonthlyQueries = monthlyQueries
		return p.newDNSAAAARecord(vals).component()
	case "azurerm_dns_caa_record":
		vals, err := decoderDNSCAARecord(tfRes.Values)
		if err != nil {
			return nil
		}
		monthlyQueries := int64(1000000000)
		vals.MonthlyQueries = monthlyQueries
		return p.newDNSCAARecord(vals).component()
	case "azurerm_dns_cname_record":
		vals, err := decoderDNSCNAMERecord(tfRes.Values)
		if err != nil {
			return nil
		}
		monthlyQueries := int64(1000000000)
		vals.MonthlyQueries = monthlyQueries
		return p.newDNSCNAMERecord(vals).component()
	case "azurerm_dns_mx_record":
		vals, err := decoderDNSMXRecord(tfRes.Values)
		if err != nil {
			return nil
		}
		monthlyQueries := int64(1000000000)
		vals.MonthlyQueries = monthlyQueries
		return p.newDNSMXRecord(vals).component()
	case "azurerm_dns_ns_record":
		vals, err := decoderDNSNSRecord(tfRes.Values)
		if err != nil {
			return nil
		}
		monthlyQueries := int64(1000000000)
		vals.MonthlyQueries = monthlyQueries
		return p.newDNSNSRecord(vals).component()
	case "azurerm_dns_ptr_record":
		vals, err := decoderDNSPTRRecord(tfRes.Values)
		if err != nil {
			return nil
		}
		monthlyQueries := int64(1000000000)
		vals.MonthlyQueries = monthlyQueries
		return p.newDNSPTRRecord(vals).component()
	case "azurerm_dns_srv_record":
		vals, err := decoderDNSARecord(tfRes.Values)
		if err != nil {
			return nil
		}
		monthlyQueries := int64(1000000000)
		vals.MonthlyQueries = &monthlyQueries
		return p.newDNSARecord(vals).component()
	case "azurerm_dns_txt_record":
		vals, err := decoderDNSTXTRecord(tfRes.Values)
		if err != nil {
			return nil
		}
		monthlyQueries := int64(1000000000)
		vals.MonthlyQueries = monthlyQueries
		return p.newDNSTXTRecord(vals).component()
	case "azurerm_dns_zone":
		vals, err := decoderRMDNSZone(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newRMDNSZone(vals).component()
	case "azurerm_container_registry":
		vals, err := decoderContainerRegistry(tfRes.Values)
		if err != nil {
			return nil
		}
		storageGb := 150.0
		monthlyBuildVCPUHrs := 150.0
		vals.StorageGB = &storageGb
		vals.MonthlyBuildVCPUHrs = &monthlyBuildVCPUHrs
		return p.newContainerRegistry(vals).component()
	case "azurerm_kubernetes_cluster":
		vals, err := decoderKubernetesCluster(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.NewAzureRMKubernetesCluster(vals).Components()
	case "azurerm_kubernetes_cluster_node_pool":
		vals, err := decoderKubernetesClusterNodePool(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newAzureRMKubernetesClusterNodePool(vals).Components()
	default:
		return nil
	}
}

// getLocationName will return the location name from the location display name (ex: UK West -> ukwest)
// if the l is not found it'll return the l again meaning is not found or already a name
func getLocationName(l string) string {
	ln, ok := locationDisplayToName[l]
	if !ok {
		return l
	}
	return ln
}
