# Available gRPC Client Interceptors

## Retry & Resilience

**github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry**
- `retry.UnaryClientInterceptor()` - Automatically retries failed requests with configurable backoff and retry codes.
- `retry.StreamClientInterceptor()` - Automatically retries failed stream requests with configurable backoff.

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

**github.com/grpc-ecosystem/go-grpc-prometheus** ⚠️ *Deprecated - Use `github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus` instead*
- `grpcprom.UnaryClientInterceptor` - Exposes Prometheus metrics for unary client calls.
- `grpcprom.StreamClientInterceptor` - Exposes Prometheus metrics for streaming client calls.

**github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus** *Replacement for go-grpc-prometheus*
- `ClientMetrics.UnaryClientInterceptor(opts...)` - Exposes Prometheus metrics for unary client calls (create with `NewClientMetrics()`).
- `ClientMetrics.StreamClientInterceptor(opts...)` - Exposes Prometheus metrics for streaming client calls (create with `NewClientMetrics()`).

**go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc**
- `otelgrpc.UnaryClientInterceptor()` - Provides OpenTelemetry tracing and metrics for unary client calls.
- `otelgrpc.StreamClientInterceptor()` - Provides OpenTelemetry tracing and metrics for streaming client calls.

**google.golang.org/grpc/stats/opentelemetry** !!! *Same as `otelgrpc` interceptors - alternative configuration method via DialOption*
- `opentelemetry.DialOption(opts)` - Configures OpenTelemetry metrics and tracing via gRPC dial options.

**go.elastic.co/apm/module/apmgrpc/v2**
- `apmgrpc.NewUnaryClientInterceptor(opts...)` - Traces gRPC client requests with Elastic APM, creating spans for external gRPC calls.
- `apmgrpc.NewStreamClientInterceptor(opts...)` - Traces gRPC streaming client requests with Elastic APM, creating spans for stream calls.

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

**github.com/grpc-ecosystem/go-grpc-prometheus** ⚠️ *Deprecated - Use `github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus` instead*
- `grpcprom.UnaryServerInterceptor` - Exposes Prometheus metrics for unary server calls.
- `grpcprom.StreamServerInterceptor` - Exposes Prometheus metrics for streaming server calls.

**github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus** *Replacement for go-grpc-prometheus*
- `ServerMetrics.UnaryServerInterceptor(opts...)` - Exposes Prometheus metrics for unary server calls (create with `NewServerMetrics()`).
- `ServerMetrics.StreamServerInterceptor(opts...)` - Exposes Prometheus metrics for streaming server calls (create with `NewServerMetrics()`).

**go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc**
- `otelgrpc.UnaryServerInterceptor()` - Provides OpenTelemetry tracing and metrics for unary server calls.
- `otelgrpc.StreamServerInterceptor()` - Provides OpenTelemetry tracing and metrics for streaming server calls.

**google.golang.org/grpc/stats/opentelemetry** !!! *Same as `otelgrpc` interceptors - alternative configuration method via ServerOption*
- `opentelemetry.ServerOption(opts)` - Configures OpenTelemetry metrics and tracing via gRPC server options.

**go.elastic.co/apm/module/apmgrpc/v2**
- `apmgrpc.NewUnaryServerInterceptor(opts...)` - Traces gRPC server requests with Elastic APM, creating transactions with optional panic recovery.
- `apmgrpc.NewStreamServerInterceptor(opts...)` - Traces gRPC streaming server requests with Elastic APM, creating transactions with optional panic recovery.

**github.com/mercari/go-grpc-interceptor/instrument**
- `instrument.UnaryServerInterceptor` - Provides instrumentation hook for tracking method calls with timing and error information.
- `instrument.StreamServerInterceptor` - Provides instrumentation hook for tracking streaming method calls with timing and error information.

## Request Processing

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

**github.com/mercari/go-grpc-interceptor/multiinterceptor**
- `multiinterceptor.NewMultiUnaryServerInterceptor(interceptors...)` - Chains multiple unary server interceptors together.
- `multiinterceptor.NewMultiStreamServerInterceptor(interceptors...)` - Chains multiple stream server interceptors together.

**github.com/codemodus/catena**  *Similar to `multiinterceptor` - provides additional append/merge operations but same core chaining functionality*
- `catena.NewUnaryServerCatena(interceptors...)` - Creates a catena for chaining unary server interceptors with append/merge operations.
- `catena.UnaryServerCatena.Interceptor()` - Returns the chained unary server interceptor from a catena.
- `catena.UnaryServerCatena.ServerOption()` - Returns a gRPC ServerOption for the chained interceptor.

---

# Dependent Interceptors

*These interceptors require specific frameworks or ecosystems to function.*

## Circuit Breaker (Client & Server)

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `HystrixClientInterceptor(opts...)` - Implements circuit breaker pattern using Hystrix to prevent cascading failures.

**github.com/zhufuyi/pkg/grpc/interceptor** *Requires zhufuyi/pkg ecosystem (github.com/zhufuyi/pkg/shield/circuitbreaker, github.com/zhufuyi/pkg/container/group)*
- `UnaryClientCircuitBreaker(opts...)` - Implements circuit breaker pattern for unary client calls using zhufuyi shield library.
- `SteamClientCircuitBreaker(opts...)` - Implements circuit breaker pattern for streaming client calls.
- `UnaryServerCircuitBreaker(opts...)` - Implements circuit breaker pattern for unary server calls.
- `SteamServerCircuitBreaker(opts...)` - Implements circuit breaker pattern for streaming server calls.

## Retry (Client)

**github.com/zhufuyi/pkg/grpc/interceptor** *Requires zhufuyi/pkg ecosystem (github.com/zhufuyi/pkg/errcode)*
- `UnaryClientRetry(opts...)` - Retries failed unary client requests with configurable options.
- `StreamClientRetry(opts...)` - Retries failed streaming client requests with configurable options.

## Timeout (Client)

**github.com/zhufuyi/pkg/grpc/interceptor** *Requires zhufuyi/pkg ecosystem*
- `UnaryTimeout(duration)` - Enforces timeout for unary client calls.
- `StreamTimeout(duration)` - Enforces timeout for streaming client calls.

## Rate Limiting (Server)

**github.com/zhufuyi/pkg/grpc/interceptor** *Requires zhufuyi/pkg ecosystem (github.com/zhufuyi/pkg/shield/ratelimit, github.com/zhufuyi/pkg/errcode)*
- `UnaryServerRateLimit(opts...)` - Rate limits unary server requests with configurable window and bucket sizes.
- `StreamServerRateLimit(opts...)` - Rate limits streaming server requests with configurable window and bucket sizes.

## Authentication (Client & Server)

**github.com/zhufuyi/pkg/grpc/interceptor** *Requires zhufuyi/pkg ecosystem (github.com/zhufuyi/pkg/jwt)*
- `UnaryServerJwtAuth(opts...)` - Validates JWT tokens for unary server requests using zhufuyi JWT library.
- `StreamServerJwtAuth(opts...)` - Validates JWT tokens for streaming server requests.

**github.com/gravitational/teleport/api/utils/grpc/interceptors** *Requires Teleport ecosystem (github.com/gravitational/teleport/api/mfa)*
- `WithMFAUnaryInterceptor(mfaCeremony)` - Prompts for MFA credentials when required and retries unary client calls with MFA authentication.

**github.com/siderolabs/go-api-signature/pkg/client/interceptor** *Requires Sidero Labs ecosystem (github.com/siderolabs/go-api-signature/pkg/message, github.com/siderolabs/go-api-signature/pkg/pgp/client)*
- `Interceptor.Unary()` - Signs unary client requests using PGP keys for API authentication.
- `Interceptor.Stream()` - Signs streaming client requests using PGP keys for API authentication.

## Logging (Client & Server)

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `DebugLoggingInterceptor()` - Logs all request/response data for debugging purposes.
- `ResponseTimeLoggingInterceptor(ff)` - Logs response time for each unary request with optional filtering.
- `ResponseTimeLoggingStreamInterceptor()` - Logs response time for streaming RPCs.

**github.com/go-admin-team/go-admin-core/server/grpc/interceptors/logging** *Requires go-admin-core ecosystem (go-admin-team/go-admin-core/logger, go-admin-team/go-admin-core/tools/utils)*
- `UnaryClientInterceptor(opts...)` - Logs unary client calls with request ID and timing information.
- `StreamClientInterceptor(opts...)` - Logs streaming client calls with request ID and timing information.
- `UnaryServerInterceptor(opts...)` - Logs unary server calls with request ID, method, status, and timing.
- `StreamServerInterceptor(opts...)` - Logs streaming server calls with request ID, method, status, and timing.

**github.com/zhufuyi/pkg/grpc/interceptor** *Requires zhufuyi/pkg ecosystem (go.uber.org/zap)*
- `UnaryClientLog(logger, opts...)` - Logs unary client calls using zap logger.
- `StreamClientLog(logger, opts...)` - Logs streaming client calls using zap logger.
- `UnaryServerLog(logger, opts...)` - Logs unary server calls using zap logger with context tags.
- `StreamServerLog(logger, opts...)` - Logs streaming server calls using zap logger with context tags.
- `UnaryServerCtxTags()` - Adds context tags to unary server requests for logging.
- `StreamServerCtxTags()` - Adds context tags to streaming server requests for logging.

**goa.design/goa/v3/grpc/middleware** *Requires Goa framework (goa.design/goa/v3/middleware)*
- `middleware.UnaryServerLog(logger)` - Logs incoming unary gRPC requests and responses with request ID, method, status, bytes, and timing.
- `middleware.UnaryServerLogContext(logFromCtx)` - Logs unary requests using logger extracted from context.
- `middleware.StreamServerLog(logger)` - Logs incoming streaming gRPC requests and responses.
- `middleware.StreamServerLogContext(logFromCtx)` - Logs streaming requests using logger extracted from context.

## Recovery & Error Handling (Client & Server)

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `PanicRecoveryInterceptor()` - Recovers from panics, logs errors, and reports to NewRelic and error notifier.
- `ServerErrorInterceptor()` - Intercepts server errors and reports them to error notifier with trace IDs.
- `ServerErrorStreamInterceptor()` - Intercepts streaming server errors and reports them to error notifier.

**github.com/zhufuyi/pkg/grpc/interceptor** *Requires zhufuyi/pkg ecosystem (grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery)*
- `UnaryServerRecovery()` - Recovers from panics in unary server handlers.
- `StreamServerRecovery()` - Recovers from panics in streaming server handlers.

**github.com/kazegusuri/grpc-panic-handler** !!! *Similar to `recovery.UnaryServerInterceptor` and `panichandler` - basic panic recovery with custom handler support*
- `UnaryPanicHandler` - Recovers from panics in unary server handlers and returns Internal error, supports custom panic handlers.
- `StreamPanicHandler` - Recovers from panics in streaming server handlers and returns Internal error, supports custom panic handlers.

**github.com/gravitational/teleport/api/utils/grpc/interceptors** *Requires Teleport ecosystem (gravitational/trace, gravitational/teleport/api/trail)*
- `GRPCClientUnaryErrorInterceptor` - Converts gRPC status errors to Teleport trace errors for unary client calls.
- `GRPCClientStreamErrorInterceptor` - Converts gRPC status errors to Teleport trace errors for streaming client calls.
- `GRPCServerUnaryErrorInterceptor` - Converts Teleport trace errors to gRPC status errors for unary server calls.
- `GRPCServerStreamErrorInterceptor` - Converts Teleport trace errors to gRPC status errors for streaming server calls.

**github.com/google/trillian/server/interceptor** *Requires Trillian ecosystem (google/trillian/storage, google/trillian/quota, google/trillian/trees)*
- `ErrorWrapper()` - Wraps errors emitted by handlers using Trillian error wrapping logic.

## Metrics & Observability (Client & Server)

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `NewRelicClientInterceptor()` - Reports client RPC calls to NewRelic APM for monitoring and performance tracking.
- `NewRelicInterceptor()` - Reports server RPC calls to NewRelic APM for monitoring and performance tracking.

**github.com/zhufuyi/pkg/grpc/interceptor** *Requires zhufuyi/pkg ecosystem (zhufuyi/pkg/grpc/metrics)*
- `UnaryClientMetrics()` - Exposes Prometheus metrics for unary client calls.
- `StreamClientMetrics()` - Exposes Prometheus metrics for streaming client calls.
- `UnaryServerMetrics(opts...)` - Exposes Prometheus metrics for unary server calls.
- `StreamServerMetrics(opts...)` - Exposes Prometheus metrics for streaming server calls.

**go.chromium.org/luci/grpc/grpcmon** *Requires Chromium tsmon ecosystem (go.chromium.org/luci/common/tsmon)*
- `grpcmon.UnaryServerInterceptor` - Gathers RPC handler metrics (count, duration) and sends them to tsmon for server-side monitoring.
- `grpcmon.StreamServerInterceptor` - Gathers streaming RPC handler metrics and sends them to tsmon for server-side monitoring.
- `grpcmon.ClientRPCStatsMonitor` - Stats handler for client-side RPC monitoring via tsmon (use with WithStatsHandler dial option).

## Tracing (Client & Server)

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `GRPCClientInterceptor(opts...)` - Adds OpenTracing support to client calls for distributed tracing.

**github.com/zhufuyi/pkg/grpc/interceptor** *Requires zhufuyi/pkg ecosystem (go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc)*
- `UnaryClientTracing()` - Provides OpenTelemetry tracing for unary client calls.
- `StreamClientTracing()` - Provides OpenTelemetry tracing for streaming client calls.
- `UnaryServerTracing()` - Provides OpenTelemetry tracing for unary server calls.
- `StreamServerTracing()` - Provides OpenTelemetry tracing for streaming server calls.

**goa.design/goa/v3/grpc/middleware** *Requires Goa framework (goa.design/goa/v3/middleware)*
- `middleware.UnaryServerTrace(opts...)` - Initializes trace information (trace ID, span ID, parent span ID) in unary request context.
- `middleware.StreamServerTrace(opts...)` - Initializes trace information in streaming request context.
- `middleware.UnaryClientTrace()` - Sets outgoing unary request metadata with trace information from context.
- `middleware.StreamClientTrace()` - Sets outgoing streaming request metadata with trace information from context.

**goa.design/goa/v3/grpc/middleware/xray** *Requires Goa framework (goa.design/goa/v3/middleware/xray)*
- `xray.NewUnaryServer(service, daemon)` - Sends AWS X-Ray segments to daemon for unary server requests.
- `xray.NewStreamServer(service, daemon)` - Sends AWS X-Ray segments to daemon for streaming server requests.
- `xray.UnaryClient(host)` - Creates X-Ray subsegments for unary client calls if segment exists in context.
- `xray.StreamClient(host)` - Creates X-Ray subsegments for streaming client calls if segment exists in context.

## Request Processing (Client & Server)

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `TraceIdInterceptor()` - Extracts and injects trace IDs from request objects into context.
- `OptionsInterceptor()` - Adds options and logging context to requests.

**github.com/go-admin-team/go-admin-core/server/grpc/interceptors/request_tag** *Requires go-admin-core ecosystem (github.com/go-admin-team/go-admin-core/tools/utils)*
- `UnaryClientInterceptor()` - Appends request ID to outgoing unary client request metadata.
- `StreamClientInterceptor()` - Appends request ID to outgoing streaming client request metadata.
- `UnaryServerInterceptor()` - Appends request ID to context for unary server requests.
- `StreamServerInterceptor()` - Appends request ID to context for streaming server requests.

**goa.design/goa/v3/grpc/middleware** *Requires Goa framework (goa.design/goa/v3/middleware)*
- `middleware.UnaryRequestID(opts...)` - Initializes request ID in unary gRPC request metadata, optionally using incoming x-request-id.
- `middleware.StreamRequestID(opts...)` - Initializes request ID in streaming gRPC request metadata, optionally using incoming x-request-id.

**github.com/mercari/go-grpc-interceptor/acceptlang/i18n** *Requires go-i18n library (github.com/nicksnyder/go-i18n/i18n)*
- `i18n.UnaryServerInterceptor` - Parses Accept-Language header and provides i18n translation function in context.

## Context Management (Server)

**goa.design/goa/v3/grpc/middleware** *Requires Goa framework (goa.design/goa/v3/middleware)*
- `middleware.StreamCanceler(ctx)` - Gracefully stops streaming requests when context is canceled, cancels all active stream contexts.

## Caching (Client)

**github.com/DeNA/cloud-datastore-interceptor/cache** *Requires Google Cloud Datastore (cloud.google.com/go/datastore)*
- `cache.UnaryClientInterceptor(cacher)` - Caches Cloud Datastore Lookup and GetMulti calls, invalidates cache on Put/Delete operations.

## Request Transformation (Client)

**github.com/DeNA/cloud-datastore-interceptor/transform** *Requires Google Cloud Datastore (cloud.google.com/go/datastore)*
- `transform.QueryToLookupWithKeysOnly()` - Transforms RunQuery requests to Lookup requests with KeysOnly query for cache compatibility.

## Tree Validation & Quota Management (Server)

**github.com/google/trillian/server/interceptor** *Requires Trillian ecosystem (github.com/google/trillian/storage, github.com/google/trillian/quota, github.com/google/trillian/trees)*
- `TrillianInterceptor.UnaryInterceptor()` - Validates tree type and state, enforces quota limits, and manages token refunds for Trillian services.

## Convenience Functions

**github.com/go-coldbrew/interceptors** *Requires coldbrew ecosystem dependencies (go-coldbrew/errors, go-coldbrew/log, go-coldbrew/options, go-coldbrew/tracing)*
- `DefaultClientInterceptor(opts...)` - Chains all default unary client interceptors (Hystrix, retry, OpenTracing, NewRelic, Prometheus).
- `DefaultClientStreamInterceptor(opts...)` - Chains all default stream client interceptors (OpenTracing, NewRelic, Prometheus).
- `DefaultInterceptors()` - Returns array of all default unary server interceptors (logging, tracing, tags, Prometheus, errors, NewRelic, panic recovery).
- `DefaultStreamInterceptors()` - Returns array of all default stream server interceptors (logging, tracing, tags, Prometheus, errors).