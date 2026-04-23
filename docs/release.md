# Release Contract

Public source of truth:

- GitHub: `https://github.com/reorc/apimux-cli`
- GitLab mirror: `https://gitlab.tool.reorc.cloud/kamay/apimux-cli`

Release contract:

- Branches follow a two-line model: `develop` for integration, `main` for stable release.
- Release candidates are promoted by merging or fast-forwarding `develop` into `main`.
- Git tags follow `vMAJOR.MINOR.PATCH`
- install script remains stable at `releases/latest/download/install.sh`
- latest manifest remains stable at `releases/latest/download/latest.json`
- assets remain stable at `releases/download/<tag>/apimux_<tag>_<os>_<arch>.tar.gz`
- checksums remain stable at `releases/download/<tag>/apimux_<tag>_checksums.txt`

Release automation:

- CI build runs on pull requests and pushes to both `develop` and `main`
- GitHub release runs on tag push matching `v*`
- Release workflow rejects tags that do not point to commits already reachable from `main`
- GoReleaser builds darwin/linux × amd64/arm64 archives
- GoReleaser uploads archives, checksums, `install.sh`, and `latest.json` to GitHub Releases

Recommended release flow:

1. Land normal development on `develop`.
2. Prepare a `develop -> main` promotion when the next release candidate is ready.
3. Merge or fast-forward `main` to the target release commit.
4. Create and push the release tag from `main`.
5. Let GitHub Actions publish the release assets.

Manual dual-push remains the repo sync model:

- GitHub is the public source of truth
- GitLab remains the internal mirror
