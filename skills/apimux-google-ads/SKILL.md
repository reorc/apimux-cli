---
name: apimux-google-ads
version: 1.0.0
description: "Google Ads Transparency Center 查询。提供广告主搜索、广告素材列表、广告详情能力，适用于广告样本检索、竞品创意分析、素材下钻等场景。"
metadata:
  source: google_ads
  requires:
    bins: ["apimux"]
  cliHelp: "apimux google_ads --help"
---

# Google Ads

**CRITICAL — 开始前 MUST 先用 Read 工具读取 [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md)，其中包含响应结构、错误处理等共享规则。**

Google Ads Transparency Center 数据查询，覆盖广告主发现、广告素材列表和广告详情。

## 快速决策

- 想先找广告主 / 域名 → `search_advertisers`
- 已有广告主 ID 或域名，想列出广告素材 → `list_ad_creatives`
- 已有 `advertiser_id + creative_id`，想看单条创意详情 → `get_ad_details`

## Capabilities 概览

| Capability | 说明 | 典型场景 |
|------------|------|----------|
| `search_advertisers` | 搜索广告主与域名 | 广告主发现、品牌确认 |
| `list_ad_creatives` | 列出广告素材 | 创意样本采集、过滤 |
| `get_ad_details` | 获取广告详情 | 定向与变体下钻 |

## Agent Journey

```
search_advertisers → list_ad_creatives → get_ad_details
```

先搜索广告主，再列素材，最后对目标 `creative_id` 拉详情。

---

## google_ads.search_advertisers

搜索 Google Ads 广告主与相关域名。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `query` | string | 是 | 搜索关键词 |
| `region` | string | 否 | ISO alpha-2 国家码 |
| `num_advertisers` | integer | 否 | 广告主数量，1-100 |
| `num_domains` | integer | 否 | 域名数量，1-100 |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `advertisers` | object[] | 广告主列表 |
| `domains` | object[] | 域名列表 |

### 规则

- `query` 必填
- `region` 必须是 ISO alpha-2
- `num_advertisers` / `num_domains` 范围 1-100

---

## google_ads.list_ad_creatives

按广告主或域名列出广告素材。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `advertiser_id` | string | 否 | 广告主 ID，必须以 `AR` 开头 |
| `domain` | string | 否 | 广告主域名 |
| `region` | string | 否 | ISO alpha-2 国家码 |
| `platform` | string | 否 | `google_play`、`google_maps`、`google_search`、`youtube`、`google_shopping` |
| `ad_format` | string | 否 | `text`、`image`、`video` |
| `time_period` | string | 否 | `last_7_days`、`last_30_days`、`last_90_days`、`last_year` |
| `page_token` | string | 否 | 分页 token |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `creative_id` | string | 素材 ID |
| `advertiser_id` | string | 广告主 ID |
| `advertiser_name` | string | 广告主名称 |
| `target_domain` | string | 落地页域名 |
| `format` | string | 素材格式 |
| `first_shown_datetime` | string | 首次展示时间 |
| `last_shown_datetime` | string | 最后展示时间 |
| `total_days_shown` | integer | 素材累计展示天数 |
| `details_link` | string | Google 详情页链接 |

### 规则

- `advertiser_id` 或 `domain` 至少提供一个
- `advertiser_id` 必须以 `AR` 开头
- `platform` / `ad_format` / `time_period` 只接受批准的字符串枚举
- 分页状态放在 `meta.page_token`

---

## google_ads.get_ad_details

获取单条 Google Ads 广告详情。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `advertiser_id` | string | 是 | 广告主 ID，必须以 `AR` 开头 |
| `creative_id` | string | 是 | 素材 ID，必须以 `CR` 开头 |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `ad_information` | object | 广告元信息与定向信息 |
| `variations` | object[] | 广告创意变体 |

### 规则

- `advertiser_id` 必须以 `AR` 开头
- `creative_id` 必须以 `CR` 开头
- 素材不存在时返回 canonical `ad_not_found`
- provider 名称不会出现在响应里

---

## 通用规则

- **service 使用 canonical envelope，CLI 默认 data-only**：详见 [apimux-shared](../apimux-shared/SKILL.md)
- **分页 token 在 `meta`**：`list_ad_creatives` 的下一页 token 在 `meta.page_token`
- **不暴露 provider 内部信息**：不会暴露真实上游系统名称
