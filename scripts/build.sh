#!/bin/bash

# Build script for Accretional CLI
# This script ensures all dependencies are fetched and builds the CLI

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo -e "${GREEN}Building Accretional CLI...${NC}"
echo "Project root: $PROJECT_ROOT"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    echo "Please install Go from https://go.dev/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "Go version: $GO_VERSION"
echo ""

# Change to project root
cd "$PROJECT_ROOT"

# Check if collector dependency exists (local replace)
COLLECTOR_PATH="../collector"
if [ ! -d "$COLLECTOR_PATH" ]; then
    echo -e "${YELLOW}Warning: Collector repository not found at $COLLECTOR_PATH${NC}"
    echo "The CLI depends on a local collector repository."
    echo "Expected structure:"
    echo "  ../collector/  (collector repository)"
    echo "  ./accretional-cli/  (this repository)"
    echo ""
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Build cancelled."
        exit 1
    fi
else
    echo -e "${GREEN}✓ Collector repository found${NC}"
fi

# Step 1: Download dependencies
echo ""
echo -e "${GREEN}[1/4] Downloading dependencies...${NC}"
go mod download
if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to download dependencies${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Dependencies downloaded${NC}"

# Step 2: Tidy modules
echo ""
echo -e "${GREEN}[2/4] Tidying modules...${NC}"
go mod tidy
if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to tidy modules${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Modules tidied${NC}"

# Step 3: Verify dependencies
echo ""
echo -e "${GREEN}[3/4] Verifying dependencies...${NC}"
go mod verify
if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Dependency verification failed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Dependencies verified${NC}"

# Step 4: Build the CLI
echo ""
echo -e "${GREEN}[4/4] Building CLI...${NC}"
OUTPUT_BINARY="$PROJECT_ROOT/accretional"
if [[ "$OS" == "Windows_NT" ]] || [[ "$(uname)" == "MINGW"* ]] || [[ "$(uname)" == "MSYS"* ]]; then
    OUTPUT_BINARY="$PROJECT_ROOT/accretional.exe"
fi

go build -o "$OUTPUT_BINARY" ./cmd/accretional
if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Build failed${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Build successful!${NC}"
echo ""
echo "Binary location: $OUTPUT_BINARY"

# Check if binary is executable
if [ -x "$OUTPUT_BINARY" ]; then
    echo -e "${GREEN}✓ Binary is executable${NC}"
    
    # Show binary info
    echo ""
    echo "Binary information:"
    ls -lh "$OUTPUT_BINARY"
    
    # Test that it runs
    echo ""
    echo "Testing binary..."
    "$OUTPUT_BINARY" --version 2>/dev/null || "$OUTPUT_BINARY" --help > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Binary runs successfully${NC}"
    fi
else
    echo -e "${YELLOW}Warning: Binary is not executable${NC}"
    chmod +x "$OUTPUT_BINARY" 2>/dev/null || true
fi

echo ""
echo -e "${GREEN}Build complete!${NC}"
echo ""
echo "To use the CLI:"
echo "  $OUTPUT_BINARY --help"
echo ""
echo "Or add to PATH:"
echo "  export PATH=\$PATH:$PROJECT_ROOT"
