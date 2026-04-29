---
name: apimux-reddit
version: 1.0.0
description: "Reddit content data. Search posts, inspect subreddit feeds, get post details, and collect comments for topic research and community monitoring."
metadata:
  source: reddit
  requires:
    bins: ["apimux"]
  cliHelp: "apimux reddit --help"
---

# Reddit

Search and inspect Reddit posts, subreddit feeds, post details, and comments. Use this for topic research, community monitoring, and discussion analysis.

**Before using:** Read [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, pagination metadata, and CLI conventions.

## What you can do

- **Find posts by keyword** → `search`
- **Monitor a subreddit feed** → `get_subreddit_feed`
- **Inspect one post** → `get_post_detail`
- **Collect comments for a post** → `get_post_comments`

## Available capabilities

| Capability | What it does | When to use |
|------------|--------------|-------------|
| `search` | Search Reddit content | Topic search and post discovery |
| `get_subreddit_feed` | Get a subreddit feed | Monitor a community or board |
| `get_post_detail` | Get one post's details | Drill into a specific post |
| `get_post_comments` | Get comments for a post | Sample opinions and discussion |

---

## reddit.search

Search Reddit content by query.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | Yes | Search keyword |
| `search_type` | string | No | `post`, `community`, `comment`, `media`, or `people`; default `post` |
| `sort` | string | No | `relevance`, `hot`, `top`, `new`, or `comments`; default `relevance` |
| `time_range` | string | No | `all`, `year`, `month`, `week`, `day`, or `hour`; default `all` |
| `after` | string | No | Pagination cursor; omit for the first page |

### CLI usage

```bash
apimux reddit search --query "standing desk"
apimux reddit search --query "standing desk" --sort "top" --time-range "month"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `post_id` | string | Post ID, usually `t3_xxxxx` |
| `title` | string | Post title |
| `subreddit` | string | Subreddit name |
| `author` | string | Author username |
| `score` | integer | Reddit score |
| `num_comments` | integer | Comment count |
| `created_at` | string | RFC3339 publish time |
| `permalink` | string | Reddit permalink |
| `url` | string | Target URL |
| `selftext` | string | Post body text |
| `thumbnail` | string | Thumbnail URL |
| `is_video` | boolean | Whether the post is a video post |

### Notes

- `query` is required.
- Prefer lowercase `sort` values. The CLI also accepts uppercase values and normalizes them.
- Pagination state is returned in `meta.cursor` and `meta.has_more`.

---

## reddit.get_subreddit_feed

Get posts from one subreddit feed.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `subreddit_name` | string | Yes | Subreddit name without the `r/` prefix |
| `sort` | string | No | `best`, `hot`, `new`, `top`, `controversial`, or `rising`; default `hot` |
| `after` | string | No | Pagination cursor; omit for the first page |

### CLI usage

```bash
apimux reddit get_subreddit_feed --subreddit-name "MechanicalKeyboards"
apimux reddit get_subreddit_feed --subreddit-name "MechanicalKeyboards" --sort "top"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `post_id` | string | Post ID, usually `t3_xxxxx` |
| `title` | string | Post title |
| `subreddit` | string | Subreddit name |
| `author` | string | Author username |
| `score` | integer | Reddit score |
| `num_comments` | integer | Comment count |
| `created_at` | string | RFC3339 publish time |
| `permalink` | string | Reddit permalink |
| `url` | string | Target URL |
| `selftext` | string | Post body text |
| `thumbnail` | string | Thumbnail URL |
| `is_video` | boolean | Whether the post is a video post |

### Notes

- `subreddit_name` is required and must not include `r/`.
- Pagination state is returned in `meta.cursor` and `meta.has_more`.

---

## reddit.get_post_detail

Get details for one Reddit post.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `post_id` | string | Yes | Post ID in `t3_xxxxx` format |

### CLI usage

```bash
apimux reddit get_post_detail --post-id "t3_abcdef"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `post_id` | string | Post ID |
| `title` | string | Post title |
| `subreddit` | string | Subreddit name |
| `author` | string | Author username |
| `score` | integer | Reddit score |
| `upvote_ratio` | number | Upvote ratio |
| `num_comments` | integer | Comment count |
| `created_at` | string | RFC3339 publish time |
| `permalink` | string | Reddit permalink |
| `url` | string | Target URL |
| `selftext` | string | Post body text |
| `thumbnail` | string | Thumbnail URL |
| `is_video` | boolean | Whether the post is a video post |
| `link_flair_text` | string | Link flair text |

### Notes

- `post_id` must use `t3_xxxxx` format.
- Missing posts return `post_not_found`.

---

## reddit.get_post_comments

Get comments for one Reddit post.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `post_id` | string | Yes | Post ID in `t3_xxxxx` format |
| `sort_type` | string | No | `confidence`, `new`, `top`, `hot`, `controversial`, `old`, or `random`; default `confidence` |
| `after` | string | No | Pagination cursor; omit for the first page |

### CLI usage

```bash
apimux reddit get_post_comments --post-id "t3_abcdef"
apimux reddit get_post_comments --post-id "t3_abcdef" --sort-type "top"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Comment ID |
| `author` | string | Comment author |
| `body` | string | Comment body |
| `score` | integer | Comment score |
| `created_at` | string | RFC3339 publish time |
| `parent_id` | string | Parent node ID |
| `depth` | integer | Comment depth |

### Notes

- `post_id` must use `t3_xxxxx` format.
- Pagination state is returned in `meta.cursor` and `meta.has_more`.

---

## General notes

- See [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure and error handling.
