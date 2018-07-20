package runscope

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceRunscopeBucket(t *testing.T) {

	teamID := os.Getenv("RUNSCOPE_TEAM_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceRunscopeBucketConfig, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRunscopeBucket("data.runscope_bucket.test"),
				),
			},
		},
	})
}

func testAccDataSourceRunscopeBucket(dataSource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[dataSource]
		a := r.Primary.Attributes

		if a["name"] != "integration-test-bucket" {
			return fmt.Errorf("expected to get 'integration-test-bucket' bucket returned from runscope data resource %v, got %v", dataSource, a["name"])
		}

		return nil
	}
}

const testAccDataSourceRunscopeBucketConfig = `
resource "runscope_bucket" "test" {
	team_uuid = "%s"
	name = "integration-test-bucket"
}
data "runscope_bucket" "test" {
	key = "${runscope_bucket.test.id}"
}
`
