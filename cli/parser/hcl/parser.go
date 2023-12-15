package hcl

import (
	"encoding/json"
	"github.com/kaytu-io/infracost/external/config"
	"github.com/kaytu-io/infracost/external/providers/terraform"
	usagePackage "github.com/kaytu-io/pennywise/cli/usage"
	"golang.org/x/net/context"
)

func ParseHclResources(path string, usage usagePackage.Usage) (string, []Resource, error) {
	var resources []Resource
	runCtx, err := config.NewRunContextFromEnv(context.Background())
	if err != nil {
		return "", nil, err
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
		return "", nil, providerErr
	}
	provider := ""
	jsons := h.LoadPlanJSONs()
	for _, j := range jsons {
		var res Project
		err := json.Unmarshal(j.JSON, &res)
		if err != nil {
			return "", nil, err
		}
		for key, _ := range res.Configuration.ProviderConfig {
			provider = key
		}
		for _, mod := range res.PlannedValues {
			resources = append(resources, mod.Resources...)
		}
	}

	for i, res := range resources {
		if makeResource, ok := makeResourceProcesses[res.Type]; ok {
			resources[i] = makeResource.setRefs(resources, res)
			resources[i] = makeResource.runFunctions(res)
		}
		// Embedded references
		if makeResource, ok := makeResourceProcesses[res.Type]; ok {
			resources[i] = makeResource.setRefs(resources, res)
		}
		resources[i] = addUsage(res, usage)
	}

	return provider, resources, nil
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
