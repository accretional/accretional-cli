package commands

import (
	"github.com/accretional/accretional-cli/internal/commands/server"
	"github.com/spf13/cobra"
)

// NewServerCmd creates the server command with subcommands
func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Server management and monitoring",
		Long: `Manage and monitor Collector server instances.

Commands:
  health  - Check server health and connectivity
  status  - Get detailed server status and information`,
	}

	cmd.AddCommand(server.NewHealthCmd())
	cmd.AddCommand(server.NewStatusCmd())

	return cmd
}
