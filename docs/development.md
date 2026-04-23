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

Use `apimux config init` to bootstrap `base_url` and `api_key` for local testing.
