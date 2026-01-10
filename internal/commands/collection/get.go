package collection

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
	getEndpoint   string
	getNamespace  string
	getCollection string
	getID         string
	getOutput     string
)

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a record from a collection",
		Long:  "Retrieve a record by ID from a collection",
		RunE:  runGet,
	}

	cmd.Flags().StringVar(&getEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&getNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&getCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&getID, "id", "", "Record ID (required)")
	cmd.Flags().StringVar(&getOutput, "output", "", "Output file (optional, prints to stdout if not provided)")

	cmd.MarkFlagRequired("collection")
	cmd.MarkFlagRequired("id")

	return cmd
}

func runGet(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: getEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Create request
	req := &pb.GetRequest{
		Namespace:      getNamespace,
		CollectionName: getCollection,
		Id:             getID,
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
	if getOutput != "" {
		if err := os.WriteFile(getOutput, jsonData, 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		cmd.Printf("✓ Record retrieved and saved to %s\n", getOutput)
	} else {
		cmd.Println(string(jsonData))
	}

	return nil
}
