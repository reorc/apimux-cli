---
name: apimux-xiaohongshu
version: 1.0.0
description: "小红书内容查询。提供笔记搜索、笔记详情和评论能力，适用于种草内容研究、笔记巡检、评论抽样等场景。"
metadata:
  source: xiaohongshu
  requires:
    bins: ["apimux"]
  cliHelp: "apimux xiaohongshu --help"
---

# Xiaohongshu

**CRITICAL — 开始前 MUST 先用 Read 工具读取 [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md)，其中包含响应结构、错误处理等共享规则。**

Xiaohongshu 数据查询，覆盖笔记搜索、详情和评论三个首批能力。

## 快速决策

- 想先找笔记样本 → `search_notes`
- 已有 `note_id`，想拉笔记详情 → `get_note_detail`
- 已有 `note_id`，想看评论 → `get_note_comments`

## Capabilities 概览

| Capability | 说明 | 典型场景 |
|------------|------|----------|
| `search_notes` | 搜索笔记 | 种草内容发现、关键词巡检 |
| `get_note_detail` | 获取单篇笔记详情 | 笔记内容与作者信息下钻 |
| `get_note_comments` | 获取笔记评论 | 评论抽样、互动分析 |

## Agent Journey

```text
search_notes → get_note_detail → get_note_comments
```

先搜索笔记，再对目标 `note_id` 拉详情和评论。

---

## xiaohongshu.search_notes

搜索 Xiaohongshu 笔记。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `keyword` | string | 是 | 搜索关键词 |
| `page` | integer | 否 | 页码，从 1 开始；默认 1 |
| `note_type` | string | 否 | `all`、`video`、`normal`、`live`；默认 `all` |
| `time_filter` | string | 否 | `all`、`1d`、`1w`、`6m`；默认 `all` |
| `sort_strategy` | string | 否 | `default`、`latest`、`likes`；默认 `default` |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `note_id` | string | 笔记 ID |
| `title` | string | 标题 |
| `description` | string | 描述 |
| `type` | string | 笔记类型 |
| `xsec_token` | string | 从搜索结果带回的安全 token，后续 `get_note_detail` 可能需要 |
| `like_count` | integer | 点赞数 |
| `collect_count` | integer | 收藏数 |
| `comment_count` | integer | 评论数 |
| `author` | object | 作者信息 |

### 规则

- `keyword` 必填
- canonical 枚举值仍是首选；CLI 也兼容 legacy provider 值：`sort_strategy` 的 `general/time_descending/popularity_descending`，以及 `note_type` 的 `0/1/2`
- 分页走 `page` 参数；返回状态在 `meta.current_page` / `meta.next_page` / `meta.has_more`

---

## xiaohongshu.get_note_detail

获取单篇 Xiaohongshu 笔记详情。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `note_id` | string | 是 | 24-char hex 笔记 ID |
| `xsec_token` | string | 否 | 从 `search_notes` 结果拿到的安全 token；部分笔记详情需要 |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `note_id` | string | 笔记 ID |
| `title` | string | 标题 |
| `description` | string | 描述 |
| `type` | string | 笔记类型 |
| `user_id` | string | 作者 ID |
| `nickname` | string | 作者昵称 |
| `avatar` | string | 作者头像 |
| `like_count` | integer | 点赞数 |
| `collect_count` | integer | 收藏数 |
| `comment_count` | integer | 评论数 |
| `share_count` | integer | 分享数 |
| `images` | string[] | 图片列表 |
| `tags` | string[] | 标签列表 |
| `time` | string | 发布时间 |
| `last_update_time` | string | 更新时间 |
| `ip_location` | string | IP 属地 |
| `video_url` | string | 视频链接（视频笔记时） |

### 规则

- `note_id` 必须是 24-char hex
- 如果 `search_notes` 返回了 `xsec_token`，调用详情时应一并传入
- share link 在 contract 层拒绝，不接受 `xhslink.com/...`
- 笔记不存在时返回 canonical `note_not_found`

---

## xiaohongshu.get_note_comments

获取 Xiaohongshu 笔记评论。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `note_id` | string | 是 | 24-char hex 笔记 ID |
| `cursor` | string | 否 | 评论分页 cursor，首页不传 |
| `sort_strategy` | string | 否 | `default`、`latest`、`likes`；默认 `default` |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `comment_id` | string | 评论 ID |
| `user_id` | string | 作者 ID |
| `nickname` | string | 作者昵称 |
| `avatar` | string | 作者头像 |
| `content` | string | 评论内容 |
| `like_count` | integer | 点赞数 |
| `reply_count` | integer | 回复数 |
| `create_time` | string | 发布时间 |
| `ip_location` | string | IP 属地 |

### 规则

- `note_id` 必须是 24-char hex
- 分页状态放在 `meta.cursor` / `meta.has_more`

---

## 通用规则

- **service 使用 canonical envelope，CLI 默认 data-only**：详见 [apimux-shared](../apimux-shared/SKILL.md)
- **不暴露 provider 内部信息**：不会暴露真实上游系统名称
