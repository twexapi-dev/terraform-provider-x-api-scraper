---
page_title: "x-api-scraper_tweet Resource"
---

# x-api-scraper_tweet

Post a tweet. Destroy deletes that tweet.

> **Note:** Public getting-started docs focus on API-key read data sources. This resource requires account session credentials and is not covered in the public getting started guide.

## Schema

### Required

- `tweet_content` (String)
- `username` (String) Acting account username, used when deleting.

### Optional

- `reply_tweet_id` (String)

### Read-Only

- `id` (String)
- `tweet_id` (String)
