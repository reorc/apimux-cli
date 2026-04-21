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

Current publishing mode is manual dual-push plus manual release asset upload. CI can replace the mechanics later without changing the public URL contract.
