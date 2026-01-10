package collection

import (
	"fmt"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	listCollectionsEndpoint  string
	listCollectionsNamespace string
	listCollectionsPageSize  int32
	listCollectionsPageToken string
)

func NewListCollectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List collections",
		Long:  "Discover and list collections in a namespace",
		RunE:  runListCollections,
	}

	cmd.Flags().StringVar(&listCollectionsEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().StringVar(&listCollectionsNamespace, "namespace", "", "Namespace (empty for all namespaces)")
	cmd.Flags().Int32Var(&listCollectionsPageSize, "page-size", 20, "Number of collections per page")
	cmd.Flags().StringVar(&listCollectionsPageToken, "page-token", "", "Page token for pagination")

	return cmd
}

func runListCollections(cmd *cobra.Command, args []string) error {
	// Create client
	cl, err := client.New(client.Config{
		Endpoint: listCollectionsEndpoint,
		Insecure: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Create request
	req := &pb.DiscoverRequest{
		Namespace: listCollectionsNamespace,
		PageSize:  listCollectionsPageSize,
	}

	if listCollectionsPageToken != "" {
		req.PageToken = listCollectionsPageToken
	}

	// Call Discover
	repoClient := pb.NewCollectionRepoClient(cl.Conn())
	resp, err := repoClient.Discover(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("list collections failed: %w", err)
	}

	// Check status
	if resp.Status != nil && resp.Status.Code != pb.Status_OK {
		return fmt.Errorf("list collections failed: %s (code: %d)", resp.Status.Message, resp.Status.Code)
	}

	// Print results
	cmd.Printf("Found %d collection(s)", len(resp.Collections))
	if listCollectionsNamespace != "" {
		cmd.Printf(" in namespace '%s'", listCollectionsNamespace)
	}
	cmd.Println()
	cmd.Println()

	for i, coll := range resp.Collections {
		cmd.Printf("[%d] %s/%s\n", i+1, coll.Namespace, coll.Name)
		if coll.MessageType != nil {
			cmd.Printf("    Type: %s/%s\n",
				coll.MessageType.Namespace,
				coll.MessageType.MessageName)
		}
		if len(coll.IndexedFields) > 0 {
			cmd.Printf("    Indexed Fields: %v\n", coll.IndexedFields)
		}
		cmd.Println()
	}

	if resp.NextPageToken != "" {
		cmd.Printf("(More results available, use --page-token %s)\n", resp.NextPageToken)
	}

	return nil
}
