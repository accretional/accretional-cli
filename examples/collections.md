# Collection Operations Examples

Commands for managing collections in Collector.

## Create a Collection

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
- `--name` - Collection name (required)
- `--namespace` - Namespace (default: `shared`)
- `--enable-fts` - Enable full-text search
- `--enable-vector` - Enable semantic/vector search
- `--vector-dimensions` - Vector dimension size (must match your embedding model)
- `--indexed-field` - Add indexed fields (can be specified multiple times)
- `--type-name` - Message type name (default: `Struct` for JSON data)
- `--type-namespace` - Message type namespace (default: `google.protobuf`)

**Example with indexed fields:**
```bash
./accretional collection create \
  --name my-collection \
  --namespace shared \
  --enable-fts \
  --indexed-field name \
  --indexed-field email \
  --endpoint localhost:50051
```

## List Collections

List all collections in a namespace:

```bash
./accretional collection list \
  --namespace shared \
  --endpoint localhost:50051
```

**Options:**
- `--namespace` - Namespace (empty for all namespaces)
- `--page-size` - Number of collections per page (default: 20)
- `--page-token` - Page token for pagination

## Describe a Collection

Get detailed information about a collection:

```bash
./accretional collection describe \
  --name my-collection \
  --namespace shared \
  --endpoint localhost:50051
```

**Shows:**
- Collection name and namespace
- Message type information
- Search configuration (FTS, vector search settings)
- Indexed fields

## Modify Collection Settings

Add indexed fields to a collection:

```bash
./accretional collection modify \
  --name my-collection \
  --namespace shared \
  --indexed-field name \
  --indexed-field email \
  --endpoint localhost:50051
```

**Options:**
- `--indexed-field` - Field to index (can be specified multiple times)

## Delete a Collection

```bash
./accretional collection delete \
  --name my-collection \
  --namespace shared \
  --endpoint localhost:50051
```

**Note:** This permanently deletes the collection and all its records.
