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
	case "azurerm_linux_virtual_machine_scale_set":
		vals, err := decodeLinuxVirtualMachineScaleSetValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newLinuxVirtualMachineScaleSet(vals).Components()
	case "azurerm_windows_virtual_machine_scale_set":
		vals, err := decodeWindowsVirtualMachineScaleSetValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newWindowsVirtualMachineScaleSet(vals).Components()
	case "azurerm_virtual_machine_scale_set":
		vals, err := decodeVirtualMachineScaleSetValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newVirtualMachineScaleSet(vals).Components()
	case "azurerm_lb":
		vals, err := decodeLoadBalancerValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newLoadBalancer(vals).Components()
	case "azurerm_lb_rule":
		vals, err := decodeLoadBalancerRuleValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newLoadBalancerRule(vals).Components()
	case "azurerm_lb_outbound_rule":
		vals, err := decodeLoadBalancerRuleValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newLoadBalancerRule(vals).Components()
	case "azurerm_application_gateway":
		vals, err := decodeApplicationGatewayValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newApplicationGateway(vals).Components()
	case "azurerm_nat_gateway":
		vals, err := decodeNatGatewayValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newNatGateway(vals).Components()
	case "azurerm_public_ip":
		vals, err := decodePublicIPValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newPublicIP(vals).Components()
	case "azurerm_public_ip_prefix":
		vals, err := decodePublicIPPrefixValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newPublicIPPrefix(vals).Components()
	case "azurerm_container_registry":
		vals, err := decodeContainerRegistry(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newContainerRegistry(vals).component()
	case "azurerm_private_endpoint":
		vals, err := decodePrivateEndpointValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newPrivateEndpoint(vals).Components()
	case "azurerm_storage_queue":
		vals, err := decodeStorageQueueValues(tfRes.Values)
		if err != nil {
			return nil
		}
		return p.newStorageQueue(vals).Components()
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
