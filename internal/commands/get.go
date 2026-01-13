package commands

import (
	"github.com/accretional/accretional-cli/internal/commands/get"
	"github.com/spf13/cobra"
)

// NewGetCmd creates the get command with subcommands
func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resources (record)",
		Long:  "Retrieve resources from Collector.",
	}

	cmd.AddCommand(get.NewRecordCmd())

	return cmd
}
