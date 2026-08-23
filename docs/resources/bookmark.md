---
page_title: "x-api-scraper_bookmark Resource"
---

# x-api-scraper_bookmark

Bookmark a tweet. Destroy removes the bookmark.

## Example Usage

```hcl
resource "x-api-scraper_bookmark" "example" {
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
