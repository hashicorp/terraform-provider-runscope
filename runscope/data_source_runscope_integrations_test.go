package runscope

import (
	"fmt"
	"os"
	"testing"

	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceRunscopeIntegrations_Basic(t *testing.T) {

	teamID := os.Getenv("RUNSCOPE_TEAM_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceRunscopeIntegrationsConfig, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRunscopeIntegrations("data.runscope_integrations.by_type"),
				),
			},
		},
	})
}

func testAccDataSourceRunscopeIntegrations(dataSource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[dataSource]
		a := r.Primary.Attributes

		if a["ids.#"] != "2" {
			return fmt.Errorf("expected to get 2 integrations ids returned from runscope data resource %v, got %v", dataSource, a["ids.#"])
		}

		return nil
	}
}

const testAccDataSourceRunscopeIntegrationsConfig = `
data "runscope_integrations" "by_type" {
	team_uuid = "%s"
	filter = {
		name = "type"
		values = ["slack"]
	}
}
`

func TestAccDataSourceRunscopeIntegrations_usage(t *testing.T) {

	teamID := os.Getenv("RUNSCOPE_TEAM_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceRunscopeIntegrationsUsageConfig, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvironmentIntegrations("runscope_environment.environment_with_integrations", true),
					testAccCheckEnvironmentIntegrations("runscope_environment.environment_no_integrations", false),
				),
			},
		},
	})
}

func testAccCheckEnvironmentIntegrations(environment string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[environment]

		if !ok {
			return fmt.Errorf("Not found: %s", environment)
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
		foundRecord, err = client.ReadSharedEnvironment(environment,
			&runscope.Bucket{Key: bucketID})

		if err != nil {
			return err
		}

		if len(foundRecord.Integrations) == 0 && expected {
			return fmt.Errorf("Expected environment to have integrations")
		} else if len(foundRecord.Integrations) != 0 && !expected {
			return fmt.Errorf("Expected environment not to have integrations, but had %d", len(foundRecord.Integrations))
		}

		return nil
	}
}

const testAccDataSourceRunscopeIntegrationsUsageConfig = `
data "runscope_integrations" "slack" {
	team_uuid = "%[1]v"
	filter = {
		name = "type"
		values = ["slack"]
	}
}

data "runscope_integrations" "empty" {
	team_uuid = "%[1]v"
	filter = {
		name = "type"
		values = ["unknown"]
	}
}

resource "runscope_environment" "environment_with_integrations" {
  bucket_id    = "${runscope_bucket.bucket.id}"
  name         = "test-environment-1"

  integrations = ["${data.runscope_integrations.slack.ids}"]
}

resource "runscope_environment" "environment_no_integrations" {
  bucket_id    = "${runscope_bucket.bucket.id}"
  name         = "test-environment-2"

  integrations = ["${data.runscope_integrations.empty.ids}"]
}

resource "runscope_bucket" "bucket" {
	name = "terraform-provider-test"
	team_uuid = "%[1]v"
}

`
