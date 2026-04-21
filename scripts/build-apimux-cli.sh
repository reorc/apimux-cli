#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST_DIR="${DIST_DIR:-$ROOT_DIR/dist/releases}"
VERSION="${VERSION:-$(git -C "$ROOT_DIR" describe --tags --always --dirty 2>/dev/null || printf 'dev')}"
GIT_COMMIT="${GIT_COMMIT:-$(git -C "$ROOT_DIR" rev-parse --short HEAD 2>/dev/null || printf 'unknown')}"
BUILD_DATE="${BUILD_DATE:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}"
RELEASE_MANIFEST_URL="${RELEASE_MANIFEST_URL:-https://github.com/reorc/apimux-cli/releases/latest/download/latest.json}"

TARGETS=(
  "darwin amd64"
  "darwin arm64"
  "linux amd64"
  "linux arm64"
)

mkdir -p "$DIST_DIR/$VERSION"
CHECKSUM_FILE="$DIST_DIR/$VERSION/apimux_${VERSION}_checksums.txt"
: > "$CHECKSUM_FILE"

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

for target in "${TARGETS[@]}"; do
  read -r GOOS GOARCH <<<"$target"
  OUT_DIR="$DIST_DIR/$VERSION/${GOOS}_${GOARCH}"
  BIN_NAME="apimux"
  ARCHIVE_NAME="apimux_${VERSION}_${GOOS}_${GOARCH}.tar.gz"
  ARCHIVE_PATH="$DIST_DIR/$VERSION/$ARCHIVE_NAME"

  rm -rf "$OUT_DIR"
  mkdir -p "$OUT_DIR"

  (
    cd "$ROOT_DIR"
    CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
      go build \
      -ldflags="-X github.com/reorc/apimux-cli/internal/buildinfo.Version=$VERSION -X github.com/reorc/apimux-cli/internal/buildinfo.Commit=$GIT_COMMIT -X github.com/reorc/apimux-cli/internal/buildinfo.BuildDate=$BUILD_DATE -X github.com/reorc/apimux-cli/internal/buildinfo.ReleaseManifestURL=$RELEASE_MANIFEST_URL" \
      -o "$OUT_DIR/$BIN_NAME" \
      ./cmd/apimux
  )

  tar -C "$OUT_DIR" -czf "$ARCHIVE_PATH" "$BIN_NAME"
  printf '%s  %s\n' "$(sha256_file "$ARCHIVE_PATH")" "$ARCHIVE_NAME" >> "$CHECKSUM_FILE"
done

cat > "$DIST_DIR/latest.json" <<EOF
{
  "latest_version": "$VERSION"
}
EOF

printf 'Built APIMux CLI release artifacts in %s\n' "$DIST_DIR/$VERSION"
