package azurerm

// RegistryName is the fully qualified name under which this provider is stored in the registry.
const RegistryName = "registry.terraform.io/hashicorp/azurerm"

// TerraformProviderInitializer is a terraform.ProviderInitializer that initializes the default GCP provider.
//var TerraformProviderInitializer = resource.ProviderInitializer{
//	MatchNames: []string{ProviderName, RegistryName},
//	Provider: func(values map[string]interface{}) (resource.Provider, error) {
//		return azurerm.NewProvider(ProviderName, nil)
//	},
//}
