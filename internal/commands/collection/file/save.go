package file

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	saveFileEndpoint    string
	saveFileNamespace   string
	saveFileCollection  string
	saveFilePath        string
	saveFileLocalPath   string
	saveFileExtractText bool
)

func NewSaveFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "save-file",
		Short: "Save a file to a collection",
		Long:  "Store a file as a standalone record in a collection using CollectionData",
		RunE:  runSaveFile,
	}

	cmd.Flags().StringVar(&saveFileEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&saveFileNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&saveFileCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&saveFilePath, "path", "", "File path in collection (required)")
	cmd.Flags().StringVar(&saveFileLocalPath, "file", "", "Local file path to upload (required)")
	cmd.Flags().BoolVar(&saveFileExtractText, "extract-text", false, "Extract text content for searchability (FTS and semantic search)")

	cmd.MarkFlagRequired("collection")
	cmd.MarkFlagRequired("path")
	cmd.MarkFlagRequired("file")

	return cmd
}

func runSaveFile(cmd *cobra.Command, args []string) error {
	// Read the local file
	fileData, err := os.ReadFile(saveFileLocalPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Get file name from local path
	fileName := filepath.Base(saveFileLocalPath)

	// Extract text content from file (only if flag is set)
	var textContent string
	var mimeType string
	if saveFileExtractText {
		var err error
		textContent, mimeType, err = extractTextFromFile(fileData, fileName)
		if err != nil {
			// Log warning but continue - file will be stored but not searchable
			cmd.Printf("Warning: Could not extract text content: %v\n", err)
			mimeType = "application/octet-stream"
		}
	} else {
		// Detect MIME type from extension without extracting text
		mimeType = detectMimeType(fileName)
	}

	// Store file as a Struct record (since collections are typically configured for Struct)
	// We'll store file metadata, base64-encoded data, and extracted text content
	fileRecord := map[string]interface{}{
		"_type":    "file",
		"name":     fileName,
		"path":     saveFilePath,
		"size":     len(fileData),
		"data":     base64.StdEncoding.EncodeToString(fileData), // Base64 encode for JSON storage
		"mimeType": mimeType,
	}

	// Add text content if extraction was successful (for FTS and semantic search)
	if textContent != "" {
		fileRecord["content"] = textContent
	}

	// Convert to protobuf Struct
	structValue, err := structpb.NewStruct(fileRecord)
	if err != nil {
		return fmt.Errorf("failed to create Struct: %w", err)
	}

	// Create Any from Struct
	anyValue, err := anypb.New(structValue)
	if err != nil {
		return fmt.Errorf("failed to create Any from Struct: %w", err)
	}

	// Create client
	cl, err := client.New(client.Config{
		Endpoint: saveFileEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Generate record ID from path (use _file: prefix to identify file records)
	recordID := fmt.Sprintf("_file:%s", saveFilePath)

	// Create request
	req := &pb.CreateRequest{
		Namespace:      saveFileNamespace,
		CollectionName: saveFileCollection,
		Id:             recordID,
		Item:           anyValue,
	}

	// Call Create
	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	resp, err := collectionClient.Create(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("save file failed: %w", err)
	}

	// Check status
	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("save file failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	cmd.Printf("✓ File saved successfully\n")
	cmd.Printf("  Path: %s\n", saveFilePath)
	cmd.Printf("  Record ID: %s\n", resp.Id)
	cmd.Printf("  Size: %d bytes\n", len(fileData))
	cmd.Printf("  MIME Type: %s\n", mimeType)
	if textContent != "" {
		cmd.Printf("  Text Content: Extracted (%d characters)\n", len(textContent))
		cmd.Printf("  Searchable: Yes (FTS and semantic search enabled)\n")
	} else {
		cmd.Printf("  Text Content: Not extracted (binary file or unsupported format)\n")
		cmd.Printf("  Searchable: No (only metadata searchable)\n")
	}
	cmd.Printf("  Collection: %s/%s\n", saveFileNamespace, saveFileCollection)

	return nil
}
