---
name: apimux-instagram
description: "Instagram public content data. Search users, hashtags, and reels; inspect profiles, posts, reels, comments, and replies for creator research, brand monitoring, and audience feedback analysis."
metadata:
  source: instagram
  requires:
    bins: ["apimux"]
  cliHelp: "apimux instagram --help"
---

# Instagram

Search and inspect public Instagram users, hashtags, reels, profiles, posts, comments, and comment replies. Use this for creator research, brand monitoring, content sampling, and audience feedback analysis.

**Before using:** Read [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, pagination metadata, and CLI conventions.

**CLI note:** Use the Instagram source subcommands for normal workflows:

```bash
apimux instagram search_reels --query "desk setup" --count 10
```

## What you can do

- **Find users by keyword** -> `instagram.search_users`
- **Find hashtags by keyword** -> `instagram.search_hashtags`
- **Search public reels** -> `instagram.search_reels`
- **Inspect one profile** -> `instagram.get_user_profile`
- **List one user's posts** -> `instagram.get_user_posts`
- **List one user's reels** -> `instagram.get_user_reels`
- **Inspect one post or reel** -> `instagram.get_post_detail`
- **Collect comments** -> `instagram.get_post_comments`
- **Collect replies to one comment** -> `instagram.get_comment_replies`

## Available capabilities

| Capability | What it does | When to use |
|------------|--------------|-------------|
| `instagram.search_users` | Search public Instagram users | Creator discovery and profile lookup |
| `instagram.search_hashtags` | Search hashtags | Topic and hashtag sizing |
| `instagram.search_reels` | Search public reels by keyword | Reels discovery and content research |
| `instagram.get_user_profile` | Get one profile by username or user ID | Creator qualification and profile summary |
| `instagram.get_user_posts` | List public posts for one user | Profile content sampling |
| `instagram.get_user_reels` | List public reels for one user | Short-video content sampling |
| `instagram.get_post_detail` | Get one post/reel by shortcode, media ID, or URL | Inspect metadata, media, owner, and engagement |
| `instagram.get_post_comments` | List comments for one post/reel | Audience feedback and objection mining |
| `instagram.get_comment_replies` | List replies for one parent comment | Thread-level discussion analysis |

## Common workflows

- Creator discovery: use `instagram.search_users`, pick `username` or `user_id`, then call `instagram.get_user_profile`, `instagram.get_user_posts`, or `instagram.get_user_reels`.
- Reels research: use `instagram.search_reels` for topic discovery, then inspect promising results with `instagram.get_post_detail`.
- Post URL known: call `instagram.get_post_detail` by `url`, then use `instagram.get_post_comments`.
- Comment mining: use `instagram.get_post_comments`, then pass a returned `comment_id` into `instagram.get_comment_replies`.

---

## instagram.search_users

Search public Instagram users by keyword.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | Yes | Instagram user search query |
| `cursor` | string | No | Pagination cursor |
| `count` | integer | No | Maximum users to return |

### CLI usage

```bash
apimux instagram search_users --query "desk setup" --count 10
apimux instagram search_users --query "desk setup" --cursor "<cursor>" --count 10
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `user_id` | string | Instagram numeric user ID |
| `username` | string | Username |
| `full_name` | string | Display name |
| `url` | string | Canonical profile URL |
| `bio` | string | Biography |
| `category` | string | Business/category label |
| `profile_pic_url` | string | Avatar URL |
| `follower_count` | integer | Follower count |
| `following_count` | integer | Following count |
| `media_count` | integer | Public media count |
| `is_business` | boolean | Business account indicator |
| `is_verified` | boolean | Verified badge |

---

## instagram.search_hashtags

Search Instagram hashtags by keyword.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | Yes | Instagram hashtag search query |

### CLI usage

```bash
apimux instagram search_hashtags --query "desksetup"
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `hashtag_id` | string | Hashtag ID |
| `name` | string | Hashtag name |
| `media_count` | integer | Public media count |

---

## instagram.search_reels

Search public Instagram reels by keyword.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | Yes | Instagram reels search query |
| `cursor` | string | No | Pagination cursor |
| `count` | integer | No | Maximum reels to return |

### CLI usage

```bash
apimux instagram search_reels --query "desk setup" --count 10
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `media_id` | string | Instagram media ID |
| `shortcode` | string | Shortcode |
| `caption` | string | Caption text |
| `description` | string | Description text when returned |
| `media_type` | string | Media type |
| `taken_at` | string | Timestamp or provider display value |
| `like_count` | integer | Like count |
| `comment_count` | integer | Comment count |
| `view_count` | integer | View/play count |
| `owner` | object | Compact owner summary: `user_id`, `username`, `full_name`, `is_verified`, `is_business` |

Default compact list output hides long CDN/media fields such as `media_url`, `thumbnail`, `profile_pic_url`, and `permalink`. Use `--output data` when raw URLs are required.

---

## instagram.get_user_profile

Fetch one public Instagram profile by username or user ID.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `username` | string | Conditional | Instagram username |
| `user_id` | string | Conditional | Instagram numeric user ID |

### CLI usage

```bash
apimux instagram get_user_profile --username "natgeo"
apimux instagram get_user_profile --user-id "123456789"
```

### Notes

- Provide at least one identity input: `username` or `user_id`.
- Output fields match `instagram.search_users` with fuller detail when the provider returns it.

---

## instagram.get_user_posts

List public posts for one Instagram user.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `username` | string | Conditional | Instagram username |
| `user_id` | string | Conditional | Instagram numeric user ID |
| `cursor` | string | No | Pagination cursor |
| `count` | integer | No | Maximum posts to return |

### CLI usage

```bash
apimux instagram get_user_posts --username "natgeo" --count 12
apimux instagram get_user_posts --user-id "123456789" --count 12
```

### Notes

- Provide at least one identity input: `username` or `user_id`.
- Response fields match the media fields shown under `instagram.search_reels`.

---

## instagram.get_user_reels

List public reels for one Instagram user.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `username` | string | Conditional | Instagram username |
| `user_id` | string | Conditional | Instagram numeric user ID |
| `cursor` | string | No | Pagination cursor |
| `count` | integer | No | Maximum reels to return |

### CLI usage

```bash
apimux instagram get_user_reels --username "natgeo" --count 12
```

### Notes

- Provide at least one identity input: `username` or `user_id`.
- Response fields match the media fields shown under `instagram.search_reels`.

---

## instagram.get_post_detail

Fetch one public Instagram post or reel by shortcode, media ID, or URL.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `shortcode` | string | Conditional | Instagram post shortcode |
| `media_id` | string | Conditional | Instagram media ID |
| `url` | string | Conditional | Instagram post/reel URL |

### CLI usage

```bash
apimux instagram get_post_detail --shortcode "C0abc123xyz"
apimux instagram get_post_detail --url "https://www.instagram.com/p/C0abc123xyz/"
```

### Notes

- Provide at least one identity input: `shortcode`, `media_id`, or `url`.
- Response fields match the media fields shown under `instagram.search_reels`.

---

## instagram.get_post_comments

List comments for one public Instagram post or reel.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `shortcode` | string | Conditional | Instagram post shortcode |
| `media_id` | string | Conditional | Instagram media ID |
| `url` | string | Conditional | Instagram post/reel URL |
| `cursor` | string | No | Pagination cursor |
| `count` | integer | No | Maximum comments to return |

### CLI usage

```bash
apimux instagram get_post_comments --shortcode "C0abc123xyz" --count 20
apimux instagram get_post_comments --url "https://www.instagram.com/p/C0abc123xyz/" --count 20
```

### Response fields

| Field | Type | Description |
|-------|------|-------------|
| `comment_id` | string | Comment ID |
| `text` | string | Comment text |
| `create_time` | string | Timestamp or provider display value |
| `like_count` | integer | Like count |
| `reply_count` | integer | Reply count |
| `author` | object | Compact author summary: `user_id`, `username`, `full_name`, `is_verified`, `is_business` when returned |

Default compact comment output hides avatar/CDN URL fields. Use `--output data` when raw author objects or URLs are required.

---

## instagram.get_comment_replies

List replies for one Instagram parent comment.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `shortcode` | string | Conditional | Instagram post/reel shortcode; provide one of `shortcode`, `url`, or `media_id` |
| `url` | string | Conditional | Instagram post/reel URL; provide one of `shortcode`, `url`, or `media_id` |
| `media_id` | string | Conditional | Instagram media ID; provide one of `shortcode`, `url`, or `media_id` |
| `comment_id` | string | Yes | Parent comment ID |
| `cursor` | string | No | Pagination cursor |
| `count` | integer | No | Maximum replies to return |

### CLI usage

```bash
apimux instagram get_comment_replies --shortcode "C0abc123xyz" --comment-id "18000000000000000" --count 20
apimux instagram get_comment_replies --url "https://www.instagram.com/p/C0abc123xyz/" --comment-id "18000000000000000" --count 20
```

### Notes

- `comment_id` is required, and the post/reel identity is also required via one of `shortcode`, `url`, or `media_id`.
- Use `meta.cursor` and `meta.has_more` for pagination when returned.
- Default compact reply output keeps short author identity fields and hides avatar/CDN URL fields. Use `--output data` for raw objects.
