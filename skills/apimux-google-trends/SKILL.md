---
name: apimux-google-trends
version: 1.0.0
description: "Google Trends interest data. Get normalized search interest over time with region, time range, category, and Google property filters."
metadata:
  source: google_trends
  requires:
    bins: ["apimux"]
  cliHelp: "apimux google_trends --help"
---

# Google Trends

Get Google search interest data for demand validation, trend discovery, and keyword comparison.

**Before using:** Read [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, and CLI conventions.

## What you can do

- **Validate product or keyword demand** → `get_interest_over_time`
- **Compare up to 5 keywords in one request** → `get_interest_over_time`
- **Inspect long-term seasonality** → use `time="today 5-y"`
- **Assess supply after demand validation** → use [`amazon.search_products`](../apimux-amazon/SKILL.md)

## Available capabilities

| Capability | What it does | When to use |
|------------|--------------|-------------|
| `get_interest_over_time` | Get normalized search-interest time series | Demand validation and trend comparison |

Phase 1 includes only `get_interest_over_time`. Future capabilities may include interest by region, related topics, and related queries.

---

## google_trends.get_interest_over_time

Get a search-interest time series for one or more Google Trends queries.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `q` | string | Yes | Search query; comma-separate multiple queries, up to 5 |
| `geo` | string | No | ISO 3166-1 alpha-2 country or subregion code, such as `US` or `US-CA`; omit for worldwide |
| `time` | string | No | Time range; default `today 12-m`; see table below |
| `cat` | string | No | Google Trends category ID; default `"0"` for all categories |
| `gprop` | string | No | Google property: `""` (web), `"images"`, `"news"`, `"froogle"` (shopping), or `"youtube"`; default `""` |
| `tz` | integer | No | Timezone offset in minutes; omit to use UTC |

### Time ranges

| Value | Meaning |
|-------|---------|
| `now 1-H` | Past 1 hour |
| `now 4-H` | Past 4 hours |
| `now 1-d` | Past 1 day |
| `now 7-d` | Past 7 days |
| `today 1-m` | Past 30 days |
| `today 3-m` | Past 90 days |
| `today 12-m` | Past 12 months (default) |
| `today 5-y` | Past 5 years |
| `all` | 2004 to present |
| `YYYY-MM-DD YYYY-MM-DD` | Custom date range |

### CLI usage

```bash
# Search interest over the past 12 months
apimux google_trends get_interest_over_time --q "wireless earbuds" --geo "US"

# Compare two queries
apimux google_trends get_interest_over_time --q "wireless earbuds,bluetooth headphones" --geo "US"

# Long-term trend
apimux google_trends get_interest_over_time --q "wireless earbuds" --geo "US" --time "today 5-y"

# YouTube search interest
apimux google_trends get_interest_over_time --q "wireless earbuds" --gprop "youtube"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `timeline_data[].date` | string | Time-series date label |
| `timeline_data[].timestamp` | string | Unix timestamp |
| `timeline_data[].values[].query` | string | Query represented in the time-series point |
| `timeline_data[].values[].value` | integer | Normalized interest value, 0-100 |
| `averages[].query` | string | Query represented by the average |
| `averages[].value` | integer | Average interest over the selected period |
| `regions[].geo` | string | Region geo code |
| `regions[].name` | string | Region display name |
| `regions[].values[].query` | string | Query represented in the region value |
| `regions[].values[].value` | integer | Regional interest value |
| `related_topics.top[].position` | integer | Top related topic rank |
| `related_topics.top[].topic_id` | string | Topic ID |
| `related_topics.top[].title` | string | Top related topic title |
| `related_topics.top[].topic_type` | string | Top related topic type |
| `related_topics.top[].value` | integer | Top related topic value |
| `related_topics.top[].link` | string | Top related topic link |
| `related_topics.rising[].position` | integer | Rising related topic rank |
| `related_topics.rising[].topic_id` | string | Topic ID |
| `related_topics.rising[].title` | string | Rising related topic title |
| `related_topics.rising[].topic_type` | string | Rising related topic type |
| `related_topics.rising[].value` | integer | Rising related topic value |
| `related_topics.rising[].link` | string | Rising related topic link |
| `related_queries.top[].position` | integer | Top related query rank |
| `related_queries.top[].query` | string | Top related query text |
| `related_queries.top[].value` | integer | Top related query value |
| `related_queries.top[].link` | string | Top related query link |
| `related_queries.rising[].position` | integer | Rising related query rank |
| `related_queries.rising[].query` | string | Rising related query text |
| `related_queries.rising[].value` | integer | Rising related query value |
| `related_queries.rising[].link` | string | Rising related query link |
| `search_metadata.status` | string | Upstream request status |
| `search_metadata.created_at` | string | Request creation time |
| `search_metadata.request_url` | string | Upstream request URL |
| `search_parameters.q` | string | Original query parameter |
| `search_parameters.geo` | string | Original geo parameter |
| `search_parameters.time` | string | Original time parameter |
| `search_parameters.tz` | integer | Original timezone offset |
| `search_parameters.data_type` | string | Upstream data type |
| `search_parameters.cat` | string | Original category parameter |
| `search_parameters.region` | string | Original region parameter |
| `search_parameters.gprop` | string | Original Google property parameter |

### Common errors

- Treating interest values as absolute search volume → values are normalized 0-100 indices, not raw search counts.
- Comparing separate requests directly → values are comparable only within the same request.
- Using natural-language time ranges such as `"last year"` → use the supported values above or `YYYY-MM-DD YYYY-MM-DD`.
- Using `gprop="shopping"` → use `gprop="froogle"` for Google Shopping.
- Omitting `geo` when analyzing a specific market → omitted `geo` means worldwide.

### Notes

- `q` is required and cannot be empty.
- Up to 5 queries can be comma-separated in one request.
- Interest values are relative, normalized 0-100 indices.
- Omitted `geo` means worldwide.
- `gprop` accepts only `""`, `"images"`, `"news"`, `"froogle"`, or `"youtube"`.

---

## General notes

- See [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure and error handling.
