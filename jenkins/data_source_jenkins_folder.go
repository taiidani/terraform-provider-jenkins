package jenkins

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type folderDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Folder      types.String `tfsdk:"folder"`
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
	Security    types.Object `tfsdk:"security"`
	Template    types.String `tfsdk:"display_name"`
}

type folderDataSource struct {
	*dataSourceHelper
}

// Ensure the implementation satisfies the desired interfaces.
var _ datasource.DataSourceWithConfigure = &folderDataSource{}

func newFolderDataSource() datasource.DataSource {
	return &folderDataSource{
		dataSourceHelper: newDataSourceHelper(),
	}
}

func (d *folderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

func (d *folderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Get the attributes of a folder within Jenkins.

~> The Jenkins installation that uses this resource is expected to have the [Cloudbees Folders Plugin](https://plugins.jenkins.io/cloudbees-folder) installed in their system.`,
		Attributes: d.schemaCredential(map[string]schema.Attribute{
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The name of the folder to display in the UI.",
				Computed:            true,
			},
			"security": schema.SingleNestedAttribute{
				MarkdownDescription: "The Jenkins project-based security configuration.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"inheritance_strategy": schema.StringAttribute{
						MarkdownDescription: "The strategy for applying these permissions sets to existing inherited permissions.",
						Computed:            true,
					},
					"permissions": schema.ListAttribute{
						MarkdownDescription: "The Jenkins permissions sets that provide access to this folder.",
						Computed:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"template": schema.StringAttribute{
				MarkdownDescription: "The configuration file template, used to communicate with Jenkins.",
				Computed:            true,
			},
		}),
	}
}

func (d *folderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data folderDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name, folders := parseCanonicalJobID(data.ID.ValueString())
	job, err := d.client.GetJob(ctx, name, folders...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Data Source",
			"An unexpected error occurred while parsing the data source read response. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}
	data.ID = types.StringValue(job.Base)

	// Extract the raw XML configuration
	config, err := job.GetConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Data Source",
			"An unexpected error occurred while extracting the job configuration. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}
	data.Template = types.StringValue(config)

	// Next, parse the properties from the config
	f, err := parseFolder(config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Folder",
			"An unexpected error occurred while parsing the folder configuration. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	data.Name = types.StringValue(name)
	data.DisplayName = types.StringValue(f.DisplayName)
	data.Folder = types.StringValue(formatFolderID(folders))
	data.Description = types.StringValue(f.Description)

	// Convert the security block
	s := folderSecurityModel{}
	s.FromXML(f.Properties.Security)
	sTF, diag := types.ObjectValueFrom(ctx, s.AttributeTypes(), &s)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Security = sTF

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
