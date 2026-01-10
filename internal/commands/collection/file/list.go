package file

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	listFilesEndpoint   string
	listFilesNamespace  string
	listFilesCollection string
	listFilesRecordID   string
	listFilesPrefix     string
	listFilesStandalone bool
)

func NewListFilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-files",
		Short: "List files in a collection",
		Long: `List files in a collection. Can list:
  - Files attached to a specific record (use --record-id)
  - Standalone files (use --standalone with optional --prefix)`,
		RunE: runListFiles,
	}

	cmd.Flags().StringVar(&listFilesEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&listFilesNamespace, "namespace", "shared", "Namespace")
	cmd.Flags().StringVar(&listFilesCollection, "collection", "", "Collection name (required)")
	cmd.Flags().StringVar(&listFilesRecordID, "record-id", "", "Record ID to list attached files for")
	cmd.Flags().StringVar(&listFilesPrefix, "prefix", "", "Path prefix for filtering standalone files")
	cmd.Flags().BoolVar(&listFilesStandalone, "standalone", false, "List standalone files instead of attached files")

	cmd.MarkFlagRequired("collection")

	return cmd
}

func runListFiles(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: listFilesEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	collectionClient := pb.NewCollectionServiceClient(cl.Conn())

	if listFilesStandalone {
		// List standalone files (records with _file: prefix)
		return listStandaloneFiles(cmd, cl, collectionClient)
	} else if listFilesRecordID != "" {
		// List files attached to a specific record
		return listAttachedFiles(cmd, cl, collectionClient)
	} else {
		return fmt.Errorf("either --record-id or --standalone must be specified")
	}
}

func listStandaloneFiles(cmd *cobra.Command, cl *client.Client, collectionClient pb.CollectionServiceClient) error {
	// Note: Since List doesn't return record IDs, we can't directly filter by _file: prefix
	// This is a limitation - we list all records and try to identify file records
	// A better implementation would use a dedicated file collection or search with filters
	
	req := &pb.ListRequest{
		Namespace:      listFilesNamespace,
		CollectionName: listFilesCollection,
		PageSize:       100,
	}

	resp, err := collectionClient.List(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("list files failed: %w", err)
	}

	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("list files failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	// Try to identify file records by checking for _type: "file" in Struct records
	// This is heuristic-based since we can't filter by ID prefix
	var files []fileInfo
	for _, item := range resp.Items {
		var fileStruct structpb.Struct
		if err := proto.Unmarshal(item.Value, &fileStruct); err == nil {
			fileMap := fileStruct.AsMap()
			// Check if this is a file record
			if fileType, ok := fileMap["_type"].(string); ok && fileType == "file" {
				filePath := "unknown"
				if pathVal, ok := fileMap["path"].(string); ok {
					filePath = pathVal
				} else if nameVal, ok := fileMap["name"].(string); ok {
					filePath = nameVal
				}
				
				size := 0
				if sizeVal, ok := fileMap["size"].(float64); ok {
					size = int(sizeVal)
				} else if dataVal, ok := fileMap["data"]; ok {
					// Try to determine size from data (base64 string)
					if dataStr, ok := dataVal.(string); ok {
						// Approximate size (base64 is ~4/3 of original)
						size = len(dataStr) * 3 / 4
					}
				}
				
				files = append(files, fileInfo{
					name: filePath,
					size: size,
				})
			}
		}
	}

	if len(files) == 0 {
		cmd.Println("No standalone files found")
		cmd.Println("Note: File listing is limited - files are identified heuristically.")
		cmd.Println("Consider using a dedicated collection for files or use --record-id to list attached files.")
		return nil
	}

	cmd.Printf("Found %d standalone file(s)\n", len(files))
	cmd.Println()

	for i, file := range files {
		if file.size >= 0 {
			cmd.Printf("[%d] %s (%d bytes)\n", i+1, file.name, file.size)
		} else {
			cmd.Printf("[%d] %s (URI reference)\n", i+1, file.name)
		}
	}

	return nil
}

func listAttachedFiles(cmd *cobra.Command, cl *client.Client, collectionClient pb.CollectionServiceClient) error {
	// Get the record
	req := &pb.GetRequest{
		Namespace:      listFilesNamespace,
		CollectionName: listFilesCollection,
		Id:             listFilesRecordID,
	}

	resp, err := collectionClient.Get(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("get record failed: %w", err)
	}

	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("get record failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	// Unmarshal the record
	var recordData structpb.Struct
	if err := proto.Unmarshal(resp.Item.Value, &recordData); err != nil {
		return fmt.Errorf("failed to unmarshal record: %w", err)
	}

	recordMap := recordData.AsMap()
	dataURI, ok := recordMap["data_uri"].(string)
	if !ok || dataURI == "" {
		cmd.Println("No files attached to this record")
		return nil
	}

	// Get the file record
	fileReq := &pb.GetRequest{
		Namespace:      listFilesNamespace,
		CollectionName: listFilesCollection,
		Id:             dataURI,
	}

	fileResp, err := collectionClient.Get(cmd.Context(), fileReq)
	if err != nil {
		return fmt.Errorf("get file record failed: %w", err)
	}

	if fileResp.Status != nil && fileResp.Status.Code != pb.Status_OK {
		return fmt.Errorf("get file record failed: %s (code: %d)", fileResp.Status.Message, fileResp.Status.Code)
	}

	// Unmarshal file record (Struct format)
	var fileStruct structpb.Struct
	if err := proto.Unmarshal(fileResp.Item.Value, &fileStruct); err != nil {
		return fmt.Errorf("failed to unmarshal file: %w", err)
	}

	fileMap := fileStruct.AsMap()

	cmd.Printf("File attached to record %s:\n", listFilesRecordID)
	cmd.Printf("  File Record ID: %s\n", dataURI)
	
	if nameVal, ok := fileMap["name"].(string); ok {
		cmd.Printf("  Name: %s\n", nameVal)
	}
	if pathVal, ok := fileMap["path"].(string); ok {
		cmd.Printf("  Path: %s\n", pathVal)
	}
	if sizeVal, ok := fileMap["size"].(float64); ok {
		cmd.Printf("  Size: %d bytes\n", int(sizeVal))
	} else if dataVal, ok := fileMap["data"]; ok {
		switch v := dataVal.(type) {
		case []byte:
			cmd.Printf("  Size: %d bytes\n", len(v))
		case string:
			cmd.Printf("  Size: %d bytes\n", len(v))
		}
	}

	return nil
}

type fileInfo struct {
	name string
	size int
}
