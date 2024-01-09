package hcl

import (
	"github.com/kaytu-io/pennywise-server/resource"
	"github.com/kaytu-io/pennywise/parser/azurerm"
)

type Resource struct {
	Address string                 `mapstructure:"address"`
	Mode    string                 `mapstructure:"mode"`
	Name    string                 `mapstructure:"name"`
	Type    string                 `mapstructure:"type"`
	Values  map[string]interface{} `mapstructure:"values"`
}

func (r Resource) ToResource(provider resource.ProviderName, defaultRegion *string) resource.ResourceDef {
	region := ""
	if defaultRegion != nil {
		region = *defaultRegion
	}
	for key, value := range r.Values {
		if provider == resource.AzureProvider && key == "location" {
			region = azurerm.GetRegionCode(value.(string))
			break
		}
	}
	return resource.ResourceDef{
		Address:      r.Address,
		Type:         r.Type,
		Name:         r.Name,
		RegionCode:   region,
		ProviderName: provider,
		Values:       r.Values,
	}
}

type ProviderConfig struct {
	Name resource.ProviderName `json:"name"`
}

type Config struct {
	ProviderConfig map[resource.ProviderName]ProviderConfig `json:"provider_config"`
}

type Project struct {
	TerraformVersion         string            `json:"terraform_version"`
	FormatVersion            string            `json:"format_version"`
	InfracostResourceChanges interface{}       `json:"infracost_resource_changes"`
	Configuration            Config            `json:"configuration"`
	PlannedValues            map[string]Module `json:"planned_values"`
	PriorState               interface{}       `json:"prior_state"`
}

type Module struct {
	Resources []Resource `mapstructure:"resources"`
}
