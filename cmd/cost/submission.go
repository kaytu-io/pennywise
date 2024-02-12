package cost

import (
	"fmt"
	"github.com/kaytu-io/pennywise/cmd/flags"
	outputCost "github.com/kaytu-io/pennywise/pkg/output/cost"
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
		err := estimateSubmission(classic, submissionId, DefaultServerAddress)
		if err != nil {
			return err
		}

		return nil
	},
}

func estimateSubmission(classic bool, submissionId string, ServerClientAddress string) error {
	serverClient, err := server.NewPennywiseServerClient(ServerClientAddress)
	if err != nil {
		return err
	}
	sub, err := schema.ReadSubmissionFile(submissionId)
	if err != nil {
		return err
	}
	state, err := serverClient.GetStateCost(*sub)
	if err != nil {
		return err
	}
	if classic {
		costString, err := state.CostString()
		if err != nil {
			return err
		}
		fmt.Println(costString)
		fmt.Println("To learn how to use usage open:\nhttps://github.com/kaytu-io/pennywise/blob/main/docs/usage.md")
	} else {
		err = outputCost.ShowStateCosts(state)
		if err != nil {
			return err
		}
	}

	return nil
}
