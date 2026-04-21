#!/usr/bin/env bash

set -euo pipefail

DEFAULT_RELEASE_BASE_URL="https://github.com/reorc/apimux-cli/releases/download"
DEFAULT_RELEASE_MANIFEST_URL="https://github.com/reorc/apimux-cli/releases/latest/download/latest.json"

APIMUX_RELEASE_MANIFEST_URL="${APIMUX_RELEASE_MANIFEST_URL:-$DEFAULT_RELEASE_MANIFEST_URL}"
APIMUX_RELEASE_BASE_URL="${APIMUX_RELEASE_BASE_URL:-$DEFAULT_RELEASE_BASE_URL}"
APIMUX_INSTALL_DIR="${APIMUX_INSTALL_DIR:-$HOME/.local/bin}"
APIMUX_VERSION="${APIMUX_VERSION:-latest}"

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    printf 'missing required command: %s\n' "$1" >&2
    exit 1
  }
}

need_cmd curl
need_cmd tar

if ! command -v python3 >/dev/null 2>&1 && ! command -v python >/dev/null 2>&1; then
  printf 'missing required command: python3 or python\n' >&2
  exit 1
fi

json_get() {
  local key="$1"
  local input="$2"
  local python_bin
  python_bin="$(command -v python3 || command -v python)"
  printf '%s' "$input" | "$python_bin" -c 'import json,sys; data=json.load(sys.stdin); print(data.get(sys.argv[1], ""))' "$key"
}

detect_os() {
  case "$(uname -s)" in
    Darwin) printf 'darwin' ;;
    Linux) printf 'linux' ;;
    *) printf 'unsupported' ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) printf 'amd64' ;;
    arm64|aarch64) printf 'arm64' ;;
    *) printf 'unsupported' ;;
  esac
}

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

OS_NAME="$(detect_os)"
ARCH_NAME="$(detect_arch)"

if [[ "$OS_NAME" == "unsupported" || "$ARCH_NAME" == "unsupported" ]]; then
  printf 'unsupported platform: %s/%s\n' "$(uname -s)" "$(uname -m)" >&2
  exit 1
fi

if [[ "${APIMUX_RELEASE_BASE_URL%/}" != "$DEFAULT_RELEASE_BASE_URL" && "${APIMUX_RELEASE_MANIFEST_URL:-}" == "$DEFAULT_RELEASE_MANIFEST_URL" ]]; then
  APIMUX_RELEASE_MANIFEST_URL="${APIMUX_RELEASE_BASE_URL%/}/latest/download/latest.json"
fi

if [[ -z "$APIMUX_RELEASE_MANIFEST_URL" ]]; then
  printf 'set APIMUX_RELEASE_MANIFEST_URL or APIMUX_RELEASE_BASE_URL before running install.sh\n' >&2
  exit 1
fi

MANIFEST="$(curl -fsSL "$APIMUX_RELEASE_MANIFEST_URL")"

if [[ "$APIMUX_VERSION" == "latest" ]]; then
  APIMUX_VERSION="$(json_get latest_version "$MANIFEST")"
fi

if [[ -z "$APIMUX_VERSION" ]]; then
  printf 'could not resolve APIMUX_VERSION from manifest %s\n' "$APIMUX_RELEASE_MANIFEST_URL" >&2
  exit 1
fi

if [[ -z "$APIMUX_RELEASE_BASE_URL" ]]; then
  APIMUX_RELEASE_BASE_URL="$DEFAULT_RELEASE_BASE_URL"
fi

ARCHIVE_NAME="apimux_${APIMUX_VERSION}_${OS_NAME}_${ARCH_NAME}.tar.gz"
ARCHIVE_URL="${APIMUX_RELEASE_BASE_URL%/}/${APIMUX_VERSION}/${ARCHIVE_NAME}"
CHECKSUM_URL="${APIMUX_RELEASE_BASE_URL%/}/${APIMUX_VERSION}/apimux_${APIMUX_VERSION}_checksums.txt"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

ARCHIVE_PATH="$TMP_DIR/$ARCHIVE_NAME"
CHECKSUM_PATH="$TMP_DIR/checksums.txt"

curl -fsSL "$ARCHIVE_URL" -o "$ARCHIVE_PATH"
curl -fsSL "$CHECKSUM_URL" -o "$CHECKSUM_PATH"

EXPECTED_SUM="$(awk -v file="$ARCHIVE_NAME" '$2 == file {print $1}' "$CHECKSUM_PATH")"
ACTUAL_SUM="$(sha256_file "$ARCHIVE_PATH")"

if [[ -z "$EXPECTED_SUM" ]]; then
  printf 'checksum entry not found for %s\n' "$ARCHIVE_NAME" >&2
  exit 1
fi
if [[ "$EXPECTED_SUM" != "$ACTUAL_SUM" ]]; then
  printf 'checksum mismatch for %s\n' "$ARCHIVE_NAME" >&2
  exit 1
fi

mkdir -p "$APIMUX_INSTALL_DIR"
tar -C "$TMP_DIR" -xzf "$ARCHIVE_PATH"
install -m 0755 "$TMP_DIR/apimux" "$APIMUX_INSTALL_DIR/apimux"

printf 'Installed apimux %s to %s/apimux\n' "$APIMUX_VERSION" "$APIMUX_INSTALL_DIR"
"$APIMUX_INSTALL_DIR/apimux" version
case ":$PATH:" in
  *":$APIMUX_INSTALL_DIR:"*) ;;
  *)
    printf 'Add %s to PATH before using the CLI.\n' "$APIMUX_INSTALL_DIR"
    ;;
esac
