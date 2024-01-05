package ingestion

import (
	"fmt"
	"github.com/kaytu-io/pennywise/cli/cmd/flags"
	"github.com/kaytu-io/pennywise/server/client"
	"github.com/spf13/cobra"
)

var add = &cobra.Command{
	Use:   "add",
	Short: `Adds an ingestion job to receive pricing and store in the database`,
	Long:  `Adds an ingestion job to receive pricing and store in the database for the specified provider and resource type and region.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := flags.ReadStringFlag(cmd, "provider")
		service := flags.ReadStringFlag(cmd, "service")
		region := flags.ReadStringFlag(cmd, "region")

		serverClient := client.NewPennywiseServerClient(ServerClientAddress)
		job, err := serverClient.Ingest(provider, service, region)
		if err != nil {
			return err
		}

		fmt.Println(*job)
		return nil
	},
}
