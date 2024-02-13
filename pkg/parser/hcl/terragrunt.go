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

func ParseTerragruntProject(path string, usage usagePackage.Usage) (*schema.ModuleDef, error) {
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
	var projectsModule schema.ModuleDef
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
				if _, ok := map[string]bool{
					"aws":     true,
					"azure":   true,
					"azurerm": true,
				}[string(key)]; ok {
					provider = key
					defaultRegion = providerConfig.Expressions.Region.ConstantValue
					break
				}
			}
			for _, mod := range res.PlannedValues {
				rootModule = mod
			}
		}

		addUsageToModule(usage, &rootModule)
		currentDir, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}

		projectName, err := filepath.Rel(currentDir, dir.ConfigDir)
		if err != nil {
			return nil, err
		}

		parsedProject := ParsedProject{
			Directory:     projectName,
			Provider:      provider,
			DefaultRegion: defaultRegion,
			RootModule:    rootModule,
		}
		projectModule := parsedProject.GetModule()
		changeResourcesId(projectName, &projectModule)
		projectsModule.ChildModules = append(projectsModule.ChildModules, schema.ModuleDef{
			Address:      projectName,
			ChildModules: projectModule.ChildModules,
			Resources:    projectModule.Resources,
		})
	}
	return &projectsModule, nil
}

func changeResourcesId(project string, mod *schema.ModuleDef) {
	for i, res := range mod.Resources {
		mod.Resources[i].Address = project + "." + res.Address
	}
	for _, childMod := range mod.ChildModules {
		changeResourcesId(project, &childMod)
	}
}
