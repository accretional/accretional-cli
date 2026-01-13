# accretional-cli

Command-line interface for interacting with Collector services.

## Installation

```bash
go install github.com/accretional/accretional-cli/cmd/accretional@latest
```

## Usage

### Create a collection

```bash
accretional create collection \
  --endpoint localhost:50051 \
  --namespace shared \
  --name users \
  --indexed-field email \
  --indexed-field name
```

### Create a record

```bash
# From JSON string
accretional create record \
  --endpoint localhost:50051 \
  --namespace shared \
  --collection users \
  --data '{"name": "Alice", "email": "alice@example.com"}'

# From JSON file
accretional create record \
  --endpoint localhost:50051 \
  --collection users \
  --file ./user.json \
  --id user-123
```

### Get a record

```bash
accretional get record \
  --endpoint localhost:50051 \
  --collection users \
  --id user-123

# Save to file
accretional get record \
  --endpoint localhost:50051 \
  --collection users \
  --id user-123 \
  --output ./user.json
```

### List records

```bash
# List all records
accretional list records \
  --endpoint localhost:50051 \
  --collection users

# With filter
accretional list records \
  --endpoint localhost:50051 \
  --collection users \
  --filter '{"status": "active"}' \
  --page-size 20
```

### Search records

```bash
accretional search records \
  --endpoint localhost:50051 \
  --collection users \
  --query "alice" \
  --limit 10
```

## Development

### Setup

```bash
# Install dependencies
go mod tidy

# Build
go build ./cmd/accretional

# Run
./accretional --help
```

### Dependencies

- `github.com/accretional/collector` - For proto definitions
- `github.com/grpc-ecosystem/go-grpc-middleware/v2` - Client middleware
- `github.com/spf13/cobra` - CLI framework

## Commands

- `create collection` - Create a new collection
- `create record` - Create a new record in a collection
- `get record` - Get a record by ID
- `list records` - List records with optional filtering
- `search records` - Full-text search records

All commands support:
- `--endpoint` - gRPC server endpoint (default: localhost:50051)
- `--namespace` - Namespace (default: shared)
- `--collection` - Collection name (required for record operations)
