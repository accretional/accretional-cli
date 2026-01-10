package collection

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

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
	searchVectorStr           string
	searchVectorFile          string
	searchSimilarityThreshold float32
	searchLimit               int32
)

func NewSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search records in a collection",
		Long:  "Search records using full-text search (FTS), vector similarity, or hybrid search",
		RunE:  runSearch,
	}

	cmd.Flags().StringVar(&searchEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&searchNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&searchCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&searchQuery, "query", "", "Full-text search query (optional, use with --vector for hybrid search)")
	cmd.Flags().StringVar(&searchVectorStr, "vector", "", "Comma or space-separated vector values for vector search")
	cmd.Flags().StringVar(&searchVectorFile, "vector-file", "", "JSON file containing vector array for vector search")
	cmd.Flags().Float32Var(&searchSimilarityThreshold, "similarity-threshold", 0.0, "Minimum similarity threshold for vector search (0.0-1.0)")
	cmd.Flags().Int32Var(&searchLimit, "limit", 10, "Maximum number of results")

	cmd.MarkFlagRequired("collection")

	return cmd
}

func runSearch(cmd *cobra.Command, args []string) error {
	// Validate that at least one search method is provided
	if searchQuery == "" && searchVectorFile == "" && searchVectorStr == "" {
		return fmt.Errorf("either --query (FTS), --vector, or --vector-file must be provided")
	}

	// Parse vector from string, file, or use empty
	var vector []float32
	if searchVectorFile != "" {
		data, err := os.ReadFile(searchVectorFile)
		if err != nil {
			return fmt.Errorf("failed to read vector file: %w", err)
		}

		// Try to parse as JSON array
		var vectorArray []float64
		if err := json.Unmarshal(data, &vectorArray); err != nil {
			// Try parsing as space/comma-separated values
			content := strings.TrimSpace(string(data))
			parts := strings.FieldsFunc(content, func(r rune) bool {
				return r == ',' || r == ' ' || r == '\n' || r == '\t'
			})
			vectorArray = make([]float64, 0, len(parts))
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				val, err := strconv.ParseFloat(part, 32)
				if err != nil {
					return fmt.Errorf("invalid vector value '%s': %w", part, err)
				}
				vectorArray = append(vectorArray, val)
			}
		}

		// Convert to float32
		vector = make([]float32, len(vectorArray))
		for i, v := range vectorArray {
			vector[i] = float32(v)
		}
	} else if searchVectorStr != "" {
		// Parse comma or space-separated values
		parts := strings.FieldsFunc(searchVectorStr, func(r rune) bool {
			return r == ',' || r == ' '
		})
		vector = make([]float32, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			val, err := strconv.ParseFloat(part, 32)
			if err != nil {
				return fmt.Errorf("invalid vector value '%s': %w", part, err)
			}
			vector = append(vector, float32(val))
		}
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

	// Add vector search if provided
	if len(vector) > 0 {
		req.Vector = vector
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
	if searchQuery != "" && len(vector) > 0 {
		searchType = "hybrid (FTS + vector)"
	} else if len(vector) > 0 {
		searchType = "vector"
	} else {
		searchType = "full-text"
	}

	cmd.Printf("Found %d result(s) (%s search)\n", len(resp.Results), searchType)
	if searchQuery != "" {
		cmd.Printf("Query: '%s'\n", searchQuery)
	}
	if len(vector) > 0 {
		cmd.Printf("Vector dimension: %d\n", len(vector))
		if searchSimilarityThreshold > 0 {
			cmd.Printf("Similarity threshold: %.3f\n", searchSimilarityThreshold)
		}
	}
	cmd.Println()

	for i, result := range resp.Results {
		cmd.Printf("[%d]", i+1)
		if result.Score > 0 {
			cmd.Printf(" Relevance Score: %.4f", result.Score)
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
