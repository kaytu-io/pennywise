package diff

import "github.com/spf13/cobra"

// DiffCmd diff commands
var DiffCmd = &cobra.Command{
	Use:   "diff",
	Short: `Shows the diff between two submissions.`,
	Long:  `Shows the diff between two submissions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	DiffCmd.AddCommand(projectCommand)
	projectCommand.Flags().String("json-path", "", "terraform plan json file path")
	projectCommand.Flags().String("project-path", ".", "path to terraform project")
	projectCommand.Flags().String("usage", "", "usage file path")
	projectCommand.Flags().Bool("classic", false, "Show results in classic view (not interactive)")
	projectCommand.Flags().String("compare-to", "", "submission id to compare other submission with (latest submission by default)")

	DiffCmd.AddCommand(submissionCommand)
	submissionCommand.Flags().String("submission-id", "", "submission id")
	submissionCommand.MarkFlagRequired("submission-id")
	submissionCommand.Flags().String("compare-to", "", "submission id to compare other submission with")
	submissionCommand.MarkFlagRequired("compare-to")

	submissionCommand.Flags().Bool("classic", false, "Show results in classic view (not interactive)")
}
