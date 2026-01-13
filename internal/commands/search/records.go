package search

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
	recordsEndpoint            string
	recordsNamespace           string
	recordsCollection          string
	recordsQuery               string
	recordsVectorStr           string
	recordsVectorFile          string
	recordsSimilarityThreshold float32
	recordsLimit               int32
)

func NewRecordsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "records",
		Short: "Search records in a collection",
		Long:  "Search records using full-text search (FTS), vector similarity, or hybrid search",
		RunE:  runSearchRecords,
	}

	cmd.Flags().StringVar(&recordsEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&recordsNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&recordsCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&recordsQuery, "query", "", "Full-text search query (optional, use with --vector for hybrid search)")
	cmd.Flags().StringVar(&recordsVectorStr, "vector", "", "Comma or space-separated vector values for vector search")
	cmd.Flags().StringVar(&recordsVectorFile, "vector-file", "", "JSON file containing vector array for vector search")
	cmd.Flags().Float32Var(&recordsSimilarityThreshold, "similarity-threshold", 0.0, "Minimum similarity threshold for vector search (0.0-1.0)")
	cmd.Flags().Int32Var(&recordsLimit, "limit", 10, "Maximum number of results")

	cmd.MarkFlagRequired("collection")

	return cmd
}

func runSearchRecords(cmd *cobra.Command, args []string) error {
	// Validate that at least one search method is provided
	if recordsQuery == "" && recordsVectorFile == "" && recordsVectorStr == "" {
		return fmt.Errorf("either --query (FTS), --vector, or --vector-file must be provided")
	}

	// Parse vector from string, file, or use empty
	var vector []float32
	if recordsVectorFile != "" {
		data, err := os.ReadFile(recordsVectorFile)
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
	} else if recordsVectorStr != "" {
		// Parse comma or space-separated values
		parts := strings.FieldsFunc(recordsVectorStr, func(r rune) bool {
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
		Endpoint: recordsEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Create request
	req := &pb.SearchRequest{
		Namespace:      recordsNamespace,
		CollectionName: recordsCollection,
		Limit:          recordsLimit,
	}

	// Add full-text search if provided
	if recordsQuery != "" {
		req.FullText = recordsQuery
	}

	// Add vector search if provided
	if len(vector) > 0 {
		req.Vector = vector
		if recordsSimilarityThreshold > 0 {
			req.SimilarityThreshold = recordsSimilarityThreshold
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
	if recordsQuery != "" && len(vector) > 0 {
		searchType = "hybrid (FTS + vector)"
	} else if len(vector) > 0 {
		searchType = "vector"
	} else {
		searchType = "full-text"
	}

	cmd.Printf("Found %d result(s) (%s search)\n", len(resp.Results), searchType)
	if recordsQuery != "" {
		cmd.Printf("Query: '%s'\n", recordsQuery)
	}
	if len(vector) > 0 {
		cmd.Printf("Vector dimension: %d\n", len(vector))
		if recordsSimilarityThreshold > 0 {
			cmd.Printf("Similarity threshold: %.3f\n", recordsSimilarityThreshold)
		}
	}
	cmd.Println()

	for i, result := range resp.Results {
		cmd.Printf("[%d]", i+1)
		if recordsQuery != "" {
			cmd.Printf(" Score: %g", result.Score)
		}
		if result.Distance > 0 {
			cmd.Printf(" Distance: %g", result.Distance)
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
