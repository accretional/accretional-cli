package commands

import (
	"github.com/accretional/accretional-cli/internal/commands/collection"
	"github.com/spf13/cobra"
)

// NewCollectionCmd creates the collection command with subcommands
func NewCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection",
		Short: "Collection operations (create, get, list, search)",
		Long:  "Manage records in Collector collections",
	}

	cmd.AddCommand(collection.NewCreateCmd())
	cmd.AddCommand(collection.NewGetCmd())
	cmd.AddCommand(collection.NewListCmd())
	cmd.AddCommand(collection.NewSearchCmd())

	return cmd
}
