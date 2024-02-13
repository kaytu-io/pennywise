package diff

import (
	"fmt"
	"github.com/kaytu-io/pennywise/cmd/flags"
	"github.com/kaytu-io/pennywise/pkg"
	outputDiff "github.com/kaytu-io/pennywise/pkg/output/diff"
	"github.com/kaytu-io/pennywise/pkg/schema"
	"github.com/kaytu-io/pennywise/pkg/server"
	"github.com/spf13/cobra"
)

var submissionCommand = &cobra.Command{
	Use:   "submission",
	Short: `Shows a submission cost.`,
	Long:  `Shows a submission cost.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		classic := flags.ReadBooleanFlag(cmd, "classic")

		submissionId := flags.ReadStringFlag(cmd, "submission-id")
		compareTo := flags.ReadStringFlag(cmd, "compare-to")

		err := submissionsDiff(classic, submissionId, compareTo, pkg.DefaultServerAddress)
		if err != nil {
			return err
		}

		return nil
	},
}

func submissionsDiff(classic bool, submissionId, compareToId string, ServerClientAddress string) error {
	serverClient, err := server.NewPennywiseServerClient(ServerClientAddress)
	if err != nil {
		return err
	}
	sub, err := schema.ReadSubmissionFileV2(submissionId)
	if err != nil {
		return err
	}
	compareTo, err := schema.ReadSubmissionFileV2(compareToId)
	if err != nil {
		return err
	}

	req := schema.SubmissionsDiffV2{
		Current:   *sub,
		CompareTo: *compareTo,
	}
	stateDiff, err := serverClient.GetSubmissionsDiffV2(req)
	if err != nil {
		return err
	}
	if classic {
		return fmt.Errorf("classic view not available for diff")
	} else {
		err = outputDiff.ShowStateCosts(stateDiff)
		if err != nil {
			return err
		}
	}
	return nil
}
