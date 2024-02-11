package hcl

import (
	"encoding/json"
	"github.com/kaytu-io/infracost/external/config"
	"github.com/kaytu-io/infracost/external/providers/terraform"
	"github.com/kaytu-io/pennywise/pkg/schema"
	usagePackage "github.com/kaytu-io/pennywise/pkg/usage"
	"golang.org/x/net/context"
)

func ParseHclResources(path string, usage usagePackage.Usage) (schema.ProviderName, string, []Resource, error) {
	var resources []Resource
	runCtx, err := config.NewRunContextFromEnv(context.Background())
	if err != nil {
		return "", "", nil, err
	}
	ctx := config.ProjectContext{
		ProjectConfig: &config.Project{
			Path: path,
		},
		RunContext: runCtx,
	}
	h, providerErr := terraform.NewHCLProvider(
		&ctx,
		nil,
	)
	if providerErr != nil {
		return "", "", nil, providerErr
	}
	var provider schema.ProviderName
	var defaultRegion string
	jsons := h.LoadPlanJSONs()
	for _, j := range jsons {
		var res Project
		err := json.Unmarshal(j.JSON, &res)
		if err != nil {
			return "", "", nil, err
		}
		for key, providerConfig := range res.Configuration.ProviderConfig {
			provider = key
			defaultRegion = providerConfig.Expressions.Region.ConstantValue
		}
		for _, mod := range res.PlannedValues {
			resources = append(resources, mod.Resources...)
			for _, childMod := range mod.ChildModules {
				resources = append(resources, childMod.Resources...)
			}
		}
	}

	for i, res := range resources {
		resources[i] = addUsage(res, usage)
	}

	return provider, defaultRegion, resources, nil
}

func addUsage(res Resource, usage usagePackage.Usage) Resource {
	newValues := res.Values

	newValues[usagePackage.Key] = usage.GetUsage(res.Type, res.Address)
	return Resource{
		Address: res.Address,
		Mode:    res.Mode,
		Name:    res.Name,
		Type:    res.Type,
		Values:  newValues,
	}
}
