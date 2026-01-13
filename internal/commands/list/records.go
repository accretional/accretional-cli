package list

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
	recordsEndpoint   string
	recordsNamespace  string
	recordsCollection string
	recordsPageSize   int32
	recordsFilter     string
)

func NewRecordsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "records",
		Short: "List records in a collection",
		Long:  "List records in a collection with optional filtering",
		RunE:  runListRecords,
	}

	cmd.Flags().StringVar(&recordsEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&recordsNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&recordsCollection, "collection", "", "Collection name (required)")
	cmd.Flags().Int32Var(&recordsPageSize, "page-size", 10, "Number of records per page")
	cmd.Flags().StringVar(&recordsFilter, "filter", "", "JSON filter (optional)")

	cmd.MarkFlagRequired("collection")

	return cmd
}

func runListRecords(cmd *cobra.Command, args []string) error {
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
	req := &pb.ListRequest{
		Namespace:      recordsNamespace,
		CollectionName: recordsCollection,
		PageSize:       recordsPageSize,
	}

	// Add filter if provided
	if recordsFilter != "" {
		var filterMap map[string]interface{}
		if err := json.Unmarshal([]byte(recordsFilter), &filterMap); err != nil {
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
