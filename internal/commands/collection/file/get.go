package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

var (
	getFileEndpoint   string
	getFileNamespace  string
	getFileCollection string
	getFilePath       string
	getFileOutput     string
)

func NewGetFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-file",
		Short: "Get a file from a collection",
		Long:  "Retrieve a standalone file stored in a collection",
		RunE:  runGetFile,
	}

	cmd.Flags().StringVar(&getFileEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&getFileNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&getFileCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&getFilePath, "path", "", "File path in collection (required)")
	cmd.Flags().StringVar(&getFileOutput, "output", "", "Output file path (optional, uses original filename if not provided)")

	cmd.MarkFlagRequired("collection")
	cmd.MarkFlagRequired("path")

	return cmd
}

func runGetFile(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: getFileEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Generate record ID from path
	recordID := fmt.Sprintf("_file:%s", getFilePath)

	// Create request
	req := &pb.GetRequest{
		Namespace:      getFileNamespace,
		CollectionName: getFileCollection,
		Id:             recordID,
	}

	// Call Get
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	resp, err := collectionClient.Get(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("get file failed: %w", err)
	}

	// Check status
	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("get file failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	// Unmarshal CollectionData from Any
	var collectionData pb.CollectionData
	if err := proto.Unmarshal(resp.Item.Value, &collectionData); err != nil {
		return fmt.Errorf("failed to unmarshal CollectionData: %w", err)
	}

	// Extract file data
	var fileData []byte
	switch content := collectionData.Content.(type) {
	case *pb.CollectionData_Data:
		fileData = content.Data
	case *pb.CollectionData_Uri:
		return fmt.Errorf("file stored as URI reference (large file), direct retrieval not yet supported")
	default:
		return fmt.Errorf("unknown content type in CollectionData")
	}

	// Determine output path
	outputPath := getFileOutput
	if outputPath == "" {
		if collectionData.Name != "" {
			outputPath = collectionData.Name
		} else {
			outputPath = filepath.Base(getFilePath)
		}
	}

	// Write file
	if err := os.WriteFile(outputPath, fileData, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	cmd.Printf("✓ File retrieved successfully\n")
	cmd.Printf("  Path: %s\n", getFilePath)
	cmd.Printf("  Output: %s\n", outputPath)
	cmd.Printf("  Size: %d bytes\n", len(fileData))
	cmd.Printf("  Collection: %s/%s\n", getFileNamespace, getFileCollection)

	return nil
}
