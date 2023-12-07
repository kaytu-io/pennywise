package cmd

import (
	"github.com/kaytu-io/pennywise/cli/cmd/flags"
	"github.com/kaytu-io/pennywise/server/client"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	ServerClientAddress = os.Getenv("SERVER_CLIENT_URL")
)

var ingest = &cobra.Command{
	Use:   "ingest",
	Short: `Store pricing data in the server database.`,
	Long:  `Store pricing data in the server database for the specified provider and resource type and region.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := flags.ReadStringFlag(cmd, "provider")
		service := flags.ReadStringFlag(cmd, "service")
		region := flags.ReadStringFlag(cmd, "region")

		serverClient := client.NewPennywiseServerClient(ServerClientAddress)
		if strings.ToLower(provider) == "aws" {
			err := serverClient.IngestAws(service, region)
			if err != nil {
				return err
			}
		} else if strings.ToLower(provider) == "azure" {
			err := serverClient.IngestAzure(service, region)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	ingest.Flags().String("provider", "", "cloud provider (aws | azure)")
	ingest.MarkFlagRequired("provider")
	ingest.Flags().String("service", "", "service")
	ingest.MarkFlagRequired("service")
	ingest.Flags().String("region", "", "region")
}
