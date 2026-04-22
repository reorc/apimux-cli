---
name: apimux-google-trends
version: 1.0.0
description: "Google Trends 搜索热度查询。获取关键词的搜索热度趋势（0-100 归一化指数），支持地区、时间范围、Google 属性过滤。适用于市场需求验证、趋势发现、关键词对比分析等场景。"
metadata:
  source: google_trends
  requires:
    bins: ["apimux"]
  cliHelp: "apimux google_trends --help"
---

# Google Trends

**CRITICAL — 开始前 MUST 先用 Read 工具读取 [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md)，其中包含响应结构、错误处理等共享规则。**

Google 搜索热度数据查询，用于验证市场需求信号和趋势分析。

## 快速决策

- 想验证某个产品/关键词是否有市场需求 → `get_interest_over_time`
- 想对比多个关键词的热度 → `get_interest_over_time`（最多 5 个关键词）
- 想了解长期趋势（季节性、增长/衰退）→ `get_interest_over_time` + `time="today 5-y"`
- 确认需求信号后想评估供给侧 → 转到 [`amazon.search_products`](../apimux-amazon/SKILL.md)

## Capabilities 概览

| Capability | 说明 | 典型场景 |
|------------|------|----------|
| `get_interest_over_time` | 搜索热度时间序列 | 需求验证、趋势对比 |

Phase 1 只包含 `get_interest_over_time`。后续可能扩展 `get_interest_by_region`、`get_related_topics`、`get_related_queries`。

## Agent Journeys

### Journey: 市场验证（跨 skill）
```
google_trends.get_interest_over_time → amazon.search_products → (optional) amazon.get_category_trend
```
本 skill 是市场验证的起点。先确认需求端信号（搜索热度是否在增长），再去 [Amazon](../apimux-amazon/SKILL.md) 评估供给侧。

典型判断逻辑：
- 热度持续上升 → 需求增长，值得进入 Amazon 供给分析
- 热度平稳 → 成熟市场，关注竞争格局
- 热度下降 → 需求萎缩，谨慎进入
- 明显季节性 → 注意入场时机

---

## google_trends.get_interest_over_time

获取关键词在 Google 上的搜索热度时间序列。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `q` | string | 是 | 关键词，多个用逗号分隔（最多 5 个） |
| `geo` | string | 否 | 地区代码，ISO 3166-1 alpha-2 或子区域（如 `US`, `US-CA`）。不传则全球 |
| `time` | string | 否 | 时间范围，默认 `today 12-m`。见下表 |
| `cat` | string | 否 | Google Trends 类目 ID，默认 `"0"`（全部） |
| `gprop` | string | 否 | Google 属性：`""` (web), `"images"`, `"news"`, `"froogle"` (shopping), `"youtube"`；默认 `""` (web) |
| `tz` | integer | 否 | 时区偏移（分钟）；不传则使用 UTC |

### 时间范围

| 值 | 含义 |
|----|------|
| `now 1-H` | 过去 1 小时 |
| `now 4-H` | 过去 4 小时 |
| `now 1-d` | 过去 1 天 |
| `now 7-d` | 过去 7 天 |
| `today 1-m` | 过去 30 天 |
| `today 3-m` | 过去 90 天 |
| `today 12-m` | 过去 12 个月（默认） |
| `today 5-y` | 过去 5 年 |
| `all` | 2004 至今 |
| `YYYY-MM-DD YYYY-MM-DD` | 自定义日期范围 |

### CLI 用法

```bash
# 查看关键词过去 12 个月的搜索趋势
apimux google_trends get_interest_over_time --q "wireless earbuds" --geo "US"

# 对比两个关键词
apimux google_trends get_interest_over_time --q "wireless earbuds,bluetooth headphones" --geo "US"

# 查看过去 5 年的长期趋势
apimux google_trends get_interest_over_time --q "wireless earbuds" --geo "US" --time "today 5-y"

# 只看 YouTube 上的搜索热度
apimux google_trends get_interest_over_time --q "wireless earbuds" --gprop "youtube"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `timeline_data[].date` | string | 时间序列日期标签 |
| `timeline_data[].timestamp` | string | Unix 时间戳 |
| `timeline_data[].values[].query` | string | 时间序列中的关键词 |
| `timeline_data[].values[].value` | integer | 时间序列中的归一化热度 0-100 |
| `averages[].query` | string | 平均值对应关键词 |
| `averages[].value` | integer | 时间段内平均热度 |
| `regions[].geo` | string | 地区 geo 编码 |
| `regions[].name` | string | 地区展示名称 |
| `regions[].values[].query` | string | 地区热度中的关键词 |
| `regions[].values[].value` | integer | 地区热度值 |
| `related_topics.top[].position` | integer | Top related topic 排名 |
| `related_topics.top[].topic_id` | string | topic id |
| `related_topics.top[].title` | string | Top related topic 标题 |
| `related_topics.top[].topic_type` | string | Top related topic 类型 |
| `related_topics.top[].value` | integer | Top related topic 热度值 |
| `related_topics.top[].link` | string | Top related topic 链接 |
| `related_topics.rising[].position` | integer | Rising related topic 排名 |
| `related_topics.rising[].topic_id` | string | topic id |
| `related_topics.rising[].title` | string | Rising related topic 标题 |
| `related_topics.rising[].topic_type` | string | Rising related topic 类型 |
| `related_topics.rising[].value` | integer | Rising related topic 热度值 |
| `related_topics.rising[].link` | string | Rising related topic 链接 |
| `related_queries.top[].position` | integer | Top related query 排名 |
| `related_queries.top[].query` | string | Top related query 文本 |
| `related_queries.top[].value` | integer | Top related query 热度值 |
| `related_queries.top[].link` | string | Top related query 链接 |
| `related_queries.rising[].position` | integer | Rising related query 排名 |
| `related_queries.rising[].query` | string | Rising related query 文本 |
| `related_queries.rising[].value` | integer | Rising related query 热度值 |
| `related_queries.rising[].link` | string | Rising related query 链接 |
| `search_metadata.status` | string | upstream request 状态 |
| `search_metadata.created_at` | string | 请求创建时间 |
| `search_metadata.request_url` | string | upstream request URL |
| `search_parameters.q` | string | 原始查询词 |
| `search_parameters.geo` | string | 原始 geo 参数 |
| `search_parameters.time` | string | 原始时间范围参数 |
| `search_parameters.tz` | integer | 原始时区偏移 |
| `search_parameters.data_type` | string | upstream data_type |
| `search_parameters.cat` | string | 原始 category 参数 |
| `search_parameters.region` | string | 原始 region 参数 |
| `search_parameters.gprop` | string | 原始 Google property 参数 |

### 常见错误

- 把热度值当绝对搜索量 → 0-100 是归一化指数，100 代表该时间段内的最高热度，不是绝对搜索量
- 跨查询比较热度值 → 不同查询之间的绝对值不可直接比较，只有同一次查询中的多个关键词才可比较
- time 格式用自然语言（如 "last year"）→ 必须使用预定义值或 `YYYY-MM-DD YYYY-MM-DD` 格式
- gprop 传错值（如 `"shopping"`）→ 必须用 `"froogle"` 代表 shopping
- 不传 geo 就下结论 → 不传 geo 是全球数据，分析特定市场建议传对应 geo

### 规则

- **关键词必填**：`q` 不能为空
- **最多 5 个关键词**：用逗号分隔
- **热度值是相对的**：0-100 是归一化指数，不是绝对搜索量
- **geo 不传 = 全球**：分析特定市场建议传对应 geo（如 `US`）
- **time 格式严格**：必须使用预定义值或 `YYYY-MM-DD YYYY-MM-DD` 格式
- **gprop 枚举**：只接受 `""`, `"images"`, `"news"`, `"froogle"`, `"youtube"`

---

## 通用规则

- **service 使用 canonical envelope，CLI 默认 data-only**：详见 [apimux-shared](../apimux-shared/SKILL.md)
- **错误不暴露 provider 内部信息**：所有错误映射到 canonical error taxonomy
