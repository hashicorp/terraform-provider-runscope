package runscope

import (
	"fmt"
	"log"
	"strings"

	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRunscopeEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Update: resourceEnvironmentUpdate,
		Delete: resourceEnvironmentDelete,

		Schema: map[string]*schema.Schema{
			"bucket_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"test_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"script": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"preserve_cookies": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"initial_variables": &schema.Schema{
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: false,
			},
			"integrations": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"regions": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"remote_agents": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional: true,
			},
			"retry_on_failure": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"verify_ssl": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"webhooks": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"emails": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notify_all": &schema.Schema{
							Type:     schema.TypeBool,
							Required: true,
						},
						"notify_on": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"notify_threshold": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"recipients": &schema.Schema{
							Type:  schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"id": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"email": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},

							},
							},
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func resourceEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	name := d.Get("name").(string)
	log.Printf("[INFO] Creating environment with name: %s", name)

	environment, err := createEnvironmentFromResourceData(d)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] environment create: %#v", environment)

	var createdEnvironment *runscope.Environment
	bucketID := d.Get("bucket_id").(string)

	if testID, ok := d.GetOk("test_id"); ok {
		createdEnvironment, err = client.CreateTestEnvironment(environment,
			&runscope.Test{ID: testID.(string), Bucket: &runscope.Bucket{Key: bucketID}})
	} else {
		createdEnvironment, err = client.CreateSharedEnvironment(environment,
			&runscope.Bucket{Key: bucketID})
	}
	if err != nil {
		return fmt.Errorf("Failed to create environment: %s", err)
	}

	d.SetId(createdEnvironment.ID)
	log.Printf("[INFO] environment ID: %s", d.Id())

	return resourceEnvironmentRead(d, meta)
}

func resourceEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	environmentFromResource, err := createEnvironmentFromResourceData(d)
	if err != nil {
		return fmt.Errorf("Failed to read environment from resource data: %s", err)
	}

	var environment *runscope.Environment
	bucketID := d.Get("bucket_id").(string)
	if testID, ok := d.GetOk("test_id"); ok {
		environment, err = client.ReadTestEnvironment(
			environmentFromResource, &runscope.Test{ID: testID.(string), Bucket: &runscope.Bucket{Key: bucketID}})
	} else {
		environment, err = client.ReadSharedEnvironment(
			environmentFromResource, &runscope.Bucket{Key: bucketID})
	}

	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "403") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Couldn't find environment: %s", err)
	}

	d.Set("bucket_id", bucketID)
	d.Set("test_id", d.Get("test_id").(string))
	d.Set("name", environment.Name)
	d.Set("script", environment.Script)
	d.Set("preserve_cookies", environment.PreserveCookies)
	d.Set("initial_variables", environment.InitialVariables)
	d.Set("integrations", readIntegrations(environment.Integrations))
	d.Set("retry_on_failure", environment.RetryOnFailure)
	d.Set("verify_ssl", environment.VerifySsl)
	d.Set("webhooks", environment.WebHooks)
	return nil
}

func resourceEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)
	environment, err := createEnvironmentFromResourceData(d)
	if err != nil {
		return fmt.Errorf("Error updating environment: %s", err)
	}

	if d.HasChange("name") ||
		d.HasChange("script") ||
		d.HasChange("preserve_cookies") ||
		d.HasChange("initial_variables") ||
		d.HasChange("integrations") ||
		d.HasChange("regions") ||
		d.HasChange("remote_agents") ||
		d.HasChange("retry_on_failure") ||
		d.HasChange("verify_ssl") ||
		d.HasChange("webhooks") {
		client := meta.(*runscope.Client)
		bucketID := d.Get("bucket_id").(string)
		if testID, ok := d.GetOk("test_id"); ok {
			_, err = client.UpdateTestEnvironment(
				environment, &runscope.Test{ID: testID.(string), Bucket: &runscope.Bucket{Key: bucketID}})
		} else {
			_, err = client.UpdateSharedEnvironment(
				environment, &runscope.Bucket{Key: bucketID})
		}
		if err != nil {
			return fmt.Errorf("Error updating environment: %s", err)
		}
	}

	return nil
}

func resourceEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	environmentFromResource, err := createEnvironmentFromResourceData(d)
	if err != nil {
		return fmt.Errorf("Failed to read environment from resource data: %s", err)
	}

	bucketID := d.Get("bucket_id").(string)
	if testID, ok := d.GetOk("test_id"); ok {
		log.Printf("[INFO] Deleting test environment with id: %s name: %s, from test %s",
			environmentFromResource.ID, environmentFromResource.Name, testID.(string))
		err = client.DeleteEnvironment(
			environmentFromResource, &runscope.Bucket{Key: bucketID})
	} else {
		log.Printf("[INFO] Deleting shared environment with id: %s name: %s",
			environmentFromResource.ID, environmentFromResource.Name)
		err = client.DeleteEnvironment(
			environmentFromResource, &runscope.Bucket{Key: bucketID})
	}

	if err != nil {
		return fmt.Errorf("Error deleting environment: %s", err)
	}

	return nil
}

func createEnvironmentFromResourceData(d *schema.ResourceData) (*runscope.Environment, error) {

	environment := runscope.NewEnvironment()
	environment.ID = d.Id()

	if attr, ok := d.GetOk("name"); ok {
		environment.Name = attr.(string)
	}

	if attr, ok := d.GetOk("test_id"); ok {
		environment.TestID = attr.(string)
	}

	if attr, ok := d.GetOk("script"); ok {
		environment.Script = attr.(string)
	}

	if attr, ok := d.GetOk("preserve_cookies"); ok {
		environment.PreserveCookies = attr.(bool)
	}

	if attr, ok := d.GetOk("initial_variables"); ok {
		variablesRaw := attr.(map[string]interface{})
		variables := map[string]string{}
		for k, v := range variablesRaw {
			variables[k] = v.(string)
		}

		environment.InitialVariables = variables
	}

	if attr, ok := d.GetOk("integrations"); ok {
		integrations := []*runscope.EnvironmentIntegration{}
		items := attr.(*schema.Set)
		for _, item := range items.List() {
			integration := runscope.EnvironmentIntegration{
				ID: item.(string),
			}

			integrations = append(integrations, &integration)
		}

		environment.Integrations = integrations
	}

	if attr, ok := d.GetOk("regions"); ok {
		regions := []string{}
		items := attr.(*schema.Set)
		for _, x := range items.List() {
			item := x.(string)
			regions = append(regions, item)
		}

		environment.Regions = regions
	}

	if attr, ok := d.GetOk("remote_agents"); ok {
		remoteAgents := []*runscope.LocalMachine{}
		items := attr.(*schema.Set)
		for _, x := range items.List() {
			item := x.(map[string]interface{})
			remoteAgent := runscope.LocalMachine{
				Name: item["name"].(string),
				UUID: item["uuid"].(string),
			}

			remoteAgents = append(remoteAgents, &remoteAgent)
		}

		environment.RemoteAgents = remoteAgents
	}

	if attr, ok := d.GetOk("retry_on_failure"); ok {
		environment.RetryOnFailure = attr.(bool)
	}

	if attr, ok := d.Get("verify_ssl").(bool); ok {
		environment.VerifySsl = attr
	}

	if attr, ok := d.GetOk("webhooks"); ok {
		webhooks := []string{}
		items := attr.(*schema.Set)
		for _, x := range items.List() {
			item := x.(string)
			webhooks = append(webhooks, item)
		}

		environment.WebHooks = webhooks
	}

	if attr, ok := d.GetOk("emails"); ok {
		contacts := []*runscope.Contact{}
		tmp := attr.([]interface{})
		items := tmp[0].(map[string]interface{})
		emailSettings := runscope.EmailSettings{
			NotifyAll: items["notify_all"].(bool),
			NotifyOn:  items["notify_on"].(string),
			NotifyThreshold: items["notify_threshold"].(int),
		}

		for _, x := range items["recipients"].([]interface{}) {
			item := x.(map[string]interface{})
			contact := runscope.Contact{
				Name:     item["name"].(string),
				Email:   item["email"].(string),
				ID: item["id"].(string),
			}

			contacts = append(contacts, &contact)
		}

		emailSettings.Recipients = contacts
		environment.EmailSettings = &emailSettings

	}


	return environment, nil
}

func readIntegrations(integrations []*runscope.EnvironmentIntegration) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(integrations))
	for _, integration := range integrations {

		item := map[string]interface{}{
			"id":               integration.ID,
			"integration_type": integration.IntegrationType,
			"description":      integration.Description,
		}

		result = append(result, item)
	}

	return result
}
