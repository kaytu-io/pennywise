package hcl

import (
	"encoding/json"
	"github.com/kaytu-io/infracost/external/config"
	"github.com/kaytu-io/infracost/external/providers/terraform"
	"github.com/kaytu-io/pennywise/pkg/schema"
	usagePackage "github.com/kaytu-io/pennywise/pkg/usage"
	"golang.org/x/net/context"
	"path/filepath"
)

func ParseTerragruntProject(path string, usage usagePackage.Usage) ([]ParsedProject, error) {
	runCtx, err := config.NewRunContextFromEnv(context.Background())
	if err != nil {
		return nil, err
	}
	ctx := config.ProjectContext{
		ProjectConfig: &config.Project{
			Path: path,
		},
		RunContext: runCtx,
	}
	tProvider := terraform.NewTerragruntHCLProvider(&ctx, false)
	dirs, err := tProvider.PrepWorkingDirs()
	if err != nil {
		return nil, err
	}
	var parsedProjects []ParsedProject
	for _, dir := range dirs {
		var rootModule Module
		var provider schema.ProviderName
		var defaultRegion string
		jsons := dir.Provider.LoadPlanJSONs()
		for _, j := range jsons {
			var res Project
			err := json.Unmarshal(j.JSON, &res)
			if err != nil {
				return nil, err
			}
			for key, providerConfig := range res.Configuration.ProviderConfig {
				provider = key
				defaultRegion = providerConfig.Expressions.Region.ConstantValue
			}
			for _, mod := range res.PlannedValues {
				rootModule = mod
			}
		}

		addUsageToModule(usage, &rootModule)
		currentDir, err := filepath.Abs(".")
		if err != nil {
			return nil, err
		}

		relativePath, err := filepath.Rel(currentDir, dir.ConfigDir)
		if err != nil {
			return nil, err
		}

		parsedProjects = append(parsedProjects, ParsedProject{
			Directory:     relativePath,
			Provider:      provider,
			DefaultRegion: defaultRegion,
			RootModule:    rootModule,
		})
	}
	return parsedProjects, nil
}
