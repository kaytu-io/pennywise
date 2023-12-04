package hcl

import (
	"encoding/json"
	"github.com/kaytu-io/infracost/external/config"
	"github.com/kaytu-io/infracost/external/providers/terraform"
	"golang.org/x/net/context"
)

func ParseHclResources(path string) (string, []Resource, error) {
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

	return provider, resources, nil
}
