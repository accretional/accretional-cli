package collection

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	createCollectionEndpoint  string
	createCollectionNamespace string
	createCollectionName      string
	createCollectionTypeNS    string
	createCollectionTypeName  string
	createCollectionIndexed   []string
)

func NewCreateCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-collection",
		Short: "Create a new collection",
		Long: `Create a new collection in the specified namespace.

By default, collections are created without a message type (untyped).
For typed collections, specify --type-namespace and --type-name.`,
		RunE: runCreateCollection,
	}

	cmd.Flags().StringVar(&createCollectionEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&createCollectionNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&createCollectionName, "name", "", "Collection name (required)")
	cmd.Flags().StringVar(&createCollectionTypeNS, "type-namespace", "", "Message type namespace (optional)")
	cmd.Flags().StringVar(&createCollectionTypeName, "type-name", "", "Message type name (optional, leave empty for untyped collections)")
	cmd.Flags().StringArrayVar(&createCollectionIndexed, "indexed-field", []string{}, "Fields to index (can be specified multiple times)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func runCreateCollection(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: createCollectionEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Build collection
	collection := &pb.Collection{
		Namespace: createCollectionNamespace,
		Name:      createCollectionName,
	}

	if createCollectionTypeName != "" {
		collection.MessageType = &pb.MessageTypeRef{
			Namespace:   createCollectionTypeNS,
			MessageName: createCollectionTypeName,
		}
	}

	// Add indexed fields if provided
	if len(createCollectionIndexed) > 0 {
		collection.IndexedFields = createCollectionIndexed
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
	cmd.Printf("  Namespace: %s\n", createCollectionNamespace)
	cmd.Printf("  Name: %s\n", createCollectionName)
	if resp.ServerEndpoint != "" {
		cmd.Printf("  Server Endpoint: %s\n", resp.ServerEndpoint)
	}

	return nil
}
