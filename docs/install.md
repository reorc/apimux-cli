# Install

Default install path:

```bash
curl -fsSL https://github.com/reorc/apimux-cli/releases/latest/download/install.sh | sh
```

Environment variables:

- `APIMUX_VERSION`: install a pinned version such as `v0.1.0`
- `APIMUX_INSTALL_DIR`: override the default install dir (`~/.local/bin`)
- `APIMUX_RELEASE_BASE_URL`: override the release asset base URL
- `APIMUX_RELEASE_MANIFEST_URL`: override the `latest.json` manifest URL

The installer expects release assets shaped like:

- `https://github.com/reorc/apimux-cli/releases/latest/download/latest.json`
- `https://github.com/reorc/apimux-cli/releases/download/v0.1.0/apimux_v0.1.0_darwin_arm64.tar.gz`
- `https://github.com/reorc/apimux-cli/releases/download/v0.1.0/apimux_v0.1.0_checksums.txt`
