---
page_title: "x-api-scraper_dm Resource"
---

# x-api-scraper_dm

Send a DM via `POST /v3/twitter/send-dm`. Destroy removes Terraform state only; the API cannot unsend.

## Example Usage

```hcl
resource "x-api-scraper_dm" "hello" {
  recipient = "elonmusk"
  text      = "hello from Terraform"
  cookie    = var.twitter_cookie
}
```

## Schema

### Required

- `recipient` (String) User id, `@handle`, or an existing group id (`g...`).
- `text` (String) Message body.
- `cookie` (String, Sensitive) Twitter cookie or `auth_token`.

### Optional

- `media_url` (String) Public image URL (max 1). Mutually exclusive with `video_url`.
- `video_url` (String) Public video URL. Mutually exclusive with `media_url`.

### Read-Only

- `id` (String)
- `message_id` (String)
- `sent_at` (String)
