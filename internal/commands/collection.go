package commands

import (
	"github.com/accretional/accretional-cli/internal/commands/collection"
	"github.com/spf13/cobra"
)

// NewCollectionCmd creates the collection command with subcommands
func NewCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection",
		Short: "Collection and record operations",
		Long: `Manage collections and records in Collector.

Collection Operations:
  create      - Create a new collection
  list        - List collections in a namespace
  describe    - Get detailed information about a collection
  modify      - Modify collection settings (indexed fields)
  delete      - Delete a collection

Record Operations:
  create-record - Create a new record in a collection
  get-record    - Get a record by ID
  list-records  - List records in a collection
  update-record - Update an existing record
  delete-record - Delete a record
  search        - Search records (FTS, semantic, or hybrid)`,
	}

	// Collection commands
	cmd.AddCommand(collection.NewCreateCollectionCmd())
	cmd.AddCommand(collection.NewListCollectionsCmd())
	cmd.AddCommand(collection.NewDescribeCollectionCmd())
	cmd.AddCommand(collection.NewModifyCollectionCmd())
	cmd.AddCommand(collection.NewDeleteCollectionCmd())

	// Record commands
	cmd.AddCommand(collection.NewCreateRecordCmd())
	cmd.AddCommand(collection.NewGetRecordCmd())
	cmd.AddCommand(collection.NewListRecordsCmd())
	cmd.AddCommand(collection.NewUpdateRecordCmd())
	cmd.AddCommand(collection.NewDeleteRecordCmd())
	cmd.AddCommand(collection.NewSearchRecordsCmd())

	return cmd
}
