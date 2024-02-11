package cost

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/infracost/external/providers"
	"github.com/kaytu-io/pennywise/cmd/cost/terraform"
	"github.com/kaytu-io/pennywise/cmd/flags"
	"github.com/kaytu-io/pennywise/pkg/output"
	"github.com/kaytu-io/pennywise/pkg/parser/hcl"
	"github.com/kaytu-io/pennywise/pkg/schema"
	"github.com/kaytu-io/pennywise/pkg/server"
	"github.com/kaytu-io/pennywise/pkg/submission"
	usagePackage "github.com/kaytu-io/pennywise/pkg/usage"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

const (
	DefaultServerUrl = "https://pennywise.kaytu.dev/kaytu"
)

// CostCmd cost commands
var CostCmd = &cobra.Command{
	Use:   "cost",
	Short: `Shows the costs for the resources with the defined usages.`,
	Long:  `Breaks down the costs for the resources with the defined usages within the next month.`,
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

		//ServerClientAddress := DefaultServerUrl
		//if os.Getenv("SERVER_CLIENT_URL") != "" {
		//	ServerClientAddress = os.Getenv("SERVER_CLIENT_URL")
		//}

		classic := flags.ReadBooleanFlag(cmd, "classic")

		jsonPath := flags.ReadStringOptionalFlag(cmd, "json-path")
		projectPath := flags.ReadStringFlag(cmd, "project-path")
		if jsonPath != nil {
			err := estimateTfPlanJson(classic, *jsonPath, usage, DefaultServerUrl)
			if err != nil {
				return err
			}
		} else {
			err := estimateTerraformProject(classic, projectPath, usage, DefaultServerUrl)
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
	sub, err := submission.CreateSubmission(resources)
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
		err = output.ShowStateCosts(state)
		if err != nil {
			return err
		}
	}
	return nil
}

func estimateTerraformProject(classic bool, projectPath string, usage usagePackage.Usage, ServerClientAddress string) error {
	var providerName schema.ProviderName
	var defaultRegion string
	var parserResources []hcl.Resource
	var err error
	if providers.IsTerragruntNestedDir(projectPath, 5) {
		fmt.Println("terragrunt project...")
		providerName, defaultRegion, parserResources, err = hcl.ParseTerragruntProject(projectPath, usage)
	} else {
		providerName, defaultRegion, parserResources, err = hcl.ParseHclResources(projectPath, usage)
	}
	if err != nil {
		return err
	}
	var resources []schema.ResourceDef
	for _, r := range parserResources {
		resources = append(resources, r.ToResource(providerName, defaultRegion))
	}
	serverClient, err := server.NewPennywiseServerClient(ServerClientAddress)
	if err != nil {
		return err
	}
	sub, err := submission.CreateSubmission(resources)
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
		err = output.ShowStateCosts(state)
		if err != nil {
			return err
		}
	}
	return nil
}
