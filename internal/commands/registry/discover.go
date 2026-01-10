package registry

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	discoverEndpoint  string
	discoverNamespace string
	discoverService   string
)

func NewDiscoverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "discover",
		Short: "Discover services via dispatcher",
		Long:  "Discover available services and their locations using the CollectiveDispatcher",
		RunE:  runDiscover,
	}

	cmd.Flags().StringVar(&discoverEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&discoverNamespace, "namespace", "", "Namespace (empty for all namespaces)")
	cmd.Flags().StringVar(&discoverService, "service", "", "Service name to discover (optional)")

	return cmd
}

func runDiscover(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: discoverEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	ctx := cmd.Context()
	repoClient := pb.NewCollectionRepoClient(cl.Conn())

	// Use Discover to find collections (which represent services)
	req := &pb.DiscoverRequest{
		Namespace: discoverNamespace,
		PageSize:  1000,
	}

	resp, err := repoClient.Discover(ctx, req)
	if err != nil {
		return fmt.Errorf("discover failed: %w", err)
	}

	// Display discovered services/collections
	if discoverNamespace != "" {
		cmd.Printf("Discovered in namespace '%s':\n", discoverNamespace)
	} else {
		cmd.Printf("Discovered (all namespaces):\n")
	}

	if discoverService != "" {
		cmd.Printf("Filtering for service: %s\n", discoverService)
	}

	cmd.Printf("Total: %d\n\n", len(resp.Collections))

	if len(resp.Collections) == 0 {
		cmd.Println("No services/collections discovered.")
		return nil
	}

	for i, coll := range resp.Collections {
		if discoverService != "" && coll.Name != discoverService {
			continue
		}
		cmd.Printf("[%d] %s/%s\n", i+1, coll.Namespace, coll.Name)
	}

	cmd.Println("\nNote: Use 'registry route' to get routing information for specific services.")

	return nil
}
