package ingestion

import (
	"fmt"
	"github.com/kaytu-io/pennywise/cmd/flags"
	"github.com/kaytu-io/pennywise/pkg/server"
	"github.com/spf13/cobra"
)

var supportedService = &cobra.Command{
	Use:   "supported_services",
	Short: `Returns a list of the available services`,
	Long:  `Returns the services that are available for the add ingestion command`,
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := flags.ReadStringFlag(cmd, "provider")
		if provider == "aws" {

			serverClient := server.NewPennywiseServerClient(flags.ReadStringFlag(cmd, "server-url"))
			listNewServices, err := serverClient.ListServices(provider)
			if err != nil {
				return err
			}

			fmt.Println("Available services in AWS: ")
			for _, v := range listNewServices {
				fmt.Println(v)
			}

			return nil
		} else if provider == "azure" {

			serverClient := server.NewPennywiseServerClient(flags.ReadStringFlag(cmd, "server-url"))
			listNewServices, err := serverClient.ListServices(provider)
			if err != nil {
				return err
			}

			fmt.Println("Available services in Azure: ")
			for _, v := range listNewServices {
				fmt.Println(v)
			}

			return nil
		} else {
			return fmt.Errorf("please enter right provider")
		}
	},
}
