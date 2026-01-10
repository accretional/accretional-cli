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
	listEndpoint   string
	listNamespace  string
	listCollection string
	listPageSize   int32
	listFilter     string
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List records in a collection",
		Long:  "List records in a collection with optional filtering",
		RunE:  runList,
	}

	cmd.Flags().StringVar(&listEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&listNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&listCollection, "collection", "", "Collection name (required)")
	cmd.Flags().Int32Var(&listPageSize, "page-size", 10, "Number of records per page")
	cmd.Flags().StringVar(&listFilter, "filter", "", "JSON filter (optional)")

	cmd.MarkFlagRequired("collection")

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: listEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Create request
	req := &pb.ListRequest{
		Namespace:      listNamespace,
		CollectionName: listCollection,
		PageSize:       listPageSize,
	}

	// Add filter if provided
	if listFilter != "" {
		var filterMap map[string]interface{}
		if err := json.Unmarshal([]byte(listFilter), &filterMap); err != nil {
			return fmt.Errorf("invalid filter JSON: %w", err)
		}
		filterStruct, err := structpb.NewStruct(filterMap)
		if err != nil {
			return fmt.Errorf("failed to create filter: %w", err)
		}
		req.Filter = filterStruct
	}

	// Call List
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	resp, err := collectionClient.List(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("list failed: %w", err)
	}

	// Print results
	cmd.Printf("Found %d record(s)\n", len(resp.Items))
	cmd.Println()

	for i, item := range resp.Items {
		var jsonData []byte
		var structVal structpb.Struct
		if err := proto.Unmarshal(item.Value, &structVal); err == nil {
			jsonData, _ = json.MarshalIndent(structVal.AsMap(), "", "  ")
		} else {
			jsonData = []byte(fmt.Sprintf("{\"type\": \"%s\"}", item.TypeUrl))
		}

		cmd.Printf("[%d]\n", i+1)
		cmd.Println(string(jsonData))
		cmd.Println()
	}

	if resp.NextPageToken != "" {
		cmd.Printf("(More results available, use --page-token %s)\n", resp.NextPageToken)
	}

	return nil
}
