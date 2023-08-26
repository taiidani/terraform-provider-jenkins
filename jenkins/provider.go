package jenkins

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider creates a new Jenkins provider.
//
// Deprecated: Use the provider-framework version of the provider for all new resources.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL of the Jenkins server to connect to. It should be fully qualified (e.g. `https://...`) and point to the root of the Jenkins server location.",
			},
			"ca_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to the Jenkins self-signed certificate. It may be required in order to authenticate to your Jenkins instance.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true, // Needs to be optional to be able to run terraform validate without providing credentials
				Description: "The username to authenticate to Jenkins.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true, // Needs to be optional to be able to run terraform validate without providing credentials
				Description: "The password to authenticate to Jenkins. If you are using the GitHub OAuth authentication method, enter your Personal Access Token here.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"jenkins_credential_vault_approle": dataSourceJenkinsCredentialVaultAppRole(),
			"jenkins_folder":                   dataSourceJenkinsFolder(),
			"jenkins_job":                      dataSourceJenkinsJob(),
			"jenkins_view":                     dataSourceJenkinsView(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"jenkins_credential_secret_file":             resourceJenkinsCredentialSecretFile(),
			"jenkins_credential_secret_text":             resourceJenkinsCredentialSecretText(),
			"jenkins_credential_ssh":                     resourceJenkinsCredentialSSH(),
			"jenkins_credential_vault_approle":           resourceJenkinsCredentialVaultAppRole(),
			"jenkins_folder":                             resourceJenkinsFolder(),
			"jenkins_job":                                resourceJenkinsJob(),
			"jenkins_credential_azure_service_principal": resourceJenkinsCredentialAzureServicePrincipal(),
			"jenkins_view":                               resourceJenkinsView(),
		},

		ConfigureContextFunc: configureProvider,
	}
}

// Deprecated: Use the provider-framework version of the provider for all new resources.
func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	serverURL := os.Getenv("JENKINS_URL")
	if d.Get("server_url").(string) != "" {
		serverURL = d.Get("server_url").(string)
	}
	if serverURL == "" {
		return nil, diag.Errorf("server_url is required and must be provided in the provider config or the JENKINS_URL environment variable")
	}

	caCert := os.Getenv("JENKINS_CA_CERT")
	if d.Get("ca_cert").(string) != "" {
		caCert = d.Get("ca_cert").(string)
	}

	username := os.Getenv("JENKINS_USERNAME")
	if d.Get("username").(string) != "" {
		username = d.Get("username").(string)
	}

	password := os.Getenv("JENKINS_PASSWORD")
	if d.Get("password").(string) != "" {
		password = d.Get("password").(string)
	}

	config := Config{
		ServerURL: serverURL,
		Username:  username,
		Password:  password,
	}

	// Read the certificate
	var err error
	if caCert != "" {
		config.CACert, err = os.Open(caCert)
		if err != nil {
			return nil, diag.Errorf("Unable to open certificate file %s: %s", caCert, err.Error())
		}
	}

	client := newJenkinsClient(&config)
	if _, err = client.Init(ctx); err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
