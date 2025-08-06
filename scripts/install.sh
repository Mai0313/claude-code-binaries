#!/bin/bash

set -e

# Parse command line arguments
TARGET="$1"  # Optional target parameter

# Validate target if provided
if [[ -n "$TARGET" ]] && [[ ! "$TARGET" =~ ^(stable|latest|[0-9]+\.[0-9]+\.[0-9]+(-[^[:space:]]+)?)$ ]]; then
    echo "Usage: $0 [stable|latest|VERSION]" >&2
    exit 1
fi

GCS_BUCKET="https://storage.googleapis.com/claude-code-dist-86c565f3-f756-42ad-8dfa-d59b1c096819/claude-code-releases"
DOWNLOAD_DIR="$(pwd)/binaries"
# DOWNLOAD_DIR="$HOME/.claude/downloads"

# Check for required dependencies
if ! command -v curl >/dev/null 2>&1; then
    echo "curl is required but not installed" >&2
    exit 1
fi

# Check if jq is available (optional)
HAS_JQ=false
if command -v jq >/dev/null 2>&1; then
    HAS_JQ=true
fi

# Simple JSON parser for extracting checksum when jq is not available
get_checksum_from_manifest() {
    local json="$1"
    local platform="$2"

    # Normalize JSON to single line and extract checksum
    json=$(echo "$json" | tr -d '\n\r\t' | sed 's/ \+/ /g')

    # Extract checksum for platform using bash regex
    if [[ $json =~ \"$platform\"[^}]*\"checksum\"[[:space:]]*:[[:space:]]*\"([a-f0-9]{64})\" ]]; then
        echo "${BASH_REMATCH[1]}"
        return 0
    fi

    return 1
}

# Detect platform
case "$(uname -s)" in
    Darwin) os="darwin" ;;
    Linux) os="linux" ;;
    *) echo "Windows is not supported" >&2; exit 1 ;;
esac

case "$(uname -m)" in
    x86_64|amd64) arch="x64" ;;
    arm64|aarch64) arch="arm64" ;;
    *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

platform="${os}-${arch}"
mkdir -p "$DOWNLOAD_DIR"

# Always download stable version (which has the most up-to-date installer)
version=$(curl -fsSL "$GCS_BUCKET/stable")

# Download manifest and extract checksum
manifest_json=$(curl -fsSL "$GCS_BUCKET/$version/manifest.json")

# Use jq if available, otherwise fall back to pure bash parsing
if [ "$HAS_JQ" = true ]; then
    checksum=$(echo "$manifest_json" | jq -r ".platforms[\"$platform\"].checksum // empty")
else
    checksum=$(get_checksum_from_manifest "$manifest_json" "$platform")
fi

# Validate checksum format (SHA256 = 64 hex characters)
if [ -z "$checksum" ] || [[ ! "$checksum" =~ ^[a-f0-9]{64}$ ]]; then
    echo "Platform $platform not found in manifest" >&2
    exit 1
fi

# Download and verify
binary_path="$DOWNLOAD_DIR/claude-$version-$platform"
if ! curl -fsSL -o "$binary_path" "$GCS_BUCKET/$version/$platform/claude"; then
    echo "Download failed" >&2
    rm -f "$binary_path"
    exit 1
fi

# Pick the right checksum tool
if [ "$os" = "darwin" ]; then
    actual=$(shasum -a 256 "$binary_path" | cut -d' ' -f1)
else
    actual=$(sha256sum "$binary_path" | cut -d' ' -f1)
fi

if [ "$actual" != "$checksum" ]; then
    echo "Checksum verification failed" >&2
    rm -f "$binary_path"
    exit 1
fi

chmod +x "$binary_path"
echo "Claude Code has been installed: $binary_path"
