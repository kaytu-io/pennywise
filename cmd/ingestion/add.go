package ingestion

import (
	"fmt"
	"github.com/kaytu-io/pennywise/cmd/flags"
	"github.com/kaytu-io/pennywise/pkg/schema"
	"github.com/kaytu-io/pennywise/pkg/server"
	"github.com/spf13/cobra"
	"time"
)

var add = &cobra.Command{
	Use:   "add",
	Short: `Adds an ingestion job to receive pricing and store in the database`,
	Long: `Adds an ingestion job to receive pricing and store in the database for the specified provider and resource type and region.
			The command will returned the ingestion job object.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := flags.ReadStringFlag(cmd, "provider")
		service := flags.ReadStringFlag(cmd, "service")
		region := flags.ReadStringFlag(cmd, "region")
		wait := flags.ReadBooleanFlag(cmd, "wait")
		if provider != "azure" && provider != "aws" {
			return fmt.Errorf("this provider is not supported")
		}
		serverClient, err := server.NewPennywiseServerClient(flags.ReadStringFlag(cmd, "server-url"))
		if err != nil {
			return err
		}
		job, err := serverClient.AddIngestion(provider, service, region)
		if err != nil {
			return err
		}

		fmt.Println(fmt.Sprintf("Job %d is created to ingest service %s in %s region for %s provider", job.ID, job.Service, job.Location, job.Provider))
		if wait {
			fmt.Println("Waiting...")
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
		return nil
	},
}
