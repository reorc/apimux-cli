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
- **Analyze comments under a video** → `list_comments`
- **List products from a TikTok Shop seller** → `shop_products`
- **Inspect one TikTok Shop product** → `shop_product_info`
- **Validate a market across sources** → use `search_videos` for content demand, then [`amazon.search_products`](../apimux-amazon/SKILL.md) for supply-side checks

## Available capabilities

| Capability | What it does | When to use |
|------------|--------------|-------------|
| `search_videos` | Search TikTok videos | Content research and competitor video discovery |
| `list_comments` | List video comments | Audience feedback and comment insights |
| `shop_products` | List seller products | Seller and product-selection analysis |
| `shop_product_info` | Get product details | Product research and cross-platform comparison |

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
| `region` | string | No | Only `US` is supported; default `US` |

### CLI usage

```bash
apimux tiktok shop_product_info --product-id "1729384756"
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
- `region` currently supports only `US`.
- Missing products return `product_not_found`.

---

## General notes

- See [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure and error handling.
- TikTok Shop capabilities currently support only the US market.
