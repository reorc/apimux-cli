---
name: apimux-shared
version: 1.0.0
description: "Shared APIMux guidance for response structure, error handling, partial-failure semantics, metadata, and CLI conventions. Read this before using any APIMux skill."
metadata:
  requires:
    bins: ["apimux"]
  cliHelp: "apimux --help"
---

# APIMux Shared Guide

Read this before using any APIMux source skill. It explains the response envelope, error shape, metadata behavior, and CLI flag conventions that all APIMux capabilities share.

## Response structure

The APIMux service contract uses a standard envelope:

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

| Field | Type | Description |
|-------|------|-------------|
| `ok` | boolean | Whether the request succeeded |
| `data` | object/array | Business data when `ok=true` |
| `error` | object | Error information when `ok=false`; mutually exclusive with `data` |
| `meta` | object | Metadata such as capability name, contract version, pagination, and partial-failure details |

**Service rules:**
- Read `data` when `ok=true`; read `error` when `ok=false`.
- Empty matches, such as a search with no results, return `ok=true, data=[]`; they are not errors.
- `meta` is always present.

**Default CLI output:**
- `apimux <source> <capability>` prints a compact data-oriented response by default.
- Errors print as `{"error": ...}`.
- Important metadata such as pagination and partial-failure details is preserved as `{"data": ..., "meta": {...}}`.
- If there is no important metadata, the CLI prints the data body directly for backward compatibility.
- Use `--debug` to inspect the full service envelope.

## Error handling

Service-side errors use this shape:

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

### Common error codes and agent actions

| Error code | Meaning | Agent action |
|------------|---------|--------------|
| `validation_error` | Input parameters are invalid | Fix the parameter format and retry |
| `upstream_timeout` | Upstream provider timed out | Retry after a short wait; report if it persists |
| `provider_unavailable` | Upstream provider is unavailable | Retry later; report persistent failures |
| `provider_invalid_request` | Upstream rejected the request | Check whether the parameter combination is valid |

Each source can define additional business errors such as `product_not_found` or `category_not_found`. See the source-specific skill for those cases.

**Rules:**
- Error codes are normalized.
- `validation_error` is raised before the request reaches the upstream provider.
- For unknown error codes, report the full error object to the user.

## Partial-failure semantics

Some capabilities aggregate multiple dimensions in one request. These can partially succeed:

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

**When `meta.partial=true`:**
- `data` contains the successful dimensions, while failed dimensions usually appear as `null`.
- Check `meta.subrequests` to identify which dimensions failed.
- Do not interpret `null` values as zero.
- Tell the user which data is complete and which dimensions are missing.

## Metadata exposed by the CLI

Default CLI output automatically preserves important metadata.

**Pagination metadata, when present:**
- `cursor` — cursor for the next page
- `has_more` — whether more data is available
- `current_page` — current page number
- `next_page` — next page number
- `total` — total record count

**Partial-failure metadata, when present:**
- `partial` — whether the response is partial
- `subrequest_count` — number of subrequests
- `subrequests` — status for each subrequest

**Output shape:**
```json
{
  "data": { ... },
  "meta": {
    "cursor": "t3_xyz789",
    "has_more": true
  }
}
```

If there is no important metadata, output remains the compact data body.

## CLI flag conventions

APIMux maps parameter names to CLI flags consistently:

- Single-word parameters use `--` directly, such as `asin` → `--asin` and `q` → `--q`.
- Snake-case parameters become kebab-case, such as `node_id` → `--node-id`, `begin_date` → `--begin-date`, and `page_size` → `--page-size`.
- Array parameters are passed as comma-separated strings, such as `--asins "B0CM5JV26D,B0D1234567"` or `--keywords "yoga mat,pilates ring"`.
- Boolean parameters use `true` or `false`, such as `--only-purchase true`.

Parameter names can differ between skills, such as `q`, `keyword`, and `query`. Follow each capability's parameter table instead of carrying naming habits across sources.

## CLI usage

```bash
# Default compact output
apimux amazon get_product --asin "B0CM5JV26D" --market "US"

# Human-readable compact output
apimux --output pretty amazon get_product --asin "B0CM5JV26D" --market "US"

# Raw data output without compact projection
apimux --output data amazon get_product --asin "B0CM5JV26D" --market "US"

# Full service envelope for debugging
apimux --debug amazon get_product --asin "B0CM5JV26D" --market "US"

# List capability names
apimux schema capabilities

# Inspect capability parameters
apimux schema show amazon.get_product
```

**CLI rules:**
- Default output is compact and suitable for programmatic use.
- `--output pretty` prints formatted compact data.
- `--output data` and `--output data-pretty` bypass compact projection and return raw data.
- `--debug` prints the full service envelope.
- Command format is always `apimux <source> <capability> [flags]`.
