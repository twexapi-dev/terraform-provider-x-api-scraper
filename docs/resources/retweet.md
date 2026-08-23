---
page_title: "x-api-scraper_retweet Resource"
---

# x-api-scraper_retweet

Retweet a tweet. Destroy deletes the retweet.

## Example Usage

```hcl
resource "x-api-scraper_retweet" "example" {
  tweet_id = "20"
  cookie   = var.twitter_cookie
}
```

## Schema

### Required

- `tweet_id` (String)
- `cookie` (String, Sensitive)

### Read-Only

- `id` (String)
