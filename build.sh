#!/bin/bash
# Build script for AuraSpeed (Bash/Zsh/Fish compatible)
# Builds for Windows (amd64), Linux (amd64), macOS (amd64, arm64)

set -e

PROJECT_NAME="auraspeed"
OUTPUT_DIR="dist"

# Get version from git tag or use "dev"
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")

# Get commit hash
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo -e "\033[36mBuilding $PROJECT_NAME...\033[0m"
echo -e "Version: \033[36m$VERSION\033[0m"
echo -e "Commit:  \033[36m$COMMIT\033[0m"
echo ""

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Build targets
declare -a TARGETS=(
    "windows:amd64:${PROJECT_NAME}-windows-amd64.exe"
    "linux:amd64:${PROJECT_NAME}-linux-amd64"
    "darwin:amd64:${PROJECT_NAME}-darwin-amd64"
    "darwin:arm64:${PROJECT_NAME}-darwin-arm64"
)

# LDFLAGS for version info
LDFLAGS="-s -w -X 'auraspeed/cmd/root.Version=${VERSION}' -X 'auraspeed/cmd/root.Commit=${COMMIT}' -X 'auraspeed/cmd/root.BuildTime=${BUILD_TIME}'"

for TARGET in "${TARGETS[@]}"; do
    IFS=':' read -r GOOS GOARCH OUTPUT <<< "$TARGET"

    echo -e "Building for \033[33m${GOOS}/${GOARCH}\033[0m..."

    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$LDFLAGS" -o "${OUTPUT_DIR}/${OUTPUT}" ./cmd/main.go

    if [ -f "${OUTPUT_DIR}/${OUTPUT}" ]; then
        SIZE=$(du -k "${OUTPUT_DIR}/${OUTPUT}" | cut -f1)
        echo -e "  \033[32m✓ Created: ${OUTPUT} (${SIZE} KB)\033[0m"
    else
        echo -e "  \033[31m✗ Failed: Output file not found\033[0m"
    fi
done

echo ""
echo -e "\033[32mBuild complete! Output directory: ${OUTPUT_DIR}\033[0m"
ls -lh "$OUTPUT_DIR"
