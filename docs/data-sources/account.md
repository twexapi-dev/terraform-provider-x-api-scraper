---
page_title: "x-api-scraper_account Data Source"
---

# x-api-scraper_account

Read the current TwexAPI account balance.

## Example Usage

```hcl
data "x-api-scraper_account" "me" {}
```

## Schema

### Read-Only

- `id` (String)
- `code` (Number)
- `message` (String)
- `data_json` (String) Raw balance payload as JSON.
