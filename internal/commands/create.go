package commands

import (
	"github.com/accretional/accretional-cli/internal/commands/create"
	"github.com/spf13/cobra"
)

// NewCreateCmd creates the create command with subcommands
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create resources (collection, record)",
		Long:  "Create collections or records in Collector.",
	}

	cmd.AddCommand(create.NewCollectionCmd())
	cmd.AddCommand(create.NewRecordCmd())

	return cmd
}
