package hcl

import (
	"fmt"
	"github.com/kaytu-io/pennywise-server/schema"
	"github.com/kaytu-io/pennywise/parser/azurerm"
	usagePackage "github.com/kaytu-io/pennywise/usage"
	"strings"
)

// Resource represents a resource in the terraform project
type Resource struct {
	Address string                 `mapstructure:"address"`
	Name    string                 `mapstructure:"name"`
	Type    string                 `mapstructure:"type"`
	Values  map[string]interface{} `mapstructure:"values"`
}

// extractResourcesFromMapStructure extracts the resources from a terraform project map structure
func extractResourcesFromMapStructure(mapStructure map[string]interface{}) (string, []Resource, error) {
	var provider string
	var resources []Resource

	for key, value := range mapStructure {
		labels := strings.Split(key, ".")
		if labels[0] == "provider" {
			provider = labels[1]
		} else if labels[0] == "resource" {
			values, err := value.(map[string]interface{})
			if !err {
				return "", nil, fmt.Errorf("resource %s value is not a map", key)
			}
			resources = append(resources, Resource{
				Address: strings.Join(labels[1:], "."),
				Name:    strings.Split(strings.Join(labels[2:], "."), "[")[0],
				Type:    labels[1],
				Values:  values,
			})
		}
	}
	return provider, resources, nil
}

// ToResourceDef convert Resource to an acceptable type for pennywise server
func (r *Resource) ToResourceDef(provider schema.ProviderName, defaultRegion *string) schema.ResourceDef {
	region := ""
	if defaultRegion != nil {
		region = *defaultRegion
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

// addUsage add usage data to resource
func (r *Resource) addUsage(usage usagePackage.Usage) {
	newValues := r.Values

	newValues[usagePackage.Key] = usage.GetUsage(r.Type, r.Address)
	r.Values = newValues
}
