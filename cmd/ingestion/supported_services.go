package ingestion

import (
	"fmt"
	"github.com/kaytu-io/pennywise/cmd/flags"
	"github.com/spf13/cobra"
)

var supportedService = &cobra.Command{
	Use:   "supported-services",
	Short: `Shows the list of services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := flags.ReadStringFlag(cmd, "provider")
		if provider == "aws" {
			fmt.Printf("Available services in AWS: \nAmazonEC2\nAmazonEFS\nAmazonEKS\nAmazonFSx\nAmazonRDS\nAmazonElastiCache\nAmazonCloudWatch\nAWSELB\n")
			return nil
		} else if provider == "azure" {
			fmt.Printf("Available services in Azure: \nKey Vault\nVirtual Machines\nStorage\nContainer Registry\nAzure DNS\nLoad Balancer\n" +
				"Application Gateway\nNAT Gateway\nVPN Gateway\nContent Delivery Network\nVirtual Network\nAzure Cosmos DB\n" +
				"Azure Database for MariaDB\nAzure Database for MySQL\nAzure Database for PostgreSQL\nSQL Database\nSQL Managed Instance\n" +
				"Azure Kubernetes Service\nFunctions\nAzure App Service\nAPI Management\n")
			return nil
		} else {
			fmt.Errorf("please enter right provider")
			return nil
		}
	},
}
