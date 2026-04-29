---
name: apimux-amazon
version: 1.0.0
description: "Amazon e-commerce data. Product details, search, reviews, category best sellers, and trends. For product research, competitive analysis, and market assessment."
metadata:
  source: amazon
  requires:
    bins: ["apimux"]
  cliHelp: "apimux amazon --help"
---

# Amazon

Query Amazon e-commerce data across products, keywords, reviews, and categories. Use this for product research, competitive analysis, and market intelligence.

**Before using:** Read [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, and partial-failure semantics.

## What you can do

**Keywords & Search Intelligence:**
- Expand seed keywords → `expand_keywords`
- Get keyword metrics (search volume, competition, ads) → `get_keyword_overview`
- View keyword search trends over time → `get_keyword_trends`
- Find keywords associated with an ASIN → `list_asin_keywords`
- Query Amazon ABA trending keywords → `query_aba_keywords`

**Sales & Trends:**
- View daily sales trend for an ASIN → `get_asin_sales_daily_trend`
- Compare monthly sales history across ASINs → `get_asins_sales_history`
- View variant sales for parent ASIN (30 days) → `get_variant_sales_30d`

**Products:**
- Get product details by ASIN → `get_product`
- Search products by keyword → `search_products`
- Get product reviews → `get_product_reviews`

**Categories:**
- Find category by name (get node_id) → `search_category`
- Get category best sellers → `get_category_best_sellers`
- View category market trends → `get_category_trend`

**Cross-source workflow:**
- Validate market demand → Use [`google_trends.get_interest_over_time`](../apimux-google-trends/SKILL.md) first, then `search_products` to assess supply

## Available capabilities

| Capability | What it does | When to use |
|------------|--------------|-------------|
| `expand_keywords` | Expand seed keyword | Build keyword pool |
| `get_keyword_overview` | Keyword metrics overview | Search volume, competition, ad landscape |
| `get_keyword_trends` | Keyword search trends | Seasonality, trend comparison |
| `list_asin_keywords` | ASIN-related keywords | Reverse-engineer competitor keywords |
| `query_aba_keywords` | ABA trending keywords | Hot keywords, category discovery |
| `search_category` | Find category by name | Get node_id for category analysis |
| `get_asin_sales_daily_trend` | Daily sales trend | Short-term monitoring |
| `get_asins_sales_history` | Monthly sales history | Batch comparison |
| `get_variant_sales_30d` | Variant sales (30 days) | Variant structure analysis |
| `get_product` | Product details | Competitive baseline |
| `search_products` | Search by keyword | Product research |
| `get_product_reviews` | Product reviews | Review analysis |
| `get_category_best_sellers` | Category best sellers | Category analysis |
| `get_category_trend` | Category trends | Market dynamics |

---

## amazon.expand_keywords

Expand a seed keyword to find related Amazon keywords with search volume signals.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `keyword` | string | Yes | Seed keyword to expand |
| `market` | string | No | Target market: `US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT` (default: `US`) |

### CLI usage

```bash
apimux amazon expand_keywords --keyword "yoga mat"
apimux amazon expand_keywords --keyword "yoga mat" --market "DE"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `keyword` | string | Expanded keyword |
| `match_types` | string[] | Match type tags |
| `est_searches_num` | integer | Estimated monthly searches |
| `searches_rank` | integer | Search volume rank |

`meta.total` indicates total number of expanded keywords.

### Common errors

- Empty keyword → returns `missing_keyword`
- Invalid market → returns `invalid_market`

### Notes

- `keyword` is required and cannot be empty
- `market` defaults to `US` if omitted
- Returns an array of keyword objects

---

## amazon.get_keyword_overview

Get comprehensive metrics for a single Amazon keyword including search volume, competition, and ad landscape.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `keyword` | string | Yes | Keyword to analyze |
| `market` | string | No | Target market: `US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT` (default: `US`) |

### CLI usage

```bash
apimux amazon get_keyword_overview --keyword "yoga mat"
apimux amazon get_keyword_overview --keyword "standing desk" --market "UK"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `keyword` | string | Queried keyword |
| `est_searches_num` | integer | Estimated monthly searches |
| `searches_rank` | integer | Search volume rank |
| `searches` | integer/null | Raw search count |
| `demand_ratio` | number/null | Demand ratio |
| `competitor_cnt` | integer/null | Number of competing products |
| `sale_num` | integer/null | Sales-related metric |
| `ac_asin_num` | integer/null | Amazon's Choice ASIN count |
| `brand_ad_asin_num` | integer/null | Brand ad ASIN count |
| `sp_ad_asin_num` | integer/null | Sponsored Product ASIN count |
| `ppc_ad_asin_num` | integer/null | PPC ad ASIN count |
| `video_ad_asin_num` | integer/null | Video ad ASIN count |
| `er_asin_num` | integer/null | Editorial recommendation ASIN count |
| `nf_asin_num` | integer/null | New Flag ASIN count |
| `tr_asin_num` | integer/null | Top Rated ASIN count |
| `search_recommend_asin_num` | integer/null | Search recommendation ASIN count |
| `global_keyword_num` | integer/null | Global keyword count |
| `update_time` | string/null | Data update time |

### Common errors

- Empty keyword → returns `missing_keyword`
- Invalid market → returns `invalid_market`

### Notes

- Returns multi-dimensional data in a single call
- Missing metrics return `null` (not fabricated defaults)

---

## amazon.get_keyword_trends

Get historical search trends for multiple Amazon keywords.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `keywords` | string[] | Yes | Keyword list (max 1000) |
| `market` | string | No | Target market: `US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT` (default: `US`) |
| `granularity` | string | No | `week` or `month` (default: `month`) |

### CLI usage

```bash
apimux amazon get_keyword_trends --keywords "yoga mat,pilates ring"
apimux amazon get_keyword_trends --keywords "desk mat,monitor stand" --market "US" --granularity "week"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `keyword` | string | Queried keyword |
| `est_searches_num_history` | object | Time → search volume mapping |
| `searches_rank_history` | object | Time → rank mapping |

### Common errors

- Empty `keywords` → returns `missing_keywords`
- More than 1000 keywords → returns `too_many_keywords`
- Invalid granularity → returns `invalid_granularity`

### Notes

- `keywords` must be a string array (CLI converts comma-separated values)
- Monthly keys: `YYYY-MM`, weekly keys: `YYYY-MM-DD`

---

## amazon.list_asin_keywords

List keywords associated with an Amazon ASIN.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `asin` | string | Yes | Amazon ASIN, 10 uppercase alphanumeric characters |
| `keyword` | string | No | Keyword filter |
| `market` | string | No | Target market: `US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT`; default `US` |

### CLI usage

```bash
apimux amazon list_asin_keywords --asin "B0CM5JV26D"
apimux amazon list_asin_keywords --asin "B0CM5JV26D" --keyword "desk" --market "US"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `keyword` | string | Related keyword |
| `kw_characters` | string[] | Keyword tags |
| `conversion_characters` | string[] | Conversion tags |
| `exposure_type` | string[] | Exposure types |
| `last_rank` | string | Organic rank |
| `ad_last_rank` | string | Ad rank |
| `est_searches_num` | integer | Estimated search volume |
| `searches_rank` | integer | Search-volume rank |
| `ratio_score` | number | Traffic share score |

`meta.total` indicates the total keyword count.

### Common errors

- ASIN has invalid format → returns `invalid_asin`
- market is not supported → returns `invalid_market`

### Notes

- `asin` must be a 10-character alphanumeric ASIN. Lowercase input is normalized before validation.
- `keyword` is optional and filters the result set.

---

## amazon.query_aba_keywords

Query Amazon Brand Analytics trending keywords. Use this to discover keyword opportunities from ABA ranking data.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `keyword` | string | No | Keyword filter |
| `node_ids` | string | No | Comma-separated category node IDs |
| `page` | integer | No | Page number; default `1` |
| `page_size` | integer | No | Rows per page, `20..200`; default `40` |
| `market` | string | No | Target market: `US`, `UK`, `DE`, `FR`, `IN`, `CA`, `JP`, `ES`, `IT`, `MX`, `AE`, `AU`, `BR`, `SA`; default `US` |

### CLI usage

```bash
apimux amazon query_aba_keywords --keyword "yoga mat"
apimux amazon query_aba_keywords --node-ids "12345,67890" --page 2 --page-size 40 --market "US"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `keyword` | string | Trending keyword |
| `keyword_cn_name` | string | Chinese translation |
| `rank` | integer | Rank |
| `search_volume` | integer | Monthly search volume |
| `word_count` | integer | Word count |
| `product_count` | integer | Product count |
| `rank_change_of_weekly` | number | Weekly rank change |
| `cpc` | integer | Cost per click |
| `search_conversion_rate` | number | Search conversion rate |
| `click_of_90d` | integer | Clicks in the last 90 days |
| `sales_volume_of_90d` | integer | Sales volume in the last 90 days |
| `top3_asin` | string[] | Top 3 ASIN |
| `top3_brand` | string[] | Top 3 brands |
| `top3_category` | string[] | Top 3 categories |
| `season` | string | Seasonality tag |

`meta.current_page` and `meta.has_more` indicate pagination state.

### Common errors

- `page_size` outside `20..200` → returns `invalid_page_size`
- market is not supported → returns `invalid_market`

### Notes

- `keyword` and `node_ids` are optional. If both are omitted, the capability returns platform-level trending keywords.
- Returns a flat object array.

---

## amazon.get_asin_sales_daily_trend

Get the daily sales trend for one Amazon ASIN.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `asin` | string | Yes | Amazon ASIN, 10 uppercase alphanumeric characters |
| `begin_date` | string | No | Start date, `YYYY-MM-DD` |
| `market` | string | No | Target market: `US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT`; default `US` |

### CLI usage

```bash
apimux amazon get_asin_sales_daily_trend --asin "B0CM5JV26D"
apimux amazon get_asin_sales_daily_trend --asin "B0CM5JV26D" --begin-date "2026-04-01"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `date` | string | Date, `YYYY-MM-DD` |
| `sales` | integer | Sales for that day |

### Common errors

- ASIN has invalid format → returns `invalid_asin`
- `begin_date` is not an ISO date → returns `invalid_begin_date`

### Notes

- `begin_date` can be omitted. If provided, it must be `YYYY-MM-DD`.
- Results are sorted by date ascending.

---

## amazon.get_asins_sales_history

Get monthly sales history for multiple Amazon ASINs.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `asins` | string[] | Yes | Non-empty ASIN list, up to 10 ASINs |
| `market` | string | No | Target market: `US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT`; default `US` |
| `month` | integer | No | Most recent N months; `-1` or omitted means all history |

### CLI usage

```bash
apimux amazon get_asins_sales_history --asins "B0CM5JV26D,B0D1234567"
apimux amazon get_asins_sales_history --asins "B0CM5JV26D,B0D1234567" --month 6
apimux amazon get_asins_sales_history --asins "B0CM5JV26D,B0D1234567" --market "UK" --month 12
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `asin` | string | Amazon ASIN |
| `month` | string | Month, `YYYY-MM` |
| `sales` | integer | Monthly sales |

### Common errors

- `asins` is empty → returns `missing_asins`
- More than 10 ASINs → returns `too_many_asins`
- Any ASIN has invalid format → returns `invalid_asin`
- `month` is not `-1` or a positive integer → returns `invalid_month`

### Notes

- This is a batch query with a maximum of 10 ASINs.
- Missing history for some ASINs does not fail the whole request; only rows with available data are returned.
- Response rows use stable ordering by `asin + month`.
- Omit `month` or pass `-1` for all history; pass a positive integer for the most recent N months.

---

## amazon.get_variant_sales_30d

Get 30-day sales for variants associated with an Amazon ASIN.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `asin` | string | Yes | Parent or variant ASIN, 10 uppercase alphanumeric characters |
| `market` | string | No | Target market: `US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT`; default `US` |

### CLI usage

```bash
apimux amazon get_variant_sales_30d --asin "B0CM5JV26D"
apimux amazon get_variant_sales_30d --asin "B0CM5JV26D" --market "DE"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `asin` | string | Variant ASIN |
| `bought_in_past_month` | integer | Sales in the past 30 days |
| `update_time` | string | Update time |

### Common errors

- ASIN has invalid format → returns `invalid_asin`
- market is not supported → returns `invalid_market`

### Notes

- Input can be a parent ASIN or a specific variant ASIN.
- Response is not guaranteed to contain only one row; it usually returns the whole variant group.
- Stability note: this capability can be intermittently unavailable. Prefer calling it separately from critical workflows.

---

## amazon.get_product

Get details for one Amazon product.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `asin` | string | Yes | Amazon product ASIN, 10 uppercase alphanumeric characters |
| `market` | string | Yes | Target market: `US`, `UK`, `DE`, `FR`, `IN`, `CA`, `JP`, `ES`, `IT`, `MX`, `AE`, `AU`, `BR`, `SA` |

### CLI usage

```bash
apimux amazon get_product --asin "B0CM5JV26D" --market "US"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `asin` | string | Product ASIN |
| `market` | string | Amazon market code |
| `title` | string | Product title |
| `brand` | string | Brand name |
| `product_url` | string | Product URL |
| `main_image` | string | Main image URL |
| `rating` | float | Average rating, 0-5 |
| `review_count` | integer | Review count |
| `price.display` | string | Current price display text |
| `price.value` | number | Current price numeric value |
| `price.currency` | string | Current price currency |
| `brand_store.id` | string | Brand store ID |
| `brand_store.text` | string | Brand store display text |
| `brand_store.link` | string | Brand store URL |
| `feature_bullets` | array<string> | Product feature bullets |
| `images[].link` | string | Product image URL |
| `images[].variant` | string | Image variant label |
| `variants[].asin` | string | Variant ASIN |
| `variants[].title` | string | Variant title |
| `variants[].link` | string | Variant URL |
| `variants[].is_current_product` | boolean | Whether this is the current variant |
| `variants[].main_image` | string | Variant main image |
| `variants[].dimensions[].name` | string | Variant dimension name |
| `variants[].dimensions[].value` | string | Variant dimension value |
| `buybox.price.display` | string | Buy Box current price display text |
| `buybox.price.value` | number | Buy Box current price numeric value |
| `buybox.price.currency` | string | Buy Box current price currency |
| `buybox.original_price.display` | string | Buy Box original price display text |
| `buybox.original_price.value` | number | Buy Box original price numeric value |
| `buybox.original_price.currency` | string | Buy Box original price currency |
| `buybox.availability` | string | Buy Box availability |
| `buybox.is_prime` | boolean | Whether the Buy Box offer is Prime |

### Common errors

- Lowercase ASIN, such as `b0cm5jv26d` → automatically normalized to uppercase
- ASIN is not 10 characters → returns `validation_error`
- Unknown ASIN → use `search_products` first

### Notes

- ASIN must be 10 alphanumeric characters. Lowercase input is normalized; non-10-character input is rejected.
- `market` is required and does not default to `US`.
- If you do not have an ASIN, use `search_products` first.
- For multiple ASINs, call this capability in parallel.

---

## amazon.search_products

Search Amazon products by keyword.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `q` | string | Yes | Search keyword |
| `market` | string | Yes | Target market: `US`, `UK`, `DE`, `FR`, `IN`, `CA`, `JP`, `ES`, `IT` |
| `page` | integer | No | Page number; default `1`, maximum `10` |

### CLI usage

```bash
apimux amazon search_products --q "wireless earbuds" --market "US"
apimux amazon search_products --q "wireless earbuds" --market "US" --page 2
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `asin` | string | Product ASIN |
| `market` | string | Market |
| `title` | string | Product title |
| `brand` | string | Brand name |
| `product_url` | string | Product URL |
| `main_image` | string | Main image URL |
| `rating` | float | Average rating |
| `review_count` | integer | Review count |
| `price` | object | Price information |
| `position` | integer | Search result position |

`meta.current_page` identifies the current page.

### Common errors

- Empty search keyword → returns `validation_error`
- `page` greater than 10 → returns an error
- Treating search rows as complete product details → search rows are summaries; call `get_product` for full details

### Notes

- `q` is required and cannot be empty.
- `page` is capped at 10.
- Search results are summaries. Call `get_product` for full details.

---

## amazon.get_product_reviews

Get product reviews with multiple filters.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `asin` | string | Yes | Amazon product ASIN, 10 uppercase alphanumeric characters |
| `market` | string | Yes | Target market excluding JP: `US`, `UK`, `DE`, `FR`, `IN`, `CA`, `ES`, `IT`, `MX`, `AE`, `AU`, `BR`, `SA` |
| `start_date` | string | No | Review start date, ISO `YYYY-MM-DD` |
| `star` | string | No | `positive` for 4-5 star reviews or `negative` for 1-3 star reviews; omit for all reviews |
| `only_purchase` | boolean | No | `true` returns only verified-purchase reviews; default `false` |
| `page_index` | integer | No | Page number; default `1`, maximum `10` |

### CLI usage

```bash
# Get all reviews
apimux amazon get_product_reviews --asin "B0CM5JV26D" --market "US" --start-date "2025-01-01"

# Negative reviews only
apimux amazon get_product_reviews --asin "B0CM5JV26D" --market "US" --star "negative"

# Verified-purchase positive reviews only
apimux amazon get_product_reviews --asin "B0CM5JV26D" --market "US" --star "positive" --only-purchase true
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `star` | float | Rating, 1-5 |
| `title` | string | Review title |
| `content` | string | Review text |
| `date` | string | Review date, `YYYY-MM-DD` |
| `is_verified_purchase` | boolean | Whether the review is from a verified purchase |
| `helpful_votes` | integer | Helpful vote count |
| `reviewer_name` | string | Reviewer name |

Empty results return `ok=true, data=[]`; they are not errors.

### Common errors

- Numeric star values such as `"1"` or `"5"` → use `positive` or `negative`, not numbers
- String `"true"` for `only_purchase` → pass boolean `true`, not a string
- Date format such as `01/01/2025` → use ISO `YYYY-MM-DD`
- Market `JP` → returns validation error because JP reviews are not supported

### Notes

- `star` accepts only `positive` or `negative`; numeric star values are not accepted.
- `only_purchase` must be a boolean `true` or `false`; string `"true"` is rejected.
- `start_date` must be `YYYY-MM-DD`.
- JP is not supported for review queries.
- Filters are combined with AND semantics: `star=negative` + `only_purchase=true` means verified-purchase negative reviews only.

---

## amazon.search_category

Search Amazon categories by name and return reusable `node_id` values for category capabilities.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Category name; supports English or Chinese and partial matching |
| `market` | string | No | Target market: `US`, `UK`, `DE`, `FR`, `IN`, `CA`, `JP`, `ES`, `IT`, `MX`, `AE`, `AU`, `BR`, `SA`; default `US` |
| `limit` | integer | No | Number of results, `1..100`; default `20` |

### CLI usage

```bash
apimux amazon search_category --name "cell phone"
apimux amazon search_category --name "phone" --market "US" --limit 5
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `node_id` | string | Amazon category node ID |
| `name` | string | English category name |
| `cn_name` | string | Chinese category name |
| `path` | object[] | Full category path; each item includes `node_id` and `name` |

`meta.total` indicates the total match count.

### Common errors

- Empty `name` → returns `missing_name`
- `limit` outside `1..100` → returns `invalid_limit`
- market is not supported → returns `invalid_market`

### Notes

- Results are sorted by relevance.
- The response includes the full category path.
- Use returned `node_id` values directly with `get_category_best_sellers` and `get_category_trend`.

---

## amazon.get_category_best_sellers

Get best-selling products for one Amazon category.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `node_id` | string | Yes | Amazon category node ID; numeric string with at least 3 digits |
| `market` | string | Yes | Target market |
| `query_start` | string | No | Historical data start date, `YYYY-MM-DD` |
| `query_date` | string | No | Historical data end date, `YYYY-MM-DD` |
| `query_days` | integer | No | Number of days before `query_date`; maximum `365` |

### CLI usage

```bash
apimux amazon get_category_best_sellers --node-id "3743561" --market "US"
apimux amazon get_category_best_sellers --node-id "3743561" --market "US" --query-date "2026-01-01" --query-days 30
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `asin` | string | Product ASIN |
| `market` | string | Market |
| `rank` | integer | Best-seller rank; use this field instead of relying on array order |
| `title` | string | Product title |
| `brand` | string | Brand name |
| `main_image` | string | Main image URL |
| `rating` | float | Average rating |
| `rating_count` | integer | Rating count |
| `price` | object | Price information |
| `category` | string[] | Category path |
| `seller_count` | integer | Seller count |
| `is_fba` | boolean | Whether the product is FBA |

### Common errors

- Non-numeric `node_id`, such as `"electronics"` → must be numeric
- Missing category node → returns `category_not_found` (HTTP 404), not an upstream error
- Unknown `node_id` → use `search_category` first

### Notes

- `node_id` must be a numeric Amazon category node ID with at least 3 digits.
- Use `search_category` first when you do not know the node ID.
- Missing nodes return `category_not_found` (HTTP 404).
- `store_name` and `product_type` are occasional upstream supplement fields; do not rely on them as stable fields.
- `query_start` and `query_date` must be ISO `YYYY-MM-DD`.
- `query_days` is capped at 365.

---

## amazon.get_category_trend

Get category trend data. This capability can query multiple trend dimensions in one request.

**Important:** This capability can partially succeed when multiple dimensions are requested. Read the partial-failure semantics in [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) before using multi-dimension requests.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `node_id` | string | Yes | Amazon category node ID; numeric string with at least 3 digits |
| `market` | string | Yes | Target market |
| `trend_types` | string[] | Yes | Trend dimension names from the table below; supports multiple dimensions |

### Trend dimensions (`trend_types`)

| Name | Dimension |
|------|-----------|
| `sales_volume` | Sales volume |
| `brand_count` | Brand count |
| `seller_count` | Seller count |
| `avg_price` | Average price |
| `avg_rating_count` | Average review count |
| `avg_star` | Average rating |
| `new_product_ratio_1m` | New-product ratio in the past 1 month |
| `new_product_ratio_3m` | New-product ratio in the past 3 months |
| `amazon_self_ratio` | Amazon retail/self-operated ratio |
| `avg_profit` | Average profit |
| `top100_share` | Top 100 concentration |
| `top3_listing_monopoly` | Top 3 listing concentration |
| `top10_brand_monopoly` | Top 10 brand concentration |

### CLI usage

```bash
# Query one dimension
apimux amazon get_category_trend --node-id "3743561" --market "US" --trend-types "sales_volume"

# Query multiple dimensions, comma-separated
apimux amazon get_category_trend --node-id "3743561" --market "US" --trend-types "sales_volume,brand_count,avg_price"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `month` | string | Month, `YYYYMM` |
| `<requested trend_type>` | number/null | Requested trend dimension value; failed dimensions are `null` |

### Response format

Time series where each row contains a month and the requested dimension values:

```json
{
  "ok": true,
  "data": [
    {"month": "202601", "sales_volume": 1234, "brand_count": 56},
    {"month": "202602", "sales_volume": 1456, "brand_count": 61}
  ]
}
```

When requesting multiple dimensions, check `meta.partial` to see whether all dimensions succeeded. Failed dimensions are `null`; see `meta.subrequests` for details.

### Common errors

- Not checking `meta.partial` → some dimensions may have failed; `null` does not mean zero.
- Unknown dimension name → returns `invalid_trend_type`; use names from the table above.
- Numeric dimension such as `0` → pass dimension names as strings.
- Treating monthly data as daily data → granularity is monthly (`YYYYMM`); daily/weekly data is not supported.
- Duplicate dimension names → returns `duplicate_trend_type`.

### Notes

- `node_id` must be numeric and at least 3 digits.
- `trend_types` is required. Pass dimension names as strings, not numeric IDs.
- Multiple dimensions can be queried in one request.
- Duplicate dimension names are not supported.
- Data granularity is monthly (`YYYYMM`); daily/weekly granularity is not supported.
- With partial failure, `meta.partial=true` and failed dimensions appear as `null` in `data`.

---

## General notes

- See [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, and partial-failure semantics.
- Product, review, and category capabilities require explicit `market`; keyword capabilities can omit it and default to `US`.
- ASINs are 10-character alphanumeric IDs. Lowercase input is normalized to uppercase.
