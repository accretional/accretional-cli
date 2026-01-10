package collection

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	deleteFileEndpoint   string
	deleteFileNamespace  string
	deleteFileCollection string
	deleteFilePath       string
)

func NewDeleteFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-file",
		Short: "Delete a file from a collection",
		Long:  "Delete a standalone file stored in a collection",
		RunE:  runDeleteFile,
	}

	cmd.Flags().StringVar(&deleteFileEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&deleteFileNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&deleteFileCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&deleteFilePath, "path", "", "File path in collection (required)")

	cmd.MarkFlagRequired("collection")
	cmd.MarkFlagRequired("path")

	return cmd
}

func runDeleteFile(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: deleteFileEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Generate record ID from path
	recordID := fmt.Sprintf("_file:%s", deleteFilePath)

	// Create request
	req := &pb.DeleteRequest{
		Namespace:      deleteFileNamespace,
		CollectionName: deleteFileCollection,
		Id:             recordID,
	}

	// Call Delete
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	resp, err := collectionClient.Delete(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("delete file failed: %w", err)
	}

	// Check status
	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("delete file failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	cmd.Printf("✓ File deleted successfully\n")
	cmd.Printf("  Path: %s\n", deleteFilePath)
	cmd.Printf("  Collection: %s/%s\n", deleteFileNamespace, deleteFileCollection)

	return nil
}
