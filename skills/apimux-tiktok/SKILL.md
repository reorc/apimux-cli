---
name: apimux-tiktok
version: 1.0.0
description: "TikTok content and TikTok Shop data. Search videos, analyze comments, list shop products, and inspect product details for content research and commerce analysis."
metadata:
  source: tiktok
  requires:
    bins: ["apimux"]
  cliHelp: "apimux tiktok --help"
---

# TikTok

Search TikTok content and inspect TikTok Shop product data. Use this for content research, creator/product analysis, shopping research, and cross-platform market validation.

**Before using:** Read [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, pagination metadata, and CLI conventions.

## What you can do

- **Find videos by topic** → `search_videos`
- **Inspect one video** → `get_video_detail`
- **Analyze comments under a video** → `list_comments`
- **List products from a TikTok Shop seller** → `shop_products`
- **Inspect one TikTok Shop product** → `shop_product_info`
- **Search TikTok Shop products by keyword** → `search_products`
- **Browse TikTok Shop product reviews** → `product_reviews`
- **Validate a market across sources** → use `search_videos` for content demand, then [`amazon.search_products`](../apimux-amazon/SKILL.md) for supply-side checks

## Available capabilities

| Capability | What it does | When to use |
|------------|--------------|-------------|
| `search_videos` | Search TikTok videos | Content research and competitor video discovery |
| `get_video_detail` | Get details for one TikTok video | Inspect metadata, author, engagement, and media URLs |
| `list_comments` | List video comments | Audience feedback and comment insights |
| `shop_products` | List seller products | Seller and product-selection analysis |
| `shop_product_info` | Get product details | Product research and cross-platform comparison |
| `search_products` | Search shop products by keyword | Keyword-level shop discovery and trend scouting |
| `product_reviews` | List product reviews | Product sentiment, rating distribution, UGC mining |

---

## tiktok.search_videos

Search TikTok videos by keyword.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `keyword` | string | Yes | Search keyword |
| `region` | string | No | ISO alpha-2 country code; omit for no region filter |
| `sort_by` | string | No | `relevance`, `likes`, or `date`; default `relevance` |
| `publish_time` | string | No | `all`, `1d`, `1w`, `1m`, `3m`, or `6m`; default `all` |
| `cursor` | integer | No | Pagination cursor; omit for the first page |
| `count` | integer | No | Number of results, `1..35` |

### CLI usage

```bash
apimux tiktok search_videos --keyword "desk setup"
apimux tiktok search_videos --keyword "desk setup" --sort-by "likes" --publish-time "1m" --region "US"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `video_id` | string | Video ID |
| `video_url` | string | Video playback URL |
| `description` | string | Video description |
| `create_time` | string | Publish time |
| `like_count` | integer | Like count |
| `comment_count` | integer | Comment count |
| `share_count` | integer | Share count |
| `play_count` | integer | Play count |
| `cover_image` | string | Cover image URL |
| `duration` | integer | Video duration in seconds |
| `region` | string | Region information |
| `is_ad` | boolean | Whether the video is an ad |
| `author` | object | Author summary |
| `music` | object | Music summary |

### Notes

- `keyword` is required.
- Prefer the string enum values above. For backward compatibility, the CLI also accepts legacy numeric values for `sort_by` (`0/1/2`) and `publish_time` (`0/1/7/30/90/180`).
- `count` must be in `1..35`.
- Pagination state is returned in `meta.cursor` and `meta.has_more`.

---

## tiktok.get_video_detail

Get details for one TikTok video by share URL or video ID.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `share_url` | string | Conditional | TikTok share URL. Use either `share_url` or `aweme_id`, not both. |
| `aweme_id` | string | Conditional | Numeric TikTok video ID. Use either `aweme_id` or `share_url`, not both. |
| `region` | string | No | ISO alpha-2 country code for `aweme_id` lookup only; default `US`. Do not pass with `share_url`. |

### CLI usage

```bash
apimux tiktok get_video_detail --share-url "https://www.tiktok.com/t/ZTFNEj8Hk/"
apimux tiktok get_video_detail --aweme-id "7350810998023949599" --region "US"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `platform` | string | Always `tiktok` |
| `video_id` | string | TikTok video ID |
| `caption` | string | Video caption |
| `create_time` | string | RFC3339 publish time |
| `share_url` | string | Share URL |
| `video_url` | string | Signed video playback URL |
| `cover` | string | Cover image URL |
| `duration` | integer | Video duration in milliseconds |
| `author` | object | Author summary |
| `stats` | object | `play`, `digg`, `comment`, and `share` counts |
| `music` | object | Music summary when available |
| `region` | string | Video region when available |

### Notes

- Provide exactly one identity input: `share_url` or `aweme_id`.
- `region` is valid only with `aweme_id`; passing it with `share_url` returns a validation error.
- `video_url` and cover/music URLs are signed provider URLs and may expire. Download or store them immediately if you need durable media access.
- Missing videos return `video_not_found`.

---

## tiktok.list_comments

List comments for one TikTok video.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `video_id` | string | Yes | Numeric video ID |
| `cursor` | integer | No | Pagination cursor; omit for the first page |
| `count` | integer | No | Number of comments, `1..50` |

### CLI usage

```bash
apimux tiktok list_comments --video-id "7489123456789012345"
apimux tiktok list_comments --video-id "7489123456789012345" --count 20
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `comment_id` | string | Comment ID |
| `video_id` | string | Video ID |
| `text` | string | Comment text |
| `create_time` | string | Comment time |
| `like_count` | integer | Like count |
| `reply_count` | integer | Reply count |
| `images` | string[] | Comment images |
| `author` | object | Comment author summary |

### Notes

- `video_id` must be a numeric string; full TikTok URLs are not accepted.
- `count` must be in `1..50`.
- Empty comment lists return `ok=true`.
- Pagination state is returned in `meta.cursor`, `meta.has_more`, and `meta.total`.

---

## tiktok.shop_products

List products from a TikTok Shop seller.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `seller_id` | string | Yes | TikTok Shop seller ID |
| `region` | string | No | Only `US` is supported; default `US` |
| `sort` | string | No | `sale` or `rec`; default `rec` |
| `top_n` | integer | No | Number of products, `1..200`; default `20` |

### CLI usage

```bash
apimux tiktok shop_products --seller-id "123456789"
apimux tiktok shop_products --seller-id "123456789" --sort "sale" --top-n 40
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `product_id` | string | Product ID |
| `product_name` | string | Product name |
| `product_cover` | string | Product cover image |
| `product_sold_count` | integer | Sold count |
| `format_available_price` | string | Current price text |
| `format_origin_price` | string | Original price text |
| `discount` | string | Discount text |

### Notes

- `seller_id` is required.
- `region` currently supports only `US`.
- `sort` accepts only `sale` or `rec`.
- `top_n` must be in `1..200`.
- Default compact output is columnar: `{columns, rows}`.

---

## tiktok.shop_product_info

Get details for one TikTok Shop product.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `product_id` | string | Yes | Product ID |
| `region` | string | No | `US`, `GB`, `SG`, `MY`, `PH`, `TH`, `VN`, or `ID`; default `US` |

### CLI usage

```bash
apimux tiktok shop_product_info --product-id "1729384756"
apimux tiktok shop_product_info --product-id "1729384756" --region "GB"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `product_id` | string | Product ID |
| `product_name` | string | Product name |
| `status` | integer | Product status |
| `seller_id` | string | Seller ID |
| `seller_name` | string | Seller name |
| `sold_count` | integer | Sold count |
| `rating` | number | Rating |
| `original_price` | string | Original price |
| `real_price` | string | Current price |
| `discount` | string | Discount text |
| `images` | string[] | Product images |
| `is_platform_product` | boolean | Whether it is a platform product |
| `review_count` | integer | Review count |

### Notes

- `product_id` is required.
- `region` defaults to `US` and accepts the 8 supported markets listed above.
- Missing products return `product_not_found`.

---

## tiktok.search_products

Search TikTok Shop products by keyword.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `keyword` | string | Yes | Search keyword |
| `region` | string | No | `US`, `GB`, `SG`, `MY`, `PH`, `TH`, `VN`, or `ID`; default `US` |
| `cursor` | string | No | Pagination token from previous `meta.cursor` |
| `offset` | integer | No | Result offset, `>= 0`; default `0` |
| `count` | integer | No | Number of products, `1..200`; default `20` |

### CLI usage

```bash
apimux tiktok search_products --keyword "labubu"
apimux tiktok search_products --keyword "labubu" --region "GB" --count 40
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `product_id` | string | Product ID |
| `product_name` | string | Product name |
| `product_cover` | string | Product cover image |
| `product_sold_count` | integer | Sold count |
| `format_available_price` | string | Current price text |
| `format_origin_price` | string | Original price text |
| `discount` | string | Discount text |
| `rating` | number | Product rating when available |
| `review_count` | integer | Review count when available |

### Notes

- `keyword` is required.
- `region` must be one of the 8 supported markets.
- Pagination state is returned in `meta.cursor` and `meta.has_more`.

---

## tiktok.product_reviews

List reviews for one TikTok Shop product.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `product_id` | string | Yes | TikTok Shop product ID |
| `region` | string | No | `US`, `GB`, `SG`, `MY`, `PH`, `TH`, `VN`, or `ID`; default `US` |
| `page` | integer | No | Page number starting at `1`; default `1` |
| `sort` | string | No | `default` or `latest`; default `default` |
| `media_filter` | string | No | `all`, `media`, or `verified`; default `all` |
| `star` | string | No | `all` or `1`..`5`; default `all` |
| `count` | integer | No | Number of reviews, `1..100`; default `20` |

### CLI usage

```bash
apimux tiktok product_reviews --product-id "1729556436942358002"
apimux tiktok product_reviews --product-id "1729556436942358002" --sort "latest" --star "5" --media-filter "media"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `review_id` | string | Review ID |
| `rating` | integer | Star rating, `1..5` |
| `content` | string | Review content |
| `create_time` | string | RFC3339 timestamp when available |
| `verified_purchase` | boolean | Verified purchase indicator |
| `like_count` | integer | Likes on this review |
| `medias` | object[] | Review media entries `{type, url}` |
| `author` | object | `user_id`, `nickname`, `avatar` |
| `seller_reply` | object | Seller reply payload when present |

### Notes

- `product_id` is required.
- The underlying provider endpoint is documented for Americas and Europe (e.g.
  `US`, `GB`). For Southeast Asia regions (`SG`, `MY`, `PH`, `TH`, `VN`, `ID`)
  results may be empty.
- Pagination state is returned in `meta.has_more`, `meta.total`, and
  `meta.cursor` (the current page number).
- Review summary is surfaced in `meta.resolved_filters.average_rating` and
  `meta.resolved_filters.star_distribution` when present.

---

## General notes

- See [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure and error handling.
- TikTok Shop capabilities support the TikHub region list: `US, GB, SG, MY, PH, TH, VN, ID`. `shop_products` currently remains US-only; the other shop capabilities accept the full list with `US` as the default.

