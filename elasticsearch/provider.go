package elasticsearch

import (
	"context"

	api "github.com/elastic/go-elasticsearch/v7"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ELASTICSEARCH_URL", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ELASTICSEARCH_USER", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ELASTICSEARCH_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"elasticsearch_user":    resourceUser(),
			"elasticsearch_role":    resourceRole(),
			"elasticsearch_api_key": resourceAPIKey(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(context context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {

	url := data.Get("url").(string)
	username := data.Get("username").(string)
	password := data.Get("password").(string)

	var diags diag.Diagnostics
	var client *api.Client
	var err error

	if (url != "") || (username != "") || (password != "") {

		config := api.Config{
			Addresses: []string{
				url,
			},
			Username: username,
			Password: password,
		}
		client, err = api.NewClient(config)

	} else {
		client, err = api.NewDefaultClient()
	}

	if err != nil {
		return nil, diag.FromErr(err)
	}

	response, err := client.Ping()

	if err != nil {
		return nil, diag.FromErr(err)
	}

	if response.IsError() {
		return nil, diag.Errorf("Error: %s", response.String())
	}

	return client, diags
}
