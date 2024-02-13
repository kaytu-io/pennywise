package cost

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/infracost/external/providers"
	"github.com/kaytu-io/pennywise/cmd/cost/terraform"
	"github.com/kaytu-io/pennywise/cmd/flags"
	"github.com/kaytu-io/pennywise/pkg"
	"github.com/kaytu-io/pennywise/pkg/cost"
	outputCost "github.com/kaytu-io/pennywise/pkg/output/cost"
	"github.com/kaytu-io/pennywise/pkg/parser/hcl"
	"github.com/kaytu-io/pennywise/pkg/schema"
	"github.com/kaytu-io/pennywise/pkg/server"
	usagePackage "github.com/kaytu-io/pennywise/pkg/usage"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

var projectCommand = &cobra.Command{
	Use:   "project",
	Short: `Shows the costs by parsing a project resources.`,
	Long:  `Shows the costs by parsing a project resources.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		usagePath := flags.ReadStringOptionalFlag(cmd, "usage")
		var usage usagePackage.Usage
		if usagePath != nil {
			usageFile, err := os.Open(*usagePath)
			if err != nil {
				return fmt.Errorf("error while reading usage file %s", err)
			}
			defer usageFile.Close()

			ext := filepath.Ext(*usagePath)
			switch ext {
			case ".json":
				err = json.NewDecoder(usageFile).Decode(&usage)
			case ".yaml", ".yml":
				err = yaml.NewDecoder(usageFile).Decode(&usage)
			default:
				return fmt.Errorf("unsupported file format %s for usage file", ext)
			}
			if err != nil {
				return fmt.Errorf("error while parsing usage file %s", err)
			}

		} else {
			usage = usagePackage.Usage{}
		}

		classic := flags.ReadBooleanFlag(cmd, "classic")

		jsonPath := flags.ReadStringOptionalFlag(cmd, "json-path")
		projectPath := flags.ReadStringFlag(cmd, "project-path")
		tfVarFiles := flags.ReadStringArrayFlag(cmd, "terraform-var-file")
		if jsonPath != nil {
			err := estimateTfPlanJson(classic, *jsonPath, usage, pkg.DefaultServerAddress)
			if err != nil {
				return err
			}
		} else {
			err := estimateTerraformProject(classic, projectPath, usage, pkg.DefaultServerAddress, tfVarFiles)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func estimateTfPlanJson(classic bool, jsonPath string, usage usagePackage.Usage, ServerClientAddress string) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		return err
	}
	resources, err := terraform.ParseTerraformPlanJson(file, usage)
	if err != nil {
		return err
	}
	serverClient, err := server.NewPennywiseServerClient(ServerClientAddress)
	if err != nil {
		return err
	}
	sub, err := schema.CreateSubmission(resources)
	if err != nil {
		return err
	}
	err = sub.StoreAsFile()
	if err != nil {
		return err
	}
	state, err := serverClient.GetStateCost(*sub)
	if err != nil {
		return err
	}
	if classic {
		costString, err := state.CostString()
		if err != nil {
			return err
		}
		fmt.Println(costString)
		fmt.Println("To learn how to use usage open:\nhttps://github.com/kaytu-io/pennywise/blob/main/docs/usage.md")
	} else {
		modularState := cost.ModularState{
			Resources: state.Resources,
		}
		err = outputCost.ShowStateCosts(&modularState)
		if err != nil {
			return err
		}
	}
	return nil
}

func estimateTerraformProject(classic bool, projectPath string, usage usagePackage.Usage, ServerClientAddress string, tfVarFiles []string) error {
	var projects *schema.ModuleDef
	var err error
	if providers.IsTerragruntNestedDir(projectPath, 5) {
		fmt.Println("terragrunt project...")
		projects, err = hcl.ParseTerragruntProject(projectPath, usage)
	} else {
		projects, err = hcl.ParseHclResources(projectPath, usage, tfVarFiles)
	}
	if err != nil {
		return err
	}
	serverClient, err := server.NewPennywiseServerClient(ServerClientAddress)
	if err != nil {
		return err
	}
	sub, err := schema.CreateSubmissionV2(*projects)
	if err != nil {
		return err
	}
	err = sub.StoreAsFile()
	if err != nil {
		return err
	}
	state, err := serverClient.GetStateCostV2(*sub)
	if err != nil {
		return err
	}
	if classic {
		costString, err := state.ToClassicState().CostString()
		if err != nil {
			return err
		}
		fmt.Println(costString)
		fmt.Println("To learn how to use usage open:\nhttps://github.com/kaytu-io/pennywise/blob/main/docs/usage.md")
	} else {
		err = outputCost.ShowStateCosts(state)
		if err != nil {
			return err
		}
	}

	return nil
}
