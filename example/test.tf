variable "access_token" {}
variable "team_uuid" {}

provider "runscope" {
  access_token = "${var.access_token}"
}

# Create a bucket
resource "runscope_bucket" "main" {
  name         = "terraform-ftw"
  team_uuid    = "${var.team_uuid}"
}

# Create a test in the bucket
resource "runscope_test" "homepage_test" {
  name         = "homepage"
  description  = "checks example.com homepage is up and running"
  bucket_id    = "${runscope_bucket.main.id}"
}

# Create a test step
resource "runscope_step" "home_page_step" {
  bucket_id      = "${runscope_bucket.main.id}"
  test_id        = "${runscope_test.homepage_test.id}"
  step_type      = "request"
  url            = "http://example.com"
  method         = "GET"
  assertions     = [
    {
      source     = "response_status"
      comparison = "equal_number"
      value      = "200"
    }
  ]
}

# Create a schedule to execute the test
resource "runscope_schedule" "hourly" {
  bucket_id      = "${runscope_bucket.main.id}"
  test_id        = "${runscope_test.homepage_test.id}"
  interval       = "1h"
  note           = "Hourly schedule"
  environment_id = "${runscope_test.homepage_test.default_environment_id}"
}


output "bucket_url"  {
  value = "https://runscope.com/radar/${runscope_bucket.main.id}"
}