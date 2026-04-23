---
name: apimux-trendcloud
version: 1.0.0
description: "TrendCloud 中国电商渠道趋势查询。覆盖抖音/京东/天猫的部分市场数据，适合做品牌/品类/系列/SKU/属性层面的份额、排行、趋势与筛选值 discovery，不适合把它当作全网或全平台市场数据源。"
metadata:
  source: trendcloud
  requires:
    bins: ["apimux"]
  cliHelp: "apimux trendcloud --help"
---

# TrendCloud

**CRITICAL — 开始前 MUST 先用 Read 工具读取 [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md)，其中包含响应结构、错误处理等共享规则。**

TrendCloud 是一个**中国电商渠道的部分市场数据源**，当前只覆盖这些平台：
- `douyin`（抖音）
- `jd`（京东）
- `tmall`（天猫）

它更适合回答这类问题：
- 某个品牌 / 品类 / 系列 / SKU 在中国主流电商渠道里的**销售额、销量、份额、排行、价格带**怎么样
- 某个类目在抖音、京东、天猫之间的**趋势差异**是什么
- 我手里的 filter 名称不确定，想先做 **category / brand / series / sku / attribute discovery**

它**不适合**直接回答这类问题：
- “中国整个市场”或“全网电商”总体规模
- Amazon / TikTok Shop / Reddit / Google Trends 之类非中国电商渠道的数据
- 用户画像、创意内容、评论语义、VOC 这类内容侧问题

如果用户的问题本质上是：
- 中国电商渠道里的**品牌/类目格局、排行、份额、趋势** → 优先考虑 TrendCloud
- 内容热度 / 搜索需求 / 评论反馈 → 应改用其他 source（如 TikTok、Xiaohongshu、Google Trends、Amazon reviews）

## 快速决策

- 不确定某个品类/品牌/系列/SKU/属性怎么写，先用 `search_filter_values`
- 想看一段时间内销售额或销量走势，使用 `get_market_trend`
- 想看品牌/品类/系列/SKU/属性排行榜，使用 `get_top_rankings`
- 如果用户说“抖音 / 京东 / 天猫 哪个平台更强”，优先把 `filters.platforms` 显式带上
- 如果用户说“中国市场”，先确认他要的是 **TrendCloud 覆盖的中国电商渠道数据**，不要自动上升成“全中国全渠道”

## 推荐工作流

最常见的可靠流程是：

1. 先用 `search_filter_values` 找到精确 filter 值
2. 再把返回的 `label/path` 填进 `filters`
3. 用 `get_market_trend` 看趋势
4. 用 `get_top_rankings` 看头部格局 / 份额 / 排行

尤其当用户给的是自然语言词汇时（如“咖啡”“瑞幸”“挂耳”），**不要直接猜 filter**；先做 discovery，命中率和解释性都会更好。

## Capabilities 概览

| Capability | 说明 | 典型场景 |
|------------|------|----------|
| `get_market_trend` | 市场趋势时序 | 中国电商渠道内的销售额 / 销量走势、平台对比 |
| `get_top_rankings` | 排行榜查询 | 品牌份额、类目头部格局、系列 / SKU 排行 |
| `search_filter_values` | 筛选值发现 | category / brand / series / sku / attribute 的 filter discovery |

## trendcloud.get_market_trend

获取 TrendCloud 月度市场趋势。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `start_month` | string | 否 | 起始月份，`YYYY-MM` |
| `end_month` | string | 否 | 结束月份，`YYYY-MM` |
| `metrics` | string[] | 否 | 指标枚举：`sales`、`volume`，默认 `["sales"]` |
| `filters` | object | 否 | 结构化过滤条件，支持 `platforms/categories/brands/series/skus/attributes` |

### CLI 用法

```bash
apimux trendcloud get_market_trend --start-month "2025-01" --end-month "2025-12"
apimux trendcloud get_market_trend --metrics "sales,volume" --filters-json '{"platforms":["douyin"],"brands":["瑞幸"]}'
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `period` | string | 月份，`YYYY-MM` |
| `sales` | number | 销售额，单位元 |
| `volume` | integer | 销量 |

## trendcloud.get_top_rankings

获取 TrendCloud 排行榜结果。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `entity` | string | 是 | 排行实体：`brand/category/series/sku/attribute` |
| `metric` | string | 否 | 主排序指标，按 entity 约束 |
| `start_month` | string | 否 | 起始月份，`YYYY-MM` |
| `end_month` | string | 否 | 结束月份，`YYYY-MM` |
| `top_n` | integer | 否 | 返回数量，1-100，默认 20 |
| `category_level` | string | 否 | 仅 entity=`category` 时可用：`category1/category2/category3` |
| `filters` | object | 否 | 结构化过滤条件 |

### CLI 用法

```bash
apimux trendcloud get_top_rankings --entity "brand" --metric "sales"
apimux trendcloud get_top_rankings --entity "category" --category-level "category2" --filters-json '{"platforms":["tmall"]}'
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `rank` | integer | 1-based 排名 |
| `label` | string | 实体展示名称 |
| `sales` | number | 销售额，单位元 |
| `volume` | integer | 销量 |
| `market_share` | number | 市占比 |
| `avg_price` | number | 平均价格，单位元 |
| `sales_change_ratio` | number | 销售同比变化 |
| `volume_change_ratio` | number | 销量同比变化 |
| `market_share_change_ratio` | number | 市占比同比变化 |

## trendcloud.search_filter_values

搜索 TrendCloud 的筛选值候选。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `kind` | string | 是 | `category/brand/series/sku/attribute` |
| `query` | string | 是 | 搜索关键词 |
| `platforms` | string[] | 否 | 平台范围：`douyin/jd/tmall` |
| `categories` | string[] | 否 | 类目 hint |
| `limit` | integer | 否 | 返回数量，1-50，默认 10 |

### CLI 用法

```bash
apimux trendcloud search_filter_values --kind "category" --query "咖啡"
apimux trendcloud search_filter_values --kind "brand" --query "luckin" --platforms "douyin,jd"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `label` | string | 显示名称 |
| `path` | string[] | 树字段返回完整层级路径；简单字段返回 `[label]` |

## 常见错误

- 平台不在 `douyin/jd/tmall` 之内 → `invalid_platform`
- 时间范围格式不是 `YYYY-MM` 或跨度超过 12 个月 → `invalid_time_range`
- 某个 filter 无法唯一匹配 → `ambiguous_filter`
- 某个 filter 完全不存在 → `invalid_filter`
- 并发已满 → `busy`，带 `retry_after_seconds`

## 规则

- 金额字段统一返回 **元**，不会暴露上游分值
- `meta.resolved_time_range` 会说明默认值或 clamp 行为
- `meta.resolved_filters` 会说明 filter 实际解析结果
- 当 tree filter 有歧义时，先用 `search_filter_values` 做 discovery，再重试正式查询
