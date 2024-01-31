package ingestion

import (
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/pennywise/cmd/flags"
	"github.com/kaytu-io/pennywise/pkg/schema"
	"github.com/kaytu-io/pennywise/pkg/server"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var get = &cobra.Command{
	Use:   "get",
	Short: `Returns an ingestion job with the provided id`,
	Long:  `Returns an ingestion job with the provided id`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id := flags.ReadStringFlag(cmd, "id")
		wait := flags.ReadBooleanFlag(cmd, "wait")

		serverClient, err := server.NewPennywiseServerClient(flags.ReadStringFlag(cmd, "server-url"))
		if err != nil {
			return err
		}
		job, err := serverClient.GetIngestionJob(id)
		if err != nil {
			if strings.Contains(err.Error(), "job ID not found") {
				return fmt.Errorf("job ID not found")
			}
			return err
		}
		if wait && !(job.Status == schema.IngestionJobSucceeded || job.Status == schema.IngestionJobFailed) {
			fmt.Println(fmt.Sprintf("Waiting for job %d to be completed...", job.ID))
			var maxWaiting int
			for maxWaiting < 3600/10 {
				job, err = serverClient.GetIngestionJob(fmt.Sprintf("%d", job.ID))
				if err != nil {
					return err
				}
				if job.Status == schema.IngestionJobSucceeded || job.Status == schema.IngestionJobFailed {
					fmt.Println("Job is finished:", job.Status)
					return nil
				}
				time.Sleep(time.Second * 10)
			}
		}
		jobJSON, err := json.MarshalIndent(*job, "", "    ")
		if err != nil {
			return err
		}

		fmt.Println(string(jobJSON))
		return nil
	},
}
