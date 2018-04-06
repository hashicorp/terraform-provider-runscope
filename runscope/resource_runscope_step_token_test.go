package runscope

import (
	"encoding/json"
	"fmt"
	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccStepToken_basic(t *testing.T) {
	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStepTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeStepTokenConfigA, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStepTokenExists("runscope_step_token.token_step"),
					testAccCheckStepTokenBodyContainsEnvironmentWithTokenPlaceHolder("runscope_step_token.token_step", "token"),
					resource.TestCheckResourceAttr(
						"runscope_step_token.token_step", "token_name", "token")),
			},
		},
	})
}

func testAccCheckStepTokenDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*runscope.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "runscope_step_token" {
			continue
		}

		var err error
		bucketId := rs.Primary.Attributes["bucket_id"]
		testId := rs.Primary.Attributes["test_id"]
		err = client.DeleteTestStep(&runscope.TestStep{ID: rs.Primary.ID}, bucketId, testId)

		if err == nil {
			return fmt.Errorf("Record step token %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckStepTokenBodyContainsEnvironmentWithTokenPlaceHolder(n string, tokenName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*runscope.Client)

		var foundRecord *runscope.TestStep
		var err error

		step := new(runscope.TestStep)
		step.ID = rs.Primary.ID
		bucketId := rs.Primary.Attributes["bucket_id"]
		testId := rs.Primary.Attributes["test_id"]

		foundRecord, err = client.ReadTestStep(step, bucketId, testId)
		environmentBody := make(map[string]interface{})
		json.Unmarshal([]byte(foundRecord.Body), &environmentBody)

		if err != nil {
			return err
		}

		initalVariablesMap := environmentBody["initial_variables"].(map[string]interface{})

		if initalVariablesMap["var1"].(string) != "true" {
			return fmt.Errorf("var1 in initial variables should be %v but is %v", "true", initalVariablesMap["var1"])
		}

		if initalVariablesMap["var2"].(string) != "value2" {
			return fmt.Errorf("var1 in initial variables should be %v but is %v", "value2", initalVariablesMap["var2"])
		}

		if initalVariablesMap[tokenName].(string) != fmt.Sprintf("{{%s}}", tokenName) {
			return fmt.Errorf("%v in initial variables should be %v but is %v", tokenName, fmt.Sprintf("{{%s}}", tokenName), initalVariablesMap["var2"])
		}

		return nil

	}
}

func testAccCheckStepTokenExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*runscope.Client)

		var foundRecord *runscope.TestStep
		var err error

		step := new(runscope.TestStep)
		step.ID = rs.Primary.ID
		bucketId := rs.Primary.Attributes["bucket_id"]
		testId := rs.Primary.Attributes["test_id"]

		foundRecord, err = client.ReadTestStep(step, bucketId, testId)

		if err != nil {
			return err
		}

		if foundRecord.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		return nil
	}
}

const testRunscopeStepTokenConfigA = `
resource "runscope_step_token" "token_step" {
  bucket_id      = "${runscope_bucket.bucket.id}"
  test_id        = "${runscope_test.test.id}"
  environment_id = "${runscope_environment.environment.id}"
  token_name     = "token"
}

resource "runscope_environment" "environment" {
  bucket_id    = "${runscope_bucket.bucket.id}"
  name         = "test-environment"


  initial_variables {
    var1 = "true",
    var2 = "value2"
  }

	regions = ["us1", "eu1"]
}

resource "runscope_test" "test" {
  bucket_id   = "${runscope_bucket.bucket.id}"
  name        = "runscope test"
  description = "This is a test test..."
}

resource "runscope_bucket" "bucket" {
  name      = "terraform-provider-test"
  team_uuid = "%s"
}
`
