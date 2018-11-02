package runscope

import (
	"fmt"
	"os"
	"testing"

	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTest_basic(t *testing.T) {
	var test runscope.Test
	teamID := os.Getenv("RUNSCOPE_TEAM_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTestDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeTestConfigA, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTestExists("runscope_test.test", &test),
					resource.TestCheckResourceAttr(
						"runscope_test.test", "name", "runscope test"),
					resource.TestCheckResourceAttr(
						"runscope_test.test", "description", "This is a test test..."),
				),
			},
		},
	})
}

func testAccCheckTestDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*runscope.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "runscope_test" {
			continue
		}

		_, err := client.ReadTest(&runscope.Test{ID: rs.Primary.ID, Bucket: &runscope.Bucket{Key: rs.Primary.Attributes["bucket_id"]}})

		if err == nil {
			return fmt.Errorf("Record %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTestExists(n string, test *runscope.Test) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*runscope.Client)

		foundRecord, err := client.ReadTest(&runscope.Test{ID: rs.Primary.ID, Bucket: &runscope.Bucket{Key: rs.Primary.Attributes["bucket_id"]}})

		if err != nil {
			return err
		}

		if foundRecord.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		test = foundRecord

		return nil
	}
}

const testRunscopeTestConfigA = `
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
	webhooks = ["https://example.com"]
}

resource "runscope_test" "test" {
  bucket_id = "${runscope_bucket.bucket.id}"
  name = "runscope test"
  description = "This is a test test..."
  default_environment_id = "${runscope_environment.environmentA.id}"
}

resource "runscope_bucket" "bucket" {
  name = "terraform-provider-test"
  team_uuid = "%s"
}
`
