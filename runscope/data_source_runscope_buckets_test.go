package runscope

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceRunscopeBuckets(t *testing.T) {

	teamID := os.Getenv("RUNSCOPE_TEAM_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceRunscopeBucketsConfig, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRunscopeBuckets("data.runscope_buckets.test"),
				),
			},
		},
	})
}

func testAccDataSourceRunscopeBuckets(dataSource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[dataSource]
		a := r.Primary.Attributes

		if a["keys.#"] != "1" {
			return fmt.Errorf("expected to get 1 bucket key returned from runscope data resource %v, got %v", dataSource, a["keys.#"])
		}

		return nil
	}
}

const testAccDataSourceRunscopeBucketsConfig = `
resource "runscope_bucket" "test" {
	team_uuid = "%s"
	name = "integration-test-bucket"
}
data "runscope_buckets" "test" {
	filter = [
		{
			name = "name"
			values = ["${runscope_bucket.test.name}"]
		}
	]
}
`

