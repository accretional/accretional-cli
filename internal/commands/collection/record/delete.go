package record

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	deleteRecordEndpoint   string
	deleteRecordNamespace  string
	deleteRecordCollection string
	deleteRecordID         string
)

func NewDeleteRecordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-record",
		Short: "Delete a record from a collection",
		Long:  "Delete a record by ID from a collection",
		RunE:  runDelete,
	}

	cmd.Flags().StringVar(&deleteRecordEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&deleteRecordNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&deleteRecordCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&deleteRecordID, "id", "", "Record ID (required)")

	cmd.MarkFlagRequired("collection")
	cmd.MarkFlagRequired("id")

	return cmd
}

func runDelete(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: deleteRecordEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Create request
	req := &pb.DeleteRequest{
		Namespace:      deleteRecordNamespace,
		CollectionName: deleteRecordCollection,
		Id:             deleteRecordID,
	}

	// Call Delete
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	resp, err := collectionClient.Delete(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	// Check status
	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("delete failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	cmd.Printf("✓ Record deleted successfully\n")
	cmd.Printf("  ID: %s\n", deleteRecordID)
	cmd.Printf("  Collection: %s/%s\n", deleteRecordNamespace, deleteRecordCollection)

	return nil
}
