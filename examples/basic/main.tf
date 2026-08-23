terraform {
  required_providers {
    x-api-scraper = {
      source = "twexapi-dev/x-api-scraper"
    }
  }
}

provider "x-api-scraper" {}

data "x-api-scraper_user_about" "example" {
  screen_name = "elonmusk"
}

data "x-api-scraper_tweet" "example" {
  tweet_id = "20"
}

output "user_id" {
  value = data.x-api-scraper_user_about.example.user_id
}

output "tweet_text" {
  value = data.x-api-scraper_tweet.example.text
}
