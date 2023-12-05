package cost

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/pennywise/cli/cmd/cost/terraform"
	"github.com/kaytu-io/pennywise/cli/cmd/flags"
	"github.com/kaytu-io/pennywise/cli/parser/hcl"
	usagePackage "github.com/kaytu-io/pennywise/cli/usage"
	"github.com/kaytu-io/pennywise/server/client"
	"github.com/spf13/cobra"
	"log"
	"os"
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
			err = json.NewDecoder(usageFile).Decode(&usage)
			if err != nil {
				return fmt.Errorf("error while parsing usage file %s", err)
			}
		} else {
			usage = usagePackage.Default
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
	provider, resources, err := hcl.ParseHclResources(projectDir, usage)
	if err != nil {
		return err
	}
	for _, rs := range resources {
		res := rs.ToResource(provider, nil)
		serverClient := client.NewPennywiseServerClient(ServerClientAddress)
		cost, err := serverClient.GetCost(res)
		if err != nil {
			return err
		}
		fmt.Println(cost.CostString())
	}
	return nil
}

func estimateTfPlanJson(jsonPath string, usage usagePackage.Usage) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	err = terraform.EstimateTerraformPlanJson(file, usage)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return nil
}
