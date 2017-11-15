package runscope

import (
	"fmt"
	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceRunscopeIntegration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRunscopeIntegrationRead,

		Schema: map[string]*schema.Schema{
			"team_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filter": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRunscopeIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	log.Printf("[INFO] Reading Runscope integration")

	searchType := d.Get("type").(string)
	filters, filtersOk := d.GetOk("filter")

	resp, err := client.ListIntegrations(d.Get("team_uuid").(string))
	if err != nil {
		return err
	}

	found := &runscope.Integration{}
	for _, integration := range resp {
		if integration.IntegrationType == searchType {
			if filtersOk {
				if !integrationFiltersTest(integration, filters.(*schema.Set)) {
					continue
				}
			}
			found = integration
			break
		}
	}

	if found == nil {
		return fmt.Errorf("Unable to locate any integrations with the type: %s", searchType)
	}

	d.SetId(found.ID)
	d.Set("id", found.ID)
	d.Set("type", found.IntegrationType)
	d.Set("description", found.Description)

	return nil
}

func integrationFiltersTest(integration *runscope.Integration, filters *schema.Set) bool {
	for _, v := range filters.List() {
		m := v.(map[string]interface{})
		passed := false

		for _, e := range m["values"].([]interface{}) {
			switch m["name"].(string) {
			case "id":
				if integration.ID == e {
					passed = true
				}
			case "type":
				if integration.IntegrationType == e {
					passed = true
				}
			default:
				if integration.Description == e {
					passed = true
				}
			}
		}

		if passed {
			continue
		} else {
			return false
		}

	}
	return true
}
