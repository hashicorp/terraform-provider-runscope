package runscope

import (
	"fmt"
	"os"
	"testing"

	"log"
	"strings"

	runscope "github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBucket_basic(t *testing.T) {
	var bucketResponse runscope.Bucket
	teamID := os.Getenv("RUNSCOPE_TEAM_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testRunscopeBucketConfigA, teamID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists("runscope_bucket.bucket", &bucketResponse),
					resource.TestCheckResourceAttr(
						"runscope_bucket.bucket", "name", "runscope-bucket"),
				),
			},
		},
	})
}

func init() {
	resource.AddTestSweepers("runscope_bucket", &resource.Sweeper{
		Name: "runscope_bucket",
		F: func(region string) error {
			println("[DEBUG] running test sweeper function runscope_bucket")

			config := config{
				AccessToken: os.Getenv("RUNSCOPE_ACCESS_TOKEN"),
				APIURL:      "https://api.runscope.com",
			}
			client, err := config.client()

			if err != nil {
				log.Fatalf("Could not create client: %v", err)
				return err
			}

			shouldDeleteBucket := func(bucket *runscope.Bucket) bool {
				if strings.HasPrefix(bucket.Name, "test") || strings.HasSuffix(bucket.Name, "-test") {
					log.Printf("[DEBUG] deleting bucket %v id: %v", bucket.Name, bucket.Key)
					return true
				}
				return false
			}

			client.DeleteBuckets(shouldDeleteBucket)
			return nil
		},
	})
}

func testAccCheckBucketDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*runscope.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "runscope_bucket" {
			continue
		}

		_, err := client.ReadBucket(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Record %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckBucketExists(n string, bucket *runscope.Bucket) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*runscope.Client)

		foundRecord, err := client.ReadBucket(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundRecord.Key != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		bucket = foundRecord

		return nil
	}
}

const testRunscopeBucketConfigA = `
resource "runscope_bucket" "bucket" {
  name = "runscope-bucket"
  team_uuid = "%s"
}`
