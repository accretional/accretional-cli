package create

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	collectionEndpoint  string
	collectionNamespace string
	collectionName      string
	collectionTypeNS    string
	collectionTypeName  string
	collectionIndexed   []string
)

func NewCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection",
		Short: "Create a new collection",
		Long: `Create a new collection in the specified namespace.`,
		RunE: runCreateCollection,
	}

	cmd.Flags().StringVar(&collectionEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&collectionNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&collectionName, "name", "", "Collection name (required)")
	cmd.Flags().StringVar(&collectionTypeNS, "type-namespace", "", "Message type namespace (optional)")
	cmd.Flags().StringVar(&collectionTypeName, "type-name", "", "Message type name (optional, leave empty for untyped collections)")
	cmd.Flags().StringArrayVar(&collectionIndexed, "indexed-field", []string{}, "Fields to index (can be specified multiple times)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func runCreateCollection(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: collectionEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Build collection
	collection := &pb.Collection{
		Namespace: collectionNamespace,
		Name:      collectionName,
	}

	if collectionTypeName != "" {
		collection.MessageType = &pb.MessageTypeRef{
			Namespace:   collectionTypeNS,
			MessageName: collectionTypeName,
		}
	}

	// Add indexed fields if provided
	if len(collectionIndexed) > 0 {
		collection.IndexedFields = collectionIndexed
	}

	// Create request
	req := &pb.CreateCollectionRequest{
		Collection: collection,
	}

	// Call CreateCollection via CollectionRepo
	repoClient := pb.NewCollectionRepoClient(cl.Conn())
	resp, err := repoClient.CreateCollection(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("create collection failed: %w", err)
	}

	// Check status - accept both 0 (proto enum OK) and 200 (HTTP-style OK)
	if resp.Status != nil && resp.Status.Code != pb.Status_OK && resp.Status.Code != 200 {
		return fmt.Errorf("create collection failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	cmd.Printf("✓ Collection created successfully\n")
	cmd.Printf("  Collection ID: %s\n", resp.CollectionId)
	cmd.Printf("  Namespace: %s\n", collectionNamespace)
	cmd.Printf("  Name: %s\n", collectionName)
	if resp.ServerEndpoint != "" {
		cmd.Printf("  Server Endpoint: %s\n", resp.ServerEndpoint)
	}

	return nil
}
