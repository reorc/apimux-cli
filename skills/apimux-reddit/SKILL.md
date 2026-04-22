---
name: apimux-reddit
version: 1.0.0
description: "Reddit 内容查询。提供搜索、subreddit feed、帖子详情和评论能力，适用于话题研究、社区巡检、帖子下钻分析等场景。"
metadata:
  source: reddit
  requires:
    bins: ["apimux"]
  cliHelp: "apimux reddit --help"
---

# Reddit

**CRITICAL — 开始前 MUST 先用 Read 工具读取 [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md)，其中包含响应结构、错误处理等共享规则。**

Reddit 数据查询，覆盖帖子搜索、subreddit feed、帖子详情和评论列表四个首批能力。

## 快速决策

- 想按关键词找帖子 → `search`
- 想看某个 subreddit 最近内容 → `get_subreddit_feed`
- 已有 `post_id`，想拉帖子详情 → `get_post_detail`
- 已有 `post_id`，想拉评论 → `get_post_comments`

## Capabilities 概览

| Capability | 说明 | 典型场景 |
|------------|------|----------|
| `search` | 搜索 Reddit 内容 | 话题搜索、帖子发现 |
| `get_subreddit_feed` | 拉取 subreddit feed | 社区巡检、板块追踪 |
| `get_post_detail` | 获取帖子详情 | 单帖下钻分析 |
| `get_post_comments` | 获取帖子评论 | 评论抽样、观点分析 |

## Agent Journey

```text
search → get_post_detail → get_post_comments
```

先搜索帖子，再对目标 `post_id` 拉详情和评论。

---

## reddit.search

搜索 Reddit 内容。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `query` | string | 是 | 搜索关键词 |
| `search_type` | string | 否 | `post`、`community`、`comment`、`media`、`people` |
| `sort` | string | 否 | `relevance`、`hot`、`top`、`new`、`comments` |
| `time_range` | string | 否 | `all`、`year`、`month`、`week`、`day`、`hour` |
| `after` | string | 否 | 分页 cursor |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `post_id` | string | 帖子 ID，canonical `t3_xxxxx` |
| `title` | string | 标题 |
| `subreddit` | string | subreddit 名称 |
| `author` | string | 作者 |
| `score` | integer | 分数 |
| `num_comments` | integer | 评论数 |
| `created_time` | string | RFC3339 发布时间 |
| `permalink` | string | Reddit permalink |
| `url` | string | 目标链接 |
| `selftext` | string | 正文 |
| `thumbnail` | string | 缩略图 |
| `is_video` | boolean | 是否视频帖 |

### 规则

- `query` 必填
- 所有枚举参数只接受 lowercase canonical 值
- 分页状态放在 `meta.cursor` / `meta.has_more`

---

## reddit.get_subreddit_feed

拉取单个 subreddit feed。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `subreddit_name` | string | 是 | 不带 `r/` 前缀的 subreddit 名称 |
| `sort` | string | 否 | `best`、`hot`、`new`、`top`、`controversial`、`rising` |
| `after` | string | 否 | 分页 cursor |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `post_id` | string | 帖子 ID，canonical `t3_xxxxx` |
| `title` | string | 标题 |
| `subreddit` | string | subreddit 名称 |
| `author` | string | 作者 |
| `score` | integer | 分数 |
| `num_comments` | integer | 评论数 |
| `created_time` | string | RFC3339 发布时间 |
| `permalink` | string | Reddit permalink |
| `url` | string | 目标链接 |
| `selftext` | string | 正文 |
| `thumbnail` | string | 缩略图 |
| `is_video` | boolean | 是否视频帖 |

### 规则

- `subreddit_name` 必填，且不能带 `r/` 前缀
- 分页状态放在 `meta.cursor` / `meta.has_more`

---

## reddit.get_post_detail

获取单条 Reddit 帖子详情。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `post_id` | string | 是 | `t3_xxxxx` 形式的帖子 ID |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `post_id` | string | 帖子 ID |
| `title` | string | 标题 |
| `subreddit` | string | subreddit 名称 |
| `author` | string | 作者 |
| `score` | integer | 分数 |
| `upvote_ratio` | number | 点赞率 |
| `num_comments` | integer | 评论数 |
| `created_time` | string | RFC3339 发布时间 |
| `permalink` | string | Reddit permalink |
| `url` | string | 目标链接 |
| `selftext` | string | 正文 |
| `thumbnail` | string | 缩略图 |
| `is_video` | boolean | 是否视频帖 |
| `link_flair_text` | string | flair 文本 |

### 规则

- `post_id` 必须是 `t3_xxxxx` 格式
- 帖子不存在时返回 canonical `post_not_found`

---

## reddit.get_post_comments

获取 Reddit 帖子评论列表。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `post_id` | string | 是 | `t3_xxxxx` 形式的帖子 ID |
| `sort_type` | string | 否 | `confidence`、`new`、`top`、`hot`、`controversial`、`old`、`random` |
| `after` | string | 否 | 分页 cursor |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | string | 评论 ID |
| `author` | string | 评论作者 |
| `body` | string | 评论内容 |
| `score` | integer | 评论分数 |
| `created_at` | string | RFC3339 发布时间 |
| `parent_id` | string | 父节点 ID |
| `depth` | integer | 评论层级 |

### 规则

- `post_id` 必须是 `t3_xxxxx` 格式
- CLI 默认 compact 输出使用 `items.{columns,rows}` compat shape
- 分页状态仍放在 `meta.cursor` / `meta.has_more`

---

## 通用规则

- **service 使用 canonical envelope，CLI 默认 data-only**：详见 [apimux-shared](../apimux-shared/SKILL.md)
- **不暴露 provider 内部信息**：不会暴露真实上游系统名称
