# Development

Requirements:

- Go 1.24+
- network access to an APIMux service for manual CLI checks

Common commands:

```bash
make build
make test
make release-build
```

Local config is stored under `APIMUX_CONFIG_DIR` when set, otherwise the default user config path.

Use `apimux config init` to bootstrap `base_url` and `api_key` for local testing.
