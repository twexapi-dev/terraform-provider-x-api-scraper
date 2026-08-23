---
page_title: "TwexAPI Provider"
description: Use Terraform to read X profiles and tweets, follow accounts, send DMs, and post, like, retweet, or bookmark.
---

# TwexAPI Provider

Use Terraform to read X profiles and tweets, follow accounts, send DMs, and post, like, retweet, or bookmark. Search and pagination stay on the language SDKs.

## Example Usage

```hcl
terraform {
  required_providers {
    x-api-scraper = {
      source  = "twexapi-dev/x-api-scraper"
      version = "~> 0.1"
    }
  }
}

provider "x-api-scraper" {}
```

## Schema

### Optional

- `bearer_auth` (String, Sensitive) TwexAPI API key. Falls back to `X_API_SCRAPER_KEY`.
