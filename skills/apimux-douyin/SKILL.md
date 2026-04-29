---
name: apimux-douyin
version: 1.0.0
description: "Douyin video and comment data. Search videos, inspect video details, collect comments, and drill into comment replies for content research and audience feedback analysis."
metadata:
  source: douyin
  requires:
    bins: ["apimux"]
  cliHelp: "apimux douyin --help"
---

# Douyin

Search and inspect Douyin videos, details, comments, and comment replies. Use this for content research, hot-video monitoring, and comment sampling.

**Before using:** Read [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, pagination metadata, and CLI conventions.

## What you can do

- **Find video examples** → `search_videos`
- **Inspect one video** → `get_video_detail`
- **Collect comments for a video** → `get_video_comments`
- **Drill into comment replies** → `get_comment_replies`

## Available capabilities

| Capability | What it does | When to use |
|------------|--------------|-------------|
| `search_videos` | Search Douyin videos by keyword | Discover videos and monitor topics |
| `get_video_detail` | Get details for one video | Inspect metadata, author, and engagement |
| `get_video_comments` | List comments for one video | Sample audience feedback |
| `get_comment_replies` | List replies under one parent comment | Analyze a comment thread |

---

## douyin.search_videos

Search Douyin videos by keyword.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `keyword` | string | Yes | Search keyword |
| `sort_type` | string | No | `comprehensive`, `likes`, or `latest`; default `comprehensive` |
| `publish_time` | string | No | `all`, `1d`, `1w`, or `6m`; default `all` |
| `filter_duration` | string | No | `all`, `under_1m`, `1m_5m`, or `over_5m`; default `all` |
| `content_type` | string | No | `all`, `video`, `image`, or `article`; default `all` |
| `cursor` | integer | No | Pagination cursor; omit for the first page |

### CLI usage

```bash
apimux douyin search_videos --keyword "desk setup"
apimux douyin search_videos --keyword "desk setup" --sort-type "likes" --publish-time "1w"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `aweme_id` | string | Video ID |
| `description` | string | Video description |
| `create_time` | string | RFC3339 publish time |
| `author` | object | Author information |
| `statistics` | object | Like/comment/share/play statistics |
| `video` | object | Video duration and aspect information |
| `share_url` | string | Share URL |

### Notes

- `keyword` is required.
- Prefer the string enum values above. For backward compatibility, the CLI also accepts legacy aliases for `sort_type` (`0/1/2`) and `publish_time` (`0/1/7/180`).
- Pagination state is returned in `meta.cursor` and `meta.has_more`.

---

## douyin.get_video_detail

Get details for one Douyin video.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `aweme_id` | string | Yes | Numeric video ID |

### CLI usage

```bash
apimux douyin get_video_detail --aweme-id "7489123456789012345"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `aweme_id` | string | Video ID |
| `description` | string | Video description |
| `create_time` | string | RFC3339 publish time |
| `author` | object | Author information |
| `statistics` | object | Like/comment/share/play statistics |
| `video` | object | Video duration and aspect information |
| `share_url` | string | Share URL |

### Notes

- `aweme_id` must be a numeric string.
- Missing videos return `video_not_found`.

---

## douyin.get_video_comments

List comments for one Douyin video.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `aweme_id` | string | Yes | Numeric video ID |
| `cursor` | integer | No | Pagination cursor; omit for the first page |
| `count` | integer | No | Page size; default `20` |

### CLI usage

```bash
apimux douyin get_video_comments --aweme-id "7489123456789012345"
apimux douyin get_video_comments --aweme-id "7489123456789012345" --count 20
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `comment_id` | string | Comment ID |
| `text` | string | Comment text |
| `like_count` | integer | Like count |
| `reply_count` | integer | Reply count |
| `create_time` | string | RFC3339 comment time |
| `author` | object | Comment author information |

### Notes

- `aweme_id` must be a numeric string.
- Pagination state is returned in `meta.cursor` and `meta.has_more`.
- Total comment count is returned in `meta.total` when available.

---

## douyin.get_comment_replies

List replies under one parent comment.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `aweme_id` | string | Yes | Numeric video ID |
| `comment_id` | string | Yes | Parent comment ID |
| `cursor` | integer | No | Pagination cursor; omit for the first page |
| `count` | integer | No | Page size; default `20` |

### CLI usage

```bash
apimux douyin get_comment_replies --aweme-id "7489123456789012345" --comment-id "1234567890"
apimux douyin get_comment_replies --aweme-id "7489123456789012345" --comment-id "1234567890" --count 20
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `comment_id` | string | Reply comment ID |
| `text` | string | Reply text |
| `like_count` | integer | Like count |
| `reply_count` | integer | Nested reply count |
| `create_time` | string | RFC3339 reply time |
| `author` | object | Reply author information |

### Notes

- `aweme_id` must be a numeric string.
- `comment_id` is required.
- Pagination state is returned in `meta.cursor` and `meta.has_more`.
- Total reply count is returned in `meta.total` when available.

---

## General notes

- See [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure and error handling.
