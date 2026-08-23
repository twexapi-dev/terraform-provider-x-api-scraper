---
page_title: "x-api-scraper_like Resource"
---

# x-api-scraper_like

Like a tweet. Destroy unlikes it.

## Example Usage

```hcl
resource "x-api-scraper_like" "example" {
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
