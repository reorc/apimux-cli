# Changelog

All notable changes to APIMux CLI are documented here.

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
