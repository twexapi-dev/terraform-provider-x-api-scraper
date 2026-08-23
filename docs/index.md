---
page_title: "TwexAPI Provider"
description: Terraform provider for TwexAPI Twitter search, followers, DMs, and X automation.
---

# TwexAPI Provider

Manage TwexAPI reads and writes as Terraform data sources and resources. It wraps the Go SDK. Search and pagination stay on the language SDKs.

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
