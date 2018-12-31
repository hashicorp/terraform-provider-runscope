package runscope

import (
	"log"

	runscope "github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceRunscopeBucket() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRunscopeBucketRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"team_uuid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRunscopeBucketRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	log.Printf("[INFO] Reading Runscope bucket")

	resp, err := client.ReadBucket(d.Get("key").(string))
	if err != nil {
		return err
	}

	d.SetId(resp.Key)
	d.Set("name", resp.Name)
	d.Set("team_uuid", resp.Team.ID)

	return nil
}
