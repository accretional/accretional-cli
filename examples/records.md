# Record Operations Examples

Commands for managing records within collections.

## Create a Record

```bash
./accretional collection create-record \
  --collection my-collection \
  --namespace shared \
  --id user-001 \
  --data '{"name":"Alice","email":"alice@example.com","age":30}' \
  --endpoint localhost:50051
```

**Options:**
- `--collection` - Collection name (required)
- `--namespace` - Namespace (default: `shared`)
- `--id` - Record ID (required)
- `--data` - JSON data as string (required)
- `--file` - JSON data from file (alternative to `--data`)

**Example using a JSON file:**
```bash
./accretional collection create-record \
  --collection my-collection \
  --namespace shared \
  --id user-001 \
  --file ./user.json \
  --endpoint localhost:50051
```

## Get a Record

```bash
./accretional collection get-record \
  --collection my-collection \
  --namespace shared \
  --id user-001 \
  --endpoint localhost:50051
```

## List Records

List all records in a collection:

```bash
./accretional collection list-records \
  --collection my-collection \
  --namespace shared \
  --endpoint localhost:50051
```

**List records with a filter:**
```bash
./accretional collection list-records \
  --collection my-collection \
  --namespace shared \
  --filter '{"age":30}' \
  --endpoint localhost:50051
```

**Options:**
- `--filter` - JSON filter object
- `--page-size` - Number of records per page
- `--page-token` - Page token for pagination

## Update a Record

```bash
./accretional collection update-record \
  --collection my-collection \
  --namespace shared \
  --id user-001 \
  --data '{"age":31}' \
  --update-mask age \
  --endpoint localhost:50051
```

**Options:**
- `--data` - JSON data with fields to update
- `--update-mask` - Comma-separated list of fields to update (optional, updates all provided fields if not specified)

**Example updating multiple fields:**
```bash
./accretional collection update-record \
  --collection my-collection \
  --namespace shared \
  --id user-001 \
  --data '{"age":31,"email":"newemail@example.com"}' \
  --update-mask age,email \
  --endpoint localhost:50051
```

## Delete a Record

```bash
./accretional collection delete-record \
  --collection my-collection \
  --namespace shared \
  --id user-001 \
  --endpoint localhost:50051
```

## Search Records

### Full-Text Search (FTS)

```bash
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --query "alice" \
  --endpoint localhost:50051
```

**Requires:** Collection with `--enable-fts` enabled

### Semantic Search

```bash
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --semantic-text "user email address" \
  --endpoint localhost:50051
```

**Requires:** Collection with `--enable-vector` enabled

**Options:**
- `--semantic-text` - Text to search by semantic meaning
- `--similarity-threshold` - Minimum similarity threshold (0.0-1.0)

### Hybrid Search (FTS + Semantic)

```bash
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --query "alice" \
  --semantic-text "contact information" \
  --endpoint localhost:50051
```

**Combines both search types for best results**

**Common Search Options:**
- `--query` - Full-text search query
- `--semantic-text` - Semantic search query
- `--limit` - Maximum number of results (default: 10)
- `--similarity-threshold` - Minimum similarity for semantic search (0.0-1.0)

**Example with custom limit:**
```bash
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --query "machine learning" \
  --limit 20 \
  --endpoint localhost:50051
```
