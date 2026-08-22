---
page_title: "x-api-scraper_follow Resource"
---

# x-api-scraper_follow

Follow an X account. Destroy unfollows.

## Example Usage

```hcl
resource "x-api-scraper_follow" "example" {
  username = "elonmusk"
  cookie   = var.twitter_cookie
}
```

## Schema

### Required

- `username` (String) Screen name to follow, without `@`.
- `cookie` (String, Sensitive) Twitter cookie or `auth_token`.

### Read-Only

- `id` (String)
