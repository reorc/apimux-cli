---
name: apimux-douyin
version: 1.0.0
description: "抖音视频搜索与评论查询。提供视频搜索、详情和评论能力，适用于内容研究、热视频巡检、评论抽样等场景。"
metadata:
  source: douyin
  requires:
    bins: ["apimux"]
  cliHelp: "apimux douyin --help"
---

# Douyin

**CRITICAL — 开始前 MUST 先用 Read 工具读取 [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md)，其中包含响应结构、错误处理等共享规则。**

Douyin 数据查询，覆盖视频搜索、视频详情、评论列表和评论回复四个能力。

## 快速决策

- 想先找视频样本 → `search_videos`
- 已有 `aweme_id`，想看单条视频详情 → `get_video_detail`
- 已有 `aweme_id`，想拉评论 → `get_video_comments`
- 已有 `aweme_id` + `comment_id`，想继续下钻评论回复 → `get_comment_replies`

## Capabilities 概览

| Capability | 说明 | 典型场景 |
|------------|------|----------|
| `search_videos` | 搜索 Douyin 视频 | 热门内容发现、关键词巡检 |
| `get_video_detail` | 获取单条视频详情 | 视频元数据与作者信息下钻 |
| `get_video_comments` | 获取视频评论列表 | 评论抽样、舆情分析 |
| `get_comment_replies` | 获取评论回复列表 | 评论线程下钻、回复语义分析 |

---

## douyin.search_videos

搜索 Douyin 视频。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `keyword` | string | 是 | 搜索关键词 |
| `sort_type` | string | 否 | `comprehensive`、`likes`、`latest`；默认 `comprehensive` |
| `publish_time` | string | 否 | `all`、`1d`、`1w`、`6m`；默认 `all` |
| `filter_duration` | string | 否 | `all`、`under_1m`、`1m_5m`、`over_5m`；默认 `all` |
| `content_type` | string | 否 | `all`、`video`、`image`、`article`；默认 `all` |
| `cursor` | integer | 否 | 分页 cursor，首页不传 |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `aweme_id` | string | 视频 ID |
| `description` | string | 视频描述 |
| `create_time` | string | RFC3339 发布时间 |
| `author` | object | 作者信息 |
| `statistics` | object | 点赞/评论/分享/播放统计 |
| `video` | object | 视频时长与比例 |
| `share_url` | string | 分享链接 |

### 规则

- `keyword` 必填
- 推荐使用标准枚举值；CLI 为兼容历史调用，也接受旧版别名：`sort_type` 的 `0/1/2` 和 `publish_time` 的 `0/1/7/180`
- 分页状态放在 `meta.cursor` / `meta.has_more`

---

## douyin.get_video_detail

获取单条 Douyin 视频详情。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `aweme_id` | string | 是 | 数字字符串形式的视频 ID |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `aweme_id` | string | 视频 ID |
| `description` | string | 视频描述 |
| `create_time` | string | RFC3339 发布时间 |
| `author` | object | 作者信息 |
| `statistics` | object | 点赞/评论/分享/播放统计 |
| `video` | object | 视频时长与比例 |
| `share_url` | string | 分享链接 |

### 规则

- `aweme_id` 必须是数字字符串
- 视频不存在时返回 `video_not_found` 错误

---

## douyin.get_video_comments

获取 Douyin 视频评论列表。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `aweme_id` | string | 是 | 数字字符串形式的视频 ID |
| `cursor` | integer | 否 | 评论分页 cursor，首页不传 |
| `count` | integer | 否 | 评论页大小；默认 20 |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `comment_id` | string | 评论 ID |
| `text` | string | 评论文本 |
| `like_count` | integer | 点赞数 |
| `reply_count` | integer | 回复数 |
| `create_time` | string | RFC3339 评论时间 |
| `author` | object | 评论作者信息 |

### 规则

- `aweme_id` 必须是数字字符串
- 分页状态放在 `meta.cursor` / `meta.has_more`
- 总评论数放在 `meta.total`

---

## douyin.get_comment_replies

获取 Douyin 某条父评论下的回复列表。

### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `aweme_id` | string | 是 | 数字字符串形式的视频 ID |
| `comment_id` | string | 是 | 父评论 ID |
| `cursor` | integer | 否 | 回复分页 cursor，首页不传 |
| `count` | integer | 否 | 回复页大小；默认 20 |

### 返回字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `comment_id` | string | 回复评论 ID |
| `text` | string | 回复文本 |
| `like_count` | integer | 点赞数 |
| `reply_count` | integer | 嵌套回复数 |
| `create_time` | string | RFC3339 回复时间 |
| `author` | object | 回复作者信息 |

### 规则

- `aweme_id` 必须是数字字符串
- `comment_id` 必填
- 分页状态放在 `meta.cursor` / `meta.has_more`
- 总回复数放在 `meta.total`

---

## 通用规则

- **响应结构与错误处理**：详见 [apimux-shared](../apimux-shared/SKILL.md)
