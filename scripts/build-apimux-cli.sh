#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GORELEASER_CMD="${GORELEASER_CMD:-go run github.com/goreleaser/goreleaser/v2@v2.12.7}"

cd "$ROOT_DIR"

eval "$GORELEASER_CMD release --snapshot --clean"
