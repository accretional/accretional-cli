# Registry and Service Discovery Examples

Commands for managing gRPC service registration, type registration, and service mesh connections.

## Proto File Registration

### Upload and Register Proto File

Upload a `.proto` file to register all message types defined in it:

```bash
./accretional registry upload-proto \
  --file path/to/your.proto \
  --namespace shared \
  --endpoint localhost:50051
```

**With custom import paths:**
```bash
./accretional registry upload-proto \
  --file path/to/your.proto \
  --namespace shared \
  --import-path /path/to/proto/dependencies \
  --endpoint localhost:50051
```

**Output:**
```
✓ Proto file registered successfully
  Proto ID: shared/your.proto
  Registered Messages (3):
    [1] YourMessage
    [2] YourRequest
    [3] YourResponse
```

## Type Management

### List Registered Types

List all registered protobuf message types:

```bash
./accretional registry list-types \
  --namespace shared \
  --endpoint localhost:50051
```

**List all types (all namespaces):**
```bash
./accretional registry list-types \
  --endpoint localhost:50051
```

**Output:**
```
Types (all namespaces):
Total: 5

Note: Type details are stored in system/types collection.
      Use 'collection get-record' to view full type information.

Registered type records:
  [1] Record ID: shared/YourMessage
  [2] Record ID: shared/YourRequest
  ...
```

### Get Type Details

Get detailed information about a specific registered type:

```bash
./accretional registry get-type \
  --name YourMessage \
  --namespace shared \
  --endpoint localhost:50051
```

**Output:**
```
Type: shared/YourMessage
Record ID: shared/YourMessage

Note: Type details are stored as RegisteredProto in system/types collection.
      Use 'collection get-record --collection types --id' to view full details.
```

## Service Registration

### Register Service from Proto File

Register a gRPC service by parsing a proto file:

```bash
./accretional registry register-service \
  --name YourService \
  --proto-file path/to/service.proto \
  --namespace shared \
  --endpoint localhost:50051
```

**With custom import paths:**
```bash
./accretional registry register-service \
  --name YourService \
  --proto-file path/to/service.proto \
  --namespace shared \
  --import-path /path/to/dependencies \
  --endpoint localhost:50051
```

### Register Service Manually

Register a service by providing method names (minimal descriptor):

```bash
./accretional registry register-service \
  --name YourService \
  --methods "Method1,Method2,Method3" \
  --namespace shared \
  --endpoint localhost:50051
```

**Note:** Using `--proto-file` is preferred as it provides complete service descriptor information including input/output types.

**Output:**
```
✓ Service registered successfully
  Service ID: shared/YourService
  Registered Methods (3):
    [1] Method1
    [2] Method2
    [3] Method3
```

## Service Discovery

### List Registered Services

List all registered gRPC services:

```bash
./accretional registry list-services \
  --namespace shared \
  --endpoint localhost:50051
```

**List all services (all namespaces):**
```bash
./accretional registry list-services \
  --endpoint localhost:50051
```

**Output:**
```
Services in namespace 'shared':
Total: 2

[1] shared/YourService
     Methods: [Method1 Method2 Method3]
     ID: shared/YourService

[2] shared/AnotherService
     Methods: [DoSomething DoSomethingElse]
     ID: shared/AnotherService
```

### Get Service Details

Get detailed information about a specific service:

```bash
./accretional registry get-service \
  --name YourService \
  --namespace shared \
  --endpoint localhost:50051
```

**Output:**
```
Service: shared/YourService
ID: shared/YourService

Methods (3):
  [1] Method1
  [2] Method2
  [3] Method3

Service Descriptor:
  Name: YourService
  Methods in descriptor: 3
```

## Connection Management

### Add Service Connection

Register a connection to a service endpoint for service mesh routing:

```bash
./accretional registry add-connection \
  --address localhost:50052 \
  --namespace shared \
  --endpoint localhost:50051
```

**With multiple namespaces:**
```bash
./accretional registry add-connection \
  --address localhost:50052 \
  --namespace "shared,production" \
  --endpoint localhost:50051
```

**Output:**
```
✓ Connection added successfully
  Connection ID: conn-12345
  Address: localhost:50052
  Target Collector ID: 019ba998-d033-7a98-b083-c2a32cf1ef01
  Shared Namespaces: [shared production]
```

### List Connections

List all registered service connections:

```bash
./accretional registry list-connections \
  --namespace shared \
  --endpoint localhost:50051
```

**List all connections:**
```bash
./accretional registry list-connections \
  --endpoint localhost:50051
```

**Output:**
```
Connections in namespace 'shared':
Total: 2

Note: Connection details are stored in system/connections collection.
      Use 'collection get-record --collection connections --id' to view full details.

Connection records:
  [1] Record ID: conn-12345
  [2] Record ID: conn-67890
```

### Remove Connection

Remove a registered service connection:

```bash
./accretional registry remove-connection \
  --address localhost:50052 \
  --endpoint localhost:50051
```

**Output:**
```
✓ Connection removed successfully
  Connection ID: conn-12345
  Address: localhost:50052
```

## Service Discovery via Dispatcher

### Discover Services

Discover available services and their locations:

```bash
./accretional registry discover \
  --namespace shared \
  --endpoint localhost:50051
```

**Discover specific service:**
```bash
./accretional registry discover \
  --namespace shared \
  --service YourService \
  --endpoint localhost:50051
```

**Output:**
```
Discovered in namespace 'shared':
Total: 3

[1] shared/YourService
[2] shared/AnotherService
[3] shared/ThirdService

Note: Use 'registry route' to get routing information for specific services.
```

### Get Routing Information

Get routing information showing which collector handles a specific service:

```bash
./accretional registry route \
  --service YourService \
  --namespace shared \
  --endpoint localhost:50051
```

**With method name:**
```bash
./accretional registry route \
  --service YourService \
  --namespace shared \
  --method Method1 \
  --endpoint localhost:50051
```

**Output:**
```
Routing Information
===================
Service: shared/YourService
Method: Method1
Server Endpoint: localhost:50052
Collection: shared/YourService

Note: This shows which collector/server handles this service.
      For method-specific routing, use the CollectiveDispatcher.Serve method.
```

## Complete Workflow

Here's a complete workflow for registering a service:

```bash
# Step 1: Upload proto file to register message types
./accretional registry upload-proto \
  --file service.proto \
  --namespace shared \
  --endpoint localhost:50051

# Step 2: Register the service
./accretional registry register-service \
  --name MyService \
  --proto-file service.proto \
  --namespace shared \
  --endpoint localhost:50051

# Step 3: Verify service registration
./accretional registry get-service \
  --name MyService \
  --namespace shared \
  --endpoint localhost:50051

# Step 4: Add connection to service endpoint
./accretional registry add-connection \
  --address localhost:50052 \
  --namespace shared \
  --endpoint localhost:50051

# Step 5: Discover and route to the service
./accretional registry discover \
  --namespace shared \
  --endpoint localhost:50051

./accretional registry route \
  --service MyService \
  --namespace shared \
  --endpoint localhost:50051
```

## Common Flags

All registry commands support:
- `--endpoint` - gRPC server endpoint (default: `localhost:50051`)
- `--namespace` - Namespace (default: `shared`)

Additional flags:
- `--import-path` - Additional import path for proto dependencies (for `upload-proto` and `register-service`)
