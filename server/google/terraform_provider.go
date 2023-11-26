package google

import (
	googletf "github.com/kaytu.io/pennywise/server/google/terraform"
	"github.com/kaytu.io/pennywise/server/resource"
)

// RegistryName is the fully qualified name under which this provider is stored in the registry.
const RegistryName = "registry.terraform.io/hashicorp/google"

// TerraformProviderInitializer is a terraform.ProviderInitializer that initializes the default GCP provider.
var TerraformProviderInitializer = resource.ProviderInitializer{
	MatchNames: []string{ProviderName, RegistryName},
	Provider: func(values map[string]interface{}) (resource.Provider, error) {
		z, ok := values["zone"]
		if !ok {
			return nil, nil
		}
		region := zoneToRegion(z.(string))
		return googletf.NewProvider(ProviderName, region)
	},
}
