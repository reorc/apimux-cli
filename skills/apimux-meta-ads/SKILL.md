---
name: apimux-meta-ads
version: 1.0.0
description: "Meta Ads Library data. Search ads and inspect ad details for creative research, competitive analysis, and campaign sampling."
metadata:
  source: meta_ads
  requires:
    bins: ["apimux"]
  cliHelp: "apimux meta_ads --help"
---

# Meta Ads

Search and analyze ads from Meta Ads Library (Facebook, Instagram). Use this for competitive research, ad creative inspiration, and campaign analysis.

**Before using:** Read [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, pagination metadata, and CLI conventions.

## What you can do

- **Search ads by keyword** → `search_ads` — find ads matching your search terms
- **Get detailed ad info** → `get_ad_detail` — view full details for a specific ad
- **Analyze ad campaigns** → search first, then get details for interesting ads

## Available capabilities

| Capability | What it does | When to use |
|------------|--------------|-------------|
| `search_ads` | Search Meta Ads Library | Find ads by keyword, filter by country/platform/date |
| `get_ad_detail` | Get full details for an ad | View complete info after finding an ad via search |

---

## meta_ads.search_ads

Search for ads in Meta Ads Library by keyword. Returns a list of ads with creative snapshots.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `q` | string | Yes | Search keyword |
| `country` | string | No | ISO alpha-2 country code (e.g., "US", "GB") |
| `ad_type` | string | No | Filter by ad type: `all`, `political_and_issue_ads`, `housing_ads`, `employment_ads`, `credit_ads` (default: `all`) |
| `active_status` | string | No | Filter by status: `active`, `inactive`, `all` (default: `all`) |
| `media_type` | string | No | Filter by media: `all`, `video`, `image`, `meme`, `image_and_meme`, `none` (default: `all`) |
| `platforms` | string | No | Comma-separated platforms: `facebook,instagram` |
| `start_date` | string | No | Start date in `YYYY-MM-DD` format |
| `end_date` | string | No | End date in `YYYY-MM-DD` format |
| `next_page_token` | string | No | Pagination token (omit for first page) |

### CLI usage

```bash
apimux meta_ads search_ads --q "fitness app"
apimux meta_ads search_ads --q "fitness app" --country "US" --media-type "video"
apimux meta_ads search_ads --q "fitness app" --platforms "facebook,instagram" --start-date "2026-01-01"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `ad_id` | string | Ad ID |
| `page_id` | string | Page ID |
| `page_name` | string | Page name |
| `start_date` | string | Ad start date |
| `end_date` | string | Ad end date |
| `is_active` | boolean | Whether ad is currently active |
| `categories` | string[] | Ad categories |
| `publisher_platform` | string[] | Platforms where ad appears |
| `snapshot` | object | Creative snapshot with text, title, links, cards, videos |

### Notes

- `q` is required.
- `country` must be an ISO alpha-2 country code when provided.
- `ad_type`, `active_status`, and `media_type` accept only the enum values listed above.
- `platforms` must be a comma-separated list of lowercase platform names.
- `start_date` and `end_date` must be `YYYY-MM-DD`.
- Use `meta.next_page_token` as `next_page_token` to fetch the next page.

---

## meta_ads.get_ad_detail

Get full details for a specific ad. Use this after `search_ads` to view complete information for ads from your search results.

**Important:** Run `search_ads` first. This capability returns the full ad record for an `ad_id` that appeared in your APIMux search results.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ad_id` | string | Yes | Ad archive ID from search results |

### CLI usage

```bash
apimux meta_ads search_ads --q "fitness app"
apimux meta_ads get_ad_detail --ad-id "477570185419072"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `ad_id` | string | Ad ID |
| `page_id` | string | Page ID |
| `page_name` | string | Page name |
| `start_date` | string | Ad start date |
| `end_date` | string | Ad end date |
| `is_active` | boolean | Whether ad is currently active |
| `categories` | string[] | Ad categories |
| `entity_type` | string | Entity type |
| `gated_type` | string | Gated status |
| `hide_data_status` | string | Data visibility status |
| `publisher_platform` | string[] | Platforms where ad appears |
| `impressions_index` | number | Impressions index |
| `total_active_time` | number | Total active time |
| `snapshot` | object | Creative snapshot with text, title, links, cards, videos |
| `collation_count` | number | Collation count |
| `collation_id` | string | Collation ID |

### Notes

- `ad_id` is required.
- Run `search_ads` first, then use an `ad_id` from the search results.
- If the ad cannot be found, the response returns `ad_not_found`; run `meta_ads.search_ads` again and use an `ad_id` from the search results.
- This response does not include legacy provider-detail-only fields such as `eu_transparency`, `political_insights`, or `verified_voice`.

---

## Common patterns

**Find and analyze ads:**
```bash
# 1. Search for ads
apimux meta_ads search_ads --q "fitness app" --country "US"

# 2. Get details for specific ad from results
apimux meta_ads get_ad_detail --ad-id "477570185419072"
```

**Filter by platform and date:**
```bash
apimux meta_ads search_ads --q "fitness app" \
  --platforms "facebook,instagram" \
  --start-date "2026-01-01" \
  --media-type "video"
```

## General notes

- See [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure and error handling.
- `search_ads` pagination uses `meta.next_page_token`.
