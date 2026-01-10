package server

import (
	"fmt"
	"time"

	"github.com/accretional/accretional-cli/internal/client"
	pb "github.com/accretional/collector/gen/collector"
	"github.com/spf13/cobra"
)

var (
	healthEndpoint string
	healthTimeout  time.Duration
)

func NewHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Check server health and connectivity",
		Long:  "Check if the Collector server is reachable and responding",
		RunE:  runHealth,
	}

	cmd.Flags().StringVar(&healthEndpoint, "endpoint", "localhost:50051", "gRPC server endpoint")
	cmd.Flags().DurationVar(&healthTimeout, "timeout", 5*time.Second, "Timeout for health check")

	return cmd
}

func runHealth(cmd *cobra.Command, args []string) error {
	startTime := time.Now()

	// Create client
	cl, err := client.New(client.Config{
		Endpoint: healthEndpoint,
		Insecure: true,
	})
	if err != nil {
		cmd.Printf("❌ UNREACHABLE\n")
		cmd.Printf("  Error: %v\n", err)
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer cl.Close()

	// Try a lightweight RPC call to verify connectivity
	// Use Discover as it's lightweight and doesn't require specific parameters
	ctx := cmd.Context()
	repoClient := pb.NewCollectionRepoClient(cl.Conn())

	req := &pb.DiscoverRequest{
		PageSize: 1, // Minimal request
	}

	_, err = repoClient.Discover(ctx, req)
	responseTime := time.Since(startTime)

	if err != nil {
		cmd.Printf("❌ UNHEALTHY\n")
		cmd.Printf("  Endpoint: %s\n", healthEndpoint)
		cmd.Printf("  Response Time: %v\n", responseTime)
		cmd.Printf("  Error: %v\n", err)
		return fmt.Errorf("health check failed: %w", err)
	}

	// Server is healthy
	cmd.Printf("✅ HEALTHY\n")
	cmd.Printf("  Endpoint: %s\n", healthEndpoint)
	cmd.Printf("  Response Time: %v\n", responseTime)

	return nil
}
