---
name: apimux-amazon
version: 1.0.0
description: "Amazon 电商数据查询。提供商品详情、商品搜索、评论分析、类目畅销榜、类目趋势等能力。适用于选品调研、竞品分析、类目分析、市场格局评估等场景。"
metadata:
  source: amazon
  requires:
    bins: ["apimux"]
  cliHelp: "apimux amazon --help"
---

# Amazon

**CRITICAL — 开始前 MUST 先用 Read 工具读取 [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md)，其中包含响应结构、错误处理、partial-failure 语义等共享规则。**

Amazon 电商数据查询，覆盖商品、搜索、评论、类目四大维度。

## 快速决策

- 想从一个种子词扩展相关关键词 → `expand_keywords`
- 想看某个关键词的搜索量、竞争和广告格局 → `get_keyword_overview`
- 想看多个关键词的历史搜索趋势 → `get_keyword_trends`
- 想看某个 ASIN 在搜索里的关联关键词 → `list_asin_keywords`
- 想查 Amazon ABA 热搜词、类目热词和转化线索 → `query_aba_keywords`
- 不知道 node_id，想先按类目名称找 Amazon 类目 → `search_category`
- 想看某个 ASIN 的日销趋势 → `get_asin_sales_daily_trend`
- 想批量对比多个 ASIN 的月销量历史 → `get_asins_sales_history`
- 想看一个 parent/variant ASIN 的近 30 天变体销量 → `get_variant_sales_30d`
- 已知 ASIN，想了解商品详情 → `get_product`
- 不知道 ASIN，想按关键词找商品 → `search_products`
- 想了解某商品的用户评价 → `get_product_reviews`
- 想了解某类目的头部商品 → `get_category_best_sellers`（需要 node_id）
- 想了解某类目的市场趋势 → `get_category_trend`（需要 node_id）
- 想验证一个市场是否有需求 → 先用 [`google_trends.get_interest_over_time`](../apimux-google-trends/SKILL.md) 确认搜索热度，再回来用 `search_products` 评估供给侧

## Capabilities 概览

| Capability | 说明 | 典型场景 |
|------------|------|----------|
| `expand_keywords` | 关键词扩展 | 扩词、词池构建 |
| `get_keyword_overview` | 关键词概览 | 搜索量、竞争度、广告格局 |
| `get_keyword_trends` | 关键词趋势 | 季节性、趋势对比 |
| `list_asin_keywords` | ASIN 关联关键词 | 反查竞品词池 |
| `query_aba_keywords` | ABA 热搜词查询 | 热词榜、类目词发现、转化线索 |
| `search_category` | 类目搜索 | 先解析类目 node_id，再进类目分析 |
| `get_asin_sales_daily_trend` | ASIN 日销趋势 | 日度波动、短期监控 |
| `get_asins_sales_history` | 多 ASIN 月销历史 | 批量对比、走势回看 |
| `get_variant_sales_30d` | 变体 30 天销量 | 变体结构分析 |
| `get_product` | 商品详情 | 竞品基线、选品验证 |
| `search_products` | 关键词搜索 | 选品调研、市场密度 |
| `get_product_reviews` | 评论查询 | 差评分析、口碑研究 |
| `get_category_best_sellers` | 类目畅销榜 | 类目分析、选品 |
| `get_category_trend` | 类目趋势 | 市场动态、趋势判断 |

## Agent Journeys

以下是 Amazon capabilities 在典型 Agent 工作流中的编排方式：

### Journey 1: 类目分析
```
(resolve node_id) → get_category_best_sellers → get_category_trend → get_product (fan-out) → get_product_reviews
```
从类目畅销榜出发，了解趋势，再对 top ASIN 做详情和评论分析。

### Journey 2: 竞品 ASIN 深度分析
```
get_product → get_product_reviews (negative) → get_product_reviews (positive)
```
建立商品基线，先拉差评分析痛点，再拉好评了解优势。

### Journey 3: 商品发现
```
search_products → get_product (fan-out) → get_product_reviews
```
关键词搜索 → 详情 → 评论抽样。

### Journey 4: 市场验证（跨 skill）
```
google_trends.get_interest_over_time → search_products → (optional) get_category_trend
```
先用 [Google Trends](../apimux-google-trends/SKILL.md) 确认需求信号，再用搜索评估 Amazon 供给侧密度。

### Journey 5: 评论驱动产品改进
```
get_product → get_product_reviews (negative, start_date=30d ago, only_purchase=true)
```
拉取近 30 天验证购买差评，提取改进方向。

### Journey 6: 关键词机会验证
```
expand_keywords → get_keyword_overview (fan-out) → get_keyword_trends
```
先扩词，再筛关键词质量，最后看趋势确认是否具备持续需求。

### Journey 7: 竞品销量与关键词联动分析
```
get_asins_sales_history → list_asin_keywords → get_keyword_overview
```
先确定哪些 ASIN 在卖，再反查它们的核心关键词，最后对关键词质量做二次判断。

### Journey 8: ABA 热词机会挖掘
```
query_aba_keywords → get_keyword_overview → get_keyword_trends
```
先看 ABA 热搜词和类目线索，再对候选词做质量与趋势验证。

---

## amazon.expand_keywords

扩展一个 Amazon 种子词，返回相关关键词及搜索量信号。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `keyword` | string | 是 | 种子关键词 |
| `market` | string | 否 | 目标市场：`US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT`；默认 `US` |

### CLI 用法

```bash
apimux amazon expand_keywords --keyword "yoga mat"
apimux amazon expand_keywords --keyword "yoga mat" --market "DE"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `keyword` | string | 扩展后的关键词 |
| `match_types` | string[] | provider 返回的匹配类型标签 |
| `est_searches_num` | integer | 估算月搜索量 |
| `searches_rank` | integer | 热度排名 |

`meta.total` 表示总扩展词数。

### 常见错误

- 传空 keyword → 返回 `missing_keyword`
- market 不在支持列表里 → 返回 `invalid_market`

### 规则

- `keyword` 必填，空字符串会被 facade 拒绝
- `market` 可省略，默认按 `US` 查询
- 返回的是 object array，不暴露 provider 原始 columnar shape

---

## amazon.get_keyword_overview

获取单个 Amazon 关键词的综合指标。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `keyword` | string | 是 | 要分析的关键词 |
| `market` | string | 否 | 目标市场：`US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT`；默认 `US` |

### CLI 用法

```bash
apimux amazon get_keyword_overview --keyword "yoga mat"
apimux amazon get_keyword_overview --keyword "standing desk" --market "UK"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `keyword` | string | 查询关键词 |
| `est_searches_num` | integer | 估算月搜索量 |
| `searches_rank` | integer | 热度排名 |
| `searches` | integer/null | 原始搜索次数（provider 提供时） |
| `demand_ratio` | number/null | 需求比率（provider 提供时） |
| `competitor_cnt` | integer/null | 竞争商品数量 |
| `sale_num` | integer/null | 销量相关指标 |
| `ac_asin_num` | integer/null | Amazon's Choice ASIN 数量 |
| `brand_ad_asin_num` | integer/null | 品牌广告 ASIN 数量 |
| `sp_ad_asin_num` | integer/null | Sponsored Product ASIN 数量 |
| `ppc_ad_asin_num` | integer/null | PPC 广告 ASIN 数量 |
| `video_ad_asin_num` | integer/null | 视频广告 ASIN 数量 |
| `er_asin_num` | integer/null | 编辑推荐 ASIN 数量 |
| `nf_asin_num` | integer/null | New Flag ASIN 数量 |
| `tr_asin_num` | integer/null | Top Rated ASIN 数量 |
| `search_recommend_asin_num` | integer/null | 搜索推荐 ASIN 数量 |
| `global_keyword_num` | integer/null | 全局关键词数量 |
| `update_time` | string/null | provider 数据更新时间 |

### 常见错误

- 传空 keyword → 返回 `missing_keyword`
- market 不在支持列表里 → 返回 `invalid_market`

### 规则

- 这是 facade 聚合能力：内部会合并多个上游接口，对外只暴露一个 capability
- 缺失的 provider 指标会返回 `null`，不会伪造默认值
- 响应不暴露 provider 名称或原始字段名

---

## amazon.get_keyword_trends

获取多个 Amazon 关键词的历史搜索趋势。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `keywords` | string[] | 是 | 关键词列表，非空，最多 1000 个 |
| `market` | string | 否 | 目标市场：`US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT`；默认 `US` |
| `granularity` | string | 否 | `week` 或 `month`，默认 `month` |

### CLI 用法

```bash
apimux amazon get_keyword_trends --keywords "yoga mat,pilates ring"
apimux amazon get_keyword_trends --keywords "desk mat,monitor stand" --market "US" --granularity "week"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `keyword` | string | 查询关键词 |
| `est_searches_num_history` | object | 时间 -> 搜索量映射 |
| `searches_rank_history` | object | 时间 -> 排名映射 |

### 常见错误

- `keywords` 为空 → 返回 `missing_keywords`
- keyword 数量超过 1000 → 返回 `too_many_keywords`
- granularity 不是 `week|month` → 返回 `invalid_granularity`

### 规则

- `keywords` 必须是字符串数组；CLI 侧通过逗号分隔转换
- 月粒度 key 形如 `YYYY-MM`，周粒度 key 形如 `YYYY-MM-DD`
- 出站请求 body 会固定传 canonical `keywords` 列表和 `granularity`

---

## amazon.list_asin_keywords

列出某个 Amazon ASIN 关联的关键词。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `asin` | string | 是 | Amazon ASIN，10 位大写字母数字 |
| `keyword` | string | 否 | 关键词过滤器 |
| `market` | string | 否 | 目标市场：`US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT`；默认 `US` |

### CLI 用法

```bash
apimux amazon list_asin_keywords --asin "B0CM5JV26D"
apimux amazon list_asin_keywords --asin "B0CM5JV26D" --keyword "desk" --market "US"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `keyword` | string | 关联关键词 |
| `kw_characters` | string[] | 关键词标签 |
| `conversion_characters` | string[] | 转化标签 |
| `exposure_type` | string[] | 曝光类型 |
| `last_rank` | string | 自然排名 |
| `ad_last_rank` | string | 广告排名 |
| `est_searches_num` | integer | 估算搜索量 |
| `searches_rank` | integer | 热度排名 |
| `ratio_score` | number | 流量占比得分 |

`meta.total` 表示总关键词数。

### 常见错误

- ASIN 格式不对 → 返回 `invalid_asin`
- market 不在支持列表里 → 返回 `invalid_market`

### 规则

- `asin` 必须是 10 位大写字母数字；小写会先标准化再校验
- `keyword` 只是 provider 侧过滤器，不是必填
- 响应已经 canonicalize，不暴露 provider column/row 结构

---

## amazon.query_aba_keywords

查询 Amazon Brand Analytics 热搜词，适合从 ABA 榜单里挖关键词机会。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `keyword` | string | 否 | 关键词过滤器 |
| `node_ids` | string | 否 | 逗号分隔的类目 node ids |
| `page` | integer | 否 | 页码，默认 `1` |
| `page_size` | integer | 否 | 每页条数，范围 `20..200`，默认 `40` |
| `market` | string | 否 | 目标市场：`US`, `UK`, `DE`, `FR`, `IN`, `CA`, `JP`, `ES`, `IT`, `MX`, `AE`, `AU`, `BR`, `SA`；默认 `US` |

### CLI 用法

```bash
apimux amazon query_aba_keywords --keyword "yoga mat"
apimux amazon query_aba_keywords --node-ids "12345,67890" --page 2 --page-size 40 --market "US"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `keyword` | string | 热搜关键词 |
| `keyword_cn_name` | string | 中文翻译 |
| `rank` | integer | 热度排名 |
| `search_volume` | integer | 月搜索量 |
| `word_count` | integer | 词数 |
| `product_count` | integer | 商品数 |
| `rank_change_of_weekly` | number | 周度排名变化 |
| `cpc` | integer | 点击成本 |
| `search_conversion_rate` | number | 搜索转化率 |
| `click_of_90d` | integer | 90 天点击量 |
| `sales_volume_of_90d` | integer | 90 天销量 |
| `top3_asin` | string[] | Top 3 ASIN |
| `top3_brand` | string[] | Top 3 品牌 |
| `top3_category` | string[] | Top 3 类目 |
| `season` | string | 季节性标签 |

`meta.current_page` / `meta.has_more` 表示分页状态。

### 常见错误

- `page_size` 不在 `20..200` 范围内 → 返回 `invalid_page_size`
- market 不在支持列表里 → 返回 `invalid_market`

### 规则

- `keyword` 和 `node_ids` 都是可选；都不传时返回平台热词
- facade 会把 provider 的 `columns + rows` 结构扁平化成 object array
- 响应保持 provider-isolated，不暴露真实上游字段命名和 provider 身份

---

## amazon.get_asin_sales_daily_trend

获取单个 Amazon ASIN 的日销趋势。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `asin` | string | 是 | Amazon ASIN，10 位大写字母数字 |
| `begin_date` | string | 否 | 起始日期，`YYYY-MM-DD` |
| `market` | string | 否 | 目标市场：`US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT`；默认 `US` |

### CLI 用法

```bash
apimux amazon get_asin_sales_daily_trend --asin "B0CM5JV26D"
apimux amazon get_asin_sales_daily_trend --asin "B0CM5JV26D" --begin-date "2026-04-01"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `date` | string | 日期，`YYYY-MM-DD` |
| `sales` | integer | 当天销量 |

### 常见错误

- ASIN 格式不对 → 返回 `invalid_asin`
- `begin_date` 不是 ISO 日期 → 返回 `invalid_begin_date`

### 规则

- `begin_date` 可省略；传了就必须是 `YYYY-MM-DD`
- 返回按日期升序排列

---

## amazon.get_asins_sales_history

批量获取多个 Amazon ASIN 的月销量历史。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `asins` | string[] | 是 | ASIN 列表，非空，最多 10 个 |
| `market` | string | 否 | 目标市场：`US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT`；默认 `US` |

### CLI 用法

```bash
apimux amazon get_asins_sales_history --asins "B0CM5JV26D,B0D1234567"
apimux amazon get_asins_sales_history --asins "B0CM5JV26D,B0D1234567" --market "UK"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `asin` | string | Amazon ASIN |
| `month` | string | 月份，`YYYY-MM` |
| `sales` | integer | 月销量 |

### 常见错误

- `asins` 为空 → 返回 `missing_asins`
- ASIN 数量超过 10 → 返回 `too_many_asins`
- 列表中任一 ASIN 格式不对 → 返回 `invalid_asin`

### 规则

- 这是 batch query，单次最多 10 个 ASIN
- 某些 ASIN 没有历史时，不会导致整个请求失败；结果只返回有数据的行
- 响应会按 `asin + month` 的稳定顺序返回

---

## amazon.get_variant_sales_30d

获取一个 Amazon ASIN 关联变体的近 30 天销量。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `asin` | string | 是 | parent 或 variant ASIN，10 位大写字母数字 |
| `market` | string | 否 | 目标市场：`US`, `UK`, `DE`, `JP`, `CA`, `FR`, `ES`, `IT`；默认 `US` |

### CLI 用法

```bash
apimux amazon get_variant_sales_30d --asin "B0CM5JV26D"
apimux amazon get_variant_sales_30d --asin "B0CM5JV26D" --market "DE"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `asin` | string | 变体 ASIN |
| `bought_in_past_month` | integer | 近 30 天销量 |
| `update_time` | string | provider 更新时间 |

### 常见错误

- ASIN 格式不对 → 返回 `invalid_asin`
- market 不在支持列表里 → 返回 `invalid_market`

### 规则

- 输入可以是 parent ASIN 或某个具体 variant ASIN
- 响应不保证只有一条；通常会返回整个变体组

---

## amazon.get_product

获取单个商品的详情信息。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `asin` | string | 是 | Amazon 商品 ASIN，10 位大写字母数字 |
| `market` | string | 是 | 目标市场：`US`, `UK`, `DE`, `FR`, `IN`, `CA`, `JP`, `ES`, `IT`, `MX`, `AE`, `AU`, `BR`, `SA` |

### CLI 用法

```bash
apimux amazon get_product --asin "B0CM5JV26D" --market "US"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `asin` | string | 商品 ASIN |
| `market` | string | Amazon 市场代码 |
| `title` | string | 商品标题 |
| `brand` | string | 品牌名 |
| `product_url` | string | 商品链接 |
| `main_image` | string | 主图 URL |
| `rating` | float | 平均评分 0-5 |
| `review_count` | integer | 评分数量 |
| `price.display` | string | 当前价格展示文本 |
| `price.value` | number | 当前价格数值 |
| `price.currency` | string | 当前价格币种 |
| `brand_store.id` | string | 品牌店铺 ID |
| `brand_store.text` | string | 品牌店铺展示文案 |
| `brand_store.link` | string | 品牌店铺链接 |
| `feature_bullets` | array<string> | 商品卖点列表 |
| `images[].link` | string | 商品图片链接 |
| `images[].variant` | string | 图片变体标识 |
| `variants[].asin` | string | 变体 ASIN |
| `variants[].title` | string | 变体标题 |
| `variants[].link` | string | 变体链接 |
| `variants[].is_current_product` | boolean | 是否当前变体 |
| `variants[].main_image` | string | 变体主图 |
| `variants[].dimensions[].name` | string | 变体维度名 |
| `variants[].dimensions[].value` | string | 变体维度值 |
| `buybox.price.display` | string | Buy Box 当前价格展示文本 |
| `buybox.price.value` | number | Buy Box 当前价格数值 |
| `buybox.price.currency` | string | Buy Box 当前价格币种 |
| `buybox.original_price.display` | string | Buy Box 原价展示文本 |
| `buybox.original_price.value` | number | Buy Box 原价数值 |
| `buybox.original_price.currency` | string | Buy Box 原价币种 |
| `buybox.availability` | string | Buy Box 库存/可售状态 |
| `buybox.is_prime` | boolean | Buy Box 是否 Prime |

### 常见错误

- 传了小写 ASIN（如 `b0cm5jv26d`）→ 会自动转大写，不会报错
- ASIN 不是 10 位 → 返回 `validation_error`
- 不知道 ASIN → 先用 `search_products` 搜索获取

### 规则

- ASIN 必须是 10 位大写字母数字（小写自动转大写，非 10 位被拒绝）
- market 必填，不会默认到 US
- 如果没有 ASIN，应先通过 `search_products` 搜索获取
- 对多个 ASIN 查询详情时可以并行调用

---

## amazon.search_products

通过关键词搜索 Amazon 商品列表。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `q` | string | 是 | 搜索关键词 |
| `market` | string | 是 | 目标市场：`US`, `UK`, `DE`, `FR`, `IN`, `CA`, `JP`, `ES`, `IT` |
| `page` | integer | 否 | 页码，默认 1，上限 10 |

### CLI 用法

```bash
apimux amazon search_products --q "wireless earbuds" --market "US"
apimux amazon search_products --q "wireless earbuds" --market "US" --page 2
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `asin` | string | 商品 ASIN |
| `market` | string | 市场 |
| `title` | string | 商品标题 |
| `brand` | string | 品牌名 |
| `product_url` | string | 商品链接 |
| `main_image` | string | 主图 URL |
| `rating` | float | 平均评分 |
| `review_count` | integer | 评分数量 |
| `price` | object | 价格信息 |
| `position` | integer | 搜索结果排名位置 |

`meta.current_page` 标识当前页码。

### 常见错误

- 空关键词 → 返回 `validation_error`
- page 超过 10 → 被 facade 拒绝
- 把搜索结果当完整信息用 → 搜索结果是摘要，完整信息需对 ASIN 调用 `get_product`

### 规则

- 关键词必填，空关键词返回 validation error
- page 上限 10，超过被 facade 拒绝
- 搜索结果是摘要，完整信息需对 ASIN 调用 `get_product`

---

## amazon.get_product_reviews

获取商品评论列表，支持多维度过滤。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `asin` | string | 是 | Amazon 商品 ASIN，10 位大写字母数字 |
| `market` | string | 是 | 目标市场（不含 JP）：`US`, `UK`, `DE`, `FR`, `IN`, `CA`, `ES`, `IT`, `MX`, `AE`, `AU`, `BR`, `SA` |
| `start_date` | string | 否 | 评论起始日期，ISO 格式 `YYYY-MM-DD` |
| `star` | string | 否 | `positive`（4-5 星）或 `negative`（1-3 星）；不传则返回全部评论 |
| `only_purchase` | boolean | 否 | `true` 只返回验证购买评论；默认 `false` |
| `page_index` | integer | 否 | 页码，默认 1，上限 10 |

### CLI 用法

```bash
# 获取所有评论
apimux amazon get_product_reviews --asin "B0CM5JV26D" --market "US" --start-date "2025-01-01"

# 只看差评
apimux amazon get_product_reviews --asin "B0CM5JV26D" --market "US" --star "negative"

# 只看验证购买的好评
apimux amazon get_product_reviews --asin "B0CM5JV26D" --market "US" --star "positive" --only-purchase true
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `star` | float | 评分 1-5 |
| `title` | string | 评论标题 |
| `content` | string | 评论正文 |
| `date` | string | 评论日期 `YYYY-MM-DD` |
| `is_verified_purchase` | boolean | 是否验证购买 |
| `helpful_votes` | integer | 有用投票数 |
| `reviewer_name` | string | 评论者名称 |

空结果返回 `ok=true, data=[]`，不是 error。

### 常见错误

- 传数字星级（如 `"1"`, `"5"`）→ 必须传 `positive` 或 `negative`，不能传数字
- 传字符串 `"true"` 给 only_purchase → 必须传 boolean `true`，不是字符串
- 日期格式错误（如 `01/01/2025`）→ 必须 ISO 格式 `YYYY-MM-DD`
- 传 `JP` 市场 → 返回 validation error，JP 不支持评论查询

### 规则

- **star 只接受 canonical 值**：`positive` 或 `negative`，不能传数字
- **only_purchase 严格 boolean**：必须传 `true` 或 `false`，字符串 `"true"` 被拒绝
- **start_date 必须 ISO 格式**：`YYYY-MM-DD`
- **JP 市场不支持**：传 `JP` 返回 validation error
- **过滤条件是 AND 关系**：`star=negative` + `only_purchase=true` = 只看验证购买的差评

---

## amazon.search_category

按类目名称搜索 Amazon 类目，返回后续类目相关 capability 可复用的 `node_id`。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 类目名称，支持英文或中文，部分匹配 |
| `market` | string | 否 | 目标市场：`US`, `UK`, `DE`, `FR`, `IN`, `CA`, `JP`, `ES`, `IT`, `MX`, `AE`, `AU`, `BR`, `SA`；默认 `US` |
| `limit` | integer | 否 | 返回结果数，范围 `1..100`，默认 `20` |

### CLI 用法

```bash
apimux amazon search_category --name "cell phone"
apimux amazon search_category --name "手机" --market "US" --limit 5
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `node_id` | string | Amazon 类目节点 ID |
| `name` | string | 英文类目名 |
| `cn_name` | string | 中文类目名 |
| `path` | object[] | 完整类目路径，每个元素含 `node_id` + `name` |

`meta.total` 表示匹配总数。服务内部会维护 7 天 TTL 的本地 category tree cache，并在本地索引上检索。

### 常见错误

- 传空 `name` → 返回 `missing_name`
- `limit` 不在 `1..100` 范围内 → 返回 `invalid_limit`
- market 不在支持列表里 → 返回 `invalid_market`

### 规则

- **本地检索，不是逐次打 upstream**：后台按 market 维护 category tree cache，查询走本地索引
- **结果按相关度排序**：优先命中名称和路径更相关的类目
- **返回 canonical path**：后续调用 `get_category_best_sellers` / `get_category_trend` 时直接复用 `node_id`

---

## amazon.get_category_best_sellers

获取类目畅销商品排行榜。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `node_id` | string | 是 | Amazon 类目节点 ID（纯数字，至少 3 位） |
| `market` | string | 是 | 目标市场 |
| `query_start` | string | 否 | 历史数据起始日期 `YYYY-MM-DD` |
| `query_date` | string | 否 | 历史数据截止日期 `YYYY-MM-DD` |
| `query_days` | integer | 否 | 从 query_date 往前回溯天数，上限 365 |

### CLI 用法

```bash
apimux amazon get_category_best_sellers --node-id "3743561" --market "US"
apimux amazon get_category_best_sellers --node-id "3743561" --market "US" --query-date "2026-01-01" --query-days 30
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `asin` | string | 商品 ASIN |
| `market` | string | 市场 |
| `rank` | integer | 畅销排名（显式字段，不依赖列表顺序） |
| `title` | string | 商品标题 |
| `brand` | string | 品牌名 |
| `main_image` | string | 主图 URL |
| `rating` | float | 平均评分 |
| `rating_count` | integer | 评分数量 |
| `price` | object | 价格信息 |
| `category` | string[] | 类目路径 |
| `seller_count` | integer | 卖家数量 |
| `is_fba` | boolean | 是否 FBA |

### 常见错误

- node_id 传了非数字（如 `"electronics"`）→ 必须是纯数字
- node_id 不存在 → 返回 `category_not_found`（HTTP 404），不是 upstream error
- 不知道 node_id → 先调用 `search_category` 搜索类目

### 规则

- **node_id 必须是数字**：至少 3 位纯数字的 Amazon 类目节点 ID
- **不知道 node_id 先搜类目**：调用 `search_category` 再进入类目相关 capability
- **不存在的 node_id 返回 not-found**：`category_not_found`（HTTP 404）
- `store_name` / `product_type` 属于上游偶发补充字段，不应作为稳定 compat 字段依赖
- query_start / query_date 必须 ISO 格式 `YYYY-MM-DD`
- query_days 上限 365

---

## amazon.get_category_trend

获取类目趋势数据，支持同时查询多个趋势维度。

**CRITICAL — 本 capability 使用 fan-out 聚合，可能出现 partial-failure。使用前 MUST 理解 [apimux-shared](../apimux-shared/SKILL.md) 中的 Partial-Failure 语义。**

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `node_id` | string | 是 | Amazon 类目节点 ID（纯数字，至少 3 位） |
| `market` | string | 是 | 目标市场 |
| `trend_types` | string[] | 是 | 趋势维度名称列表（见下表），支持同时查询多个 |

### 趋势维度（trend_types）

| 名称 | 维度 |
|------|------|
| `sales_volume` | 销量 |
| `brand_count` | 品牌数 |
| `seller_count` | 卖家数 |
| `avg_price` | 均价 |
| `avg_rating_count` | 平均评论数 |
| `avg_star` | 平均评分 |
| `new_product_ratio_1m` | 近 1 月新品占比 |
| `new_product_ratio_3m` | 近 3 月新品占比 |
| `amazon_self_ratio` | 亚马逊自营占比 |
| `avg_profit` | 平均利润 |
| `top100_share` | Top 100 集中度 |
| `top3_listing_monopoly` | Top 3 Listing 垄断度 |
| `top10_brand_monopoly` | Top 10 品牌垄断度 |

### CLI 用法

```bash
# 查询单个维度
apimux amazon get_category_trend --node-id "3743561" --market "US" --trend-types "sales_volume"

# 同时查询多个维度（逗号分隔）
apimux amazon get_category_trend --node-id "3743561" --market "US" --trend-types "sales_volume,brand_count,avg_price"
```

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `month` | string | 月份，格式 `YYYYMM` |
| `<requested trend_type>` | number/null | 请求的趋势维度列；失败维度值为 `null` |

### 返回格式

时间序列，每个数据点包含月份和各维度的值：

```json
{
  "ok": true,
  "data": [
    {"month": "202601", "sales_volume": 1234, "brand_count": 56},
    {"month": "202602", "sales_volume": 1456, "brand_count": 61}
  ]
}
```

当请求多个维度时，检查 `meta.partial` 判断是否所有维度都成功返回。失败维度的值为 `null`，详见 `meta.subrequests`。

### 常见错误

- 不检查 `meta.partial` → 部分维度可能失败，`null` 值不代表"零"
- 传了不存在的维度名 → 返回 `invalid_trend_type`，必须使用上表中的名称
- 传了数字（如 `0`）而非字符串名称 → 必须传维度名称字符串
- 把月度数据当日度数据用 → 数据粒度是月（YYYYMM），不支持日/周
- 传了重复的维度名 → 返回 `duplicate_trend_type`

### 规则

- node_id 必须是数字，至少 3 位
- **trend_types 必填**，传维度名称字符串列表，不是数字
- **支持多维度同时查询**：一次请求可传多个维度，facade 内部并行 fan-out
- 不支持重复维度名
- 数据粒度是月（YYYYMM），不支持日/周粒度
- **partial-failure**：部分维度失败时 `meta.partial=true`，失败维度在 data 中为 `null`

---

## 通用规则

所有 Amazon capabilities 共享以下规则：

- **service 使用 canonical envelope，CLI 默认 data-only**：详见 [apimux-shared](../apimux-shared/SKILL.md)
- **market 规则按 capability 区分**：商品/评论/类目能力通常要求显式 market；keyword intelligence 这批可省略并默认 `US`
- **ASIN 格式**：10 位大写字母数字（`^[A-Z0-9]{10}$`），小写自动转大写
- **错误不暴露 provider 内部信息**：所有错误映射到 canonical error taxonomy
