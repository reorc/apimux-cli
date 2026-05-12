# Development

Requirements:

- Go 1.24+
- network access to an APIMux service for manual CLI checks

Branch model:

- `develop` is the default integration branch for feature work and bug fixes.
- `main` is the stable release branch and should only advance via `develop -> main`.
- Release tags must be created from commits already contained in `main`.

Common commands:

```bash
make build
make test
make release-build
```

Local config is stored under `APIMUX_CONFIG_DIR` when set, otherwise `~/.apimux/config.json`. Older configs under the platform user config directory are still read as a fallback and are migrated on the next write. Fresh installs default to the production APIMux service at `https://apimux.io/api/core`.

Use `apimux auth login --web-url http://localhost:<port>` to exercise browser-assisted CLI auth against a local web app without persisting the web URL. Use `APIMUX_BASE_URL`, `--base-url`, or `apimux config set --base-url ...` to point CLI capability calls at a local service such as `http://127.0.0.1:8081`. For CI or manual API key setup against production, `apimux config set --api-key ...` is enough unless you intentionally need a non-production service.
