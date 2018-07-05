package runscope

import (
	"log"
	"time"

	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceRunscopeBuckets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRunscopeBucketsRead,

		Schema: map[string]*schema.Schema{
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
			"keys": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceRunscopeBucketsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	log.Printf("[INFO] Reading Runscope buckets")

	filters, filtersOk := d.GetOk("filter")

	resp, err := client.ListBuckets()
	if err != nil {
		return err
	}

	var keys []string
	for _, bucket := range resp {
		if filtersOk {
			if !bucketFiltersTest(bucket, filters.(*schema.Set)) {
				continue
			}
		}

		keys = append(keys, bucket.Key)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("keys", keys)

	return nil
}

func bucketFiltersTest(bucket *runscope.Bucket, filters *schema.Set) bool {
	for _, v := range filters.List() {
		m := v.(map[string]interface{})
		passed := false

		for _, e := range m["values"].(*schema.Set).List() {
			switch m["name"].(string) {
			case "key":
				if bucket.Key == e {
					passed = true
				}
			default:
				if bucket.Name == e {
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
