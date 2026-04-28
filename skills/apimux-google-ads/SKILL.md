---
name: apimux-google-ads
version: 1.0.0
description: "Google Ads Transparency Center data. Search advertisers, list ad creatives, and inspect ad details for competitive creative research."
metadata:
  source: google_ads
  requires:
    bins: ["apimux"]
  cliHelp: "apimux google_ads --help"
---

# Google Ads

Search Google Ads Transparency Center data across advertisers, domains, creatives, and creative details.

**Before using:** Read [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, pagination metadata, and CLI conventions.

## What you can do

- **Find advertisers or domains** → `search_advertisers`
- **List creatives for an advertiser or domain** → `list_ad_creatives`
- **Inspect one creative** → `get_ad_details`

## Available capabilities

| Capability | What it does | When to use |
|------------|--------------|-------------|
| `search_advertisers` | Search advertisers and domains | Discover advertiser IDs and confirm brands |
| `list_ad_creatives` | List ad creatives | Collect creative samples and filter by region/platform/format |
| `get_ad_details` | Get creative details | Inspect one ad's information and variations |

---

## google_ads.search_advertisers

Search Google Ads advertisers and related domains.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | Yes | Search keyword |
| `region` | string | No | ISO alpha-2 country code; omit for no region filter |
| `num_advertisers` | integer | No | Number of advertisers to return, `1..100`; default `10` |
| `num_domains` | integer | No | Number of domains to return, `1..100`; default `10` |

### CLI usage

```bash
apimux google_ads search_advertisers --query "Nike"
apimux google_ads search_advertisers --query "Nike" --region "US" --num-advertisers 20
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `advertisers` | object[] | Matching advertisers |
| `domains` | object[] | Matching domains |

### Notes

- `query` is required.
- `region` must be an ISO alpha-2 country code when provided.
- `num_advertisers` and `num_domains` must be in `1..100`.

---

## google_ads.list_ad_creatives

List ad creatives by advertiser ID or domain.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `advertiser_id` | string | No | Advertiser ID; must start with `AR` |
| `domain` | string | No | Advertiser domain |
| `region` | string | No | ISO alpha-2 country code; omit for no region filter |
| `platform` | string | No | `google_play`, `google_maps`, `google_search`, `youtube`, or `google_shopping`; omit for no platform filter |
| `ad_format` | string | No | `text`, `image`, or `video`; omit for no format filter |
| `time_period` | string | No | `last_7_days`, `last_30_days`, `last_90_days`, or `last_year`; omit for no time filter |
| `page_token` | string | No | Pagination token; omit for the first page |

### CLI usage

```bash
apimux google_ads list_ad_creatives --advertiser-id "AR123456789"
apimux google_ads list_ad_creatives --domain "nike.com" --region "US" --ad-format "video"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `creative_id` | string | Creative ID |
| `advertiser_id` | string | Advertiser ID |
| `advertiser_name` | string | Advertiser name |
| `target_domain` | string | Landing page domain |
| `format` | string | Creative format |
| `first_shown_datetime` | string | First shown time |
| `last_shown_datetime` | string | Last shown time |
| `total_days_shown` | integer | Total days shown |
| `details_link` | string | Google details page URL |

### Notes

- Provide at least one of `advertiser_id` or `domain`.
- `advertiser_id` must start with `AR`.
- `platform`, `ad_format`, and `time_period` accept only the enum values listed above.
- Next-page state is returned in `meta.page_token` when available.

---

## google_ads.get_ad_details

Get details for one Google Ads creative.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `advertiser_id` | string | Yes | Advertiser ID; must start with `AR` |
| `creative_id` | string | Yes | Creative ID; must start with `CR` |

### CLI usage

```bash
apimux google_ads get_ad_details --advertiser-id "AR123456789" --creative-id "CR987654321"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `ad_information` | object | Ad metadata and targeting information |
| `variations` | object[] | Creative variations |

### Notes

- `advertiser_id` must start with `AR`.
- `creative_id` must start with `CR`.
- Missing creatives return `ad_not_found`.

---

## General notes

- See [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure and error handling.
- `list_ad_creatives` pagination uses `meta.page_token`.
