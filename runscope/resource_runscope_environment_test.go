package runscope

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccEnvironment_basic(t *testing.T) {
	teamID := os.Getenv("RUNSCOPE_TEAM_ID")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeEnvrionmentConfigA, teamID, teamID),
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
		bucketID := rs.Primary.Attributes["bucket_id"]
		testID := rs.Primary.Attributes["test_id"]
		if testID != "" {
			err = client.DeleteEnvironment(&runscope.Environment{ID: rs.Primary.ID},
				&runscope.Bucket{Key: bucketID})
		} else {
			err = client.DeleteEnvironment(&runscope.Environment{ID: rs.Primary.ID},
				&runscope.Bucket{Key: bucketID})
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
		bucketID := rs.Primary.Attributes["bucket_id"]
		testID := rs.Primary.Attributes["test_id"]
		if testID != "" {
			foundRecord, err = client.ReadTestEnvironment(environment,
				&runscope.Test{
					ID:     testID,
					Bucket: &runscope.Bucket{Key: bucketID}})
		} else {
			foundRecord, err = client.ReadSharedEnvironment(environment,
				&runscope.Bucket{Key: bucketID})
		}

		if err != nil {
			return err
		}

		if foundRecord.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		if len(foundRecord.Integrations) != 1 {
			return fmt.Errorf("Expected %d integrations, actual %d", 1, len(foundRecord.Integrations))
		}

		if len(foundRecord.Regions) != 2 {
			return fmt.Errorf("Expected %d regions, actual %d", 2, len(foundRecord.Regions))
		}

		if !contains(foundRecord.Regions, "us1") {
			return fmt.Errorf("Expected %s, actual %s", "us1", strings.Join(foundRecord.Regions, ","))
		}

		if !contains(foundRecord.Regions, "eu1") {
			return fmt.Errorf("Expected %s, actual %s", "eu1", strings.Join(foundRecord.Regions, ","))
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
