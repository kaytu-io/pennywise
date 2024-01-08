package cost

import "github.com/spf13/cobra"

// CostCmd cost commands
var CostCmd = &cobra.Command{
	Use:   "cost",
	Short: `Shows the costs for the resources with the defined usages.`,
	Long:  `Breaks down the costs for the resources with the defined usages within the next month.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	CostCmd.AddCommand(terraformCommand)
	terraformCommand.Flags().String("json-path", "", "terraform plan json file path")
	terraformCommand.Flags().String("project", "", "terraform project directory (recommended)")
	terraformCommand.Flags().String("usage", "", "usage file path")
}
