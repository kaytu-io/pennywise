package new_hcl

import (
	"github.com/kaytu-io/pennywise/pkg/parser/azurerm"
	"github.com/kaytu-io/pennywise/pkg/schema"
)

type Resource struct {
	Address string                 `mapstructure:"address"`
	Mode    string                 `mapstructure:"mode"`
	Name    string                 `mapstructure:"name"`
	Type    string                 `mapstructure:"type"`
	Values  map[string]interface{} `mapstructure:"values"`
}

func (r Resource) ToResource(provider schema.ProviderName, defaultRegion ...string) schema.ResourceDef {
	region := ""
	if len(defaultRegion) > 0 {
		region = defaultRegion[0]
	}
	for key, value := range r.Values {
		if provider == schema.AzureProvider && key == "location" {
			region = azurerm.GetRegionCode(value.(string))
			break
		}
	}
	return schema.ResourceDef{
		Address:      r.Address,
		Type:         r.Type,
		Name:         r.Name,
		RegionCode:   region,
		ProviderName: provider,
		Values:       r.Values,
	}
}

type ProviderConfig struct {
	Name schema.ProviderName `json:"name"`
}

type Config struct {
	ProviderConfig map[schema.ProviderName]ProviderConfig `json:"provider_config"`
}

type Project struct {
	TerraformVersion         string                `json:"terraform_version"`
	FormatVersion            string                `json:"format_version"`
	InfracostResourceChanges interface{}           `json:"infracost_resource_changes"`
	Configuration            Config                `json:"configuration"`
	PlannedValues            map[string]RootModule `json:"planned_values"`
	PriorState               interface{}           `json:"prior_state"`
}

type RootModule struct {
	Resources    []Resource    `json:"resources"`
	ChildModules []ChildModule `json:"child_modules"`
}

type ChildModule struct {
	Address   string     `json:"address"`
	Resources []Resource `json:"resources"`
}
