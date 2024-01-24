package ingestion

import (
	"github.com/spf13/cobra"
)

var IngestCmd = &cobra.Command{
	Use:   "ingestion",
	Short: `Store pricing data in the server database.`,
	Long:  `Store pricing data in the server database for the specified provider and resource type and region.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	IngestCmd.AddCommand(supportedService)
	supportedService.Flags().String("provider", "", "cloud provider (aws | azure)")
	supportedService.MarkFlagRequired("provider")

	IngestCmd.AddCommand(add)
	add.Flags().String("provider", "", "cloud provider (aws | azure)")
	add.MarkFlagRequired("provider")
	add.Flags().String("service", "", "service")
	add.MarkFlagRequired("service")
	add.Flags().String("region", "", "region")
	add.Flags().Bool("wait", false, "wait")

	IngestCmd.AddCommand(list)
	list.Flags().String("provider", "", "cloud provider (aws | azure)")
	list.Flags().String("service", "", "service")
	list.Flags().String("status", "", "status")
	list.Flags().String("region", "", "region (all for jobs with not defined region)")

	IngestCmd.AddCommand(get)
	get.Flags().String("id", "", "id")
}
