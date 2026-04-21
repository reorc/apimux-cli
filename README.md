# APIMux CLI

Public command-line client for APIMux capabilities, plus the paired skill docs used by agents.

## Install

Install the latest release:

```bash
curl -fsSL https://github.com/reorc/apimux-cli/releases/latest/download/install.sh | sh
```

Install a specific version:

```bash
curl -fsSL https://github.com/reorc/apimux-cli/releases/latest/download/install.sh | \
  APIMUX_VERSION=v0.1.0 sh
```

The installer:

- detects `darwin/linux` and `amd64/arm64`
- downloads the matching `.tar.gz`
- verifies SHA-256 checksums
- installs `apimux` to `~/.local/bin/apimux`
- runs `apimux version` after install

If you publish release artifacts somewhere else, override:

```bash
APIMUX_RELEASE_BASE_URL=https://github.com/reorc/apimux-cli/releases/download \
APIMUX_RELEASE_MANIFEST_URL=https://github.com/reorc/apimux-cli/releases/latest/download/latest.json \
curl -fsSL https://github.com/reorc/apimux-cli/releases/latest/download/install.sh | sh
```

## Quickstart

Point the CLI at your APIMux service:

```bash
apimux config init
```

Run a capability:

```bash
apimux amazon get_product --asin B0CM5JV26D --market US
apimux google_trends get_interest_over_time --q "AI" --geo US
```

Useful flags:

- `--debug`: print the sanitized response envelope
- `--output compact`: default compact agent-facing body
- `--output pretty`: compact body with indented JSON
- `--output data`: raw `data` payload without projection
- `--output data-pretty`: raw `data` payload with indented JSON

## Repo Layout

- `cmd/apimux`: entrypoint
- `internal/`: CLI internals, projection layer, config, update checks
- `skills/apimux-*`: paired capability docs for agent usage
- `scripts/build-apimux-cli.sh`: release artifact builder
- `scripts/install.sh`: binary installer
- `docs/`: install, release, and development notes

## Development

Build a local binary:

```bash
make build
./dist/apimux version
```

Run tests:

```bash
make test
```

Build snapshot release archives and checksums:

```bash
make release-build
```

Version metadata is injected via git tags and commit info. The intended first public release line starts at `v0.1.0`.

## Skills

The `skills/apimux-*` directories ship alongside the CLI because both surfaces are versioned together. If you update a capability contract, update the corresponding skill doc in the same change.
