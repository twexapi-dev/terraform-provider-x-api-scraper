---
page_title: "x-api-scraper_tweet Resource"
---

# x-api-scraper_tweet

Post a tweet. Destroy deletes that tweet.

## Example Usage

```hcl
resource "x-api-scraper_tweet" "announcement" {
  username      = "example"
  cookie        = var.twitter_cookie
  tweet_content = "Published through the TwexAPI Terraform provider."
}
```

## Schema

### Required

- `tweet_content` (String)
- `username` (String) Acting account username, used when deleting.
- `cookie` (String, Sensitive) Twitter cookie or `auth_token`.

### Optional

- `reply_tweet_id` (String)

### Read-Only

- `id` (String)
- `tweet_id` (String)
