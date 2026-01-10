# Server Management Examples

Commands for checking server health and status.

## Health Check

Check if the Collector server is reachable and responding:

```bash
./accretional server health --endpoint localhost:50051
```

**Output:**
- `✅ HEALTHY` - Server is reachable and responding
- `❌ UNHEALTHY` - Server is reachable but not responding correctly
- `❌ UNREACHABLE` - Server cannot be reached

**Options:**
- `--endpoint` - gRPC server endpoint (default: `localhost:50051`)
- `--timeout` - Timeout for health check (default: `5s`)

**Example with custom timeout:**
```bash
./accretional server health \
  --endpoint localhost:50051 \
  --timeout 10s
```

## Server Status

Get detailed information about the Collector server:

```bash
./accretional server status \
  --endpoint localhost:50051 \
  --namespace shared
```

**Shows:**
- Endpoint information
- Namespace
- Collection count
- List of all collections

**Options:**
- `--endpoint` - gRPC server endpoint (default: `localhost:50051`)
- `--namespace` - Namespace to query (empty for all namespaces)

**Example for all namespaces:**
```bash
./accretional server status \
  --endpoint localhost:50051
```

**Example for specific namespace:**
```bash
./accretional server status \
  --endpoint localhost:50051 \
  --namespace shared
```
