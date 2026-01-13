package interceptors

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	fileLogger     *log.Logger
	fileLoggerOnce sync.Once
)

func getFileLogger() *log.Logger {
	fileLoggerOnce.Do(func() {
		logDir := ".log"
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return
		}

		logFile := filepath.Join(logDir, "accretional.log")
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}

		fileLogger = log.New(f, "", log.LstdFlags)
	})
	return fileLogger
}

type Config struct {
	Timeout       time.Duration
	MaxRetries    uint
	RetryTimeout  time.Duration
	EnableLogging bool
	EnableMetrics bool
	CLIVersion    string
}

func DefaultConfig() Config {
	return Config{
		Timeout:       30 * time.Second,
		MaxRetries:    3,
		RetryTimeout:  2 * time.Second,
		EnableLogging: true,
		EnableMetrics: true,
		CLIVersion:    "dev",
	}
}

func NewUnaryInterceptors(cfg Config, metrics *Metrics) []grpc.UnaryClientInterceptor {
	interceptors := []grpc.UnaryClientInterceptor{
		NewErrorInterceptor(),
		timeout.UnaryClientInterceptor(cfg.Timeout),

		retry.UnaryClientInterceptor(
			retry.WithMax(cfg.MaxRetries),
			retry.WithPerRetryTimeout(cfg.RetryTimeout),
			retry.WithBackoff(retry.BackoffExponential(100*time.Millisecond)),
		),

		NewCLIVersionInterceptor(cfg.CLIVersion),
	}

	if cfg.EnableMetrics && metrics != nil {
		interceptors = append(interceptors, NewMetricsInterceptor(metrics))
	}

	if cfg.EnableLogging {
		interceptors = append(interceptors, logging.UnaryClientInterceptor(
			newLogger(),
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		))
	}

	return interceptors
}

func NewStreamInterceptors(cfg Config) []grpc.StreamClientInterceptor {
	interceptors := []grpc.StreamClientInterceptor{
		retry.StreamClientInterceptor(
			retry.WithMax(cfg.MaxRetries),
			retry.WithPerRetryTimeout(cfg.RetryTimeout),
		),
	}

	if cfg.EnableLogging {
		interceptors = append(interceptors, logging.StreamClientInterceptor(
			newLogger(),
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		))
	}

	return interceptors
}

func NewCLIVersionInterceptor(version string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := metadata.Pairs("x-cli-version", version)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func newLogger() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		if logger := getFileLogger(); logger != nil {
			logger.Printf("[%s] %s %v", lvl, msg, fields)
		}
	})
}

type Metrics struct {
	TotalCalls     atomic.Int64
	SuccessfulCalls atomic.Int64
	FailedCalls    atomic.Int64
	TotalLatency   atomic.Int64
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) RecordCall(duration time.Duration, err error) {
	m.TotalCalls.Add(1)
	m.TotalLatency.Add(duration.Milliseconds())
	if err != nil {
		m.FailedCalls.Add(1)
	} else {
		m.SuccessfulCalls.Add(1)
	}
}

func (m *Metrics) Stats() (total, successful, failed int64, avgLatency time.Duration) {
	total = m.TotalCalls.Load()
	successful = m.SuccessfulCalls.Load()
	failed = m.FailedCalls.Load()
	if total > 0 {
		avgLatency = time.Duration(m.TotalLatency.Load()/total) * time.Millisecond
	}
	return
}

func NewMetricsInterceptor(metrics *Metrics) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		metrics.RecordCall(time.Since(start), err)
		return err
	}
}

type CLIError struct {
	Method string
	Code   codes.Code
	Cause  error
}

func (e *CLIError) Error() string {
	switch e.Code {
	case codes.Unavailable:
		return fmt.Sprintf("server unavailable: unable to connect for %s", e.Method)
	case codes.DeadlineExceeded:
		return fmt.Sprintf("request timed out: %s", e.Method)
	case codes.PermissionDenied:
		return fmt.Sprintf("permission denied: %s", e.Method)
	case codes.Unauthenticated:
		return fmt.Sprintf("authentication required: %s", e.Method)
	case codes.NotFound:
		return fmt.Sprintf("not found: %s", e.Method)
	case codes.InvalidArgument:
		return fmt.Sprintf("invalid argument: %s - %s", e.Method, status.Convert(e.Cause).Message())
	default:
		return fmt.Sprintf("%s failed: %s", e.Method, status.Convert(e.Cause).Message())
	}
}

func (e *CLIError) Unwrap() error {
	return e.Cause
}

func NewErrorInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			return &CLIError{
				Method: method,
				Code:   status.Code(err),
				Cause:  err,
			}
		}
		return nil
	}
}
