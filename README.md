# TwexAPI Terraform Provider: Twitter API for search, followers, DMs & X automation

Use the TwexAPI Terraform provider to read X profiles and tweets, follow accounts, send DMs, and post, like, retweet, or bookmark.
It wraps the [Go SDK](https://github.com/twexapi-dev/x-api-scraper-go). Search, timelines, communities, and other reads stay on the REST API or language SDKs.

[REST API](https://docs.twexapi.io) | [Go SDK](https://github.com/twexapi-dev/x-api-scraper-go) | [TypeScript SDK](https://github.com/twexapi-dev/x-api-scraper-typescript) | [Dashboard](https://twexapi.io/dashboard)

This provider is a thin Terraform Plugin Framework wrapper over the Go SDK, the same approach as the competitor provider.

## Common Twitter & X tasks

| Task                            | REST Route                          | Terraform                                          |
| ------------------------------- | ----------------------------------- | -------------------------------------------------- |
| Read an X profile               | `GET /twitter/{screen_name}/about`  | `data.x-api-scraper_user_about`                    |
| Read a tweet                    | `POST /twitter/tweets/lookup`       | `data.x-api-scraper_tweet`                         |
| Read TwexAPI account balance    | `GET /balance`                      | `data.x-api-scraper_account`                       |
| Follow an account               | `POST /twitter/user/follow`         | `x-api-scraper_follow`                             |
| Post or reply                   | `POST /twitter/tweets/create`       | `x-api-scraper_tweet`                              |
| Like / retweet / bookmark       | tweet action routes                 | `x-api-scraper_like`, `_retweet`, `_bookmark`      |
| Send a DM                       | `POST /v3/twitter/send-dm`          | `x-api-scraper_dm`                                 |
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

Or set `bearer_auth` on the provider. Write resources also need a Twitter cookie or `auth_token`.

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

## Follow an account

```hcl
resource "x-api-scraper_follow" "example" {
  username = "elonmusk"
  cookie   = var.twitter_cookie
}
```

## Post a tweet

```hcl
resource "x-api-scraper_tweet" "announcement" {
  username      = "example"
  cookie        = var.twitter_cookie
  tweet_content = "Published through the TwexAPI Terraform provider."
}
```

Changing `tweet_content` replaces the resource. Destroy deletes the tweet.

## Like, retweet, or bookmark

```hcl
resource "x-api-scraper_like" "example" {
  tweet_id = data.x-api-scraper_tweet.example.tweet_id
  cookie   = var.twitter_cookie
}
```

Destroy reverses the action.

## Send a DM

```hcl
resource "x-api-scraper_dm" "hello" {
  recipient = "elonmusk"
  text      = "hello from Terraform"
  cookie    = var.twitter_cookie
}
```

Destroy removes Terraform state only. The API cannot unsend a DM.

Not affiliated with X Corp.
