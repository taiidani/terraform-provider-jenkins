package jenkins

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type credentialVaultAppRoleDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Domain      types.String `tfsdk:"domain"`
	Folder      types.String `tfsdk:"folder"`
	Scope       types.String `tfsdk:"scope"`
	Description types.String `tfsdk:"description"`
	Namespace   types.String `tfsdk:"namespace"`
	Path        types.String `tfsdk:"path"`
	RoleID      types.String `tfsdk:"role_id"`
}

type credentialVaultAppRoleDataSource struct {
	client *jenkinsAdapter
}

// Ensure the implementation satisfies the desired interfaces.
var _ datasource.DataSourceWithConfigure = &credentialVaultAppRoleDataSource{}

func newCredentialVaultAppRoleDataSource() datasource.DataSource {
	return &credentialVaultAppRoleDataSource{}
}

func (d *credentialVaultAppRoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_vault_approle"
}

// Configure should register the client for the resource.
func (d *credentialVaultAppRoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *credentialVaultAppRoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get the attributes of a vault approle credential within Jenkins.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The full canonical job path, e.g. `/job/job-name`",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the resource being read.",
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The domain store containing this resource.",
				Optional:            true,
			},
			"folder": schema.StringAttribute{
				MarkdownDescription: "The folder namespace containing this resource.",
				Optional:            true,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: `The visibility of the credentials to Jenkins agents. This will be either "GLOBAL" or "SYSTEM".`,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A human readable description of the credentials being stored.",
				Computed:            true,
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The Vault namespace of the approle credential.",
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "The unique name of the approle auth backend.",
				Computed:            true,
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: "The role_id associated with the credentials.",
				Computed:            true,
			},
		},
	}
}

func (d *credentialVaultAppRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data credentialVaultAppRoleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cm := d.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	if data.Domain.IsNull() {
		data.Domain = basetypes.NewStringValue(defaultValueDomain)
	}

	cred := VaultAppRoleCredentials{}
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
	data.Namespace = types.StringValue(cred.Namespace)
	data.Path = types.StringValue(cred.Path)
	data.RoleID = types.StringValue(cred.RoleID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
