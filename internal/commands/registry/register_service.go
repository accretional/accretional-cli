package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	registerServiceEndpoint   string
	registerServiceNamespace  string
	registerServiceName       string
	registerServiceMethods    string
	registerServiceProtoFile  string
	registerServiceImportPath string
)

func NewRegisterServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-service",
		Short: "Register a gRPC service definition",
		Long: `Register a gRPC service with the Collector registry.

This command registers service metadata including service name, namespace, and methods.
The service descriptor can be provided via a proto file or built manually from method names.`,
		RunE: runRegisterService,
	}

	cmd.Flags().StringVar(&registerServiceEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&registerServiceNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&registerServiceName, "name", "", "Service name (required)")
	cmd.Flags().StringVar(&registerServiceMethods, "methods", "", "Comma-separated list of method names (required if proto-file not provided)")
	cmd.Flags().StringVar(&registerServiceProtoFile, "proto-file", "", "Path to .proto file containing service definition (optional, preferred)")
	cmd.Flags().StringVar(&registerServiceImportPath, "import-path", "", "Additional import path for proto dependencies")

	cmd.MarkFlagRequired("name")

	return cmd
}

func runRegisterService(cmd *cobra.Command, args []string) error {
	var serviceDesc *descriptorpb.ServiceDescriptorProto
	var fileDesc *descriptorpb.FileDescriptorProto

	// If proto file is provided, parse it to extract service descriptor
	if registerServiceProtoFile != "" {
		// Verify proto file exists
		if _, err := os.Stat(registerServiceProtoFile); err != nil {
			return fmt.Errorf("proto file not found: %w", err)
		}

		// Parse proto file
		absProtoFile, err := filepath.Abs(registerServiceProtoFile)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute path: %w", err)
		}

		protoDir := filepath.Dir(absProtoFile)
		protoBase := filepath.Base(absProtoFile)

		parser := protoparse.Parser{
			ImportPaths: []string{protoDir},
		}

		// Add custom import paths if provided
		if registerServiceImportPath != "" {
			parser.ImportPaths = append(parser.ImportPaths, registerServiceImportPath)
		}

		// ParseFiles expects relative paths when ImportPaths is set
		fileDescs, err := parser.ParseFiles(protoBase)
		if err != nil {
			return fmt.Errorf("failed to parse proto file: %w", err)
		}

		if len(fileDescs) == 0 {
			return fmt.Errorf("no proto files parsed")
		}

		fileDesc = fileDescs[0].AsFileDescriptorProto()

		// Find the service in the parsed file
		var foundService *descriptorpb.ServiceDescriptorProto
		for _, svc := range fileDesc.Service {
			if svc.Name != nil && *svc.Name == registerServiceName {
				foundService = svc
				break
			}
		}

		if foundService == nil {
			return fmt.Errorf("service '%s' not found in proto file", registerServiceName)
		}

		serviceDesc = foundService
	} else {
		// Build service descriptor from method names
		if registerServiceMethods == "" {
			return fmt.Errorf("either --proto-file or --methods must be provided")
		}

		// Parse methods
		methodNames := strings.Split(registerServiceMethods, ",")
		for i := range methodNames {
			methodNames[i] = strings.TrimSpace(methodNames[i])
		}

		// Build minimal ServiceDescriptorProto from method names
		// Note: This creates a minimal descriptor without full type information
		// Using proto file is preferred for complete service registration
		methods := make([]*descriptorpb.MethodDescriptorProto, 0, len(methodNames))
		for _, methodName := range methodNames {
			if methodName == "" {
				continue
			}
			methods = append(methods, &descriptorpb.MethodDescriptorProto{
				Name: &methodName,
				// Input/Output types are not set - proto file parsing is preferred
			})
		}

		serviceDesc = &descriptorpb.ServiceDescriptorProto{
			Name:   &registerServiceName,
			Method: methods,
		}
	}

	// Create client
	cl, err := client.New(client.Config{
		Endpoint: registerServiceEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	ctx := cmd.Context()
	registryClient := pb.NewCollectorRegistryClient(cl.Conn())

	// Register service
	req := &pb.RegisterServiceRequest{
		Namespace:         registerServiceNamespace,
		ServiceDescriptor: serviceDesc,
		FileDescriptor:    fileDesc,
	}

	resp, err := registryClient.RegisterService(ctx, req)
	if err != nil {
		return fmt.Errorf("register service failed: %w", err)
	}

	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("register service failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	cmd.Printf("✓ Service registered successfully\n")
	cmd.Printf("  Service ID: %s\n", resp.ServiceId)
	if len(resp.RegisteredMethods) > 0 {
		cmd.Printf("  Registered Methods (%d):\n", len(resp.RegisteredMethods))
		for i, method := range resp.RegisteredMethods {
			cmd.Printf("    [%d] %s\n", i+1, method)
		}
	} else {
		cmd.Printf("  Registered Methods: (none)\n")
	}

	return nil
}
