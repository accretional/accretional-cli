package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	attachFileEndpoint   string
	attachFileNamespace  string
	attachFileCollection string
	attachFileRecordID   string
	attachFileLocalPath   string
	attachFilePath       string
)

func NewAttachFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach-file",
		Short: "Attach a file to a record",
		Long:  "Attach a file to an existing record by storing it and setting the record's data_uri field",
		RunE:  runAttachFile,
	}

	cmd.Flags().StringVar(&attachFileEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&attachFileNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&attachFileCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&attachFileRecordID, "record-id", "", "Record ID to attach file to (required)")
	cmd.Flags().StringVar(&attachFileLocalPath, "file", "", "Local file path to attach (required)")
	cmd.Flags().StringVar(&attachFilePath, "path", "", "File path in collection (optional, auto-generated if not provided)")

	cmd.MarkFlagRequired("collection")
	cmd.MarkFlagRequired("record-id")
	cmd.MarkFlagRequired("file")

	return cmd
}

func runAttachFile(cmd *cobra.Command, args []string) error {
	// Read the local file
	fileData, err := os.ReadFile(attachFileLocalPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Get file name from local path
	fileName := filepath.Base(attachFileLocalPath)

	// Generate file path if not provided
	if attachFilePath == "" {
		attachFilePath = fmt.Sprintf("attachments/%s/%s", attachFileRecordID, fileName)
	}

	// Create client
	cl, err := client.New(client.Config{
		Endpoint: attachFileEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Step 1: Store the file as a record
	collectionData := &pb.CollectionData{
		Name: fileName,
		Content: &pb.CollectionData_Data{
			Data: fileData,
		},
	}

	anyValue, err := anypb.New(collectionData)
	if err != nil {
		return fmt.Errorf("failed to create Any from CollectionData: %w", err)
	}

	fileRecordID := fmt.Sprintf("_file:%s", attachFilePath)
	createReq := &pb.CreateRequest{
		Namespace:      attachFileNamespace,
		CollectionName: attachFileCollection,
		Id:             fileRecordID,
		Item:           anyValue,
	}

	collectionClient := pb.NewCollectionServiceClient(cl.Conn())
	createResp, err := collectionClient.Create(cmd.Context(), createReq)
	if err != nil {
		return fmt.Errorf("failed to store file: %w", err)
	}

	if createResp.Status != nil && createResp.Status.Code != pb.Status_OK {
		return fmt.Errorf("failed to store file: %s (code: %d)", createResp.Status.Message, createResp.Status.Code)
	}

	// Step 2: Get the existing record
	getReq := &pb.GetRequest{
		Namespace:      attachFileNamespace,
		CollectionName: attachFileCollection,
		Id:             attachFileRecordID,
	}

	getResp, err := collectionClient.Get(cmd.Context(), getReq)
	if err != nil {
		return fmt.Errorf("failed to get record: %w", err)
	}

	if getResp.Status != nil && getResp.Status.Code != pb.Status_OK {
		return fmt.Errorf("failed to get record: %s (code: %d)", getResp.Status.Message, getResp.Status.Code)
	}

	// Step 3: Update the record with data_uri
	// Try to unmarshal as Struct first (most common case for JSON records)
	var existingData structpb.Struct
	var recordMap map[string]interface{}
	
	if err := proto.Unmarshal(getResp.Item.Value, &existingData); err == nil {
		// Successfully unmarshaled as Struct
		recordMap = existingData.AsMap()
	} else {
		// Not a Struct - we need to handle this differently
		// For now, create a new Struct with just the data_uri
		// Note: This will replace the entire record content
		cmd.Printf("Warning: Record is not a Struct type. Creating new record with data_uri.\n")
		recordMap = make(map[string]interface{})
	}
	
	// Add data_uri to the record
	recordMap["data_uri"] = fileRecordID

	updatedStruct, err := structpb.NewStruct(recordMap)
	if err != nil {
		return fmt.Errorf("failed to create updated struct: %w", err)
	}

	updatedAny, err := anypb.New(updatedStruct)
	if err != nil {
		return fmt.Errorf("failed to create Any from updated struct: %w", err)
	}

	// Step 4: Update the record
	updateReq := &pb.UpdateRequest{
		Namespace:      attachFileNamespace,
		CollectionName: attachFileCollection,
		Id:             attachFileRecordID,
		Item:           updatedAny,
	}

	updateResp, err := collectionClient.Update(cmd.Context(), updateReq)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	if updateResp.Status != nil && updateResp.Status.Code != pb.Status_OK {
		return fmt.Errorf("failed to update record: %s (code: %d)", updateResp.Status.Message, updateResp.Status.Code)
	}

	cmd.Printf("✓ File attached successfully\n")
	cmd.Printf("  Record ID: %s\n", attachFileRecordID)
	cmd.Printf("  File Path: %s\n", attachFilePath)
	cmd.Printf("  File Record ID: %s\n", fileRecordID)
	cmd.Printf("  Size: %d bytes\n", len(fileData))
	cmd.Printf("  Collection: %s/%s\n", attachFileNamespace, attachFileCollection)

	return nil
}
