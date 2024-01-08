package ingestion

import (
	"fmt"
	"github.com/kaytu-io/pennywise/cli/cmd/flags"
	"github.com/kaytu-io/pennywise/server/client"
	"github.com/spf13/cobra"
)

var list = &cobra.Command{
	Use:   "list",
	Short: `Returns list of ingestion jobs`,
	Long:  `Returns list of ingestion jobs by the provided filters`,
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := flags.ReadStringFlag(cmd, "provider")
		service := flags.ReadStringFlag(cmd, "service")
		region := flags.ReadStringFlag(cmd, "region")
		status := flags.ReadStringFlag(cmd, "status")

		serverClient := client.NewPennywiseServerClient(flags.ReadStringFlag(cmd, "server-url"))
		jobs, err := serverClient.ListIngestionJobs(provider, service, region, status)
		if err != nil {
			return err
		}

		for _, job := range jobs {
			fmt.Println(job)
		}
		return nil
	},
}
