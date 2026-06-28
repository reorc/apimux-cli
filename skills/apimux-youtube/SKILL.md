---
name: apimux-youtube
description: "YouTube public content data. Search videos, inspect video and channel details, fetch transcripts, list channel videos, and collect comments or replies for brand monitoring, creator research, and content analysis."
metadata:
  source: youtube
  requires:
    bins: ["apimux"]
  cliHelp: "apimux youtube --help"
---

# YouTube

Search and inspect public YouTube videos, channels, transcripts, comments, and comment replies. Use this for brand monitoring, creator/content research, topic validation, and audience feedback analysis.

**Before using:** Read [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, pagination metadata, and CLI conventions.

**CLI note:** Use the YouTube source subcommands for normal workflows:

```bash
apimux youtube search_videos --query "standing desk" --count 10
```

## What you can do

- **Find videos by topic** -> `youtube.search_videos`
- **Inspect one video** -> `youtube.get_video_detail`
- **Fetch video captions/transcripts** -> `youtube.get_video_transcript`
- **Inspect one channel** -> `youtube.get_channel_detail`
- **List channel videos or Shorts** -> `youtube.get_channel_videos`
- **Collect video comments** -> `youtube.get_video_comments`
- **Collect replies to one comment** -> `youtube.get_comment_replies`

## Available capabilities

| Capability | What it does | When to use |
|------------|--------------|-------------|
| `youtube.search_videos` | Search public YouTube videos | Topic discovery, brand monitoring, competitor video discovery |
| `youtube.get_video_detail` | Get one video by `video_id` or URL | Inspect metadata, channel, engagement, description, and thumbnail |
| `youtube.get_video_transcript` | Get caption/transcript text | Summarize video content, extract talking points, compare claims |
| `youtube.get_channel_detail` | Get one channel by channel ID, handle, or URL | Creator/channel qualification and audience sizing |
| `youtube.get_channel_videos` | List videos or Shorts from one channel | Creator auditing and recent-content sampling |
| `youtube.get_video_comments` | List comments for one video | Audience feedback and objection mining |
| `youtube.get_comment_replies` | List replies for one parent comment | Thread-level discussion analysis |

## Common workflows

- Unknown video: use `youtube.search_videos`, pick a `video_id` or `url`, then call `youtube.get_video_detail`, `youtube.get_video_transcript`, or `youtube.get_video_comments`.
- Channel research: call `youtube.get_channel_detail` by `--handle`, `--channel-id`, or `--url`, then sample recent content with `youtube.get_channel_videos`.
- Shorts vs videos: set `content_type` to `shorts`, `video`, or `all` in `youtube.get_channel_videos`.
- Comment mining: use `youtube.get_video_comments`, then pass a returned `comment_id` into `youtube.get_comment_replies`.

---

## youtube.search_videos

Search public YouTube videos by query.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | Yes | YouTube search query |
| `cursor` | string | No | Continuation cursor from a previous page |
| `count` | integer | No | Maximum number of videos to return |

### CLI usage

```bash
apimux youtube search_videos --query "desk setup" --count 10
apimux youtube search_videos --query "desk setup" --cursor "<cursor>" --count 10
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `video_id` | string | YouTube video ID |
| `title` | string | Video title |
| `description` | string | Video description or snippet |
| `published_time` | string | Publish timestamp or display text |
| `duration` | string | Duration display when available |
| `view_count` | integer | View count |
| `like_count` | integer | Like count |
| `comment_count` | integer | Comment count |
| `channel_id` | string | Channel ID |
| `channel_handle` | string | Channel handle |
| `channel_title` | string | Channel title |
| `channel_subscriber_count` | integer | Channel subscriber count when returned |
| `channel_is_verified` | boolean | Whether the channel is verified when returned |

Default compact list output hides canonical watch URLs and thumbnail/CDN URLs. Use `--output data` when raw URLs are required.

---

## youtube.get_video_detail

Fetch one public YouTube video detail by `video_id` or URL.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `video_id` | string | Conditional | YouTube video ID. Use either `video_id` or `url`. |
| `url` | string | Conditional | YouTube video URL. Use either `video_id` or `url`. |

### CLI usage

```bash
apimux youtube get_video_detail --video-id "dQw4w9WgXcQ"
apimux youtube get_video_detail --url "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
```

### Notes

- Provide at least one identity input: `video_id` or `url`.
- Output fields match `youtube.search_videos` with fuller detail when the provider returns it.

---

## youtube.get_video_transcript

Fetch caption or transcript text for one public YouTube video.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `video_id` | string | Yes | YouTube video ID |
| `language` | string | No | Caption language code |
| `format` | string | No | `segments` or `plain`; default `segments` |

### CLI usage

```bash
apimux youtube get_video_transcript --video-id "dQw4w9WgXcQ" --format "plain"
apimux youtube get_video_transcript --video-id "dQw4w9WgXcQ" --language "en" --format "segments"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `video_id` | string | YouTube video ID |
| `language` | string | Selected language code |
| `segments` | object[] | Timed caption segments when `format=segments` |
| `text` | string | Joined transcript text when `format=plain` |

---

## youtube.get_channel_detail

Fetch one YouTube channel summary by channel ID, handle, or URL.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `channel_id` | string | Conditional | Canonical YouTube channel ID |
| `handle` | string | Conditional | YouTube handle without `@` |
| `url` | string | Conditional | YouTube channel URL |

### CLI usage

```bash
apimux youtube get_channel_detail --handle "mkbhd"
apimux youtube get_channel_detail --url "https://www.youtube.com/@mkbhd"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `channel_id` | string | Canonical channel ID |
| `title` | string | Channel title |
| `handle` | string | Public handle |
| `url` | string | Canonical channel URL |
| `description` | string | Channel description |
| `avatar` | string | Avatar URL |
| `subscriber_count` | integer | Subscriber count |
| `video_count` | integer | Video count |
| `view_count` | integer | Total view count |

---

## youtube.get_channel_videos

List videos or Shorts from one YouTube channel.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `channel_id` | string | Conditional | Canonical YouTube channel ID |
| `handle` | string | Conditional | YouTube handle without `@` |
| `url` | string | Conditional | YouTube channel URL |
| `content_type` | string | No | `video`, `shorts`, or `all`; default `video` |
| `cursor` | string | No | Continuation cursor |
| `count` | integer | No | Maximum items to return |

### CLI usage

```bash
apimux youtube get_channel_videos --handle "mkbhd" --content-type "video" --count 10
apimux youtube get_channel_videos --handle "mkbhd" --content-type "shorts" --count 10
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `video_id` | string | YouTube video ID |
| `title` | string | Video title |
| `description` | string | Video description or snippet |
| `published_time` | string | Publish timestamp or display text |
| `duration` | string | Duration display when available |
| `view_count` | integer | View count |
| `like_count` | integer | Like count |
| `comment_count` | integer | Comment count |
| `channel_id` | string | Channel ID |
| `channel_handle` | string | Channel handle |
| `channel_title` | string | Channel title |
| `channel_subscriber_count` | integer | Channel subscriber count when returned |
| `channel_is_verified` | boolean | Whether the channel is verified when returned |

Default compact list output hides canonical watch URLs and thumbnail/CDN URLs. Use `--output data` when raw URLs are required.

---

## youtube.get_video_comments

List comments for one public YouTube video by `video_id` or URL.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `video_id` | string | Conditional | YouTube video ID |
| `url` | string | Conditional | YouTube video URL |
| `cursor` | string | No | Continuation cursor |
| `count` | integer | No | Maximum comments to return |

### CLI usage

```bash
apimux youtube get_video_comments --video-id "dQw4w9WgXcQ" --count 20
apimux youtube get_video_comments --url "https://www.youtube.com/watch?v=dQw4w9WgXcQ" --count 20
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `comment_id` | string | YouTube comment ID |
| `text` | string | Comment text |
| `create_time` | string | RFC3339 timestamp or provider display text |
| `like_count` | integer | Like count |
| `reply_count` | integer | Reply count |
| `author_channel_id` | string | Author channel ID |
| `author_handle` | string | Author handle |
| `author_name` | string | Author display name |

Default compact comment output hides author avatar/channel URL fields. Use `--output data` when raw author objects or URLs are required.

---

## youtube.get_comment_replies

List replies for one YouTube parent comment.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `comment_id` | string | Yes | Parent comment ID |
| `video_id` | string | No | YouTube video ID |
| `cursor` | string | No | Continuation cursor |
| `count` | integer | No | Maximum replies to return |

### CLI usage

```bash
apimux youtube get_comment_replies --comment-id "Ug..." --count 20
apimux youtube get_comment_replies --video-id "dQw4w9WgXcQ" --comment-id "Ug..." --count 20
```

### Notes

- `comment_id` is required.
- Use `meta.cursor` and `meta.has_more` for pagination when returned.
- Default compact reply output keeps short author identity fields and hides author avatar/channel URL fields. Use `--output data` for raw objects.
