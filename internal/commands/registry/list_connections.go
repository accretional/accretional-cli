package registry

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	listConnectionsEndpoint  string
	listConnectionsNamespace string
)

func NewListConnectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-connections",
		Short: "List service connections",
		Long:  "List all registered service connections/endpoints",
		RunE:  runListConnections,
	}

	cmd.Flags().StringVar(&listConnectionsEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&listConnectionsNamespace, "namespace", "", "Namespace (empty for all namespaces)")

	return cmd
}

func runListConnections(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: listConnectionsEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	ctx := cmd.Context()

	// Connections are stored in system/connections collection
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())

	// Build filter if namespace is specified
	var filter *structpb.Struct
	if listConnectionsNamespace != "" {
		filterMap := map[string]interface{}{
			"shared_namespaces": listConnectionsNamespace,
		}
		var err error
		filter, err = structpb.NewStruct(filterMap)
		if err != nil {
			return fmt.Errorf("failed to create filter: %w", err)
		}
	}

	listReq := &pb.ListRequest{
		Namespace:      "system",
		CollectionName: "connections",
		PageSize:       1000,
		Filter:         filter,
	}

	listResp, err := collectionClient.List(ctx, listReq)
	if err != nil {
		return fmt.Errorf("list connections failed: %w", err)
	}

	if listResp.Status != nil && listResp.Status.Code != pb.Status_OK {
		return fmt.Errorf("list connections failed: %s (code: %d)", listResp.Status.Message, listResp.Status.Code)
	}

	// Display connections
	if listConnectionsNamespace != "" {
		cmd.Printf("Connections in namespace '%s':\n", listConnectionsNamespace)
	} else {
		cmd.Printf("Connections (all namespaces):\n")
	}
	cmd.Printf("Total: %d\n\n", len(listResp.Items))

	if len(listResp.Items) == 0 {
		cmd.Println("No connections registered.")
		return nil
	}

	cmd.Println("Connection records:")
	for i, item := range listResp.Items {
		// Unmarshal *any.Any to Connection
		conn := &pb.Connection{}
		if err := anypb.UnmarshalTo(item, conn, proto.UnmarshalOptions{}); err != nil {
			cmd.Printf("  [%d] Failed to unmarshal connection: %v\n", i+1, err)
			continue
		}

		cmd.Printf("  [%d] %s\n", i+1, conn.Id)
		cmd.Printf("       Address: %s\n", conn.Address)
		cmd.Printf("       Status: %s\n", conn.Status.String())
		if len(conn.SharedNamespaces) > 0 {
			cmd.Printf("       Namespaces: %v\n", conn.SharedNamespaces)
		}
		if conn.TargetCollectorId != "" {
			cmd.Printf("       Target Collector: %s\n", conn.TargetCollectorId)
		}
		cmd.Println()
	}

	return nil
}
