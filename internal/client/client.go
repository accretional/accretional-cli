package client

import (
	"github.com/accretional/accretional-cli/internal/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config holds client configuration
type Config struct {
	Endpoint string
	Insecure bool
	Interceptors interceptors.Config
}

// DefaultConfig returns a Config with sensible defaults for the given endpoint.
func DefaultConfig(endpoint string) Config {
	return Config{
		Endpoint:     endpoint,
		Insecure:     true,
		Interceptors: interceptors.DefaultConfig(),
	}
}

// Client wraps a gRPC connection to Collector
type Client struct {
	conn    *grpc.ClientConn
	config  Config
	metrics *interceptors.Metrics
}

func New(config Config) (*Client, error) {
	// Use default interceptor config if not specified
	if config.Interceptors.Timeout == 0 {
		config.Interceptors = interceptors.DefaultConfig()
	}

	metrics := interceptors.NewMetrics()

	// Build interceptor arrays
	unaryInterceptors := interceptors.NewUnaryInterceptors(config.Interceptors, metrics)
	streamInterceptors := interceptors.NewStreamInterceptors(config.Interceptors)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		grpc.WithChainUnaryInterceptor(unaryInterceptors...),
		grpc.WithChainStreamInterceptor(streamInterceptors...),
	}

	conn, err := grpc.NewClient(config.Endpoint, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		config:  config,
		metrics: metrics,
	}, nil
}

// Conn returns the underlying gRPC connection
func (c *Client) Conn() *grpc.ClientConn {
	return c.conn
}

// Metrics returns the client's metrics collector
func (c *Client) Metrics() *interceptors.Metrics {
	return c.metrics
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	return c.conn.Close()
}
