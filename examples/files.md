# File Operations Examples

Commands for managing files in collections, both as standalone files and as attachments to records.

## Standalone Files

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

**Save a file with text extraction (for searchability):**

```bash
./accretional collection save-file \
  --collection my-collection \
  --namespace shared \
  --path "docs/readme.txt" \
  --file ./readme.txt \
  --extract-text \
  --endpoint localhost:50051
```

**Options:**
- `--collection` - Collection name (required)
- `--namespace` - Namespace (default: `shared`)
- `--path` - File path in collection (required)
- `--file` - Local file path to upload (required)
- `--extract-text` - Extract text content for searchability (FTS and semantic search)

**Supported file types for text extraction:**
- Text files: `.txt`, `.md`, `.log`
- Data formats: `.json`, `.xml`, `.html`, `.csv`, `.yaml`, `.yml`
- Code files: `.go`, `.js`, `.ts`, `.py`, `.java`, `.cpp`, `.c`, `.h`, `.rs`, `.rb`, `.php`, `.sh`, `.bash`, `.proto`
- **Not yet implemented**: PDF, DOCX, DOC, XLSX

### Get a File

```bash
./accretional collection get-file \
  --collection my-collection \
  --namespace shared \
  --path "docs/readme.txt" \
  --output ./downloaded-readme.txt \
  --endpoint localhost:50051
```

**Options:**
- `--path` - File path in collection (required)
- `--output` - Output file path (optional, uses original filename if not provided)

### Delete a File

```bash
./accretional collection delete-file \
  --collection my-collection \
  --namespace shared \
  --path "docs/readme.txt" \
  --endpoint localhost:50051
```

## Attached Files

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

**Attach with text extraction (for searchability):**

```bash
./accretional collection attach-file \
  --collection my-collection \
  --namespace shared \
  --record-id user-001 \
  --file ./document.txt \
  --extract-text \
  --endpoint localhost:50051
```

**Options:**
- `--collection` - Collection name (required)
- `--namespace` - Namespace (default: `shared`)
- `--record-id` - Record ID to attach file to (required)
- `--file` - Local file path to attach (required)
- `--path` - File path in collection (optional, auto-generated if not provided)
- `--extract-text` - Extract text content for searchability

**How it works:**
1. Stores the file as a record with `_type: "file"`
2. Updates the target record's `data_uri` field to reference the file record
3. If `--extract-text` is used, extracts and stores text content in a `content` field

### List Files

**List files attached to a record:**
```bash
./accretional collection list-files \
  --collection my-collection \
  --namespace shared \
  --record-id user-001 \
  --endpoint localhost:50051
```

**List standalone files (with optional path prefix filter):**
```bash
./accretional collection list-files \
  --collection my-collection \
  --namespace shared \
  --prefix "docs/" \
  --endpoint localhost:50051
```

**Options:**
- `--record-id` - List files attached to this record
- `--prefix` - Path prefix for filtering standalone files

### Detach a File

**Detach a file from a record (keeps the file record):**
```bash
./accretional collection detach-file \
  --collection my-collection \
  --namespace shared \
  --record-id user-001 \
  --endpoint localhost:50051
```

**Detach a file and delete the file record:**
```bash
./accretional collection detach-file \
  --collection my-collection \
  --namespace shared \
  --record-id user-001 \
  --delete-file \
  --endpoint localhost:50051
```

**Options:**
- `--record-id` - Record ID to detach file from (required)
- `--file-path` - Specific file path to detach (optional)
- `--delete-file` - Also delete the file record (not just detach)

## File Storage Details

### How Files Are Stored

Files are stored as Struct records with the following structure:
```json
{
  "_type": "file",
  "name": "filename.txt",
  "path": "docs/filename.txt",
  "size": 1234,
  "data": "<base64-encoded-file-content>",
  "mimeType": "text/plain",
  "content": "<extracted-text-if-extract-text-flag-used>"
}
```

### Text Extraction

- **With `--extract-text`**: File content is extracted and stored in the `content` field, making it searchable via FTS and semantic search
- **Without `--extract-text`**: Only file metadata (name, path, size) is stored; content is not searchable

### File Record IDs

- Standalone files: `_file:<path>` (e.g., `_file:docs/readme.txt`)
- Attached files: `_file:attachments/<record-id>/<filename>`

### Searching File Content

Once files are stored with `--extract-text`, you can search their content:

```bash
# Search for text in file contents
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --query "neural networks" \
  --endpoint localhost:50051

# Semantic search on file contents
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --semantic-text "artificial intelligence" \
  --endpoint localhost:50051
```

This will return both regular records and file records that match the search query.
