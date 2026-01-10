package client

import (
	"context"
	"log"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config holds client configuration
type Config struct {
	Endpoint string
	Insecure bool
}

// Client wraps a gRPC connection to Collector
type Client struct {
	conn   *grpc.ClientConn
	config Config
}

// interceptorLogger adapts standard log to interceptor logger
func interceptorLogger() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		log.Printf("[%s] %s %v", lvl, msg, fields)
	})
}

// New creates a new Collector client with middleware
func New(config Config) (*Client, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		// Client middleware
		grpc.WithUnaryInterceptor(logging.UnaryClientInterceptor(
			interceptorLogger(),
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		)),
		grpc.WithUnaryInterceptor(retry.UnaryClientInterceptor(
			retry.WithMax(3),
			retry.WithPerRetryTimeout(2*time.Second),
			retry.WithBackoff(retry.BackoffExponential(100*time.Millisecond)),
		)),
		grpc.WithUnaryInterceptor(timeout.UnaryClientInterceptor(30 * time.Second)),

		// Stream interceptors
		grpc.WithStreamInterceptor(logging.StreamClientInterceptor(
			interceptorLogger(),
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		)),
		grpc.WithStreamInterceptor(retry.StreamClientInterceptor(
			retry.WithMax(3),
			retry.WithPerRetryTimeout(2*time.Second),
		)),
	}

	conn, err := grpc.NewClient(config.Endpoint, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		config: config,
	}, nil
}

// Conn returns the underlying gRPC connection
func (c *Client) Conn() *grpc.ClientConn {
	return c.conn
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	return c.conn.Close()
}
