---
name: apimux-tiktok
version: 1.0.0
description: "TikTok 内容与 TikTok Shop 数据查询。提供视频搜索、评论分析、店铺商品列表、商品详情等能力。适用于内容研究、达人分析、带货选品、跨平台市场验证等场景。"
metadata:
  source: tiktok
  requires:
    bins: ["apimux"]
  cliHelp: "apimux tiktok --help"
---

# TikTok

**CRITICAL — 开始前 MUST 先用 Read 工具读取 [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md)，其中包含响应结构、错误处理等共享规则。**

TikTok 数据查询，覆盖内容侧和 TikTok Shop 两条分析路径。

## 快速决策

- 想找某个主题下的热视频 → `search_videos`
- 想分析某个视频下的用户反馈 → `list_comments`
- 想看某个 TikTok Shop 卖家的带货商品 → `shop_products`
- 想看某个商品的详细信息 → `shop_product_info`
- 想做跨平台市场验证 → 先 `search_videos` 看内容热度，再转到 [`amazon.search_products`](../apimux-amazon/SKILL.md) 看供给侧

## Capabilities 概览

| Capability | 说明 | 典型场景 |
|------------|------|----------|
| `search_videos` | 视频搜索 | 内容研究、竞品视频发现 |
| `list_comments` | 视频评论 | 用户反馈、评论洞察 |
| `shop_products` | 店铺商品列表 | 达人带货分析、选品 |
| `shop_product_info` | 商品详情 | 商品研究、跨平台对比 |

## Agent Journeys

### Journey 1: 内容分析
```
search_videos → list_comments
```
先找热视频，再对目标视频拉评论。

### Journey 2: TikTok Shop 选品
```
shop_products → shop_product_info
```
先看卖家商品列表，再下钻单品详情。

### Journey 3: 跨平台市场验证
```
google_trends.get_interest_over_time → search_videos → amazon.search_products
```
先看需求信号，再比对 TikTok 内容热度和 Amazon 供给侧。

---

## tiktok.search_videos

按关键词搜索 TikTok 视频。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `keyword` | string | 是 | 搜索关键词 |
| `region` | string | 否 | 地区代码，ISO 两位国家码 |
| `sort_by` | string | 否 | `relevance`、`likes`、`date` |
| `publish_time` | string | 否 | `all`、`1d`、`1w`、`1m`、`3m`、`6m` |
| `cursor` | integer | 否 | 分页 cursor |
| `count` | integer | 否 | 返回数量，范围 1-35 |

### CLI 用法

```bash
apimux tiktok search_videos --keyword "desk setup"
apimux tiktok search_videos --keyword "desk setup" --sort-by "likes" --publish-time "1m" --region "US"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `video_id` | string | 视频 ID |
| `video_url` | string | 视频播放地址 |
| `description` | string | 视频描述 |
| `create_time` | string | 发布时间 |
| `like_count` | integer | 点赞数 |
| `comment_count` | integer | 评论数 |
| `share_count` | integer | 分享数 |
| `play_count` | integer | 播放数 |
| `cover_image` | string | 封面图 |
| `duration` | integer | 视频时长（秒） |
| `region` | string | 地区信息（若 provider 返回） |
| `is_ad` | boolean | 是否广告（若 provider 返回） |
| `author` | object | 作者摘要信息 |
| `music` | object | 背景音乐摘要信息 |

### 规则

- `keyword` 必填
- `sort_by` 只接受 `relevance`、`likes`、`date`
- `publish_time` 只接受 `all`、`1d`、`1w`、`1m`、`3m`、`6m`
- `count` 范围 1-35
- provider 的数值枚举不会暴露到 contract
- 分页状态放在 `meta.cursor` 和 `meta.has_more`，不是 `data`

---

## tiktok.list_comments

获取 TikTok 视频评论列表。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `video_id` | string | 是 | 视频 ID，必须是纯数字字符串 |
| `cursor` | integer | 否 | 分页 cursor |
| `count` | integer | 否 | 返回数量，范围 1-50 |

### CLI 用法

```bash
apimux tiktok list_comments --video-id "7489123456789012345"
apimux tiktok list_comments --video-id "7489123456789012345" --count 20
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `comment_id` | string | 评论 ID |
| `video_id` | string | 视频 ID |
| `text` | string | 评论正文 |
| `create_time` | string | 评论时间 |
| `like_count` | integer | 点赞数 |
| `reply_count` | integer | 回复数 |
| `images` | string[] | 评论图片 |
| `author` | object | 评论作者摘要信息 |

### 规则

- `video_id` 必须是纯数字字符串，不接受完整 TikTok URL
- `count` 范围 1-50
- 空评论列表返回 `ok=true`
- 分页状态放在 `meta.cursor`、`meta.has_more`、`meta.total`，不是 `data`

---

## tiktok.shop_products

获取某个 TikTok Shop 卖家的商品列表。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `seller_id` | string | 是 | TikTok Shop 卖家 ID |
| `region` | string | 否 | 仅支持 `US`，默认 `US` |
| `sort` | string | 否 | `sale` 或 `rec` |
| `top_n` | integer | 否 | 返回商品数量，范围 1-200，默认 20 |

### CLI 用法

```bash
apimux tiktok shop_products --seller-id "123456789"
apimux tiktok shop_products --seller-id "123456789" --sort "sale" --top-n 40
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `product_id` | string | 商品 ID |
| `product_name` | string | 商品名 |
| `product_cover` | string | 商品主图 |
| `product_sold_count` | integer | 销量 |
| `format_available_price` | string | 当前价格 |
| `format_origin_price` | string | 原价 |
| `discount` | string | 折扣信息 |

### 规则

- `seller_id` 必填
- `region` 仅支持 `US`
- `sort` 只接受 `sale` 或 `rec`
- `top_n` 范围 1-200
- 内部分页 fan-out 对调用方透明
- 返回的 `data` 直接是商品列表，不包一层 `items`

---

## tiktok.shop_product_info

获取单个 TikTok Shop 商品详情。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `product_id` | string | 是 | 商品 ID |
| `region` | string | 否 | 仅支持 `US`，默认 `US` |

### CLI 用法

```bash
apimux tiktok shop_product_info --product-id "1729384756"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `product_id` | string | 商品 ID |
| `product_name` | string | 商品名 |
| `status` | integer | 商品状态 |
| `seller_id` | string | 卖家 ID |
| `seller_name` | string | 卖家名称 |
| `sold_count` | integer | 销量 |
| `rating` | number | 评分 |
| `original_price` | string | 原价 |
| `real_price` | string | 当前价格 |
| `discount` | string | 折扣信息 |
| `images` | string[] | 商品图片 |
| `is_platform_product` | boolean | 是否平台商品 |

### 规则

- `product_id` 必填
- `region` 仅支持 `US`
- 商品不存在时返回 canonical not-found error

---

## 通用规则

- **service 使用 canonical envelope，CLI 默认 data-only**：详见 [apimux-shared](../apimux-shared/SKILL.md)
- **不暴露 provider 内部信息**：不会暴露 TikLiveAPI / PrimeAPI / YimianData 名称
- **TikTok Shop 第一阶段仅支持 US**：不要对非 US 市场调用 `shop_products` 或 `shop_product_info`
