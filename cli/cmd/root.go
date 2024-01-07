package cmd

import (
	"errors"
	"github.com/kaytu-io/pennywise/cli/cmd/cost"
	"github.com/kaytu-io/pennywise/cli/cmd/ingestion"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "pennywise",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().ParseErrorsWhitelist.UnknownFlags {
			return errors.New("invalid flags")
		}
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(ingestion.IngestCmd)
	rootCmd.AddCommand(cost.CostCmd)
	rootCmd.PersistentFlags().String("server-url", "http://localhost:8080", "define the server http address")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
