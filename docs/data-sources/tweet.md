---
page_title: "x-api-scraper_tweet Data Source"
---

# x-api-scraper_tweet

Read a tweet by id. Maps to the competitor `x_tweet` data source.

## Example Usage

```hcl
data "x-api-scraper_tweet" "example" {
  tweet_id = "20"
}
```

## Schema

### Required

- `tweet_id` (String)

### Read-Only

- `text` (String)
- `created_at` (String)
- `screen_name` (String)
- `user_id` (String)
- `favorite_count` (Number)
- `retweet_count` (Number)
- `reply_count` (Number)
- `data_json` (String) Raw tweet payload as JSON.
