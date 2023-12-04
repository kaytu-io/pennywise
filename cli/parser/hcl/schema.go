package hcl

import (
	"github.com/kaytu-io/pennywise/cli/parser/azurerm"
	"github.com/kaytu-io/pennywise/server/resource"
)

type Resource struct {
	Address string                 `mapstructure:"address"`
	Mode    string                 `mapstructure:"mode"`
	Name    string                 `mapstructure:"name"`
	Type    string                 `mapstructure:"type"`
	Values  map[string]interface{} `mapstructure:"values"`
}

func (r Resource) ToResource(provider string, defaultRegion *string) resource.Resource {
	region := ""
	if defaultRegion != nil {
		region = *defaultRegion
	}
	for key, value := range r.Values {
		if provider == "azurerm" && key == "location" {
			region = azurerm.GetRegionCode(value.(string))
			break
		}
	}
	return resource.Resource{
		Address:      r.Address,
		Type:         r.Type,
		Name:         r.Name,
		RegionCode:   region,
		ProviderName: provider,
		Values:       r.Values,
	}
}

type ProviderConfig struct {
	Name string `json:"name"`
}

type Config struct {
	ProviderConfig map[string]ProviderConfig `json:"provider_config"`
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
