package runscope

import (
	"fmt"
	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccEnvironment_basic(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeEnvrionmentConfigA, teamId, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists("runscope_environment.environmentA"),
					resource.TestCheckResourceAttr(
						"runscope_environment.environmentA", "name", "test-environment"),
					resource.TestCheckResourceAttr(
						"runscope_environment.environmentA", "verify_ssl", "true")),
			},
		},
	})
}
func TestAccEnvironment_do_not_verify_ssl(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeEnvrionmentConfigB, teamId, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentExists("runscope_environment.environmentB"),
					resource.TestCheckResourceAttr(
						"runscope_environment.environmentB", "name", "test-no-ssl"),
					resource.TestCheckResourceAttr(
						"runscope_environment.environmentB", "verify_ssl", "false")),
			},
		},
	})
}

func testAccCheckEnvironmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*runscope.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "runscope_environment" {
			continue
		}

		var err error
		bucketId := rs.Primary.Attributes["bucket_id"]
		testId := rs.Primary.Attributes["test_id"]
		if testId != "" {
			err = client.DeleteEnvironment(&runscope.Environment{ID: rs.Primary.ID},
				&runscope.Bucket{Key: bucketId})
		} else {
			err = client.DeleteEnvironment(&runscope.Environment{ID: rs.Primary.ID},
				&runscope.Bucket{Key: bucketId})
		}

		if err == nil {
			return fmt.Errorf("Record %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckEnvironmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*runscope.Client)

		var foundRecord *runscope.Environment
		var err error

		environment := new(runscope.Environment)
		environment.ID = rs.Primary.ID
		bucketId := rs.Primary.Attributes["bucket_id"]
		testId := rs.Primary.Attributes["test_id"]
		if testId != "" {
			foundRecord, err = client.ReadTestEnvironment(environment,
				&runscope.Test{
					ID:     testId,
					Bucket: &runscope.Bucket{Key: bucketId}})
		} else {
			foundRecord, err = client.ReadSharedEnvironment(environment,
				&runscope.Bucket{Key: bucketId})
		}

		if err != nil {
			return err
		}

		if foundRecord.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		if len(foundRecord.Integrations) != 1 {
			return fmt.Errorf("Expected %d integrations, actual %d", 1, len(environment.Integrations))
		}

		if len(foundRecord.Regions) != 2 {
			return fmt.Errorf("Expected %d regions, actual %d", 2, len(environment.Regions))
		}

		if foundRecord.Regions[0] != "us1" {
			return fmt.Errorf("Expected %s, actual %s", "us1", environment.Regions[0])
		}

		if foundRecord.Regions[1] != "eu1" {
			return fmt.Errorf("Expected %s, actual %s", "eu1", environment.Regions[1])
		}

		if !foundRecord.RetryOnFailure {
			return fmt.Errorf("Expected retry_on_failure to be set to true")
		}

		return nil
	}
}

const testRunscopeEnvrionmentConfigA = `
resource "runscope_environment" "environmentA" {
  bucket_id    = "${runscope_bucket.bucket.id}"
  name         = "test-environment"

  integrations = [
		"${data.runscope_integration.slack.id}"
  ]

  initial_variables {
    var1 = "true",
    var2 = "value2"
  }

	regions = ["us1", "eu1"]
	
	remote_agents = [
		{
			name = "test agent"
			uuid = "arbitrary-string"
		}
	]

	retry_on_failure = true
}

resource "runscope_test" "test" {
  bucket_id = "${runscope_bucket.bucket.id}"
  name = "runscope test"
  description = "This is a test test..."
}

resource "runscope_bucket" "bucket" {
  name = "terraform-provider-test"
  team_uuid = "%s"
}

data "runscope_integration" "slack" {
  team_uuid = "%s"
  type = "slack"
}
`

const testRunscopeEnvrionmentConfigB = `
resource "runscope_environment" "environmentB" {
  bucket_id    = "${runscope_bucket.bucket.id}"
  name         = "test-no-ssl"

  integrations = [
		"${data.runscope_integration.slack.id}"
  ]

  initial_variables {
    var1 = "true",
    var2 = "value2"
  }

  regions = ["us1", "eu1"]
	
  remote_agents = [
    {
      name = "test agent"
	  uuid = "arbitrary-string"
	}
  ]

  retry_on_failure = true
  verify_ssl = false
}

resource "runscope_test" "test" {
  bucket_id = "${runscope_bucket.bucket.id}"
  name = "runscope test"
  description = "This is a test test..."
}

resource "runscope_bucket" "bucket" {
  name = "terraform-provider-test"
  team_uuid = "%s"
}

data "runscope_integration" "slack" {
  team_uuid = "%s"
  type = "slack"
}
`
