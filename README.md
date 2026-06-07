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
  APIMUX_VERSION=v1.0.0 sh
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

Authorize the CLI:

```bash
apimux auth login
```

The CLI defaults to the production APIMux service at `https://apimux.io/api/core`.

On SSH/headless servers, print the authorization URL instead of opening a browser:

```bash
apimux auth login --no-browser
```

`auth login` also auto-detects common headless environments such as SSH, CI, `NO_BROWSER=1`, and Linux without `$DISPLAY`.

For CI or manual API key setup, use `apimux config set --api-key ...`. Override the service URL only for local development or private deployments:

```bash
APIMUX_BASE_URL=http://127.0.0.1:8081 apimux schema capabilities
apimux --base-url http://127.0.0.1:8081 schema capabilities
apimux config set --base-url http://127.0.0.1:8081
```

Run a capability:

```bash
apimux amazon get_product --asin B0CM5JV26D --market US
apimux google_trends get_interest_over_time --q "AI" --geo US
```

Discover available capabilities:

```bash
apimux schema capabilities
apimux schema show amazon.get_product
apimux schema list
```

Check for updates or upgrade direct binary installs:

```bash
apimux upgrade --check
apimux upgrade
```

If your `apimux` executable is managed by a package manager or is installed as a symlink, use that package manager or rerun the install script instead.

Useful flags:

- `--debug`: print the sanitized response envelope
- `--output compact`: default compact agent-facing body
- `--output pretty`: compact body with indented JSON
- `--output data`: raw `data` payload without projection
- `--output data-pretty`: raw `data` payload with indented JSON

Shell completion:

```bash
# bash
apimux completion bash > ~/.local/share/bash-completion/completions/apimux

# zsh
apimux completion zsh > "${fpath[1]}/_apimux"

# fish
apimux completion fish > ~/.config/fish/completions/apimux.fish

# PowerShell
apimux completion powershell | Out-String | Invoke-Expression
```

## Caller Context (for billing webhook callers)

API keys that have a billing webhook configured on the apimux-service side
need each request to carry caller-provided attribution context so the
receiver can settle the call against the right subject. `apimux` scans the
process environment for any variable prefixed `APIMUX_CALLER_CTX_` and
forwards it transparently as the matching `X-Apimux-Caller-Ctx-*` HTTP
header — the prefix convention is generic by design so apimux never has
to know which keys a particular downstream uses.

The naming rule:

```
APIMUX_CALLER_CTX_<SUFFIX>   →   X-Apimux-Caller-Ctx-<Suffix-With-Dashes>
```

For Kamay's apimux receiver, that looks like:

```bash
APIMUX_CALLER_CTX_USER_ID="usr_abc" \
APIMUX_CALLER_CTX_WORKSPACE_ID="ws_def" \
APIMUX_CALLER_CTX_CONVERSATION_ID="conv_xyz" \
APIMUX_CALLER_CTX_MESSAGE_ID="msg_qrs" \
  apimux reddit search --query "headphones"
```

A different downstream might use entirely different keys (the CLI does not
care):

```bash
APIMUX_CALLER_CTX_ACCOUNT_ID="acct_abc" \
APIMUX_CALLER_CTX_PROJECT_ID="proj_def" \
  apimux reddit search --query "..."
```

Unset / empty / whitespace-only values are treated as absent — no header
is sent and apimux-service records SQL NULL on `api_call_logs.caller_context`.
If the apikey has no webhook configured the headers are still persisted
for reconciliation but no webhook event is enqueued.

## Repo Layout

- `cmd/apimux`: entrypoint
- `internal/`: CLI internals, projection layer, config, update checks
- `skills/apimux-*`: paired capability docs for agent usage
- `scripts/build-apimux-cli.sh`: release artifact builder
- `scripts/install.sh`: binary installer
- `docs/`: install, release, and development notes

## Development

Branch model:

- `develop`: integration branch for day-to-day work
- `main`: release branch; keep it aligned with the latest stable tag
- release tags: create only from commits already contained in `main`

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

Version metadata is injected via git tags and commit info. The current stable release line starts at `v1.0.0`.

## Skills

The `skills/apimux-*` directories ship alongside the CLI because both surfaces are versioned together. If you update a capability contract, update the corresponding skill doc in the same change.
