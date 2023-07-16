package jenkins

import (
	"context"
	"fmt"

	jenkins "github.com/bndr/gojenkins"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type credentialUsernameDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Domain      types.String `tfsdk:"domain"`
	Folder      types.String `tfsdk:"folder"`
	Scope       types.String `tfsdk:"scope"`
	Description types.String `tfsdk:"description"`
	Username    types.String `tfsdk:"username"`
}

type credentialUsernameDataSource struct {
	client *jenkinsAdapter
}

// Ensure the implementation satisfies the desired interfaces.
var _ datasource.DataSourceWithConfigure = &credentialUsernameDataSource{}

func newCredentialUsernameDataSource() datasource.DataSource {
	return &credentialUsernameDataSource{}
}

func (d *credentialUsernameDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_username"
}

// Configure should register the client for the resource.
func (d *credentialUsernameDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jenkinsAdapter)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *jenkinsAdapter, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *credentialUsernameDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Service generated identifier.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The identifier assigned to the credentials.",
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The domain namespace that the credentials will be added to.",
				Optional:            true,
			},
			"folder": schema.StringAttribute{
				MarkdownDescription: "The folder namespace that the credentials will be added to.",
				Optional:            true,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "The Jenkins scope assigned to the credentials.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The credentials descriptive text.",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The credentials user username.",
				Computed:            true,
			},
		},
	}
}

func (d *credentialUsernameDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data credentialUsernameDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cm := d.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	cred := jenkins.UsernameCredentials{}
	err := cm.GetSingle(ctx, data.Domain.ValueString(), data.Name.ValueString(), &cred)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Data Source",
			"An unexpected error occurred while parsing the data source read response. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)

		return
	}

	data.ID = types.StringValue(generateCredentialID(data.Folder.ValueString(), cred.ID))
	data.Scope = types.StringValue(cred.Scope)
	data.Description = types.StringValue(cred.Description)
	data.Username = types.StringValue(cred.Username)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
