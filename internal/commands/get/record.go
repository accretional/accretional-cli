package get

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	recordEndpoint   string
	recordNamespace  string
	recordCollection string
	recordID         string
	recordOutput     string
)

func NewRecordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "Get a record from a collection",
		Long:  "Retrieve a record by ID from a collection",
		RunE:  runGetRecord,
	}

	cmd.Flags().StringVar(&recordEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&recordNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&recordCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&recordID, "id", "", "Record ID (required)")
	cmd.Flags().StringVar(&recordOutput, "output", "", "Output file (optional, prints to stdout if not provided)")

	cmd.MarkFlagRequired("collection")
	cmd.MarkFlagRequired("id")

	return cmd
}

func runGetRecord(cmd *cobra.Command, args []string) error {
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
	req := &pb.GetRequest{
		Namespace:      recordNamespace,
		CollectionName: recordCollection,
		Id:             recordID,
	}

	// Call Get
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	resp, err := collectionClient.Get(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("get failed: %w", err)
	}

	// Convert Any to JSON
	var jsonData []byte
	if resp.Item != nil {
		// Try to unmarshal the raw bytes as a Struct (ignoring TypeUrl)
		var structVal structpb.Struct
		if err := proto.Unmarshal(resp.Item.Value, &structVal); err == nil {
			// Convert Struct to JSON
			jsonData, err = json.MarshalIndent(structVal.AsMap(), "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal to JSON: %w", err)
			}
		} else {
			// If not a Struct, just print the raw value info
			jsonData = []byte(fmt.Sprintf("{\"type\": \"%s\", \"raw_length\": %d}", resp.Item.TypeUrl, len(resp.Item.Value)))
		}
	}

	// Output
	if recordOutput != "" {
		if err := os.WriteFile(recordOutput, jsonData, 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		cmd.Printf("✓ Record retrieved and saved to %s\n", recordOutput)
	} else {
		cmd.Println(string(jsonData))
	}

	return nil
}
