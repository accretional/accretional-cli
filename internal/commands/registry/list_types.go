package registry

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	listTypesEndpoint  string
	listTypesNamespace string
)

func NewListTypesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-types",
		Short: "List registered message types",
		Long:  "List all registered protobuf message types in a namespace",
		RunE:  runListTypes,
	}

	cmd.Flags().StringVar(&listTypesEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&listTypesNamespace, "namespace", "", "Namespace (empty for all namespaces)")

	return cmd
}

func runListTypes(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: listTypesEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	ctx := cmd.Context()

	// Types are stored in system/types collection
	// We need to query that collection to list types
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())

	// Query system/types collection
	// Note: This is a simplified approach - we list records from system/types
	// A more complete implementation would parse the RegisteredProto records
	listReq := &pb.ListRequest{
		Namespace:      "system",
		CollectionName: "types",
		PageSize:       1000,
	}

	listResp, err := collectionClient.List(ctx, listReq)
	if err != nil {
		return fmt.Errorf("list types failed: %w", err)
	}

	if listResp.Status != nil && listResp.Status.Code != pb.Status_OK {
		return fmt.Errorf("list types failed: %s (code: %d)", listResp.Status.Message, listResp.Status.Code)
	}

	// Display types
	if listTypesNamespace != "" {
		cmd.Printf("Types in namespace '%s':\n", listTypesNamespace)
	} else {
		cmd.Printf("Types (all namespaces):\n")
	}
	cmd.Printf("Total: %d\n\n", len(listResp.Items))

	if len(listResp.Items) == 0 {
		cmd.Println("No types registered.")
		return nil
	}

	cmd.Println("Registered type records:")
	for i, item := range listResp.Items {
		// Unmarshal *any.Any to RegisteredProto
		regProto := &pb.RegisteredProto{}
		if err := anypb.UnmarshalTo(item, regProto, proto.UnmarshalOptions{}); err != nil {
			cmd.Printf("  [%d] Failed to unmarshal type: %v\n", i+1, err)
			continue
		}

		cmd.Printf("  [%d] %s\n", i+1, regProto.Id)
		cmd.Printf("       Namespace: %s\n", regProto.Namespace)
		if len(regProto.MessageNames) > 0 {
			cmd.Printf("       Messages (%d): %v\n", len(regProto.MessageNames), regProto.MessageNames)
		}
		if len(regProto.Dependencies) > 0 {
			cmd.Printf("       Dependencies: %v\n", regProto.Dependencies)
		}
		cmd.Println()
	}

	return nil
}
