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

output "user_id" {
  value = data.x-api-scraper_user_about.example.user_id
}
