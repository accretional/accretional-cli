# accretional-cli

Command-line interface for interacting with Collector services.

## Installation

```bash
go install github.com/accretional/accretional-cli/cmd/accretional@latest
```

## Usage

### Collection Operations

#### Create a record

```bash
# From JSON string
accretional collection create \
  --endpoint localhost:50051 \
  --namespace shared \
  --collection users \
  --data '{"name": "Alice", "email": "alice@example.com"}'

# From JSON file
accretional collection create \
  --endpoint localhost:50051 \
  --collection users \
  --file ./user.json \
  --id user-123
```

#### Get a record

```bash
accretional collection get \
  --endpoint localhost:50051 \
  --collection users \
  --id user-123

# Save to file
accretional collection get \
  --endpoint localhost:50051 \
  --collection users \
  --id user-123 \
  --output ./user.json
```

#### List records

```bash
# List all records
accretional collection list \
  --endpoint localhost:50051 \
  --collection users

# With filter
accretional collection list \
  --endpoint localhost:50051 \
  --collection users \
  --filter '{"status": "active"}' \
  --page-size 20
```

#### Search records

```bash
accretional collection search \
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
./accretional file --help
```

### Dependencies

- `github.com/accretional/collector` - For proto definitions
- `github.com/grpc-ecosystem/go-grpc-middleware/v2` - Client middleware
- `github.com/spf13/cobra` - CLI framework

## Commands

- `collection create` - Create a new record
- `collection get` - Get a record by ID
- `collection list` - List records with optional filtering
- `collection search` - Full-text search records

All commands support:
- `--endpoint` - gRPC server endpoint (default: localhost:50051)
- `--namespace` - Namespace (default: shared)
- `--collection` - Collection name (required)
