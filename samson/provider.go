package samson

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	_samson "github.com/tolgaakyuz/samson-go"
)

// Provider returns a schema.Provider for Samson.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SAMSON_TOKEN", nil),
				Description: "The auth token for the Samson api.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"samson_project": resourceSamsonProjects(),
			"samson_command": resourceSamsonCommands(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := _samson.New(d.Get("token").(string))

	return client, nil
}
