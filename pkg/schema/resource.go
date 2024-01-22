package schema

// ResourceDef is a single resource definition.
type ResourceDef struct {
	Address      string                 `json:"address"`
	Type         string                 `json:"type"`
	Name         string                 `json:"name"`
	RegionCode   string                 `json:"region_code"`
	ProviderName ProviderName           `json:"provider_name"`
	Values       map[string]interface{} `json:"values"`
}
