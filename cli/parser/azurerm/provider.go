package azurerm

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

// Provider is an implementation of the terraform.Provider, used to extract component queries from
// terraform resources.
type Provider struct {
	key string
}

// NewProvider initializes a new Google provider with key and region
func NewProvider(key string) (*Provider, error) {
	return &Provider{
		key: key,
	}, nil
}

// Name returns the Provider's common name.
func (p *Provider) Name() string { return p.key }
