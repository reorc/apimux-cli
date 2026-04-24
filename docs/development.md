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

Local config is stored under `APIMUX_CONFIG_DIR` when set, otherwise the default user config path.

Use `apimux auth login --web-url http://localhost:<port>` to exercise browser-assisted CLI auth against a local web app without persisting the web URL. For CI or manual API key setup, use `apimux config set --base-url ... --api-key ...`.
