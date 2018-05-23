package runscope

import (
	"fmt"
	"log"
	"strings"

	"github.com/ewilde/go-runscope"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRunscopeSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceScheduleCreate,
		Read:   resourceScheduleRead,
		Delete: resourceScheduleDelete,

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
			"interval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"note": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	schedule, bucketID, testID, err := createScheduleFromResourceData(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] schedule create: %#v", schedule)

	createdSchedule, err := client.CreateSchedule(schedule, bucketID, testID)
	if err != nil {
		return fmt.Errorf("Failed to create schedule: %s", err)
	}

	d.SetId(createdSchedule.ID)
	log.Printf("[INFO] schedule ID: %s", d.Id())

	return resourceScheduleRead(d, meta)
}

func resourceScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	scheduleFromResource, bucketID, testID, err := createScheduleFromResourceData(d)
	if err != nil {
		return fmt.Errorf("Failed to read schedule from resource data: %s", err)
	}

	schedule, err := client.ReadSchedule(scheduleFromResource, bucketID, testID)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "403") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Couldn't find schedule: %s", err)
	}

	d.Set("bucket_id", bucketID)
	d.Set("test_id", testID)
	d.Set("environment_id", schedule.EnvironmentID)
	d.Set("interval", schedule.Interval)
	d.Set("note", schedule.Note)
	return nil
}

func resourceScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*runscope.Client)

	scheduleFromResource, bucketID, testID, err := createScheduleFromResourceData(d)
	if err != nil {
		return fmt.Errorf("Failed to read schedule from resource data: %s", err)
	}

	err = client.DeleteSchedule(scheduleFromResource, bucketID, testID)
	if err != nil {
		return fmt.Errorf("Error deleting schedule: %s", err)
	}

	return nil
}

func createScheduleFromResourceData(d *schema.ResourceData) (*runscope.Schedule, string, string, error) {

	schedule := runscope.NewSchedule()
	bucketID := d.Get("bucket_id").(string)
	testID := d.Get("test_id").(string)
	environmentID := d.Get("environment_id").(string)
	interval := d.Get("interval").(string)
	note := ""

	if v, ok := d.GetOk("note"); ok {
		note = v.(string)
	}

	schedule.ID = d.Id()
	schedule.Interval = interval
	schedule.Note = note
	schedule.EnvironmentID = environmentID
	return schedule, bucketID, testID, nil
}
