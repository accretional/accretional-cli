package collection

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	searchEndpoint   string
	searchNamespace  string
	searchCollection string
	searchQuery      string
	searchLimit      int32
)

func NewSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search records in a collection",
		Long:  "Full-text search records in a collection",
		RunE:  runSearch,
	}

	cmd.Flags().StringVar(&searchEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&searchNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&searchCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&searchQuery, "query", "", "Search query (required)")
	cmd.Flags().Int32Var(&searchLimit, "limit", 10, "Maximum number of results")

	cmd.MarkFlagRequired("collection")
	cmd.MarkFlagRequired("query")

	return cmd
}

func runSearch(cmd *cobra.Command, args []string) error {
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
		FullText:       searchQuery,
		Limit:          searchLimit,
	}

	// Call Search
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	resp, err := collectionClient.Search(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("search failed: %s", resp.Status.Message)
	}

	// Print results
	cmd.Printf("Found %d result(s) for '%s'\n", resp.TotalCount, searchQuery)
	cmd.Println()

	for i, result := range resp.Results {
		cmd.Printf("[%d] Score: %.2f\n", i+1, result.Score)
		if result.Item != nil {
			cmd.Printf("  Type: %s\n", result.Item.TypeUrl)
		}
		cmd.Println()
	}

	return nil
}
