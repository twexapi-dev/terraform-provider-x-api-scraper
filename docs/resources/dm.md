---
page_title: "x-api-scraper_dm Resource"
---

# x-api-scraper_dm

Destroy removes Terraform state only; the API cannot unsend.

> **Note:** Public getting-started docs focus on API-key read data sources. This resource requires account session credentials and is not covered in the public getting started guide.

## Schema

### Required

- `recipient` (String) User id, `@handle`, or an existing group id (`g...`).
- `text` (String) Message body.

### Optional

- `media_url` (String) Public image URL (max 1). Mutually exclusive with `video_url`.
- `video_url` (String) Public video URL. Mutually exclusive with `media_url`.

### Read-Only

- `id` (String)
- `message_id` (String)
- `sent_at` (String)
