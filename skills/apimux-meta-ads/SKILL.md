---
name: apimux-meta-ads
version: 1.0.0
description: "Meta Ads Library 查询。提供广告搜索与广告详情能力，适用于广告创意研究、竞品分析、投放样本采集等场景。"
metadata:
  source: meta_ads
  requires:
    bins: ["apimux"]
  cliHelp: "apimux meta_ads --help"
---

# Meta Ads

**CRITICAL — 开始前 MUST 先用 Read 工具读取 [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md)，其中包含响应结构、错误处理等共享规则。**

Meta Ads Library 数据查询，覆盖广告搜索与广告详情两个核心能力。

## 快速决策

- 想按关键词找广告创意样本 → `search_ads`
- 想下钻单条广告详情 → `get_ad_detail`
- 想做广告样本 fan-out 分析 → 先 `search_ads`，再对目标 `ad_id` 调 `get_ad_detail`

## Capabilities 概览

| Capability | 说明 | 典型场景 |
|------------|------|----------|
| `search_ads` | 搜索 Meta Ads Library 广告 | 广告创意研究、竞品广告发现 |
| `get_ad_detail` | 获取单条广告详情 | EU 透明度信息、政治广告详情下钻 |

## Agent Journey

```
search_ads → get_ad_detail
```

先搜索广告，再对感兴趣的 `ad_id` 拉详情。

---

## meta_ads.search_ads

按关键词搜索 Meta Ads Library 广告。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `q` | string | 是 | 搜索关键词 |
| `country` | string | 否 | ISO alpha-2 国家码；不传则不限地区 |
| `ad_type` | string | 否 | `all`、`political_and_issue_ads`、`housing_ads`、`employment_ads`、`credit_ads`；默认 `all` |
| `active_status` | string | 否 | `active`、`inactive`、`all`；默认 `all` |
| `media_type` | string | 否 | `all`、`video`、`image`、`meme`、`image_and_meme`、`none`；默认 `all` |
| `platforms` | string | 否 | 逗号分隔平台名：`facebook,instagram` 等；不传则不限平台 |
| `start_date` | string | 否 | 开始日期，`YYYY-MM-DD`；不传则不限起始时间 |
| `end_date` | string | 否 | 结束日期，`YYYY-MM-DD`；不传则不限截止时间 |
| `next_page_token` | string | 否 | 分页 token，首页不传 |

### CLI 用法

```bash
apimux meta_ads search_ads --q "fitness app"
apimux meta_ads search_ads --q "fitness app" --country "US" --media-type "video"
apimux meta_ads search_ads --q "fitness app" --platforms "facebook,instagram" --start-date "2026-01-01"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `ad_id` | string | 广告 ID |
| `page_id` | string | Page ID |
| `page_name` | string | Page 名称 |
| `start_date` | string | 广告开始时间 |
| `end_date` | string | 广告结束时间 |
| `is_active` | boolean | 是否活跃 |
| `categories` | string[] | 广告类别 |
| `publisher_platforms` | string[] | 广告投放平台，小写 canonical 值 |
| `snapshot` | object | 创意摘要，包括正文、标题、链接、卡片、视频等 |

### 规则

- `q` 必填
- `country` 必须是 ISO alpha-2
- `ad_type` / `active_status` / `media_type` 只接受批准的字符串枚举
- `platforms` 必须是逗号分隔的小写平台名
- `start_date` / `end_date` 必须是 `YYYY-MM-DD`
- 分页状态放在 `meta.next_page_token`
- provider 名称不会出现在响应里

---

## meta_ads.get_ad_detail

获取单条 Meta 广告详情。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `ad_id` | string | 是 | 广告 archive ID |

### CLI 用法

```bash
apimux meta_ads get_ad_detail --ad-id "477570185419072"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `ad_id` | string | 广告 ID |
| `eu_transparency` | object | EU 透明度信息 |
| `political_insights` | object | 政治广告洞察信息（若 provider 返回） |
| `verified_voice` | object | verified voice 信息（若 provider 返回） |

### 规则

- `ad_id` 必填
- 广告不存在时返回 canonical `ad_not_found`
- contract 不要求必须先调用 `search_ads`
- provider 名称不会出现在响应里

---

## 通用规则

- **service 使用 canonical envelope，CLI 默认 data-only**：详见 [apimux-shared](../apimux-shared/SKILL.md)
- **分页 token 在 `meta`**：`search_ads` 的下一页 token 在 `meta.next_page_token`
- **不暴露 provider 内部信息**：不会暴露真实上游系统名称
