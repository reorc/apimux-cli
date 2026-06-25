---
name: apimux-instagram
version: 1.0.0
description: "Instagram public profile and media data. Search users, reels, hashtags, inspect profiles/posts/reels, and collect comments and replies for brand monitoring and social research."
metadata:
  source: instagram
  requires:
    bins: ["apimux"]
  cliHelp: "apimux instagram --help"
---

# Instagram

Search and inspect public Instagram users, reels, hashtags, posts/reels, comments, and replies. Use this for brand monitoring, creator discovery, UGC research, and public social feedback analysis.

**Before using:** Read [`../apimux-shared/SKILL.md`](../apimux-shared/SKILL.md) for response structure, error handling, pagination metadata, and CLI conventions.

## What you can do

- **Find accounts by brand/name** -> `search_users`
- **Find reels by topic** -> `search_reels`
- **Find hashtags by topic** -> `search_hashtags`
- **Inspect one public profile** -> `get_user_profile`
- **List a user's posts** -> `get_user_posts`
- **List a user's reels** -> `get_user_reels`
- **Inspect one post or reel** -> `get_post_detail`
- **Collect comments** -> `get_post_comments`
- **Collect replies under one comment** -> `get_comment_replies`

## Common workflows

- Brand monitoring: `search_users --query pimax` -> `get_user_profile --username ...` -> `get_user_posts`.
- Reel research: `search_reels --query "VR headset"` -> `get_post_detail` -> `get_post_comments`.
- Hashtag scouting: `search_hashtags --query pimax` to discover public topic clusters.
- Comment mining: use `get_post_comments` and `get_comment_replies` to collect customer language and objections.

## CLI examples

```bash
apimux instagram search_users --query "pimax"
apimux instagram search_reels --query "VR headset" --count 20
apimux instagram search_hashtags --query "pimax"
apimux instagram get_user_profile --username "pimaxofficial"
apimux instagram get_user_posts --username "pimaxofficial" --count 30
apimux instagram get_user_reels --username "pimaxofficial"
apimux instagram get_post_detail --url "https://www.instagram.com/p/CODE/"
apimux instagram get_post_comments --url "https://www.instagram.com/p/CODE/" --count 50
apimux instagram get_comment_replies --comment-id "COMMENT_ID" --media-id "MEDIA_ID"
```

## Capability notes

| Capability | Required identity | Notes |
|------------|-------------------|-------|
| `search_users` | `query` | Account discovery by keyword. |
| `search_reels` | `query` | Reel discovery by keyword. |
| `search_hashtags` | `query` | Hashtag discovery by keyword. |
| `get_user_profile` | `username` or `user_id` | Provide exactly one user identity. |
| `get_user_posts` | `username` or `user_id` | Paginates through public posts. |
| `get_user_reels` | `username` or `user_id` | Paginates through public reels. |
| `get_post_detail` | `shortcode`, `media_id`, or `url` | Provide exactly one media identity. |
| `get_post_comments` | `shortcode`, `media_id`, or `url` | Returns canonical comments plus pagination metadata when available. |
| `get_comment_replies` | `comment_id` | Pass `media_id` when available to improve provider resolution. |

## Errors

- `missing_user_identity`, `ambiguous_user_identity`, `invalid_username`, `invalid_user_id`
- `missing_media_identity`, `ambiguous_media_identity`, `invalid_shortcode`, `invalid_media_id`, `invalid_url`
- `instagram_failed`, `invalid_provider_payload`
