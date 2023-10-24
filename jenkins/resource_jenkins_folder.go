package jenkins

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type folderResourceModel struct {
	ID          types.String         `tfsdk:"id"`
	Name        types.String         `tfsdk:"name"`
	Folder      types.String         `tfsdk:"folder"`
	Description types.String         `tfsdk:"description"`
	DisplayName types.String         `tfsdk:"display_name"`
	Security    *folderSecurityModel `tfsdk:"security"`
	Template    types.String         `tfsdk:"template"`
}

type folderResource struct {
	*resourceHelper
}

// Ensure the implementation satisfies the desired interfaces.
var _ resource.ResourceWithConfigure = &folderResource{}

func newFolderResource() resource.Resource {
	return &folderResource{
		resourceHelper: newResourceHelper(),
	}
}

// Metadata should return the full name of the resource.
func (r *folderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

// Schema should return the schema for this resource.
func (r *folderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Manages a username credential within Jenkins. This username may then be referenced within jobs that are created.

~> The "password" property may leave plain-text passwords in your state file. If using the property to manage the password in Terraform, ensure that your state file is properly secured and encrypted at rest.`,
		Attributes: r.schema(map[string]schema.Attribute{
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The name of the folder to display in the UI.",
				Optional:            true,
			},
			"template": schema.StringAttribute{
				MarkdownDescription: "The configuration file template, used to communicate with Jenkins.",
				Computed:            true,
			},
		}),
		Blocks: map[string]schema.Block{
			"security": schema.SingleNestedBlock{
				MarkdownDescription: "The Jenkins project-based security configuration.",
				Attributes: map[string]schema.Attribute{
					"inheritance_strategy": schema.StringAttribute{
						MarkdownDescription: "The strategy for applying these permissions sets to existing inherited permissions.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("org.jenkinsci.plugins.matrixauth.inheritance.InheritParentStrategy"),
					},
					"permissions": schema.ListAttribute{
						MarkdownDescription: "The Jenkins permissions sets that provide access to this folder.",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
		},
	}
}

// Create is called when the provider must create a new resource. Config
// and planned state values should be read from the
// CreateRequest and new state values set on the CreateResponse.
func (r *folderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data folderResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that the folder exists
	if err := folderExists(ctx, r.client, data.Folder.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Folder",
			fmt.Sprintf("An invalid folder name %q was specified. ", data.Folder.ValueString())+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	f := folder{
		Description: data.Description.ValueString(),
		DisplayName: data.DisplayName.ValueString(),
	}

	// Set up the security block
	// s := folderSecurityModel{}
	// resp.Diagnostics.Append(data.Security.As(ctx, &s, basetypes.ObjectAsOptions{
	// 	UnhandledUnknownAsEmpty: true,
	// })...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	f.Properties.Security = data.Security.ToXML() // securityToXML(data.Security)

	xml, err := f.Render()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Render Folder Template",
			"An unexpected error occurred while rendering the folder XML. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	folders := extractFolders(data.Folder.ValueString())
	_, err = r.client.CreateJobInFolder(ctx, string(xml), data.Name.ValueString(), folders...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while creating the resource. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	// Extract the raw XML configuration
	data.Template = types.StringValue("")

	// Convert from the API data model to the Terraform data model
	// and set any unknown attribute values.
	data.ID = types.StringValue(formatFolderName(data.Folder.ValueString() + "/" + data.Name.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read is called when the provider must read resource values in order
// to update state. Planned state values should be read from the
// ReadRequest and new state values set on the ReadResponse.
func (r *folderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data folderResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name, folders := parseCanonicalJobID(data.ID.ValueString())
	job, err := r.client.GetJob(ctx, name, folders...)
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			// Job does not exist
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to Refresh Resource",
			"An unexpected error occurred while parsing the resource read response. "+
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
			"Unable to Read Resource",
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
	// s := folderSecurityModel{}
	// s.FromXML(f.Properties.Security)
	// sTF, diag := types.ObjectValueFrom(ctx, s.AttributeTypes(), &s)
	// resp.Diagnostics.Append(diag...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// resp.Diagnostics.AddAttributeWarning(path.Root("security")).AddWarning("WTF", fmt.Sprintf("%#v", sTF))
	data.Security.FromXML(f.Properties.Security)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update is called to update the state of the resource. Config, planned
// state, and prior state values should be read from the
// UpdateRequest and new state values set on the UpdateResponse.
func (r *folderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data folderResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name, folders := parseCanonicalJobID(data.ID.ValueString())
	job, err := r.client.GetJob(ctx, name, folders...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Find Resource",
			"An unexpected error occurred while attempting to update the resource. "+
				"Please retry the operation or report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	// Extract the raw XML configuration
	config, err := job.GetConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Resource",
			"An unexpected error occurred while extracting the job configuration. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

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

	// Then update the values
	f.Description = data.Description.ValueString()
	f.DisplayName = data.DisplayName.ValueString()

	// s := folderSecurityModel{}
	// resp.Diagnostics.Append(data.Security.As(ctx, &s, basetypes.ObjectAsOptions{
	// 	UnhandledUnknownAsEmpty: true,
	// })...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	f.Properties.Security = data.Security.ToXML()

	// And send it back to Jenkins
	xml, err := f.Render()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Bind Resource XML",
			"An unexpected error occurred while extracting the job configuration. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	err = job.UpdateConfig(ctx, string(xml))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Resource",
			"An unexpected error occurred while extracting the job configuration. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete is called when the provider must delete the resource. Config
// values may be read from the DeleteRequest.
//
// If execution completes without error, the framework will automatically
// call DeleteResponse.State.RemoveResource(), so it can be omitted
// from provider logic.
func (r *folderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data folderResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	name, folders := parseCanonicalJobID(data.ID.ValueString())
	_, err := r.client.DeleteJobInFolder(ctx, name, folders...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Resource",
			"An unexpected error occurred while deleting the resource. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)

		return
	}
}

type folderSecurityModel struct {
	InheritanceStrategy types.String `tfsdk:"inheritance_strategy"`
	Permissions         types.List   `tfsdk:"permissions"`
}

func (f *folderSecurityModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"inheritance_strategy": types.StringType,
		"permissions":          types.ListType{ElemType: types.StringType},
	}
}

func (f *folderSecurityModel) FromXML(config *folderSecurity) {
	if config == nil {
		return
	}

	f.InheritanceStrategy = types.StringValue(config.InheritanceStrategy.Class)

	values := []attr.Value{}
	for _, p := range config.Permission {
		values = append(values, types.StringValue(p))
	}
	f.Permissions = types.ListValueMust(types.StringType, values)
}

func (f *folderSecurityModel) ToXML() *folderSecurity {
	ret := &folderSecurity{}
	ret.InheritanceStrategy = folderPermissionInheritanceStrategy{
		Class: f.InheritanceStrategy.ValueString(),
	}
	ret.Permission = []string{}
	for _, permission := range f.Permissions.Elements() {
		ret.Permission = append(ret.Permission, permission.String())
	}
	return ret
}

// func securityToXML(data types.Object) folderSecurity {

// }
