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
	CostCmd.AddCommand(projectCommand)
	projectCommand.Flags().String("json-path", "", "terraform plan json file path")
	projectCommand.Flags().String("project-path", ".", "path to terraform project")
	projectCommand.Flags().String("usage", "", "usage file path")
	projectCommand.Flags().Bool("classic", false, "Show results in classic view (not interactive)")

	CostCmd.AddCommand(submissionCommand)
	submissionCommand.Flags().String("submission-id", "", "submission id")
	submissionCommand.MarkFlagRequired("submission-id")
	submissionCommand.Flags().Bool("classic", false, "Show results in classic view (not interactive)")
}
