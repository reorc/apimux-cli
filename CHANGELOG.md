# Changelog

All notable changes to APIMux CLI are documented here.

## [1.1.13] - 2026-06-29

### Changed

- Adjusted YouTube and Instagram compact list projections for agent workflows:
  default list output now hides long CDN/media URL fields while keeping
  descriptions/captions and short structured creator fields for ranking and
  drill-down decisions.
- Kept detail-style endpoints richer by default so raw/canonical URLs remain
  available when agents intentionally inspect one item.

### Documentation

- Updated bundled YouTube and Instagram skill docs to describe compact output
  field choices and when to use `--output data` for raw URL fields.

## [1.1.5] - 2026-05-12

### Added

- Implemented `apimux upgrade` for direct binary installs. The command now
  checks the release manifest, downloads the matching platform archive, verifies
  the checksum, extracts the `apimux` binary, and atomically replaces the
  current executable.

### Fixed

- Replaced the previous `cli_upgrade_not_implemented` response with either a
  successful upgrade result or a clear package-manager/install-script guidance
  error when the current executable cannot be safely replaced.
- Made schema-bound source command help degrade gracefully when the service
  schema cannot be reached, instead of returning a transport error for
  `--help`.

## [1.1.4] - 2026-05-12

### Fixed

- Changed the fresh-install default APIMux service URL from local development
  (`http://127.0.0.1:8081`) to production (`https://apimux.io/api/core`).
- Kept local/private service usage available through `APIMUX_BASE_URL`,
  `--base-url`, or `apimux config set --base-url ...`.

### Documentation

- Clarified that localhost is a development override, not the default product
  endpoint.

## [1.1.3] - 2026-05-12

### Added

- Added the `apimux tiktok get_video_detail` shortcut for TikTok video detail
  lookup by share URL or `aweme_id`.
- Added compact output projection for TikTok and Douyin video detail responses.

### Documentation

- Updated bundled TikTok and Douyin skills with video detail examples, expiring
  media URL notes, and `video_not_found` guidance.

## [1.1.2] - 2026-04-29

### Fixed

- Improved dynamic capability help output for schema-backed commands.
- Improved CLI handling of non-envelope HTTP error responses.

### Documentation

- Converted bundled `apimux-*` skill docs to English, user-facing guidance.

## [1.1.1] - 2026-04-28

### Fixed

- Updated `meta_ads.get_ad_detail` compact output to match the service's
  workflow-backed `CoreMetaAd` contract. Default CLI output now preserves
  `page_id`, `page_name`, dates, activity status, `publisher_platform`,
  `snapshot`, and collation fields instead of reducing the response to `ad_id`.
- Updated `meta_ads.search_ads` compact output to use the canonical
  `publisher_platform` field name.

### Documentation

- Updated `skills/apimux-meta-ads` to explain that `get_ad_detail` reads ad
  items previously populated by `search_ads` and no longer calls the provider
  detail endpoint.

## [1.1.0] - 2026-04-24

### Added

- Added the web-assisted `apimux auth login` flow and made it the onboarding
  path.
- Added explicit TrendCloud compact output rules.

### Changed

- Migrated CLI config storage to the user's home directory.

## [1.0.0] - 2026-04-23

### Added

- Initial stable release of the APIMux CLI and bundled `apimux-*` skills.
