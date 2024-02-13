package hcl

import (
	"github.com/kaytu-io/pennywise/pkg/parser/azurerm"
	"github.com/kaytu-io/pennywise/pkg/schema"
	"github.com/kaytu-io/pennywise/pkg/usage"
)

type ParsedProject struct {
	Directory     string
	Provider      schema.ProviderName
	DefaultRegion string
	RootModule    Module
}

type Resource struct {
	Address string                 `mapstructure:"address"`
	Mode    string                 `mapstructure:"mode"`
	Name    string                 `mapstructure:"name"`
	Type    string                 `mapstructure:"type"`
	Values  map[string]interface{} `mapstructure:"values"`
}

func addUsageToModule(usage usage.Usage, module *Module) {
	for i, res := range module.Resources {
		module.Resources[i] = addUsage(res, usage)
	}
	for _, childModule := range module.ChildModules {
		addUsageToModule(usage, &childModule)
	}
}

func (pp ParsedProject) GetModule() schema.ModuleDef {
	return pp.buildModuleDef(pp.RootModule)
}

func (pp ParsedProject) buildModuleDef(module Module) schema.ModuleDef {
	var moduleDef = schema.ModuleDef{
		Address: module.Address,
	}
	for _, childModule := range module.ChildModules {
		moduleDef.ChildModules = append(moduleDef.ChildModules, pp.buildModuleDef(childModule))
	}
	for _, resource := range module.Resources {
		moduleDef.Resources = append(moduleDef.Resources, resource.ToResource(pp.Provider, pp.DefaultRegion))
	}
	return moduleDef
}

func (pp ParsedProject) GetResources() []schema.ResourceDef {
	var resources []schema.ResourceDef
	resources = append(resources, pp.getModuleResources(pp.RootModule)...)

	return resources
}

func (pp ParsedProject) getModuleResources(module Module) []schema.ResourceDef {
	var resources []schema.ResourceDef
	for _, res := range module.Resources {
		resources = append(resources, res.ToResource(pp.Provider, pp.DefaultRegion))
	}
	for _, childModule := range module.ChildModules {
		resources = append(resources, pp.getModuleResources(childModule)...)
	}
	return resources
}

func (r Resource) ToResource(provider schema.ProviderName, defaultRegion string) schema.ResourceDef {
	var region string
	if defaultRegion != "" {
		region = defaultRegion
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
	Name        schema.ProviderName `json:"name"`
	Expressions struct {
		Region struct {
			ConstantValue string `json:"constant_value"`
		} `json:"region"`
	} `json:"expressions"`
}

type Config struct {
	ProviderConfig map[schema.ProviderName]ProviderConfig `json:"provider_config"`
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
	Address      string     `json:"address"`
	Resources    []Resource `json:"resources"`
	ChildModules []Module   `json:"child_modules"'`
}
