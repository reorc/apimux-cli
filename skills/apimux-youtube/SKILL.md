---
name: apimux-youtube
version: 1.0.0
description: "YouTube public video and channel data. Search videos, inspect video/channel metadata, collect comments and replies, and fetch transcripts for brand monitoring and social research."
metadata:
  source: youtube
  requires:
    bins: ["apimux"]
  cliHelp: "apimux youtube --help"
---

# YouTube

Search and inspect public YouTube videos, channels, comments, replies, and transcripts. Use this for brand monitoring, creator research, competitor tracking, and video/comment summarization.

**Before using:** Read [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, pagination metadata, and CLI conventions.

## What you can do

- **Find videos by topic or brand** -> `search_videos`
- **Inspect one video** -> `get_video_detail`
- **Collect video comments** -> `get_video_comments`
- **Collect replies under one comment** -> `get_comment_replies`
- **Inspect a channel** -> `get_channel_detail`
- **List a creator/channel's videos or Shorts** -> `get_channel_videos`
- **Fetch captions/transcript text** -> `get_video_transcript`

## Common workflows

- Brand monitoring: `search_videos --query "Pimax VR"` -> `get_video_detail` -> `get_video_comments`.
- Creator research: `get_channel_detail --handle @PimaxOfficial` -> `get_channel_videos --handle @PimaxOfficial`.
- Summarization: `get_video_transcript --video-id ... --format plain` before summarizing long videos.
- Comment mining: use `get_video_comments --count 50` for audience feedback, then `get_comment_replies` for discussion threads.

## CLI examples

```bash
apimux youtube search_videos --query "Pimax VR" --count 20
apimux youtube get_video_detail --url "https://www.youtube.com/watch?v=VIDEO_ID"
apimux youtube get_video_comments --url "https://www.youtube.com/watch?v=VIDEO_ID" --count 50
apimux youtube get_comment_replies --comment-id "COMMENT_ID" --video-id "VIDEO_ID"
apimux youtube get_channel_detail --handle "@PimaxOfficial"
apimux youtube get_channel_videos --handle "@PimaxOfficial" --content-type all --count 30
apimux youtube get_video_transcript --video-id "VIDEO_ID" --language en --format plain
```

## Capability notes

| Capability | Required identity | Notes |
|------------|-------------------|-------|
| `search_videos` | `query` | First wave supports `query`, `cursor`, and `count`. |
| `get_video_detail` | `video_id` or `url` | Provide exactly one identity. |
| `get_video_comments` | `video_id` or `url` | Returns canonical comments plus pagination metadata when available. |
| `get_comment_replies` | `comment_id` | Pass `video_id` when available to improve provider resolution. |
| `get_channel_detail` | `channel_id`, `handle`, or `url` | Provide exactly one channel identity. |
| `get_channel_videos` | `channel_id`, `handle`, or `url` | Use `content_type=video`, `shorts`, or `all`. |
| `get_video_transcript` | `video_id` | Use `format=plain` for summarization input, `segments` for timestamps. |

## Errors

- `missing_video_identity`, `ambiguous_video_identity`, `invalid_video_id`
- `missing_channel_identity`, `ambiguous_channel_identity`, `invalid_channel_id`, `invalid_handle`
- `invalid_content_type`, `invalid_format`
- `youtube_failed`, `video_not_found`, `invalid_provider_payload`
