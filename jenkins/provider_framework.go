package jenkins

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func New() provider.Provider {
	return &JenkinsProvider{}
}

// Ensure the implementation satisfies the provider.Provider interface.
var _ provider.Provider = &JenkinsProvider{}

type JenkinsProvider struct{}

// Metadata satisfies the provider.Provider interface for JenkinsProvider
func (p *JenkinsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jenkins"
}

// Schema satisfies the provider.Provider interface for JenkinsProvider.
func (p *JenkinsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				Optional:    true,
				Description: "The URL of the Jenkins server to connect to. It should be fully qualified (e.g. `https://...`) and point to the root of the Jenkins server location.",
			},
			"ca_cert": schema.StringAttribute{
				Optional:    true,
				Description: "The path to the Jenkins self-signed certificate. It may be required in order to authenticate to your Jenkins instance.",
			},
			"username": schema.StringAttribute{
				Optional:    true, // Needs to be optional to be able to run terraform validate without providing credentials
				Description: "The username to authenticate to Jenkins.",
			},
			"password": schema.StringAttribute{
				Optional:    true, // Needs to be optional to be able to run terraform validate without providing credentials
				Description: "The password to authenticate to Jenkins. If you are using the GitHub OAuth authentication method, enter your Personal Access Token here.",
			},
		},
	}
}

type JenkinsProviderModel struct {
	ServerURL types.String `tfsdk:"server_url"`
	CACert    types.String `tfsdk:"ca_cert"`
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
}

// Configure satisfies the provider.Provider interface for JenkinsProvider.
func (p *JenkinsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data JenkinsProviderModel

	// Read configuration data into model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	serverURL := os.Getenv("JENKINS_URL")
	if data.ServerURL.ValueString() != "" {
		serverURL = data.ServerURL.ValueString()
	}
	if serverURL == "" {
		resp.Diagnostics.AddError(
			"server_url is required",
			"server_url is required and must be provided in the provider config or the JENKINS_URL environment variable",
		)
	}

	caCert := os.Getenv("JENKINS_CA_CERT")
	if data.CACert.ValueString() != "" {
		caCert = data.CACert.ValueString()
	}

	username := os.Getenv("JENKINS_USERNAME")
	if data.Username.ValueString() != "" {
		username = data.Username.ValueString()
	}

	password := os.Getenv("JENKINS_PASSWORD")
	if data.Password.ValueString() != "" {
		password = data.Password.ValueString()
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
			resp.Diagnostics.AddError(
				"Unable to open certificate file",
				fmt.Sprintf("Unable to open certificate file %s: %s", caCert, err.Error()),
			)
		}
	}

	client := newJenkinsClient(&config)
	if _, err = client.Init(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Unable to initialize client",
			err.Error(),
		)
	}
	resp.ResourceData = client
	resp.DataSourceData = client
}

// DataSources satisfies the provider.Provider interface for JenkinsProvider.
func (p *JenkinsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newCredentialUsernameDataSource,
	}
}

// Resources satisfies the provider.Provider interface for JenkinsProvider.
func (p *JenkinsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newCredentialUsernameResource,
	}
}
