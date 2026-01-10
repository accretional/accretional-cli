package commands

import (
	"github.com/accretional/accretional-cli/internal/commands/collection"
	"github.com/spf13/cobra"
)

// NewCollectionCmd creates the collection command with subcommands
func NewCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection",
		Short: "Collection operations (create, create-record, get-record, list-records, search)",
		Long:  "Manage collections and records in Collector. Use 'create' to create a collection, then 'create-record' to add records.",
	}

	cmd.AddCommand(collection.NewCreateCollectionCmd())
	cmd.AddCommand(collection.NewCreateCmd())
	cmd.AddCommand(collection.NewGetCmd())
	cmd.AddCommand(collection.NewListCmd())
	cmd.AddCommand(collection.NewSearchCmd())

	return cmd
}
