package registry

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	getTypeEndpoint  string
	getTypeNamespace string
	getTypeName      string
)

func NewGetTypeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-type",
		Short: "Get details about a registered type",
		Long:  "Get detailed information about a registered protobuf message type",
		RunE:  runGetType,
	}

	cmd.Flags().StringVar(&getTypeEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&getTypeNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&getTypeName, "name", "", "Type name (required)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func runGetType(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: getTypeEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	ctx := cmd.Context()

	// Types are stored in system/types collection
	// The record ID format is: namespace/type_name
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())

	recordID := fmt.Sprintf("%s/%s", getTypeNamespace, getTypeName)

	getReq := &pb.GetRequest{
		Namespace:      "system",
		CollectionName: "types",
		Id:             recordID,
	}

	getResp, err := collectionClient.Get(ctx, getReq)
	if err != nil {
		return fmt.Errorf("get type failed: %w", err)
	}

	if getResp.Status != nil && getResp.Status.Code != pb.Status_OK {
		return fmt.Errorf("get type failed: %s (code: %d)", getResp.Status.Message, getResp.Status.Code)
	}

	cmd.Printf("Type: %s/%s\n", getTypeNamespace, getTypeName)
	cmd.Printf("Record ID: %s\n", recordID)
	cmd.Println("\nNote: Type details are stored as RegisteredProto in system/types collection.")
	cmd.Println("      Use 'collection get-record --collection types --id' to view full details.")

	return nil
}
