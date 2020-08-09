package jenkins

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider creates a new Jenkins provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JENKINS_URL", nil),
				Description: "The URL of the Jenkins server to connect to.",
			},
			"ca_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JENKINS_CA_CERT", nil),
				Description: "The path to the Jenkins self-signed certificate.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JENKINS_USERNAME", nil),
				Description: "Username to authenticate to Jenkins.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JENKINS_PASSWORD", nil),
				Description: "Password to authenticate to Jenkins.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"jenkins_folder": resourceJenkinsFolder(),
			"jenkins_job":    resourceJenkinsJob(),
		},

		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		ServerURL: d.Get("server_url").(string),
		CACert:    d.Get("ca_cert").(string),
		Username:  d.Get("username").(string),
		Password:  d.Get("password").(string),
	}

	client, err := newJenkinsClient(&config)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
