package collection

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

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
