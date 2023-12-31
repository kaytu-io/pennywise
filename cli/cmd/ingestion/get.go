package ingestion

import (
	"fmt"
	"github.com/kaytu-io/pennywise/cli/cmd/flags"
	"github.com/kaytu-io/pennywise/server/client"
	"github.com/spf13/cobra"
)

var get = &cobra.Command{
	Use:   "get",
	Short: `Returns an ingestion job with the provided id`,
	Long:  `Returns an ingestion job with the provided id`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := flags.ReadStringFlag(cmd, "id")

		serverClient := client.NewPennywiseServerClient(ServerClientAddress)
		job, err := serverClient.GetIngestionJob(id)
		if err != nil {
			return err
		}

		fmt.Println(*job)
		return nil
	},
}
