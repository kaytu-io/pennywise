package cost

import (
	"fmt"
	"github.com/kaytu-io/pennywise/cli/cmd/cost/terraform"
	"github.com/kaytu-io/pennywise/cli/cmd/flags"
	"github.com/kaytu-io/pennywise/cli/usage"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var terraformCommand = &cobra.Command{
	Use:   "terraform",
	Short: `Store pricing data in the server database.`,
	Long:  `Store pricing data in the server database for the specified provider and resource type and region.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonPath := flags.ReadStringFlag(cmd, "json-path")
		file, err := os.Open(jsonPath)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		usage := usage.Default
		err = terraform.EstimateTerraformPlan(file, usage)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		return nil
	},
}
