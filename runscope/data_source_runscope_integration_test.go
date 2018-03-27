package runscope

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
)

func TestAccDataSourceRunscopeIntegration_Basic(t *testing.T) {

	teamId := os.Getenv("RUNSCOPE_TEAM_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceRunscopeIntegrationConfig, teamId),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRunscopeIntegration("data.runscope_integration.by_type"),
				),
			},
		},
	})
}

func testAccDataSourceRunscopeIntegration(dataSource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[dataSource]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get an integration ID from runscope data resource")
		}

		if a["type"] != "slack" {
			return fmt.Errorf("Expected to get an integration type slack from runscope data resource")
		}

		if a["description"] == "" {
			return fmt.Errorf("Expected to get an integration description from runscope data resource")
		}

		return nil
	}
}

const testAccDataSourceRunscopeIntegrationConfig = `
data "runscope_integration" "by_type" {
	team_uuid = "%s"
	type      = "slack"
}
`

func TestAccDataSourceRunscopeIntegration_Filter(t *testing.T) {

	teamId := os.Getenv("RUNSCOPE_TEAM_ID")
	integrationDesc := os.Getenv("RUNSCOPE_INTEGRATION_DESC")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceRunscopeIntegrationFilterConfig, teamId, integrationDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRunscopeIntegrationFilter(t, "data.runscope_integration.by_type"),
				),
			},
		},
	})
}

func testAccDataSourceRunscopeIntegrationFilter(t *testing.T, dataSource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		integrationDesc := os.Getenv("RUNSCOPE_INTEGRATION_DESC")

		r := s.RootModule().Resources[dataSource]
		if r == nil {
			t.Fatalf("Integration not found matching '%s'. Please check you have an integration with that description.\n", integrationDesc)
		}

		a := r.Primary.Attributes

		if a["description"] != integrationDesc {
			return fmt.Errorf("Expected integration description %s to be %s", a["description"], integrationDesc)
		}

		return nil
	}
}

const testAccDataSourceRunscopeIntegrationFilterConfig = `
data "runscope_integration" "by_type" {
	team_uuid = "%s"
	type      = "slack"
	filter = {
		name = "description"
		values = ["%s","other test description"]
	}
}
`
