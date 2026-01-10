package registry

import (
	"fmt"
	"strings"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	addConnectionEndpoint  string
	addConnectionService   string
	addConnectionAddress   string
	addConnectionNamespace string
)

func NewAddConnectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-connection",
		Short: "Add a service connection/endpoint",
		Long:  "Register a connection to a service endpoint for service mesh routing",
		RunE:  runAddConnection,
	}

	cmd.Flags().StringVar(&addConnectionEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&addConnectionService, "service", "", "Service name (optional)")
	cmd.Flags().StringVar(&addConnectionAddress, "address", "", "Service endpoint address (required)")
	cmd.Flags().StringVar(&addConnectionNamespace, "namespace", "shared", "Namespace")

	cmd.MarkFlagRequired("address")

	return cmd
}

func runAddConnection(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: addConnectionEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	ctx := cmd.Context()
	dispatcherClient := pb.NewCollectiveDispatcherClient(cl.Conn())

	// Parse namespaces from flag (comma-separated)
	namespaces := []string{addConnectionNamespace}
	if strings.Contains(addConnectionNamespace, ",") {
		namespaces = strings.Split(addConnectionNamespace, ",")
		for i := range namespaces {
			namespaces[i] = strings.TrimSpace(namespaces[i])
		}
	}

	// Connect to the service endpoint
	req := &pb.ConnectRequest{
		Address:    addConnectionAddress,
		Namespaces: namespaces,
	}

	resp, err := dispatcherClient.Connect(ctx, req)
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}

	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("connect failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	cmd.Printf("✓ Connection added successfully\n")
	cmd.Printf("  Connection ID: %s\n", resp.ConnectionId)
	cmd.Printf("  Address: %s\n", addConnectionAddress)
	cmd.Printf("  Target Collector ID: %s\n", resp.TargetCollectorId)
	if len(resp.SharedNamespaces) > 0 {
		cmd.Printf("  Shared Namespaces: %v\n", resp.SharedNamespaces)
	}

	return nil
}
