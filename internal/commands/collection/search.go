package collection

import (
	"encoding/json"
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	searchEndpoint            string
	searchNamespace           string
	searchCollection          string
	searchQuery               string
	searchSemanticText        string
	searchSimilarityThreshold float32
	searchLimit               int32
)

func NewSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search records in a collection",
		Long:  "Search records using full-text search (FTS), semantic similarity, or hybrid search",
		RunE:  runSearch,
	}

	cmd.Flags().StringVar(&searchEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&searchNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&searchCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&searchQuery, "query", "", "Full-text search query (optional, use with --semantic-text for hybrid search)")
	cmd.Flags().StringVar(&searchSemanticText, "semantic-text", "", "Text to search by semantic meaning (embeddings-based similarity)")
	cmd.Flags().Float32Var(&searchSimilarityThreshold, "similarity-threshold", 0.0, "Minimum similarity threshold for semantic search (0.0-1.0)")
	cmd.Flags().Int32Var(&searchLimit, "limit", 10, "Maximum number of results")

	cmd.MarkFlagRequired("collection")

	return cmd
}

func runSearch(cmd *cobra.Command, args []string) error {
	// Validate that at least one search method is provided
	if searchQuery == "" && searchSemanticText == "" {
		return fmt.Errorf("either --query (FTS) or --semantic-text must be provided")
	}

	// Create client
	cl, err := client.New(client.Config{
		Endpoint: searchEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Create request
	req := &pb.SearchRequest{
		Namespace:      searchNamespace,
		CollectionName: searchCollection,
		Limit:          searchLimit,
	}

	// Add full-text search if provided
	if searchQuery != "" {
		req.FullText = searchQuery
	}

	// Add semantic search if provided
	if searchSemanticText != "" {
		req.SemanticText = searchSemanticText
		if searchSimilarityThreshold > 0 {
			req.SimilarityThreshold = searchSimilarityThreshold
		}
	}

	// Call Search
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	resp, err := collectionClient.Search(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// Print results
	searchType := "search"
	if searchQuery != "" && searchSemanticText != "" {
		searchType = "hybrid (FTS + semantic)"
	} else if searchSemanticText != "" {
		searchType = "semantic"
	} else {
		searchType = "full-text"
	}

	cmd.Printf("Found %d result(s) (%s search)\n", len(resp.Results), searchType)
	if searchQuery != "" {
		cmd.Printf("Query: '%s'\n", searchQuery)
	}
	if searchSemanticText != "" {
		cmd.Printf("Semantic text: '%s'\n", searchSemanticText)
		if searchSimilarityThreshold > 0 {
			cmd.Printf("Similarity threshold: %.3f\n", searchSimilarityThreshold)
		}
	}
	cmd.Println()

	for i, result := range resp.Results {
		cmd.Printf("[%d]", i+1)
		// Always show score if it's available (even if 0, it indicates relevance)
		if searchQuery != "" {
			cmd.Printf(" Score: %.4f", result.Score)
		}
		if result.Distance > 0 {
			cmd.Printf(" Distance: %.4f", result.Distance)
		}
		cmd.Println()

		if result.Item != nil {
			var structVal structpb.Struct
			if err := proto.Unmarshal(result.Item.Value, &structVal); err == nil {
				jsonData, _ := json.MarshalIndent(structVal.AsMap(), "", "  ")
				cmd.Println(string(jsonData))
			} else {
				cmd.Printf("  Type: %s\n", result.Item.TypeUrl)
			}
		}
		cmd.Println()
	}

	return nil
}
