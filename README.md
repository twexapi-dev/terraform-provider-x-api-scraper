# TwexAPI Terraform Provider: read X profiles and tweets

Use the TwexAPI Terraform provider to look up X profiles, tweets, and TwexAPI account balance with an API key.
It wraps the [Go SDK](https://github.com/twexapi-dev/x-api-scraper-go). Search, timelines, communities, and other reads also stay available on the REST API or language SDKs.

Public getting-started docs focus on API-key read data sources. Account-session write resources are not covered here.

[REST API](https://docs.twexapi.io) | [Go SDK](https://github.com/twexapi-dev/x-api-scraper-go) | [TypeScript SDK](https://github.com/twexapi-dev/x-api-scraper-typescript) | [Dashboard](https://twexapi.io/dashboard)

This provider is a thin Terraform Plugin Framework wrapper over the Go SDK.

## Common Twitter & X tasks

| Task                            | REST Route                          | Terraform                                          |
| ------------------------------- | ----------------------------------- | -------------------------------------------------- |
| Read an X profile               | `GET /twitter/{screen_name}/about`  | `data.x-api-scraper_user_about`                    |
| Read a tweet                    | `POST /twitter/tweets/lookup`       | `data.x-api-scraper_tweet`                         |
| Read TwexAPI account balance    | `GET /balance`                      | `data.x-api-scraper_account`                       |
| Search tweets without the X API | `POST /twitter/advanced_search/page` | Use the [Go SDK](https://github.com/twexapi-dev/x-api-scraper-go) |

## Package & registry trust

- Source: `twexapi-dev/x-api-scraper`
- Repository: [twexapi-dev/terraform-provider-x-api-scraper](https://github.com/twexapi-dev/terraform-provider-x-api-scraper)
- Go client: [github.com/twexapi-dev/x-api-scraper-go](https://github.com/twexapi-dev/x-api-scraper-go)
- Docs: [docs.twexapi.io](https://docs.twexapi.io)
- License: MIT
- Dashboard: [twexapi.io/dashboard](https://twexapi.io/dashboard)
- Release signing key: [`gpg-public.asc`](./gpg-public.asc)

## Install

Terraform CLI 1.0 or newer.

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

The Terraform Registry listing is not published yet. Until then, build locally and use a `dev_overrides` block in `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "twexapi-dev/x-api-scraper" = "/absolute/path/to/terraform-provider-x-api-scraper"
  }
  direct {}
}
```

Then:

```bash
go build -o terraform-provider-x-api-scraper
export X_API_SCRAPER_KEY="your-api-key"
terraform plan
```

## Authenticate

Prefer environment variables:

```bash
export X_API_SCRAPER_KEY="your-api-key"
```

Or set `bearer_auth` on the provider.

Never commit credentials or Terraform state.

## Read a profile

```hcl
data "x-api-scraper_user_about" "elon" {
  screen_name = "elonmusk"
}

output "user_id" {
  value = data.x-api-scraper_user_about.elon.user_id
}
```

## Read a tweet

```hcl
data "x-api-scraper_tweet" "example" {
  tweet_id = "20"
}

output "tweet_text" {
  value = data.x-api-scraper_tweet.example.text
}
```

## Read account balance

```hcl
data "x-api-scraper_account" "me" {}
```

Not affiliated with X Corp.
