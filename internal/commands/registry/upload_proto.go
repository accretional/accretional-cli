package registry

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	uploadProtoEndpoint   string
	uploadProtoNamespace  string
	uploadProtoFile       string
	uploadProtoImportPath string
)

func NewUploadProtoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-proto",
		Short: "Upload and register proto file definitions",
		Long: `Upload a .proto file and register its message types with the Collector registry.

This command parses the proto file and registers all message types defined in it.
Dependencies are automatically handled if they are already registered.`,
		RunE: runUploadProto,
	}

	cmd.Flags().StringVar(&uploadProtoEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&uploadProtoNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&uploadProtoFile, "file", "", "Path to .proto file (required)")
	cmd.Flags().StringVar(&uploadProtoImportPath, "import-path", "", "Additional import path for proto dependencies (can be specified multiple times)")

	cmd.MarkFlagRequired("file")

	return cmd
}

func runUploadProto(cmd *cobra.Command, args []string) error {
	// Verify proto file exists
	if _, err := os.Stat(uploadProtoFile); err != nil {
		return fmt.Errorf("proto file not found: %w", err)
	}

	// Parse proto file
	parser := protoparse.Parser{
		ImportPaths: []string{filepath.Dir(uploadProtoFile)},
	}

	// Add custom import paths if provided
	if uploadProtoImportPath != "" {
		parser.ImportPaths = append(parser.ImportPaths, uploadProtoImportPath)
	}

	fileDescs, err := parser.ParseFiles(uploadProtoFile)
	if err != nil {
		return fmt.Errorf("failed to parse proto file: %w", err)
	}

	if len(fileDescs) == 0 {
		return fmt.Errorf("no proto files parsed")
	}

	// Get FileDescriptorProto
	fileDesc := fileDescs[0]
	fileDescProto := fileDesc.AsFileDescriptorProto()

	// Extract dependencies (imported files)
	var dependencies []*descriptorpb.FileDescriptorProto
	for _, dep := range fileDesc.GetDependencies() {
		dependencies = append(dependencies, dep.AsFileDescriptorProto())
	}

	// Create client
	cl, err := client.New(client.Config{
		Endpoint: uploadProtoEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	ctx := cmd.Context()
	registryClient := pb.NewCollectorRegistryClient(cl.Conn())

	// Register proto
	req := &pb.RegisterProtoRequest{
		Namespace:      uploadProtoNamespace,
		FileDescriptor: fileDescProto,
		Dependencies:   dependencies,
	}

	resp, err := registryClient.RegisterProto(ctx, req)
	if err != nil {
		return fmt.Errorf("register proto failed: %w", err)
	}

	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("register proto failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	cmd.Printf("✓ Proto file registered successfully\n")
	cmd.Printf("  Proto ID: %s\n", resp.ProtoId)
	if len(resp.RegisteredMessages) > 0 {
		cmd.Printf("  Registered Messages (%d):\n", len(resp.RegisteredMessages))
		for i, msg := range resp.RegisteredMessages {
			cmd.Printf("    [%d] %s\n", i+1, msg)
		}
	} else {
		cmd.Printf("  Registered Messages: (none)\n")
	}

	return nil
}
