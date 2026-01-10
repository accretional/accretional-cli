#!/bin/bash
# Accretional CLI Examples
# This script demonstrates all available CLI commands with examples
# Note: Replace localhost:50051 with your actual collector endpoint

ENDPOINT="localhost:50051"
NAMESPACE="shared"
COLLECTION="example-collection"

echo "=========================================="
echo "Accretional CLI Examples"
echo "=========================================="
echo ""

# ============================================
# COLLECTION OPERATIONS
# ============================================
echo "=== COLLECTION OPERATIONS ==="
echo ""

echo "# 1. Create a collection with FTS and vector search enabled"
echo "./accretional collection create \\"
echo "  --name $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --enable-fts \\"
echo "  --enable-vector \\"
echo "  --vector-dimensions 384 \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 2. List all collections in a namespace"
echo "./accretional collection list \\"
echo "  --namespace $NAMESPACE \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 3. Describe a collection (get detailed information)"
echo "./accretional collection describe \\"
echo "  --name $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 4. Modify collection settings (add indexed fields)"
echo "./accretional collection modify \\"
echo "  --name $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --indexed-field name \\"
echo "  --indexed-field email \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 5. Delete a collection"
echo "./accretional collection delete \\"
echo "  --name $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo ""
echo "=========================================="
echo ""

# ============================================
# RECORD OPERATIONS
# ============================================
echo "=== RECORD OPERATIONS ==="
echo ""

echo "# 1. Create a record"
echo "./accretional collection create-record \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --id user-001 \\"
echo "  --data '{\"name\":\"Alice\",\"email\":\"alice@example.com\",\"age\":30}' \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 2. Get a record by ID"
echo "./accretional collection get-record \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --id user-001 \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 3. List records (with optional filter)"
echo "./accretional collection list-records \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 3b. List records with filter"
echo "./accretional collection list-records \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --filter '{\"age\":30}' \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 4. Update a record"
echo "./accretional collection update-record \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --id user-001 \\"
echo "  --data '{\"age\":31}' \\"
echo "  --update-mask age \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 5. Delete a record"
echo "./accretional collection delete-record \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --id user-001 \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 6. Search records - Full-text search (FTS)"
echo "./accretional collection search \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --query \"alice\" \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 7. Search records - Semantic search"
echo "./accretional collection search \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --semantic-text \"user email address\" \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 8. Search records - Hybrid search (FTS + semantic)"
echo "./accretional collection search \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --query \"alice\" \\"
echo "  --semantic-text \"contact information\" \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo ""
echo "=========================================="
echo ""

# ============================================
# FILE OPERATIONS - Standalone Files
# ============================================
echo "=== FILE OPERATIONS - Standalone Files ==="
echo ""

echo "# 1. Save a file (without text extraction - faster, not searchable)"
echo "./accretional collection save-file \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --path \"docs/readme.txt\" \\"
echo "  --file ./readme.txt \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 2. Save a file with text extraction (for searchability)"
echo "./accretional collection save-file \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --path \"docs/readme.txt\" \\"
echo "  --file ./readme.txt \\"
echo "  --extract-text \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 3. Get a file"
echo "./accretional collection get-file \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --path \"docs/readme.txt\" \\"
echo "  --output ./downloaded-readme.txt \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 4. Delete a file"
echo "./accretional collection delete-file \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --path \"docs/readme.txt\" \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo ""
echo "=========================================="
echo ""

# ============================================
# FILE OPERATIONS - Attached Files
# ============================================
echo "=== FILE OPERATIONS - Attached Files ==="
echo ""

echo "# 1. Attach a file to a record (without text extraction)"
echo "./accretional collection attach-file \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --record-id user-001 \\"
echo "  --file ./document.pdf \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 2. Attach a file to a record with text extraction (for searchability)"
echo "./accretional collection attach-file \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --record-id user-001 \\"
echo "  --file ./document.txt \\"
echo "  --extract-text \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 3. List files attached to a record"
echo "./accretional collection list-files \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --record-id user-001 \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 4. List standalone files (with optional path prefix filter)"
echo "./accretional collection list-files \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --prefix \"docs/\" \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 5. Detach a file from a record (keep the file record)"
echo "./accretional collection detach-file \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --record-id user-001 \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo "# 6. Detach a file and delete the file record"
echo "./accretional collection detach-file \\"
echo "  --collection $COLLECTION \\"
echo "  --namespace $NAMESPACE \\"
echo "  --record-id user-001 \\"
echo "  --delete-file \\"
echo "  --endpoint $ENDPOINT"
echo ""

echo ""
echo "=========================================="
echo ""

# ============================================
# COMPLETE WORKFLOW EXAMPLE
# ============================================
echo "=== COMPLETE WORKFLOW EXAMPLE ==="
echo ""

cat << 'WORKFLOW_EOF'
# Complete workflow: Create collection, add records, attach files, and search

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

# Step 5: Search for the content
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --query "neural networks" \
  --endpoint localhost:50051

# Step 6: Semantic search
./accretional collection search \
  --collection my-collection \
  --namespace shared \
  --semantic-text "computer vision and image recognition" \
  --endpoint localhost:50051
WORKFLOW_EOF

echo ""
echo "=========================================="
echo ""

# ============================================
# NOTES AND TIPS
# ============================================
echo "=== NOTES AND TIPS ==="
echo ""
echo "1. Text Extraction:"
echo "   - Use --extract-text flag to make file contents searchable"
echo "   - Without this flag, only file metadata (name, path) is searchable"
echo "   - Text extraction works for: .txt, .md, .json, .xml, .html, .csv, .yaml, .proto, code files"
echo "   - PDF, DOCX, DOC, XLSX extraction not yet implemented"
echo ""
echo "2. Search Types:"
echo "   - FTS (--query): Fast keyword-based search"
echo "   - Semantic (--semantic-text): Meaning-based search using embeddings"
echo "   - Hybrid: Combine both for best results"
echo ""
echo "3. File Storage:"
echo "   - Standalone files: Use save-file/get-file/delete-file"
echo "   - Attached files: Use attach-file/detach-file/list-files"
echo "   - Files are stored as Struct records with _type: \"file\""
echo ""
echo "4. Collection Configuration:"
echo "   - --enable-fts: Enables full-text search"
echo "   - --enable-vector: Enables semantic/vector search"
echo "   - --vector-dimensions: Must match your embedding model (e.g., 384, 768, 1536)"
echo ""
echo "5. Namespaces:"
echo "   - Default namespace is 'shared'"
echo "   - Collections are scoped by namespace"
echo ""
echo "=========================================="
echo "End of examples"
echo "=========================================="
