#!/usr/bin/env bash

set -euo pipefail

VERSION="${1:-$(git describe --tags --always --dirty 2>/dev/null || printf 'dev')}"
DIST_DIR="${2:-.goreleaser-artifacts}"

mkdir -p "$DIST_DIR"
cat > "$DIST_DIR/latest.json" <<EOF
{
  "latest_version": "$VERSION"
}
EOF
