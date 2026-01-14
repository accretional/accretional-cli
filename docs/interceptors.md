# Available gRPC Client Interceptors

## Retry & Resilience

**github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry**
- `retry.UnaryClientInterceptor()` - Automatically retries failed requests with configurable backoff and retry codes.
- `retry.StreamClientInterceptor()` - Automatically retries failed stream requests with configurable backoff.

## Circuit Breaker

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `HystrixClientInterceptor(opts...)` - Implements circuit breaker pattern using Hystrix to prevent cascading failures.

## Timeout

**github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout**
- `timeout.UnaryClientInterceptor(duration)` - Enforces a maximum timeout for unary RPC calls.

**github.com/nhatthm/go-grpc-middleware/timeout**
- `timeout.UnaryClientTimeoutInterceptor(duration)` - Automatically creates a context with timeout if none exists in the current context.
- `timeout.StreamClientTimeoutInterceptor(duration)` - Automatically creates a context with timeout for streaming calls if none exists.
- `timeout.UnaryClientSleepInterceptor(duration)` - Sleeps for a duration before executing the unary call.
- `timeout.StreamClientSleepInterceptor(duration)` - Sleeps for a duration before executing the streaming call.
- `timeout.WithUnaryClientTimeoutInterceptor(duration)` - Dial option that appends timeout interceptor.
- `timeout.WithStreamClientTimeoutInterceptor(duration)` - Dial option that appends timeout interceptor for streams.
- `timeout.WithUnaryClientSleepInterceptor(duration)` - Dial option that appends sleep interceptor.
- `timeout.WithStreamClientSleepInterceptor(duration)` - Dial option that appends sleep interceptor for streams.

## Rate Limiting

**github.com/tommy-sho/rate-limiter-grpc-go**
- `ratelimit.UnaryClientInterceptor(limiter)` - Limits the number of unary requests per second using a rate limiter.
- `ratelimit.StreamClientInterceptor(limiter)` - Limits the number of streaming requests per second using a rate limiter.

## Logging

**github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging**
- `logging.UnaryClientInterceptor(logger, opts...)` - Logs unary RPC calls with configurable events and fields.
- `logging.StreamClientInterceptor(logger, opts...)` - Logs streaming RPC calls with configurable events and fields.

**github.com/nhatthm/go-grpc-middleware/logging/ctxd**
- `ctxd.UnaryClientInterceptor(logger, opts...)` - Logs unary client calls using ctxd logger with structured fields and timing.
- `ctxd.StreamClientInterceptor(logger, opts...)` - Logs streaming client calls using ctxd logger with structured fields and timing.

## Metrics & Observability

**github.com/grpc-ecosystem/go-grpc-prometheus**
- `grpcprom.UnaryClientInterceptor` - Exposes Prometheus metrics for unary client calls.
- `grpcprom.StreamClientInterceptor` - Exposes Prometheus metrics for streaming client calls.

**go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc**
- `otelgrpc.UnaryClientInterceptor()` - Provides OpenTelemetry tracing and metrics for unary client calls.
- `otelgrpc.StreamClientInterceptor()` - Provides OpenTelemetry tracing and metrics for streaming client calls.

**google.golang.org/grpc/stats/opentelemetry** !!! *Same as `otelgrpc` interceptors - alternative configuration method via DialOption*
- `opentelemetry.DialOption(opts)` - Configures OpenTelemetry metrics and tracing via gRPC dial options.

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `NewRelicClientInterceptor()` - Reports client RPC calls to NewRelic APM for monitoring and performance tracking.
- `GRPCClientInterceptor(opts...)` - Adds OpenTracing support to client calls for distributed tracing.
- `DefaultClientInterceptor(opts...)` - Chains all default unary client interceptors (Hystrix, retry, OpenTracing, NewRelic, Prometheus).
- `DefaultClientStreamInterceptor(opts...)` - Chains all default stream client interceptors (OpenTracing, NewRelic, Prometheus).

---

# Available gRPC Server Interceptors

## Rate Limiting

**github.com/tommy-sho/rate-limiter-grpc-go**
- `ratelimit.UnaryServerInterceptor(limiter)` - Limits the number of unary requests per second using a rate limiter.
- `ratelimit.StreamServerInterceptor(limiter)` - Limits the number of streaming requests per second using a rate limiter.

## Logging

**github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging**
- `logging.UnaryServerInterceptor(logger, opts...)` - Logs unary RPC calls with configurable events and fields.
- `logging.StreamServerInterceptor(logger, opts...)` - Logs streaming RPC calls with configurable events and fields.

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `DebugLoggingInterceptor()` - Logs all request/response data for debugging purposes.
- `ResponseTimeLoggingInterceptor(ff)` - Logs response time for each unary request with optional filtering.
- `ResponseTimeLoggingStreamInterceptor()` - Logs response time for streaming RPCs.

**github.com/mercari/go-grpc-interceptor/zap**
- `zap.UnaryServerInterceptor(logger)` - Attaches zap logger to each unary request context with method name.
- `zap.StreamServerInterceptor(logger)` - Attaches zap logger to each streaming request context with method name.
- `zap.UnaryServerInterceptorWithRequestID(logger)` - Attaches zap logger with request ID to each unary request context.
- `zap.StreamServerInterceptorWithRequestID(logger)` - Attaches zap logger with request ID to each streaming request context.

**github.com/mercari/go-grpc-interceptor/requestdump**
- `requestdump.UnaryServerInterceptor(opts...)` - Dumps request and response messages for debugging purposes.

**github.com/nhatthm/go-grpc-middleware/logging/ctxd**
- `ctxd.UnaryServerInterceptor(logger, opts...)` - Logs unary server calls using ctxd logger with structured fields and timing.
- `ctxd.StreamServerInterceptor(logger, opts...)` - Logs streaming server calls using ctxd logger with structured fields and timing.

## Recovery & Error Handling

**github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery**
- `recovery.UnaryServerInterceptor(opts...)` - Recovers from panics and converts them to gRPC errors.

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `PanicRecoveryInterceptor()` - Recovers from panics, logs errors, and reports to NewRelic and error notifier.
- `ServerErrorInterceptor()` - Intercepts server errors and reports them to error notifier with trace IDs.
- `ServerErrorStreamInterceptor()` - Intercepts streaming server errors and reports them to error notifier.

**github.com/mercari/go-grpc-interceptor/panichandler** !!! *Similar to `recovery.UnaryServerInterceptor` - basic panic recovery with less configuration options*
- `panichandler.UnaryServerInterceptor` - Protects process from aborting by panic and returns Internal error as status code.
- `panichandler.StreamServerInterceptor` - Protects streaming handlers from panic and returns Internal error as status code.

## Authentication & Authorization

**github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth**
- `auth.UnaryServerInterceptor(authFunc)` - Validates authentication tokens from request metadata.

**google.golang.org/grpc/authz**
- `authz.NewStatic(policy)` - Provides RBAC authorization with static policy configuration.
- `authz.NewFileWatcher(path, interval)` - Provides RBAC authorization with file-based policy that auto-reloads.

## Validation

**github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator**
- `validator.UnaryServerInterceptor()` - Validates request messages using protobuf validation rules.
- `validator.StreamServerInterceptor()` - Validates streaming request messages using protobuf validation rules.

**github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate**
- `protovalidate.UnaryServerInterceptor(validator)` - Validates requests using protovalidate library.
- `protovalidate.StreamServerInterceptor(validator)` - Validates streaming requests using protovalidate library.

## Metrics & Observability

**github.com/grpc-ecosystem/go-grpc-prometheus**
- `grpcprom.UnaryServerInterceptor` - Exposes Prometheus metrics for unary server calls.
- `grpcprom.StreamServerInterceptor` - Exposes Prometheus metrics for streaming server calls.

**go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc**
- `otelgrpc.UnaryServerInterceptor()` - Provides OpenTelemetry tracing and metrics for unary server calls.
- `otelgrpc.StreamServerInterceptor()` - Provides OpenTelemetry tracing and metrics for streaming server calls.

**google.golang.org/grpc/stats/opentelemetry** !!! *Same as `otelgrpc` interceptors - alternative configuration method via ServerOption*
- `opentelemetry.ServerOption(opts)` - Configures OpenTelemetry metrics and tracing via gRPC server options.

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `NewRelicInterceptor()` - Reports server RPC calls to NewRelic APM for monitoring and performance tracking.

**github.com/mercari/go-grpc-interceptor/instrument**
- `instrument.UnaryServerInterceptor` - Provides instrumentation hook for tracking method calls with timing and error information.
- `instrument.StreamServerInterceptor` - Provides instrumentation hook for tracking streaming method calls with timing and error information.

## Request Processing

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `TraceIdInterceptor()` - Extracts and injects trace IDs from request objects into context.
- `OptionsInterceptor()` - Adds options and logging context to requests.

**github.com/mercari/go-grpc-interceptor/xrequestid**
- `xrequestid.UnaryServerInterceptor(opts...)` - Generates or extracts X-Request-Id from metadata and injects into context.
- `xrequestid.StreamServerInterceptor(opts...)` - Generates or extracts X-Request-Id from metadata for streaming requests.

**github.com/mercari/go-grpc-interceptor/acceptlang**
- `acceptlang.UnaryServerInterceptor` - Parses Accept-Language header from metadata and makes available via context.
- `acceptlang.StreamServerInterceptor` - Parses Accept-Language header from metadata for streaming requests.

## Conditional Application

**github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector**
- `selector.UnaryServerInterceptor(interceptor, matchFunc)` - Conditionally applies an interceptor based on method matching.
- `selector.StreamServerInterceptor(interceptor, matchFunc)` - Conditionally applies a stream interceptor based on method matching.

## Real IP Extraction

**github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/realip**
- `realip.UnaryServerInterceptor(trustedPeers, headers)` - Extracts real client IP from proxy headers (X-Real-IP, X-Forwarded-For, etc.).

## Convenience Functions

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `DefaultInterceptors()` - Returns array of all default unary server interceptors (logging, tracing, tags, Prometheus, errors, NewRelic, panic recovery).
- `DefaultStreamInterceptors()` - Returns array of all default stream server interceptors (logging, tracing, tags, Prometheus, errors).

**github.com/mercari/go-grpc-interceptor/multiinterceptor**
- `multiinterceptor.NewMultiUnaryServerInterceptor(interceptors...)` - Chains multiple unary server interceptors together.
- `multiinterceptor.NewMultiStreamServerInterceptor(interceptors...)` - Chains multiple stream server interceptors together.

**github.com/codemodus/catena**  *Similar to `multiinterceptor` - provides additional append/merge operations but same core chaining functionality*
- `catena.NewUnaryServerCatena(interceptors...)` - Creates a catena for chaining unary server interceptors with append/merge operations.
- `catena.UnaryServerCatena.Interceptor()` - Returns the chained unary server interceptor from a catena.
- `catena.UnaryServerCatena.ServerOption()` - Returns a gRPC ServerOption for the chained interceptor.
