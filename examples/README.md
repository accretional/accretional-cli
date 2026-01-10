# Accretional CLI Examples

This document provides comprehensive examples of all available CLI commands.

**Note:** Replace `localhost:50051` with your actual collector endpoint, and adjust collection/namespace names as needed.

---

## Server Management

### Health Check

Check if the Collector server is reachable and responding:

```bash
./accretional server health --endpoint localhost:50051
```

**Output:**
- `✅ HEALTHY` - Server is reachable and responding
- `❌ UNHEALTHY` - Server is reachable but not responding correctly
- `❌ UNREACHABLE` - Server cannot be reached

### Server Status

Get detailed information about the Collector server:

```bash
./accretional server status \
  --endpoint localhost:50051 \
  --namespace shared
```

Shows:
- Endpoint information
- Namespace
- Collection count
- List of all collections

---

## Collection Operations

### Create a Collection

Create a new collection with FTS and vector search enabled:

```bash
./accretional collection create \
  --name my-collection \
  --namespace shared \
  --enable-fts \
  --enable-vector \
  --vector-dimensions 384 \
  --endpoint localhost:50051
```

**Flags:**
- `--enable-fts` - Enable full-text search
- `--enable-vector` - Enable semantic/vector search
- `--vector-dimensions` - Vector dimension size (must match your embedding model)
- `--indexed-field` - Add indexed fields (can be specified multiple times)

### List Collections

List all collections in a namespace:

```bash
./accretional collection list \
  --namespace shared \
  --endpoint localhost:50051
```

### Describe a Collection

Get detailed information about a collection:

```bash
./accretional collection describe \
  --name my-collection \
  --namespace shared \
  --endpoint localhost:50051
```

### Modify Collection Settings

Add indexed fields to a collection:

```bash
./accretional collection modify \
  --name my-collection \
  --namespace shared \
  --indexed-field name \
  --indexed-field email \
  --endpoint localhost:50051
```

### Delete a Collection

```bash
./accretional collection delete \
  --name my-collection \
  --namespace shared \
  --endpoint localhost:50051
```

---

## Record Operations

### Create a Record

```bash
./accretional collection create-record \
  --collection my-collection \
  --namespace shared \
  --id user-001 \
  --data '{"name":"Alice","email":"alice@example.com","age":30}' \
  --endpoint localhost:50051
```

### Get a Record

```bash
./accretional collection get-record \
  --collection my-collection \
  --namespace shared \
  --id user-001 \
  --endpoint localhost:50051
```

### List Records

List all records in a collection:

```bash
./accretional collection list-records \
  --collection my-collection \
  --namespace shared \
  --endpoint localhost:50051
```

List records with a filter:

```bash
./accretional collection list-records \
  --collection my-collection \
  --namespace shared \
  --filter '{"age":30}' \
  --endpoint localhost:50051
```

### Update a Record

```bash
./accretional collection update-record \
  --collection my-collection \
  --namespace shared \
  --id user-001 \
  --data '{"age":31}' \
  --update-mask age \
  --endpoint localhost:50051
```

### Delete a Record

```bash
./accretional collection delete-record \
  --collection my-collection \
  --namespace shared \
  --id user-001 \
  --endpoint localhost:50051
```

### Search Records

#### Full-Text Search (FTS)

```bash
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --query "alice" \
  --endpoint localhost:50051
```

#### Semantic Search

```bash
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --semantic-text "user email address" \
  --endpoint localhost:50051
```

#### Hybrid Search (FTS + Semantic)

```bash
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --query "alice" \
  --semantic-text "contact information" \
  --endpoint localhost:50051
```

**Search Flags:**
- `--query` - Full-text search query
- `--semantic-text` - Semantic search query (requires vector search enabled)
- `--limit` - Maximum number of results (default: 10)
- `--similarity-threshold` - Minimum similarity for semantic search (0.0-1.0)

---

## File Operations - Standalone Files

### Save a File

Save a file without text extraction (faster, not searchable):

```bash
./accretional collection save-file \
  --collection my-collection \
  --namespace shared \
  --path "docs/readme.txt" \
  --file ./readme.txt \
  --endpoint localhost:50051
```

Save a file with text extraction (for searchability):

```bash
./accretional collection save-file \
  --collection my-collection \
  --namespace shared \
  --path "docs/readme.txt" \
  --file ./readme.txt \
  --extract-text \
  --endpoint localhost:50051
```

### Get a File

```bash
./accretional collection get-file \
  --collection my-collection \
  --namespace shared \
  --path "docs/readme.txt" \
  --output ./downloaded-readme.txt \
  --endpoint localhost:50051
```

### Delete a File

```bash
./accretional collection delete-file \
  --collection my-collection \
  --namespace shared \
  --path "docs/readme.txt" \
  --endpoint localhost:50051
```

---

## File Operations - Attached Files

### Attach a File to a Record

Attach without text extraction:

```bash
./accretional collection attach-file \
  --collection my-collection \
  --namespace shared \
  --record-id user-001 \
  --file ./document.pdf \
  --endpoint localhost:50051
```

Attach with text extraction (for searchability):

```bash
./accretional collection attach-file \
  --collection my-collection \
  --namespace shared \
  --record-id user-001 \
  --file ./document.txt \
  --extract-text \
  --endpoint localhost:50051
```

### List Files

List files attached to a record:

```bash
./accretional collection list-files \
  --collection my-collection \
  --namespace shared \
  --record-id user-001 \
  --endpoint localhost:50051
```

List standalone files (with optional path prefix filter):

```bash
./accretional collection list-files \
  --collection my-collection \
  --namespace shared \
  --prefix "docs/" \
  --endpoint localhost:50051
```

### Detach a File

Detach a file from a record (keeps the file record):

```bash
./accretional collection detach-file \
  --collection my-collection \
  --namespace shared \
  --record-id user-001 \
  --endpoint localhost:50051
```

Detach a file and delete the file record:

```bash
./accretional collection detach-file \
  --collection my-collection \
  --namespace shared \
  --record-id user-001 \
  --delete-file \
  --endpoint localhost:50051
```

---

## Complete Workflow Example

Here's a complete end-to-end workflow:

```bash
# Step 1: Create a collection with search capabilities
./accretional collection create \
  --name my-collection \
  --namespace shared \
  --enable-fts \
  --enable-vector \
  --vector-dimensions 384 \
  --endpoint localhost:50051

# Step 2: Create a record
./accretional collection create-record \
  --collection my-collection \
  --namespace shared \
  --id doc-001 \
  --data '{"title":"AI Research","author":"Dr. Smith","topic":"machine learning"}' \
  --endpoint localhost:50051

# Step 3: Create a text file
cat > research-notes.txt << EOF
This document discusses recent advances in artificial intelligence.
Topics include neural networks, deep learning, and natural language processing.
The research shows promising results in computer vision applications.
EOF

# Step 4: Attach the file with text extraction (makes it searchable)
./accretional collection attach-file \
  --collection my-collection \
  --namespace shared \
  --record-id doc-001 \
  --file research-notes.txt \
  --extract-text \
  --endpoint localhost:50051

# Step 5: Search for the content using FTS
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --query "neural networks" \
  --endpoint localhost:50051

# Step 6: Search using semantic search
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --semantic-text "computer vision and image recognition" \
  --endpoint localhost:50051
```

---

## Notes and Tips

### Text Extraction

- **Use `--extract-text` flag** to make file contents searchable
- **Without this flag**, only file metadata (name, path) is searchable
- **Supported formats**: `.txt`, `.md`, `.json`, `.xml`, `.html`, `.csv`, `.yaml`, `.proto`, code files (`.go`, `.js`, `.py`, etc.)
- **Not yet implemented**: PDF, DOCX, DOC, XLSX extraction

### Search Types

- **FTS (`--query`)**: Fast keyword-based search
- **Semantic (`--semantic-text`)**: Meaning-based search using embeddings
- **Hybrid**: Combine both for best results

### File Storage

- **Standalone files**: Use `save-file`/`get-file`/`delete-file`
- **Attached files**: Use `attach-file`/`detach-file`/`list-files`
- Files are stored as Struct records with `_type: "file"`

### Collection Configuration

- `--enable-fts`: Enables full-text search
- `--enable-vector`: Enables semantic/vector search
- `--vector-dimensions`: Must match your embedding model (e.g., 384, 768, 1536)

### Namespaces

- Default namespace is `shared`
- Collections are scoped by namespace
- Leave namespace empty to query all namespaces (where supported)

### Common Flags

Most commands support:
- `--endpoint` - gRPC server endpoint (default: `localhost:50051`)
- `--namespace` - Namespace (default: `shared`)

---

## Quick Reference

| Command | Purpose |
|---------|---------|
| `server health` | Check server connectivity |
| `server status` | Get server information |
| `collection create` | Create a new collection |
| `collection list` | List collections |
| `collection describe` | Get collection details |
| `collection modify` | Modify collection settings |
| `collection delete` | Delete a collection |
| `collection create-record` | Create a record |
| `collection get-record` | Get a record by ID |
| `collection list-records` | List records |
| `collection update-record` | Update a record |
| `collection delete-record` | Delete a record |
| `collection search` | Search records (FTS/semantic/hybrid) |
| `collection save-file` | Save a standalone file |
| `collection get-file` | Get a standalone file |
| `collection delete-file` | Delete a standalone file |
| `collection attach-file` | Attach a file to a record |
| `collection detach-file` | Detach a file from a record |
| `collection list-files` | List files (attached or standalone) |
