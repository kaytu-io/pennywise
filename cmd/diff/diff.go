package diff

import "github.com/spf13/cobra"

const (
	DefaultServerAddress string = "http://localhost:8080"
)

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
	DiffCmd.AddCommand(submissionCommand)
	submissionCommand.Flags().String("submission-id", "", "submission id")
	submissionCommand.MarkFlagRequired("submission-id")
	submissionCommand.Flags().String("compare-to", "", "submission id to compare other submission with")
	submissionCommand.MarkFlagRequired("compare-to")

	submissionCommand.Flags().Bool("classic", false, "Show results in classic view (not interactive)")
}
