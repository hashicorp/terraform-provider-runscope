// Note this source file ends in an '_'; otherwise the compiler
// will treat is as a test file.

package runscope

import (
	"fmt"
	"log"
	"strings"

	runscope "github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRunscopeTest() *schema.Resource {
	return &schema.Resource{
		Create: resourceTestCreate,
		Read:   resourceTestRead,
		Update: resourceTestUpdate,
		Delete: resourceTestDelete,

		Schema: map[string]*schema.Schema{
			"bucket_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"default_environment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceRequiresTwoPhase(d *schema.ResourceData) bool {

	if _, ok := d.GetOk("default_environment_id"); ok {
		return true
	}
	return false
}

func resourceTestCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	name := d.Get("name").(string)
	log.Printf("[INFO] Creating test with name: %s", name)

	test, err := createTestFromResourceData(d, true)
	if err != nil {
		return fmt.Errorf("Failed to create test: %s", err)
	}

	log.Printf("[DEBUG] test create: %#v", test)

	createdTest, err := client.CreateTest(test)
	if err != nil {
		return fmt.Errorf("Failed to create test: %s", err)
	}
	if resourceRequiresTwoPhase(d) {
		testUpdate, err := createTestFromResourceData(d, false)
		if err != nil {
			return fmt.Errorf("Failed to create test: %s", err)
		}
		updatedTest, err := client.UpdateTest(testUpdate)
		if err != nil {
			return fmt.Errorf("Failed to create test: %s", err)
		}
		createdTest = updatedTest
	}

	d.SetId(createdTest.ID)
	log.Printf("[INFO] test ID: %s", d.Id())

	return resourceTestRead(d, meta)
}

func resourceTestRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	testFromResource, err := createTestFromResourceData(d, false)
	if err != nil {
		return fmt.Errorf("Error reading test: %s", err)
	}

	test, err := client.ReadTest(testFromResource)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "403") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Couldn't find test: %s", err)
	}

	d.Set("name", test.Name)
	d.Set("description", test.Description)
	d.Set("default_environment_id", test.DefaultEnvironmentID)
	return nil
}

func resourceTestUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)
	testFromResource, err := createTestFromResourceData(d, true)
	if err != nil {
		return fmt.Errorf("Error updating test: %s", err)
	}

	if d.HasChange("description") ||
		d.HasChange("default_environment_id") {
		client := meta.(*runscope.Client)
		_, err = client.UpdateTest(testFromResource)

		if err != nil {
			return fmt.Errorf("Error updating test: %s", err)
		}
	}

	return nil
}

func resourceTestDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	test, err := createTestFromResourceData(d, false)
	if err != nil {
		return fmt.Errorf("Error deleting test: %s", err)
	}
	log.Printf("[INFO] Deleting test with id: %s name: %s", test.ID, test.Name)

	if err := client.DeleteTest(test); err != nil {
		return fmt.Errorf("Error deleting test: %s", err)
	}

	return nil
}

func createTestFromResourceData(d *schema.ResourceData, init bool) (*runscope.Test, error) {

	test := runscope.NewTest()
	test.ID = d.Id()
	if attr, ok := d.GetOk("bucket_id"); ok {
		test.Bucket.Key = attr.(string)
	}

	if attr, ok := d.GetOk("name"); ok {
		test.Name = attr.(string)
	}

	if attr, ok := d.GetOk("description"); ok {
		test.Description = attr.(string)
	}
	// attributes that require create-then-update cycle
	if !init {
		if attr, ok := d.GetOk("default_environment_id"); ok {
			test.DefaultEnvironmentID = attr.(string)
		}
	}

	return test, nil
}
