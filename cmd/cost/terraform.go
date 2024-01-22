package cost

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/pennywise-server/client"
	"github.com/kaytu-io/pennywise-server/schema"
	"github.com/kaytu-io/pennywise/cmd/cost/terraform"
	"github.com/kaytu-io/pennywise/cmd/flags"
	"github.com/kaytu-io/pennywise/parser/hcl"
	"github.com/kaytu-io/pennywise/submission"
	usagePackage "github.com/kaytu-io/pennywise/usage"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

var (
	ServerClientAddress = os.Getenv("SERVER_CLIENT_URL")
)

var terraformCommand = &cobra.Command{
	Use:   "terraform",
	Short: `Shows the costs by parsing terraform resources.`,
	Long:  `Shows the costs by parsing terraform resources.`,
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
		projectDirectory := flags.ReadStringOptionalFlag(cmd, "project")
		if projectDirectory != nil {
			paths, err := os.ReadDir(*projectDirectory)
			if err != nil {
				return err
			}
			for _, p := range paths {
				if !p.IsDir() && strings.HasSuffix(p.Name(), ".tf") {
					err := estimateTfProject(*projectDirectory, usage)
					if err != nil {
						return err
					}
					return nil
				}
			}
		}

		jsonPath := flags.ReadStringOptionalFlag(cmd, "json-path")
		if jsonPath != nil {
			err := estimateTfPlanJson(*jsonPath, usage)
			if err != nil {
				return err
			}
			return nil
		}
		fmt.Println("Please provide a terraform project or a terraform plan json file")
		return nil
	},
}

func estimateTfProject(projectDir string, usage usagePackage.Usage) error {
	provider, hclResources, err := hcl.ParseHclResources(projectDir, usage)
	if err != nil {
		return err
	}
	var resources []schema.ResourceDef
	for _, res := range hclResources {
		resources = append(resources, res.ToResourceDef(provider, nil))
	}
	sub, err := submission.CreateSubmission(resources)
	if err != nil {
		return err
	}
	err = sub.StoreAsFile()
	if err != nil {
		return err
	}
	serverClient := client.NewPennywiseServerClient(ServerClientAddress)
	cost, err := serverClient.GetStateCost(*sub)
	if err != nil {
		return err
	}
	costString, err := cost.CostString()
	if err != nil {
		return err
	}
	fmt.Println(costString)
	return nil
}

func estimateTfPlanJson(jsonPath string, usage usagePackage.Usage) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		return err
	}
	stateCost, err := terraform.EstimateTerraformPlanJson(file, usage)
	if err != nil {
		return err
	}
	costString, err := stateCost.CostString()
	if err != nil {
		return err
	}
	fmt.Println(costString)
	return nil
}
