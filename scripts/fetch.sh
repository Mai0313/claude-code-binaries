#!/usr/bin/env bash
# Everything this mirror knows about its upstream lives in this file. To mirror
# another project, rewrite these three commands; the updater workflow publishes
# whatever `download` leaves in the directory.
#
#   fetch.sh version               print the newest upstream version
#   fetch.sh download VERSION DIR  download and verify every file of VERSION into DIR
#   fetch.sh notes VERSION         print the release notes of VERSION

set -euo pipefail

BASE_URL="https://downloads.claude.ai/claude-code-releases"
KEY_URL="https://downloads.claude.ai/keys/claude-code.asc"
CHANGELOG_URL="https://raw.githubusercontent.com/anthropics/claude-code/main/CHANGELOG.md"

version() {
    local latest
    latest=$(curl -fsSL "$BASE_URL/latest")
    # An error page served with 200 must not become a tag.
    if [[ ! $latest =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[^[:space:]]+)?$ ]]; then
        echo "Unexpected upstream version: $latest" >&2
        return 1
    fi
    echo "$latest"
}

download() {
    local version=$1 dir=$2 file platform binary checksum out
    mkdir -p "$dir"
    curl -fsSL -o "$dir/claude-code.asc" "$KEY_URL"
    for file in manifest.json manifest.json.sig manifest.json.raw-sig.json \
        manifest.zst.json manifest.zst.json.sig manifest.zst.json.raw-sig.json; do
        curl -fsSL -o "$dir/$file" "$BASE_URL/$version/$file"
    done
    # manifest.json lists the plain binaries, manifest.zst.json their zstd copies.
    for file in manifest.json manifest.zst.json; do
        jq -r '.platforms | to_entries[] | "\(.key) \(.value.binary) \(.value.checksum)"' "$dir/$file" |
            while read -r platform binary checksum; do
                out="$dir/claude-$version-$platform${binary#claude}"
                echo "Downloading $platform/$binary"
                curl -fsSL -o "$out" "$BASE_URL/$version/$platform/$binary"
                echo "$checksum  $out" | sha256sum -c --quiet
            done
    done
}

notes() {
    local version=$1 section
    # Reads the whole file: stopping early would kill curl with SIGPIPE.
    section=$(curl -fsSL "$CHANGELOG_URL" | awk -v heading="## $version" '/^## / { found = ($0 == heading); next } found')
    echo "${section:-See https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md}"
}

case "${1:-}" in
    version | download | notes) "$@" ;;
    *)
        echo "Usage: $0 version | download VERSION DIR | notes VERSION" >&2
        exit 1
        ;;
esac
