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
	updateRecordEndpoint   string
	updateRecordNamespace  string
	updateRecordCollection string
	updateRecordID         string
	updateRecordData       string
	updateRecordFile       string
	updateRecordMask       []string
)

func NewUpdateRecordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-record",
		Short: "Update a record in a collection",
		Long:  "Update an existing record in a collection with new JSON data",
		RunE:  runUpdate,
	}

	cmd.Flags().StringVar(&updateRecordEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&updateRecordNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&updateRecordCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&updateRecordID, "id", "", "Record ID (required)")
	cmd.Flags().StringVar(&updateRecordData, "data", "", "JSON data for the record")
	cmd.Flags().StringVar(&updateRecordFile, "file", "", "JSON file containing record data")
	cmd.Flags().StringArrayVar(&updateRecordMask, "update-mask", []string{}, "Field paths to update (optional, updates all fields if not specified)")

	cmd.MarkFlagRequired("collection")
	cmd.MarkFlagRequired("id")

	return cmd
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// Get data from either --data or --file
	var jsonData string
	if updateRecordFile != "" {
		data, err := os.ReadFile(updateRecordFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		jsonData = string(data)
	} else if updateRecordData != "" {
		jsonData = updateRecordData
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
		Endpoint: updateRecordEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Create request
	req := &pb.UpdateRequest{
		Namespace:      updateRecordNamespace,
		CollectionName: updateRecordCollection,
		Id:             updateRecordID,
		Item:           anyValue,
	}

	if len(updateRecordMask) > 0 {
		req.UpdateMask = updateRecordMask
	}

	// Call Update
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	resp, err := collectionClient.Update(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	// Check status
	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("update failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	cmd.Printf("✓ Record updated successfully\n")
	cmd.Printf("  ID: %s\n", updateRecordID)
	cmd.Printf("  Collection: %s/%s\n", updateRecordNamespace, updateRecordCollection)

	return nil
}
