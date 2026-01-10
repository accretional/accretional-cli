package registry

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	listServicesEndpoint  string
	listServicesNamespace string
)

func NewListServicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-services",
		Short: "List registered services",
		Long:  "List all registered gRPC services in a namespace",
		RunE:  runListServices,
	}

	cmd.Flags().StringVar(&listServicesEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&listServicesNamespace, "namespace", "", "Namespace (empty for all namespaces)")

	return cmd
}

func runListServices(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: listServicesEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	ctx := cmd.Context()
	registryClient := pb.NewCollectorRegistryClient(cl.Conn())

	// List services
	req := &pb.ListServicesRequest{
		Namespace: listServicesNamespace,
	}

	resp, err := registryClient.ListServices(ctx, req)
	if err != nil {
		return fmt.Errorf("list services failed: %w", err)
	}

	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("list services failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	// Display services
	if listServicesNamespace != "" {
		cmd.Printf("Services in namespace '%s':\n", listServicesNamespace)
	} else {
		cmd.Printf("Services (all namespaces):\n")
	}
	cmd.Printf("Total: %d\n\n", len(resp.Services))

	if len(resp.Services) == 0 {
		cmd.Println("No services registered.")
		return nil
	}

	for i, svc := range resp.Services {
		cmd.Printf("[%d] %s/%s\n", i+1, svc.Namespace, svc.ServiceName)
		if len(svc.MethodNames) > 0 {
			cmd.Printf("     Methods: %v\n", svc.MethodNames)
		}
		if svc.Id != "" {
			cmd.Printf("     ID: %s\n", svc.Id)
		}
		cmd.Println()
	}

	return nil
}
