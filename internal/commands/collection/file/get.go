package file

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
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

	// Unmarshal Struct from Any
	var fileStruct structpb.Struct
	if err := proto.Unmarshal(resp.Item.Value, &fileStruct); err != nil {
		return fmt.Errorf("failed to unmarshal file record: %w", err)
	}

	fileMap := fileStruct.AsMap()

	// Extract file data
	var fileData []byte
	var fileName string

	// Handle different storage formats
	if dataVal, ok := fileMap["data"]; ok {
		switch v := dataVal.(type) {
		case []byte:
			fileData = v
		case string:
			// Try base64 decode first
			if decoded, err := base64.StdEncoding.DecodeString(v); err == nil {
				fileData = decoded
			} else {
				// If not base64, treat as raw string
				fileData = []byte(v)
			}
		default:
			return fmt.Errorf("unexpected data type in file record")
		}
	} else {
		return fmt.Errorf("file record missing 'data' field")
	}

	// Get file name
	if nameVal, ok := fileMap["name"].(string); ok {
		fileName = nameVal
	} else if pathVal, ok := fileMap["path"].(string); ok {
		fileName = filepath.Base(pathVal)
	} else {
		fileName = filepath.Base(getFilePath)
	}

	// Determine output path
	outputPath := getFileOutput
	if outputPath == "" {
		outputPath = fileName
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
