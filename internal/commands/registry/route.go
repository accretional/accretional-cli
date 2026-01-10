package registry

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	routeEndpoint  string
	routeNamespace string
	routeService   string
	routeMethod    string
)

func NewRouteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route",
		Short: "Get routing information for a service method",
		Long:  "Get routing information showing which collector handles a specific service method",
		RunE:  runRoute,
	}

	cmd.Flags().StringVar(&routeEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&routeNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&routeService, "service", "", "Service name (required)")
	cmd.Flags().StringVar(&routeMethod, "method", "", "Method name (optional)")

	cmd.MarkFlagRequired("service")

	return cmd
}

func runRoute(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: routeEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	ctx := cmd.Context()
	repoClient := pb.NewCollectionRepoClient(cl.Conn())

	// Use Route to get routing information
	// Note: RouteRequest uses NamespacedName for collection
	req := &pb.RouteRequest{
		Collection: &pb.NamespacedName{
			Namespace: routeNamespace,
			Name:      routeService,
		},
	}

	resp, err := repoClient.Route(ctx, req)
	if err != nil {
		return fmt.Errorf("route failed: %w", err)
	}

	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("route failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	// Display routing information
	cmd.Printf("Routing Information\n")
	cmd.Printf("===================\n")
	cmd.Printf("Service: %s/%s\n", routeNamespace, routeService)

	if routeMethod != "" {
		cmd.Printf("Method: %s\n", routeMethod)
	}

	if resp.ServerEndpoint != "" {
		cmd.Printf("Server Endpoint: %s\n", resp.ServerEndpoint)
	}

	if resp.Collection != nil {
		cmd.Printf("Collection: %s/%s\n", resp.Collection.Namespace, resp.Collection.Name)
	}

	cmd.Println("\nNote: This shows which collector/server handles this service.")
	cmd.Println("      For method-specific routing, use the CollectiveDispatcher.Serve method.")

	return nil
}
