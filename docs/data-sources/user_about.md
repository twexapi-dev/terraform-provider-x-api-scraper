---
page_title: "x-api-scraper_user_about Data Source"
---

# x-api-scraper_user_about

Read an X profile by screen name. Maps to the competitor `x_user` data source.

## Example Usage

```hcl
data "x-api-scraper_user_about" "example" {
  screen_name = "elonmusk"
}
```

## Schema

### Required

- `screen_name` (String) X/Twitter screen name without `@`.

### Read-Only

- `user_id` (String)
- `name` (String)
- `avatar` (String)
- `created_at` (String)
- `is_blue_verified` (Boolean)
- `verification` (Boolean)
- `account_based_in` (String)
- `affiliate_username` (String)
- `username_changes_count` (Number)
