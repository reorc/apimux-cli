---
name: apimux-trendcloud
version: 1.0.0
description: "TrendCloud China e-commerce channel data. Analyze sales, volume, share, rankings, trends, and filter values across Douyin, JD, and Tmall."
metadata:
  source: trendcloud
  requires:
    bins: ["apimux"]
  cliHelp: "apimux trendcloud --help"
---

# TrendCloud

Analyze partial China e-commerce channel data across `douyin`, `jd`, and `tmall`. Use this for brand/category/series/SKU/attribute share, rankings, trends, and filter discovery inside the TrendCloud coverage area.

**Before using:** Read [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, partial-failure semantics, and CLI conventions.

TrendCloud is best for questions about China e-commerce channel performance, such as sales, volume, market share, rankings, price bands, and platform differences across Douyin/JD/Tmall.

TrendCloud is not a full-market source for:
- All of China or all online retail channels
- Amazon, TikTok Shop, Reddit, Google Trends, or other non-China e-commerce sources
- User profiles, creative content, comments, or VOC analysis

## What you can do

- **Discover exact filter values** → `search_filter_values`
- **View sales or volume trends over time** → `get_market_trend`
- **Rank brands, categories, series, SKUs, or attributes** → `get_top_rankings`
- **Compare Douyin/JD/Tmall** → pass `filters.platforms` explicitly

## Recommended workflow

1. Use `search_filter_values` to find exact filter labels or paths.
2. Put returned `label` or `path` values into `filters`.
3. Use `get_market_trend` for time-series analysis.
4. Use `get_top_rankings` for leaderboards, share, and category structure.

When the user gives natural-language terms such as "coffee", "Luckin", or "drip coffee", do not guess filters directly. Run discovery first.

## Available capabilities

| Capability | What it does | When to use |
|------------|--------------|-------------|
| `get_market_trend` | Get monthly trend time series | Sales/volume trends and platform comparison |
| `get_top_rankings` | Get rankings | Brand share, category structure, series/SKU ranking |
| `search_filter_values` | Discover filter values | Category, brand, series, SKU, and attribute discovery |

---

## trendcloud.get_market_trend

Get monthly market trends from TrendCloud.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `start_month` | string | No | Start month, `YYYY-MM` |
| `end_month` | string | No | End month, `YYYY-MM` |
| `metrics` | string[] | No | Metric enum: `sales`, `volume`; default `["sales"]` |
| `filters` | object | No | Structured filters supporting `platforms`, `categories`, `brands`, `series`, `skus`, and `attributes` |

### CLI usage

```bash
apimux trendcloud get_market_trend --start-month "2025-01" --end-month "2025-12"
apimux trendcloud get_market_trend --metrics "sales,volume" --filters-json '{"platforms":["douyin"],"brands":["Luckin"]}'
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `period` | string | Month, `YYYY-MM` |
| `sales` | number | Sales amount in CNY |
| `volume` | integer | Sales volume |

### Notes

- Amount fields are returned in CNY.
- `meta.resolved_time_range` explains defaults or clamping.
- `meta.resolved_filters` explains how filters were resolved.

---

## trendcloud.get_top_rankings

Get TrendCloud ranking results.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `entity` | string | Yes | Ranking entity: `brand`, `category`, `series`, `sku`, or `attribute` |
| `metric` | string | No | Primary sort metric; valid values depend on `entity` |
| `start_month` | string | No | Start month, `YYYY-MM` |
| `end_month` | string | No | End month, `YYYY-MM` |
| `top_n` | integer | No | Number of rows, `1..100`; default `20` |
| `category_level` | string | No | Only for `entity=category`: `category1`, `category2`, or `category3` |
| `filters` | object | No | Structured filters |

### CLI usage

```bash
apimux trendcloud get_top_rankings --entity "brand" --metric "sales"
apimux trendcloud get_top_rankings --entity "category" --category-level "category2" --filters-json '{"platforms":["tmall"]}'
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `rank` | integer | 1-based rank |
| `label` | string | Entity display name |
| `sales` | number | Sales amount in CNY |
| `volume` | integer | Sales volume |
| `market_share` | number | Market share |
| `avg_price` | number | Average price in CNY |
| `sales_change_ratio` | number | Sales year-over-year change ratio |
| `volume_change_ratio` | number | Volume year-over-year change ratio |
| `market_share_change_ratio` | number | Market-share year-over-year change ratio |

### Notes

- Amount fields are returned in CNY.
- Use `search_filter_values` first when the exact entity/filter label is uncertain.

---

## trendcloud.search_filter_values

Search for TrendCloud filter-value candidates.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `kind` | string | Yes | `category`, `brand`, `series`, `sku`, or `attribute` |
| `query` | string | Yes | Search keyword |
| `platforms` | string[] | No | Platform scope: `douyin`, `jd`, `tmall` |
| `categories` | string[] | No | Category hints |
| `limit` | integer | No | Number of candidates, `1..50`; default `10` |

### CLI usage

```bash
apimux trendcloud search_filter_values --kind "category" --query "coffee"
apimux trendcloud search_filter_values --kind "brand" --query "luckin" --platforms "douyin,jd"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `label` | string | Display name |
| `path` | string[] | Full hierarchy path for tree fields; `[label]` for simple fields |

### Common errors

- Platform outside `douyin/jd/tmall` → `invalid_platform`
- Time range is not `YYYY-MM` or spans more than 12 months → `invalid_time_range`
- A filter cannot be resolved uniquely → `ambiguous_filter`
- A filter does not exist → `invalid_filter`
- Concurrency is full → `busy` with `retry_after_seconds`

### Notes

- If a tree filter is ambiguous, use `search_filter_values` for discovery, then retry the analysis call with the resolved filter.
- `meta.resolved_filters` explains the actual filter resolution.

---

## General notes

- See [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, and partial-failure semantics.
- TrendCloud covers only the configured China e-commerce channels; do not present it as all-channel China market data.
