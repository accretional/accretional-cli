# Accretional CLI Examples

This directory contains comprehensive examples of all available CLI commands, organized by functionality.

**Note:** Replace `localhost:50051` with your actual collector endpoint, and adjust collection/namespace names as needed.

## Documentation Files

- **[Server Management](server.md)** - Health checks and server status
- **[Collection Operations](collections.md)** - Create, list, describe, modify, and delete collections
- **[Record Operations](records.md)** - CRUD operations and search for records
- **[File Operations](files.md)** - Standalone files and file attachments
- **[Registry & Service Discovery](registry.md)** - Service registration, type management, and service mesh connections

## Quick Start

1. **Check server health:**
   ```bash
   ./accretional server health --endpoint localhost:50051
   ```

2. **Create a collection:**
   ```bash
   ./accretional collection create \
     --name my-collection \
     --namespace shared \
     --enable-fts \
     --enable-vector \
     --vector-dimensions 384 \
     --endpoint localhost:50051
   ```

3. **Create a record:**
   ```bash
   ./accretional collection create-record \
     --collection my-collection \
     --namespace shared \
     --id doc-001 \
     --data '{"title":"Example","content":"This is a test"}' \
     --endpoint localhost:50051
   ```

## Complete Workflow

See the individual documentation files for detailed examples, or follow this complete workflow:

```bash
# Step 1: Check server health
./accretional server health --endpoint localhost:50051

# Step 2: Create a collection with search capabilities
./accretional collection create \
  --name my-collection \
  --namespace shared \
  --enable-fts \
  --enable-vector \
  --vector-dimensions 384 \
  --endpoint localhost:50051

# Step 3: Create a record
./accretional collection create-record \
  --collection my-collection \
  --namespace shared \
  --id doc-001 \
  --data '{"title":"AI Research","author":"Dr. Smith","topic":"machine learning"}' \
  --endpoint localhost:50051

# Step 4: Attach a file with text extraction
./accretional collection attach-file \
  --collection my-collection \
  --namespace shared \
  --record-id doc-001 \
  --file research-notes.txt \
  --extract-text \
  --endpoint localhost:50051

# Step 5: Search for content
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --query "neural networks" \
  --endpoint localhost:50051
```

## Common Flags

Most commands support:
- `--endpoint` - gRPC server endpoint (default: `localhost:50051`)
- `--namespace` - Namespace (default: `shared`)

## Quick Reference

| Category | Commands |
|----------|----------|
| **Server** | `health`, `status` |
| **Collections** | `create`, `list`, `describe`, `modify`, `delete` |
| **Records** | `create-record`, `get-record`, `list-records`, `update-record`, `delete-record`, `search` |
| **Files** | `save-file`, `get-file`, `delete-file`, `attach-file`, `detach-file`, `list-files` |
| **Registry** | `upload-proto`, `list-types`, `get-type`, `register-service`, `list-services`, `get-service`, `add-connection`, `list-connections`, `remove-connection`, `discover`, `route` |

For detailed examples, see the individual documentation files linked above.
