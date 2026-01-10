package collection

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	createCollectionEndpoint     string
	createCollectionNamespace    string
	createCollectionName         string
	createCollectionTypeNS       string
	createCollectionTypeName     string
	createCollectionIndexed      []string
	createCollectionEnableFTS    bool
	createCollectionEnableJSON   bool
	createCollectionEnableVector bool
	createCollectionVectorDims   int32
)

func NewCreateCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new collection",
		Long: `Create a new collection in the specified namespace.

For JSON data (using 'create-record' command with --data), you should specify:
  --type-namespace "google.protobuf"
  --type-name "Struct"

This enables JSON conversion for search functionality. Without a message type,
search may not work correctly.`,
		RunE: runCreateCollection,
	}

	cmd.Flags().StringVar(&createCollectionEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&createCollectionNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&createCollectionName, "name", "", "Collection name (required)")
	cmd.Flags().StringVar(&createCollectionTypeNS, "type-namespace", "google.protobuf", "Message type namespace (default: google.protobuf for JSON data)")
	cmd.Flags().StringVar(&createCollectionTypeName, "type-name", "Struct", "Message type name (default: Struct for JSON data)")
	cmd.Flags().StringArrayVar(&createCollectionIndexed, "indexed-field", []string{}, "Fields to index (can be specified multiple times)")

	// Search configuration flags
	cmd.Flags().BoolVar(&createCollectionEnableFTS, "enable-fts", true, "Enable full-text search (FTS5)")
	cmd.Flags().BoolVar(&createCollectionEnableJSON, "enable-json", true, "Enable JSON field extraction and filtering")
	cmd.Flags().BoolVar(&createCollectionEnableVector, "enable-vector", false, "Enable vector/semantic search")
	cmd.Flags().Int32Var(&createCollectionVectorDims, "vector-dimensions", 384, "Vector dimensions for semantic search (required if enable-vector is true)")

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

	// Always set message type (defaults to google.protobuf.Struct for JSON data)
	// This is required for JSON conversion to work, which enables search functionality
	if createCollectionTypeName != "" {
		collection.MessageType = &pb.MessageTypeRef{
			Namespace:   createCollectionTypeNS,
			MessageName: createCollectionTypeName,
		}
	} else {
		// Default to google.protobuf.Struct for JSON data
		collection.MessageType = &pb.MessageTypeRef{
			Namespace:   "google.protobuf",
			MessageName: "Struct",
		}
	}

	// Add indexed fields if provided
	if len(createCollectionIndexed) > 0 {
		collection.IndexedFields = createCollectionIndexed
	}

	// Set SearchConfig if any search options are specified
	if createCollectionEnableFTS || createCollectionEnableJSON || createCollectionEnableVector {
		searchConfig := &pb.SearchConfig{
			EnableFts:        createCollectionEnableFTS,
			EnableJson:       createCollectionEnableJSON,
			EnableVector:     createCollectionEnableVector,
			VectorDimensions: createCollectionVectorDims,
			EmbedderType:     pb.EmbedderType_EMBEDDER_DETERMINISTIC,
		}

		// Validate vector config
		if createCollectionEnableVector && createCollectionVectorDims <= 0 {
			return fmt.Errorf("vector-dimensions must be > 0 when enable-vector is true")
		}

		collection.CollectionConfig = &pb.CollectionConfig{
			SearchConfig: searchConfig,
		}
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

	// Check status - OK = 0, so any non-zero code is an error
	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
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
