#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INSTALL_SCRIPT="$ROOT_DIR/scripts/install.sh"

run_install_case() {
  local requested_version="$1"
  local manifest_version="$2"
  local expected_archive="$3"
  local expected_checksums="$4"

  local tmpdir
  tmpdir="$(mktemp -d)"
  trap 'rm -rf "$tmpdir"' RETURN

  cat >"$tmpdir/manifest.json" <<EOF
{"latest_version":"$manifest_version"}
EOF

  cat >"$tmpdir/apimux" <<'EOF'
#!/usr/bin/env bash
printf '{"version":"1.0.0","commit":"test"}\n'
EOF
  chmod +x "$tmpdir/apimux"

  tar -C "$tmpdir" -czf "$tmpdir/$expected_archive" apimux

  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$tmpdir/$expected_archive" | awk -v file="$expected_archive" '{print $1 "  " file}' >"$tmpdir/$expected_checksums"
  else
    shasum -a 256 "$tmpdir/$expected_archive" | awk -v file="$expected_archive" '{print $1 "  " file}' >"$tmpdir/$expected_checksums"
  fi

  local curl_stub="$tmpdir/curl"
  cat >"$curl_stub" <<EOF
#!/usr/bin/env bash
set -euo pipefail
args=("\$@")
url=""
output=""
for ((i=0; i<\${#args[@]}; i++)); do
  if [[ "\${args[i]}" == "-o" ]]; then
    output="\${args[i+1]}"
    ((i++))
    continue
  fi
  if [[ "\${args[i]}" != -* ]]; then
    url="\${args[i]}"
  fi
done
case "\$url" in
  http://example.test/latest.json)
    if [[ -n "\$output" ]]; then cp "$tmpdir/manifest.json" "\$output"; else cat "$tmpdir/manifest.json"; fi
    ;;
  http://example.test/releases/download/v1.0.0/$expected_archive)
    cp "$tmpdir/$expected_archive" "\$output"
    ;;
  http://example.test/releases/download/v1.0.0/$expected_checksums)
    cp "$tmpdir/$expected_checksums" "\$output"
    ;;
  *)
    printf 'unexpected curl url: %s\n' "\$url" >&2
    exit 1
    ;;
esac
EOF
  chmod +x "$curl_stub"

  local install_dir="$tmpdir/install"
  mkdir -p "$install_dir"
  PATH="$tmpdir:$PATH" \
    APIMUX_RELEASE_MANIFEST_URL="http://example.test/latest.json" \
    APIMUX_RELEASE_BASE_URL="http://example.test/releases/download" \
    APIMUX_INSTALL_DIR="$install_dir" \
    APIMUX_VERSION="$requested_version" \
    sh "$INSTALL_SCRIPT" >/dev/null

  test -x "$install_dir/apimux"
}

run_install_case "latest" "1.0.0" "apimux_1.0.0_darwin_arm64.tar.gz" "apimux_1.0.0_checksums.txt"
run_install_case "v1.0.0" "1.0.0" "apimux_1.0.0_darwin_arm64.tar.gz" "apimux_1.0.0_checksums.txt"

echo "install.sh tests passed"
