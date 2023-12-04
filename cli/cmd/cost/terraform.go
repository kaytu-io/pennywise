package cost

import (
	"fmt"
	"github.com/kaytu-io/pennywise/cli/cmd/cost/terraform"
	"github.com/kaytu-io/pennywise/cli/cmd/flags"
	"github.com/kaytu-io/pennywise/cli/parser/hcl"
	"github.com/kaytu-io/pennywise/cli/usage"
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
		projectDirectory := flags.ReadStringFlag(cmd, "project")
		if projectDirectory != "" {
			paths, err := os.ReadDir(projectDirectory)
			if err != nil {
				return err
			}
			for _, p := range paths {
				if !p.IsDir() && strings.HasSuffix(p.Name(), ".tf") {
					err := estimateTfProject(projectDirectory)
					if err != nil {
						return err
					}
					return nil
				}
			}
		}

		jsonPath := flags.ReadStringFlag(cmd, "json-path")
		if jsonPath != "" {
			err := estimateTfPlanJson(jsonPath)
			if err != nil {
				return err
			}
			return nil
		}
		fmt.Println("Please provide a terraform project or a terraform plan json file")
		return nil
	},
}

func estimateTfProject(projectDir string) error {
	provider, resources, err := hcl.ParseHclResources(projectDir)
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
		fmt.Println(rs.Address, ":", cost)
	}
	return nil
}

func estimateTfPlanJson(jsonPath string) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	usage := usage.Default
	err = terraform.EstimateTerraformPlanJson(file, usage)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return nil
}
