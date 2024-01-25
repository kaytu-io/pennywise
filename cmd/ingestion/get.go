package ingestion

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/pennywise/cmd/flags"
	"github.com/kaytu-io/pennywise/pkg/server"
	"github.com/spf13/cobra"
	"strings"
)

var get = &cobra.Command{
	Use:   "get",
	Short: `Returns an ingestion job with the provided id`,
	Long:  `Returns an ingestion job with the provided id`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := flags.ReadStringFlag(cmd, "id")

		serverClient := server.NewPennywiseServerClient(flags.ReadStringFlag(cmd, "server-url"))
		job, err := serverClient.GetIngestionJob(id)
		if err != nil {
			if strings.Contains(err.Error(), "job ID not found") {
				return fmt.Errorf("job ID not found")
			}
			return err
		}

		jobJSON, err := json.MarshalIndent(*job, "", "    ")
		if err != nil {
			return err
		}

		fmt.Println(string(jobJSON))
		return nil
	},
}
