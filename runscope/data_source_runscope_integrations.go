package runscope

import (
	"log"
	"time"

	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceRunscopeIntegrations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRunscopeIntegrationsRead,

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
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"ids": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceRunscopeIntegrationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	log.Printf("[INFO] Reading Runscope integration")

	filters, filtersOk := d.GetOk("filter")

	resp, err := client.ListIntegrations(d.Get("team_uuid").(string))
	if err != nil {
		return err
	}

	var ids []string
	for _, integration := range resp {
		if filtersOk {
			if !integrationFiltersTest(integration, filters.(*schema.Set)) {
				continue
			}
		}

		ids = append(ids, integration.ID)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("ids", ids)

	return nil
}
