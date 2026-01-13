package commands

import (
	"github.com/accretional/accretional-cli/internal/commands/search"
	"github.com/spf13/cobra"
)

// NewSearchCmd creates the search command with subcommands
func NewSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search resources (records)",
		Long:  "Search resources in Collector using full-text search (FTS), vector similarity, or hybrid search.",
	}

	cmd.AddCommand(search.NewRecordsCmd())

	return cmd
}
