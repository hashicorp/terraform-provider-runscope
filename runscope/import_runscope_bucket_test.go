package runscope

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccImportRunscopeBucket(t *testing.T) {

	teamID := os.Getenv("RUNSCOPE_TEAM_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccImportRunscopeBucketConfig, teamID),
			},
			{
				ResourceName:      "runscope_bucket.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccImportRunscopeBucketConfig = `
resource "runscope_bucket" "test" {
	team_uuid = "%s"
	name = "integration-test-bucket"
}
`

