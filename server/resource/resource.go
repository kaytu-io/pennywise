package resource

// Resource is a single Terraform resource definition.
type Resource struct {
	Address      string                 `json:"address"`
	Type         string                 `json:"type"`
	Name         string                 `json:"name"`
	RegionCode   string                 `json:"region_code"`
	ProviderName string                 `json:"provider_name"`
	Values       map[string]interface{} `json:"values"`
}
