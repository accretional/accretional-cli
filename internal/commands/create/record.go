package create

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
	recordEndpoint   string
	recordNamespace  string
	recordCollection string
	recordID         string
	recordData       string
	recordFile       string
)

func NewRecordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "Create a record in a collection",
		Long:  "Create a new record in a collection with JSON data",
		RunE:  runCreateRecord,
	}

	cmd.Flags().StringVar(&recordEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&recordNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&recordCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&recordID, "id", "", "Record ID (optional, auto-generated if not provided)")
	cmd.Flags().StringVar(&recordData, "data", "", "JSON data for the record")
	cmd.Flags().StringVar(&recordFile, "file", "", "JSON file containing record data")

	cmd.MarkFlagRequired("collection")

	return cmd
}

func runCreateRecord(cmd *cobra.Command, args []string) error {
	// Get data from either --data or --file
	var jsonData string
	if recordFile != "" {
		data, err := os.ReadFile(recordFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		jsonData = string(data)
	} else if recordData != "" {
		jsonData = recordData
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
		Endpoint: recordEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Create request
	req := &pb.CreateRequest{
		Namespace:      recordNamespace,
		CollectionName: recordCollection,
		Item:           anyValue,
	}

	if recordID != "" {
		req.Id = recordID
	}

	// Call Create
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	resp, err := collectionClient.Create(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("create failed: %w", err)
	}

	cmd.Printf("✓ Record created successfully\n")
	cmd.Printf("  ID: %s\n", resp.Id)
	cmd.Printf("  Collection: %s/%s\n", recordNamespace, recordCollection)

	return nil
}
