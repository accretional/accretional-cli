package commands

import (
	"github.com/accretional/accretional-cli/internal/commands/list"
	"github.com/spf13/cobra"
)

// NewListCmd creates the list command with subcommands
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List resources (records)",
		Long:  "List resources from Collector.",
	}

	cmd.AddCommand(list.NewRecordsCmd())

	return cmd
}
