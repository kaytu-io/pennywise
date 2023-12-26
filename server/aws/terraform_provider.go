package aws

const (
	// RegistryName is the fully qualified name under which this provider is stored in the registry.
	RegistryName = "registry.terraform.io/hashicorp/aws"

	// DefaultRegion is the region used by default when none is defined on the provider
	DefaultRegion = "us-east-1"
)

// TerraformProviderInitializer is a terraform.ProviderInitializer that initializes the default AWS provider.
//var TerraformProviderInitializer = resource.ProviderInitializer{
//	MatchNames: []string{ProviderName, RegistryName},
//	Provider: func(values map[string]interface{}) (resource.Provider, error) {
//		logger, err := zap.NewProduction()
//		if err != nil {
//			return nil, fmt.Errorf("Error : %v ", err)
//		}
//
//		r, ok := values["region"]
//		// If no region is defined it means it was passed via ENV variables
//		// and it's not tracked on the Plan or HCL so we'll assume the
//		// region to be the DefaultRegion
//		if !ok {
//			r = DefaultRegion
//		}
//		regCode := region.Code(r.(string))
//		return aws.NewProvider(ProviderName, regCode, logger)
//	},
//}
