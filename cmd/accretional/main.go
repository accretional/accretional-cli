package main

import (
	"github.com/accretional/accretional-cli/internal/commands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "accretional",
		Short: "Accretional CLI for Collector",
		Long:  "Command-line interface for interacting with Collector services",
	}

	// Add collection subcommand
	rootCmd.AddCommand(commands.NewCollectionCmd())

	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrln(err)
	}
}
