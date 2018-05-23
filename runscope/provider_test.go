package runscope

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		// provider is called terraform-provider-runscope ie runscope
		"runscope": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}

}

func TestProviderImpl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("RUNSCOPE_ACCESS_TOKEN"); v == "" {
		t.Fatal("RUNSCOPE_ACCESS_TOKEN must be set for acceptance tests")
	}

	if v := os.Getenv("RUNSCOPE_TEAM_ID"); v == "" {
		t.Fatal("RUNSCOPE_TEAM_ID must be set for acceptance tests")
	}

	if v := os.Getenv("RUNSCOPE_INTEGRATION_DESC"); v == "" {
		t.Fatal("RUNSCOPE_INTEGRATION_DESC must be set for acceptance tests")
	}
}

func TestMain(m *testing.M) {

	config := config{
		AccessToken: os.Getenv("RUNSCOPE_ACCESS_TOKEN"),
		APIURL:      "https://api.runscope.com",
	}
	client, err := config.client()

	if err != nil {
		log.Fatalf("Could not create client: %v", err)
		os.Exit(-1)
		return
	}

	shouldDeleteBucket := func(bucket *runscope.Bucket) bool {
		if strings.HasPrefix(bucket.Name, "test") || strings.HasSuffix(bucket.Name, "-test") {
			log.Printf("[DEBUG] deleting bucket %v id: %v", bucket.Name, bucket.Key)
			return true
		}
		return false
	}

	client.DeleteBuckets(shouldDeleteBucket)

	code := m.Run()

	client.DeleteBuckets(shouldDeleteBucket)

	os.Exit(code)
}
