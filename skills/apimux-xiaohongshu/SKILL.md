---
name: apimux-xiaohongshu
version: 1.0.0
description: "Xiaohongshu content data. Search notes, inspect note details, and collect comments for content research and audience feedback analysis."
metadata:
  source: xiaohongshu
  requires:
    bins: ["apimux"]
  cliHelp: "apimux xiaohongshu --help"
---

# Xiaohongshu

Search Xiaohongshu notes, inspect note details, and collect comments. Use this for content research, seed-content discovery, note monitoring, and comment sampling.

**Before using:** Read [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, pagination metadata, and CLI conventions.

## What you can do

- **Find note examples** → `search_notes`
- **Inspect one note** → `get_note_detail`
- **Collect note comments** → `get_note_comments`

## Available capabilities

| Capability | What it does | When to use |
|------------|--------------|-------------|
| `search_notes` | Search Xiaohongshu notes | Discover content examples and monitor keywords |
| `get_note_detail` | Get one note's details | Inspect content, author, engagement, images, and tags |
| `get_note_comments` | Get note comments | Sample audience feedback |

---

## xiaohongshu.search_notes

Search Xiaohongshu notes by keyword.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `keyword` | string | Yes | Search keyword |
| `page` | integer | No | Page number, starting at 1; default `1` |
| `note_type` | string | No | `all`, `video`, `normal`, or `live`; default `all` |
| `time_filter` | string | No | `all`, `1d`, `1w`, or `6m`; default `all` |
| `sort_strategy` | string | No | `default`, `latest`, or `likes`; default `default` |

### CLI usage

```bash
apimux xiaohongshu search_notes --keyword "desk setup"
apimux xiaohongshu search_notes --keyword "desk setup" --sort-strategy "likes" --note-type "video"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `note_id` | string | Note ID |
| `title` | string | Note title |
| `description` | string | Note description |
| `type` | string | Note type |
| `xsec_token` | string | Security token from search results; may be needed by `get_note_detail` |
| `like_count` | integer | Like count; may be `null` in search results |
| `collect_count` | integer | Collection count; may be `null` in search results |
| `comment_count` | integer | Comment count; may be `null` in search results |
| `author` | object | Author information |

### Notes

- `keyword` is required.
- Prefer the string enum values above. For backward compatibility, the CLI also accepts legacy values for `sort_strategy` (`general/time_descending/popularity_descending`) and `note_type` (`0/1/2`).
- Pagination uses the `page` parameter. Page state is returned in `meta.current_page`, `meta.next_page`, and `meta.has_more`.

---

## xiaohongshu.get_note_detail

Get details for one Xiaohongshu note.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `note_id` | string | Yes | 24-character hex note ID |
| `xsec_token` | string | No | Security token from `search_notes`; required for some notes |

### CLI usage

```bash
apimux xiaohongshu get_note_detail --note-id "64f1a2b3c4d5e6f789abcdef"
apimux xiaohongshu get_note_detail --note-id "64f1a2b3c4d5e6f789abcdef" --xsec-token "..."
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `note_id` | string | Note ID |
| `title` | string | Note title |
| `description` | string | Note description |
| `type` | string | Note type |
| `user_id` | string | Author user ID |
| `nickname` | string | Author nickname |
| `avatar` | string | Author avatar URL |
| `like_count` | integer | Like count |
| `collect_count` | integer | Collection count |
| `comment_count` | integer | Comment count |
| `share_count` | integer | Share count |
| `images` | string[] | Image URLs |
| `tags` | string[] | Tags |
| `time` | string | Publish time |
| `last_update_time` | string | Last update time |
| `ip_location` | string | IP location |
| `video_url` | string | Video URL for video notes |

### Notes

- `note_id` must be a 24-character hex string.
- If `search_notes` returns `xsec_token`, pass it when calling details.
- Share links such as `xhslink.com/...` are not accepted; pass the note ID.
- Missing notes return `note_not_found`.

---

## xiaohongshu.get_note_comments

Get comments for one Xiaohongshu note.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `note_id` | string | Yes | 24-character hex note ID |
| `cursor` | string | No | Comment pagination cursor; omit for the first page |
| `sort_strategy` | string | No | `default`, `latest`, or `likes`; default `default` |

### CLI usage

```bash
apimux xiaohongshu get_note_comments --note-id "64f1a2b3c4d5e6f789abcdef"
apimux xiaohongshu get_note_comments --note-id "64f1a2b3c4d5e6f789abcdef" --sort-strategy "latest"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `comment_id` | string | Comment ID |
| `user_id` | string | Author user ID |
| `nickname` | string | Author nickname |
| `avatar` | string | Author avatar URL |
| `content` | string | Comment text |
| `like_count` | integer | Like count |
| `reply_count` | integer | Reply count |
| `create_time` | string | Publish time |
| `ip_location` | string | IP location |

### Notes

- `note_id` must be a 24-character hex string.
- Pagination state is returned in `meta.cursor` and `meta.has_more`.

---

## General notes

- See [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure and error handling.
