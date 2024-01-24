package ingestion

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/pennywise/cmd/flags"
	"github.com/kaytu-io/pennywise/pkg/server"
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

		serverClient := server.NewPennywiseServerClient(flags.ReadStringFlag(cmd, "server-url"))
		jobs, err := serverClient.ListIngestionJobs(provider, service, region, status)
		if err != nil {
			return err
		}

		jobJSON, err := json.MarshalIndent(jobs, "", "    ")
		if err != nil {
			return err
		}

		fmt.Println(string(jobJSON))
		return nil
	},
}
