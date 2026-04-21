# Release Contract

Public source of truth:

- GitHub: `https://github.com/reorc/apimux-cli`
- GitLab mirror: `https://gitlab.tool.reorc.cloud/kamay/apimux-cli`

Release contract:

- Git tags follow `vMAJOR.MINOR.PATCH`
- install script remains stable at `releases/latest/download/install.sh`
- latest manifest remains stable at `releases/latest/download/latest.json`
- assets remain stable at `releases/download/<tag>/apimux_<tag>_<os>_<arch>.tar.gz`
- checksums remain stable at `releases/download/<tag>/apimux_<tag>_checksums.txt`

Release automation:

- CI build runs on pull requests and pushes to `main`
- GitHub release runs on tag push matching `v*`
- GoReleaser builds darwin/linux × amd64/arm64 archives
- GoReleaser uploads archives, checksums, `install.sh`, and `latest.json` to GitHub Releases

Manual dual-push remains the repo sync model:

- GitHub is the public source of truth
- GitLab remains the internal mirror
