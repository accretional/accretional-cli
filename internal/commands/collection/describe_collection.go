package collection

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	describeCollectionEndpoint  string
	describeCollectionNamespace string
	describeCollectionName      string
)

func NewDescribeCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe a collection",
		Long:  "Get detailed information about a collection including record count and storage size",
		RunE:  runDescribeCollection,
	}

	cmd.Flags().StringVar(&describeCollectionEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&describeCollectionNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&describeCollectionName, "name", "", "Collection name (required)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func runDescribeCollection(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: describeCollectionEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Create request
	req := &pb.DescribeRequest{
		Namespace:      describeCollectionNamespace,
		CollectionName: describeCollectionName,
	}

	// Call Describe
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	resp, err := collectionClient.Describe(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("describe failed: %w", err)
	}

	// Check status
	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("describe failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	// Display collection information
	cmd.Printf("Collection: %s/%s\n", describeCollectionNamespace, describeCollectionName)
	if resp.CollectionDefinition != nil {
		if resp.CollectionDefinition.MessageType != nil {
			cmd.Printf("  Message Type: %s/%s\n",
				resp.CollectionDefinition.MessageType.Namespace,
				resp.CollectionDefinition.MessageType.MessageName)
		}
		if len(resp.CollectionDefinition.IndexedFields) > 0 {
			cmd.Printf("  Indexed Fields: %v\n", resp.CollectionDefinition.IndexedFields)
		}
		if resp.CollectionDefinition.CollectionConfig != nil &&
			resp.CollectionDefinition.CollectionConfig.SearchConfig != nil {
			cfg := resp.CollectionDefinition.CollectionConfig.SearchConfig
			cmd.Printf("  Search Config:\n")
			cmd.Printf("    FTS: %v\n", cfg.EnableFts)
			cmd.Printf("    JSON: %v\n", cfg.EnableJson)
			cmd.Printf("    Vector: %v\n", cfg.EnableVector)
			if cfg.EnableVector {
				cmd.Printf("    Vector Dimensions: %d\n", cfg.VectorDimensions)
			}
		}
	}
	cmd.Printf("  Record Count: %d\n", resp.RecordCount)
	cmd.Printf("  Storage Size: %d bytes\n", resp.StorageSizeBytes)

	return nil
}
