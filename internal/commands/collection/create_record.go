package collection

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	createEndpoint   string
	createNamespace  string
	createCollection string
	createID         string
	createData       string
	createFile       string
)

func NewCreateRecordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-record",
		Short: "Create a record in a collection",
		Long:  "Create a new record in a collection with JSON data",
		RunE:  runCreate,
	}

	cmd.Flags().StringVar(&createEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&createNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&createCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&createID, "id", "", "Record ID (optional, auto-generated if not provided)")
	cmd.Flags().StringVar(&createData, "data", "", "JSON data for the record")
	cmd.Flags().StringVar(&createFile, "file", "", "JSON file containing record data")

	cmd.MarkFlagRequired("collection")

	return cmd
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Get data from either --data or --file
	var jsonData string
	if createFile != "" {
		data, err := os.ReadFile(createFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		jsonData = string(data)
	} else if createData != "" {
		jsonData = createData
	} else {
		return fmt.Errorf("either --data or --file must be provided")
	}

	// Parse JSON to validate it
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &jsonObj); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Convert to protobuf Struct
	structValue, err := structpb.NewValue(jsonObj)
	if err != nil {
		return fmt.Errorf("failed to convert JSON: %w", err)
	}

	// Create Any from Struct
	anyValue, err := anypb.New(structValue.GetStructValue())
	if err != nil {
		return fmt.Errorf("failed to create Any: %w", err)
	}

	// Create client
	cl, err := client.New(client.Config{
		Endpoint: createEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Create request
	req := &pb.CreateRequest{
		Namespace:      createNamespace,
		CollectionName: createCollection,
		Item:           anyValue,
	}

	if createID != "" {
		req.Id = createID
	}

	// Call Create
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	resp, err := collectionClient.Create(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("create failed: %w", err)
	}

	cmd.Printf("✓ Record created successfully\n")
	cmd.Printf("  ID: %s\n", resp.Id)
	cmd.Printf("  Collection: %s/%s\n", createNamespace, createCollection)

	return nil
}
