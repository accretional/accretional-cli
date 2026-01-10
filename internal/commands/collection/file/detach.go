package file

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
	detachFileEndpoint   string
	detachFileNamespace  string
	detachFileCollection string
	detachFileRecordID   string
	detachFilePath       string
	detachFileDeleteFile bool
)

func NewDetachFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach-file",
		Short: "Detach a file from a record",
		Long:  "Remove a file attachment from a record by clearing the data_uri field",
		RunE:  runDetachFile,
	}

	cmd.Flags().StringVar(&detachFileEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&detachFileNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&detachFileCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&detachFileRecordID, "record-id", "", "Record ID to detach file from (required)")
	cmd.Flags().StringVar(&detachFilePath, "file-path", "", "File path to detach (optional, removes any attached file if not specified)")
	cmd.Flags().BoolVar(&detachFileDeleteFile, "delete-file", false, "Also delete the file record after detaching")

	cmd.MarkFlagRequired("collection")
	cmd.MarkFlagRequired("record-id")

	return cmd
}

func runDetachFile(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: detachFileEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	collectionClient := pb.NewCollectionServiceClient(cl.Conn())

	// Step 1: Get the existing record
	getReq := &pb.GetRequest{
		Namespace:      detachFileNamespace,
		CollectionName: detachFileCollection,
		Id:             detachFileRecordID,
	}

	getResp, err := collectionClient.Get(cmd.Context(), getReq)
	if err != nil {
		return fmt.Errorf("failed to get record: %w", err)
	}

	if getResp.Status != nil && getResp.Status.Code != pb.Status_OK {
		return fmt.Errorf("failed to get record: %s (code: %d)", getResp.Status.Message, getResp.Status.Code)
	}

	// Step 2: Unmarshal and update the record
	// Try to unmarshal as Struct first (most common case for JSON records)
	var existingData structpb.Struct
	var recordMap map[string]interface{}
	
	if err := proto.Unmarshal(getResp.Item.Value, &existingData); err == nil {
		// Successfully unmarshaled as Struct
		recordMap = existingData.AsMap()
	} else {
		// Not a Struct - cannot update data_uri for non-Struct records
		return fmt.Errorf("record is not a Struct type, cannot update data_uri. Only JSON records (google.protobuf.Struct) support file attachments")
	}

	dataURI, hasDataURI := recordMap["data_uri"].(string)

	if !hasDataURI || dataURI == "" {
		cmd.Println("No file attached to this record")
		return nil
	}

	// If specific file path provided, verify it matches
	if detachFilePath != "" {
		expectedFileID := fmt.Sprintf("_file:%s", detachFilePath)
		if dataURI != expectedFileID {
			return fmt.Errorf("file path mismatch: record has %s, expected %s", dataURI, expectedFileID)
		}
	}

	// Step 3: Remove data_uri from record
	delete(recordMap, "data_uri")

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
		Namespace:      detachFileNamespace,
		CollectionName: detachFileCollection,
		Id:             detachFileRecordID,
		Item:           updatedAny,
	}

	updateResp, err := collectionClient.Update(cmd.Context(), updateReq)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	if updateResp.Status != nil && updateResp.Status.Code != pb.Status_OK {
		return fmt.Errorf("failed to update record: %s (code: %d)", updateResp.Status.Message, updateResp.Status.Code)
	}

	cmd.Printf("✓ File detached successfully\n")
	cmd.Printf("  Record ID: %s\n", detachFileRecordID)
	cmd.Printf("  File Record ID: %s\n", dataURI)

	// Step 5: Optionally delete the file record
	if detachFileDeleteFile {
		deleteReq := &pb.DeleteRequest{
			Namespace:      detachFileNamespace,
			CollectionName: detachFileCollection,
			Id:             dataURI,
		}

		deleteResp, err := collectionClient.Delete(cmd.Context(), deleteReq)
		if err != nil {
			cmd.Printf("  Warning: Failed to delete file record: %v\n", err)
		} else if deleteResp.Status != nil && deleteResp.Status.Code != pb.Status_OK {
			cmd.Printf("  Warning: Failed to delete file record: %s\n", deleteResp.Status.Message)
		} else {
			cmd.Printf("  File record deleted\n")
		}
	}

	cmd.Printf("  Collection: %s/%s\n", detachFileNamespace, detachFileCollection)

	return nil
}
