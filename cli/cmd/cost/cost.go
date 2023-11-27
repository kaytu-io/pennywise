package cost

import "github.com/spf13/cobra"

// CostCmd cost commands
var CostCmd = &cobra.Command{
	Use: "cost",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	CostCmd.AddCommand(terraformCommand)
	terraformCommand.Flags().String("json-path", "", "terraform plan json file path")
	terraformCommand.MarkFlagRequired("json-path")
}
