---
page_title: "TwexAPI Provider"
description: Use Terraform to read X profiles and tweets, follow accounts, send DMs, and post, like, retweet, or bookmark.
---

# TwexAPI Provider

Use Terraform to read X profiles and tweets, follow accounts, send DMs, and post, like, retweet, or bookmark. Search and pagination stay on the language SDKs.

Get an API key at [twexapi.io](https://twexapi.io) ([dashboard](https://twexapi.io/dashboard)). REST API: [docs.twexapi.io](https://docs.twexapi.io).

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
