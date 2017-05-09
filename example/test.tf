variable "access_token" {}
variable "team_uuid" {}

provider "runscope" {
  access_token = "${var.access_token}"
}

# Create a bucket
resource "runscope_bucket" "main" {
  name         = "terraform-ftw"
  team_uuid    = "${var.team_uuid}}"
}

# Create a test in the bucket
resource "runscope_test" "api" {
  name         = "api-test"
  description  = "checks the api is up and running"
  bucket_id    = "${runscope_bucket.main.id}"
}