package main

import (
	"os"

	"github.com/accretional/accretional-cli/internal/commands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "accretional",
		Short: "CLI for interacting with Collector services",
		Long: `accretional is a command-line interface for interacting with Collector services.

Use it to create collections and records, retrieve data, list records, and perform searches.`,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	// Add commands
	rootCmd.AddCommand(commands.NewCreateCmd())
	rootCmd.AddCommand(commands.NewGetCmd())
	rootCmd.AddCommand(commands.NewListCmd())
	rootCmd.AddCommand(commands.NewSearchCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
