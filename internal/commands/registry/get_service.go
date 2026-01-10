package registry

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	getServiceEndpoint  string
	getServiceNamespace string
	getServiceName       string
)

func NewGetServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-service",
		Short: "Get details about a registered service",
		Long:  "Get detailed information about a registered gRPC service",
		RunE:  runGetService,
	}

	cmd.Flags().StringVar(&getServiceEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&getServiceNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&getServiceName, "name", "", "Service name (required)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func runGetService(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: getServiceEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	ctx := cmd.Context()
	registryClient := pb.NewCollectorRegistryClient(cl.Conn())

	// Lookup service
	req := &pb.LookupServiceRequest{
		Namespace:   getServiceNamespace,
		ServiceName: getServiceName,
	}

	resp, err := registryClient.LookupService(ctx, req)
	if err != nil {
		return fmt.Errorf("lookup service failed: %w", err)
	}

	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("lookup service failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	if resp.Service == nil {
		cmd.Printf("Service '%s' not found in namespace '%s'\n", getServiceName, getServiceNamespace)
		return nil
	}

	svc := resp.Service

	// Display service details
	cmd.Printf("Service: %s/%s\n", svc.Namespace, svc.ServiceName)
	cmd.Printf("ID: %s\n", svc.Id)

	if len(svc.MethodNames) > 0 {
		cmd.Printf("\nMethods (%d):\n", len(svc.MethodNames))
		for i, method := range svc.MethodNames {
			cmd.Printf("  [%d] %s\n", i+1, method)
		}
	}

	if svc.ServiceDescriptor != nil {
		if svc.ServiceDescriptor.Name != nil {
			cmd.Printf("\nService Descriptor:\n")
			cmd.Printf("  Name: %s\n", *svc.ServiceDescriptor.Name)
		}
		if len(svc.ServiceDescriptor.Method) > 0 {
			cmd.Printf("  Methods in descriptor: %d\n", len(svc.ServiceDescriptor.Method))
		}
	}

	return nil
}
