package commands

import (
	"github.com/accretional/accretional-cli/internal/commands/registry"
	"github.com/spf13/cobra"
)

// NewRegistryCmd creates the registry command with subcommands
func NewRegistryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "Service registry and discovery",
		Long: `Manage gRPC service registration and discovery for service mesh functionality.

Commands:
  upload-proto      - Upload and register proto file definitions
  list-types        - List registered message types
  get-type          - Get details about a registered type
  register-service  - Register a gRPC service definition
  list-services     - List registered services
  get-service       - Get details about a registered service
  add-connection    - Add a service connection/endpoint
  list-connections  - List service connections
  remove-connection - Remove a service connection
  discover          - Discover services via dispatcher
  route             - Get routing information for a service method`,
	}

	// Proto and type commands
	cmd.AddCommand(registry.NewUploadProtoCmd())
	cmd.AddCommand(registry.NewListTypesCmd())
	cmd.AddCommand(registry.NewGetTypeCmd())

	// Service registration commands
	cmd.AddCommand(registry.NewRegisterServiceCmd())
	cmd.AddCommand(registry.NewListServicesCmd())
	cmd.AddCommand(registry.NewGetServiceCmd())

	// Connection management commands
	cmd.AddCommand(registry.NewAddConnectionCmd())
	cmd.AddCommand(registry.NewListConnectionsCmd())
	cmd.AddCommand(registry.NewRemoveConnectionCmd())

	// Discovery commands
	cmd.AddCommand(registry.NewDiscoverCmd())
	cmd.AddCommand(registry.NewRouteCmd())

	return cmd
}
