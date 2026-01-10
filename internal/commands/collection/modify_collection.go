package collection

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	modifyCollectionEndpoint  string
	modifyCollectionNamespace string
	modifyCollectionName      string
	modifyCollectionIndexed   []string
)

func NewModifyCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify",
		Short: "Modify a collection",
		Long:  "Modify collection settings (currently supports changing indexed fields)",
		RunE:  runModifyCollection,
	}

	cmd.Flags().StringVar(&modifyCollectionEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&modifyCollectionNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&modifyCollectionName, "name", "", "Collection name (required)")
	cmd.Flags().StringArrayVar(&modifyCollectionIndexed, "indexed-field", []string{}, "Fields to index (can be specified multiple times)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func runModifyCollection(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: modifyCollectionEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Create request
	req := &pb.ModifyRequest{
		Namespace:      modifyCollectionNamespace,
		CollectionName: modifyCollectionName,
		IndexedFields:  modifyCollectionIndexed,
	}

	// Call Modify
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	resp, err := collectionClient.Modify(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("modify failed: %w", err)
	}

	// Check status
	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("modify failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	cmd.Printf("✓ Collection modified successfully\n")
	cmd.Printf("  Collection: %s/%s\n", modifyCollectionNamespace, modifyCollectionName)
	if len(modifyCollectionIndexed) > 0 {
		cmd.Printf("  Indexed Fields: %v\n", modifyCollectionIndexed)
	}

	return nil
}
