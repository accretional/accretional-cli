package collection

import (
	"fmt"
	"strings"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

var (
	deleteCollectionEndpoint  string
	deleteCollectionNamespace string
	deleteCollectionName      string
)

func NewDeleteCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a collection",
		Long:  "Delete a collection and all its records",
		RunE:  runDeleteCollection,
	}

	cmd.Flags().StringVar(&deleteCollectionEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&deleteCollectionNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&deleteCollectionName, "name", "", "Collection name (required)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func runDeleteCollection(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: deleteCollectionEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Create request
	req := &pb.DeleteCollectionRequest{
		Collection: &pb.NamespacedName{
			Namespace: deleteCollectionNamespace,
			Name:      deleteCollectionName,
		},
	}

	// Call DeleteCollection
	repoClient := pb.NewCollectionRepoClient(cl.Conn())
	resp, err := repoClient.DeleteCollection(cmd.Context(), req)
	if err != nil {
		// Check if it's an "Unimplemented" error (method not registered)
		if err.Error() != "" && (contains(err.Error(), "not registered") || contains(err.Error(), "Unimplemented")) {
			return fmt.Errorf("delete collection failed: DeleteCollection method is not registered on the server.\n"+
				"Original error: %w", err)
		}
		return fmt.Errorf("delete collection failed: %w", err)
	}

	// Check status
	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("delete collection failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	cmd.Printf("✓ Collection deleted successfully\n")
	cmd.Printf("  Collection: %s/%s\n", deleteCollectionNamespace, deleteCollectionName)
	if resp.BytesFreed > 0 {
		cmd.Printf("  Bytes freed: %d\n", resp.BytesFreed)
	}

	return nil
}
