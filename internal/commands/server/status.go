package server

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	statusEndpoint  string
	statusNamespace string
)

func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get detailed server status and information",
		Long:  "Display comprehensive information about the Collector server",
		RunE:  runStatus,
	}

	cmd.Flags().StringVar(&statusEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&statusNamespace, "namespace", "", "Namespace to query (empty for all namespaces)")

	return cmd
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: statusEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	ctx := cmd.Context()
	repoClient := pb.NewCollectionRepoClient(cl.Conn())

	// Get collections
	discoverReq := &pb.DiscoverRequest{
		Namespace: statusNamespace,
		PageSize:  1000, // Get as many as possible
	}

	discoverResp, err := repoClient.Discover(ctx, discoverReq)
	if err != nil {
		return fmt.Errorf("failed to discover collections: %w", err)
	}

	collections := discoverResp.Collections
	collectionCount := len(collections)

	// Display status
	cmd.Printf("Server Status\n")
	cmd.Printf("=============\n")
	cmd.Printf("Endpoint: %s\n", statusEndpoint)

	if statusNamespace != "" {
		cmd.Printf("Namespace: %s\n", statusNamespace)
	} else {
		cmd.Printf("Namespace: (all namespaces)\n")
	}

	cmd.Printf("\nCollections:\n")
	cmd.Printf("  Total: %d\n", collectionCount)

	if collectionCount > 0 {
		cmd.Printf("\nCollection List:\n")
		for i, coll := range collections {
			cmd.Printf("  [%d] %s/%s\n", i+1, coll.Namespace, coll.Name)
		}
		cmd.Printf("\nNote: Use 'collection describe' to get detailed statistics\n")
		cmd.Printf("      (record counts, storage size) for individual collections.\n")
	} else {
		cmd.Printf("\nNo collections found.\n")
	}

	return nil
}
