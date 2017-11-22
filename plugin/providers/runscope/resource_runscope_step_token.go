package runscope

import (
	"encoding/json"
	"fmt"
	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceRunscopeStepToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceStepTokenCreate,
		Read:   resourceStepTokenRead,
		Update: resourceStepTokenUpdate,
		Delete: resourceStepTokenDelete,
		Schema: map[string]*schema.Schema{
			"bucket_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"test_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"token_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceStepTokenCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	step, bucketId, testId, environmentId, err := createStepTokenFromResourceData(d, client)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] step token create: %#v", step)

	environment, err := client.ReadSharedEnvironment(&runscope.Environment{ID: environmentId}, &runscope.Bucket{Key: bucketId})
	if err != nil {
		return err
	}

	tokenName := d.Get("token_name").(string)

	environment.InitialVariables[tokenName] = fmt.Sprintf("{{%s}}", tokenName)

	bytes, err := json.Marshal(environment)
	if err != nil {
		return err
	}

	step.Body = string(bytes)
	createdStep, err := client.CreateTestStep(step, bucketId, testId)
	if err != nil {
		return fmt.Errorf("Failed to create step token: %s", err)
	}

	d.SetId(createdStep.ID)
	log.Printf("[INFO] step token ID: %s", d.Id())

	return resourceStepTokenRead(d, meta)
}

func resourceStepTokenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	stepFromResource, bucketId, testId, environmentId, err := createStepTokenFromResourceData(d, client)
	if err != nil {
		return fmt.Errorf("Failed to read step token from resource data: %s", err)
	}

	step, err := client.ReadTestStep(stepFromResource, bucketId, testId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Couldn't find step token: %s", err)
	}

	d.Set("bucket_id", bucketId)
	d.Set("test_id", testId)
	d.Set("environment_id", environmentId)
	d.Set("body", step.Body)
	return nil
}

func resourceStepTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)
	client := meta.(*runscope.Client)
	stepFromResource, bucketId, testId, _, err := createStepTokenFromResourceData(d, client)
	if err != nil {
		return fmt.Errorf("Error updating step token: %s", err)
	}

	if d.HasChange("token") {
		_, err = client.UpdateTestStep(stepFromResource, bucketId, testId)

		if err != nil {
			return fmt.Errorf("Error updating step token: %s", err)
		}
	}

	return nil
}

func resourceStepTokenDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	stepFromResource, bucketId, testId, _, err := createStepTokenFromResourceData(d, client)
	if err != nil {
		return fmt.Errorf("Failed to read step token from resource data: %s", err)
	}

	err = client.DeleteTestStep(stepFromResource, bucketId, testId)
	if err != nil {
		return fmt.Errorf("Error deleting step token: %s", err)
	}

	return nil
}

func createStepTokenFromResourceData(d *schema.ResourceData, client *runscope.Client) (*runscope.TestStep, string, string, string, error) {

	step := runscope.NewTestStep()
	bucketId := d.Get("bucket_id").(string)
	testId := d.Get("test_id").(string)
	environmentId := d.Get("environment_id").(string)
	step.ID = d.Id()
	step.StepType = "request"
	step.Body = "environment json goes here"
	step.Method = "PUT"
	step.URL = fmt.Sprintf("%s/buckets/%s/environments/%s", client.APIURL, bucketId, environmentId)
	step.Headers = map[string][]string{
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", client.AccessToken)},
	}
	step.Assertions = []*runscope.Assertion{
		{
			Comparison: "equal_number",
			Value:      200,
			Source:     "response_status",
		},
	}

	return step, bucketId, testId, environmentId, nil
}
