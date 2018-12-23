package jenkins

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider creates a new Jenkins provider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JENKINS_URL", nil),
				Description: "The URL of the Jenkins server to connect to.",
			},
			"ca_cert": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JENKINS_CA_CERT", nil),
				Description: "The path to the Jenkins self-signed certificate.",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JENKINS_USERNAME", nil),
				Description: "Username to authenticate to Jenkins.",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JENKINS_PASSWORD", nil),
				Description: "Password to authenticate to Jenkins.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"jenkins_folder": resourceJenkinsFolder(),
			"jenkins_job": resourceJenkinsJob(),
		},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ServerURL: d.Get("server_url").(string),
		CACert:    d.Get("ca_cert").(string),
		Username:  d.Get("username").(string),
		Password:  d.Get("password").(string),
	}

	client, err := newJenkinsClient(&config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
