package runscope

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUNSCOPE_ACCESS_TOKEN", nil),
				Description: "A runscope access token.",
			},
			"api_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RUNSCOPE_API_URL", nil),
				Description: "A runscope api url i.e. https://api.runscope.com.",
				Default:     "https://api.runscope.com",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"runscope_integration":  dataSourceRunscopeIntegration(),
			"runscope_integrations": dataSourceRunscopeIntegrations(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"runscope_bucket":      resourceRunscopeBucket(),
			"runscope_test":        resourceRunscopeTest(),
			"runscope_environment": resourceRunscopeEnvironment(),
			"runscope_schedule":    resourceRunscopeSchedule(),
			"runscope_step":        resourceRunscopeStep(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := config{
		AccessToken: d.Get("access_token").(string),
		APIURL:      d.Get("api_url").(string),
	}
	return config.client()
}
