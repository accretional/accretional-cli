# Accretional CLI Examples

This directory contains example scripts demonstrating how to use the Accretional CLI.

## Files

- `cli_examples.sh` - Comprehensive documentation of all CLI commands with examples

## Usage

### View Examples

To view all available command examples:

```bash
./examples/cli_examples.sh
```

### Run Examples

The script is primarily for documentation. To actually run commands, copy the relevant command from the script and execute it, or modify the script to set your actual endpoint and collection names.

### Quick Start

1. **View the examples:**
   ```bash
   cat examples/cli_examples.sh
   ```

2. **Set your endpoint:**
   ```bash
   export ENDPOINT="localhost:50051"  # or your collector endpoint
   ```

3. **Copy and run commands** from the script output, replacing placeholders with your actual values.

## Command Categories

The examples script covers:

1. **Collection Operations**
   - Create, list, describe, modify, delete collections

2. **Record Operations**
   - Create, get, list, update, delete records
   - Search (FTS, semantic, hybrid)

3. **File Operations - Standalone**
   - Save, get, delete standalone files
   - Text extraction for searchability

4. **File Operations - Attached**
   - Attach, detach, list files attached to records

5. **Complete Workflow**
   - End-to-end example combining all operations

## Notes

- Replace `localhost:50051` with your actual collector endpoint
- Replace `example-collection` with your collection name
- Replace `shared` with your namespace if different
- The `--extract-text` flag is required for file content to be searchable
- Vector search requires `--enable-vector` and appropriate `--vector-dimensions`
