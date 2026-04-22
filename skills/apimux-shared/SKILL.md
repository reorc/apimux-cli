---
name: apimux-shared
version: 1.0.0
description: "APIMux 共享基础：响应结构、错误处理、partial-failure 语义。所有 APIMux skill 的前置依赖，使用任何 APIMux capability 前必须先读取本文件。"
metadata:
  requires:
    bins: ["apimux"]
  cliHelp: "apimux --help"
---

# APIMux 共享规则

本文件定义所有 APIMux capability 共享的契约和规则。使用任何 APIMux skill 前必须先理解这些基础概念。

## 响应结构

APIMux service contract 使用标准 envelope：

```json
{
  "ok": true,
  "data": { ... },
  "meta": {
    "capability": "amazon.get_product",
    "contract_version": "2025-04-01"
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `ok` | boolean | 请求是否成功 |
| `data` | object/array | 成功时的业务数据 |
| `error` | object | 失败时的错误信息（与 `data` 互斥） |
| `meta` | object | 元信息：capability 标识、契约版本等 |

**Service 规则**：
- `ok=true` 时读 `data`，`ok=false` 时读 `error`
- 空结果（如搜索无匹配）返回 `ok=true, data=[]`，不是 error
- `meta` 始终存在

**CLI 默认输出**：
- `apimux <source> <capability>` 默认输出精简的数据体
- 错误时输出 `{"error": ...}`
- 自动包含关键 `meta` 字段（分页、partial-failure 等），格式为 `{"data": ..., "meta": {...}}`
- 无关键 metadata 时，输出仅为 data（向后兼容）
- 使用 `--debug` 查看完整响应结构

## 错误处理（Error Taxonomy）

Service 侧错误响应结构：

```json
{
  "ok": false,
  "error": {
    "type": "provider",
    "code": "product_not_found",
    "message": "No product found for the given ASIN"
  },
  "meta": { ... }
}
```

### 通用错误码与 Agent 应对策略

| 错误码 | 含义 | Agent 应对 |
|--------|------|-----------|
| `validation_error` | 输入参数不合法 | 检查参数格式，修正后重试 |
| `upstream_timeout` | 上游超时 | 等待片刻后重试；持续超时说明上游服务异常 |
| `provider_unavailable` | 上游不可用 | 等待后重试；持续不可用需报告 |
| `provider_invalid_request` | 上游拒绝请求 | 检查参数组合是否合理 |

各 source 还有自己的业务错误码（如 `product_not_found`、`category_not_found`），详见对应 skill 文档。

**规则**：
- 所有错误码都是标准化的
- `validation_error` 在请求到达上游前就被拦截
- 遇到未知错误码时，向用户报告完整 error 对象

## Partial-Failure 语义

部分能力支持多维度聚合查询，可能出现部分成功：

```json
{
  "ok": true,
  "data": [ ... ],
  "meta": {
    "partial": true,
    "subrequest_count": 3,
    "subrequests": [
      {"name": "sales_volume", "ok": true},
      {"name": "brand_count", "ok": true},
      {"name": "avg_price", "ok": false, "error": {"code": "upstream_timeout"}}
    ]
  }
}
```

**CRITICAL — 当 `meta.partial=true` 时：**
- `data` 中包含成功维度的数据，失败维度的值为 `null`
- 必须检查 `meta.subrequests` 了解哪些维度失败
- 不要把 `null` 值当作"该维度数据为零"
- 向用户说明哪些数据是完整的、哪些缺失

## 关键 Metadata 自动暴露

CLI 默认输出会自动包含以下关键 metadata：

**分页 metadata（当存在时）：**
- `cursor` — 下一页的游标
- `has_more` — 是否还有更多数据
- `current_page` — 当前页码
- `next_page` — 下一页页码
- `total` — 总记录数

**Partial-failure metadata（当 `partial=true` 时）：**
- `partial` — 标识部分成功
- `subrequest_count` — 子请求总数
- `subrequests` — 各子请求的状态详情

**输出格式：**
```json
{
  "data": { ... },
  "meta": {
    "cursor": "t3_xyz789",
    "has_more": true
  }
}
```

当无关键 metadata 时，输出仍为纯 data（向后兼容）。

## CLI Flag 命名约定

APIMux CLI 的参数名到 flag 的映射规则：

- **单词参数**：直接加 `--` 前缀，如 `asin` → `--asin`，`q` → `--q`
- **多词参数（snake_case）**：转为 kebab-case，如 `node_id` → `--node-id`，`begin_date` → `--begin-date`，`page_size` → `--page-size`
- **数组参数**：传逗号分隔的字符串，如 `--asins "B0CM5JV26D,B0D1234567"`，`--keywords "yoga mat,pilates ring"`
- **boolean 参数**：传 `true` 或 `false`，如 `--only-purchase true`

**注意**：不同 skill 的参数名可能不同（如 `q` vs `keyword` vs `query`），这是 capability 语义差异，不是命名错误。跨 skill 调用时，请以各 skill 文档的参数表为准，不要从其他 skill 迁移参数名习惯。

## CLI 通用用法

```bash
# 默认输出
apimux amazon get_product --asin "B0CM5JV26D" --market "US"

# 人类可读格式
apimux --output pretty amazon get_product --asin "B0CM5JV26D" --market "US"

# 原始数据（跳过精简处理）
apimux --output data amazon get_product --asin "B0CM5JV26D" --market "US"

# 完整响应结构（排查问题时使用）
apimux --debug amazon get_product --asin "B0CM5JV26D" --market "US"

# 列出所有能力名称
apimux schema capabilities

# 查看能力参数
apimux schema show amazon.get_product
```

**CLI 规则**：
- 默认输出精简数据体，适合程序消费
- `--output pretty` 输出格式化的精简数据
- `--output data` / `data-pretty` 跳过精简处理，返回原始数据
- `--debug` 输出完整响应结构
- 所有命令格式：`apimux <source> <capability> [flags]`
