package jenkins

import (
	"context"
	"fmt"
	"strings"

	jenkins "github.com/bndr/gojenkins"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type credentialUsernameResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Domain      types.String `tfsdk:"domain"`
	Folder      types.String `tfsdk:"folder"`
	Scope       types.String `tfsdk:"scope"`
	Description types.String `tfsdk:"description"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
}

type credentialUsernameResource struct {
	client *jenkinsAdapter
}

var _ resource.ResourceWithConfigure = &credentialUsernameResource{}

func newCredentialUsernameResource() resource.Resource {
	return &credentialUsernameResource{}
}

// Metadata should return the full name of the resource.
func (r *credentialUsernameResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_username"
}

// Configure should register the client for the resource.
func (r *credentialUsernameResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Schema should return the schema for this resource.
func (r *credentialUsernameResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Username credentials",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Service generated identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The identifier assigned to the credentials.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The domain namespace that the credentials will be added to.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("_"),
				PlanModifiers: []planmodifier.String{
					// In-place updates should be possible, but gojenkins does not support move operations
					stringplanmodifier.RequiresReplace(),
				},
			},
			"folder": schema.StringAttribute{
				MarkdownDescription: "The folder namespace that the credentials will be added to.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "The Jenkins scope assigned to the credentials.",
				Computed:            true,
				Default:             stringdefault.StaticString("GLOBAL"),
				Validators: []validator.String{
					stringvalidator.OneOf(supportedCredentialScopes...),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The credentials descriptive text.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Managed by Terraform"),
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The credentials user username.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The credentials user password. If left empty will be unmanaged.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

// Create is called when the provider must create a new resource. Config
// and planned state values should be read from the
// CreateRequest and new state values set on the CreateResponse.
func (r *credentialUsernameResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data credentialUsernameResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cm := r.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	// Validate that the folder exists
	if err := folderExists(ctx, r.client, cm.Folder); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Folder",
			fmt.Sprintf("An invalid folder name %q was specified. ", cm.Folder)+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)

		return
	}

	cred := jenkins.UsernameCredentials{
		ID:          data.Name.ValueString(),
		Scope:       data.Scope.ValueString(),
		Description: data.Description.ValueString(),
		Username:    data.Username.ValueString(),
		Password:    data.Password.ValueString(),
	}

	err := cm.Add(ctx, data.Domain.ValueString(), cred)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while creating the resource. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)

		return
	}

	// Convert from the API data model to the Terraform data model
	// and set any unknown attribute values.
	data.ID = types.StringValue(generateCredentialID(data.Folder.ValueString(), cred.ID))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read is called when the provider must read resource values in order
// to update state. Planned state values should be read from the
// ReadRequest and new state values set on the ReadResponse.
func (r *credentialUsernameResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data credentialUsernameResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cm := r.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	cred := jenkins.UsernameCredentials{}
	err := cm.GetSingle(ctx, data.Domain.ValueString(), data.Name.ValueString(), &cred)
	if err != nil {
		if strings.HasSuffix(err.Error(), "404") {
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

	data.ID = types.StringValue(generateCredentialID(data.Folder.ValueString(), cred.ID))
	data.Scope = types.StringValue(cred.Scope)
	data.Description = types.StringValue(cred.Description)
	data.Username = types.StringValue(cred.Username)
	// NOTE: We are NOT setting the password here, as the password returned by GetSingle is garbage
	// Password only applies to Create/Update operations if the "password" property is non-empty

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update is called to update the state of the resource. Config, planned
// state, and prior state values should be read from the
// UpdateRequest and new state values set on the UpdateResponse.
func (r *credentialUsernameResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data credentialUsernameResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cm := r.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	cred := jenkins.UsernameCredentials{
		ID:          data.Name.ValueString(),
		Scope:       data.Scope.ValueString(),
		Description: data.Description.ValueString(),
		Username:    data.Username.ValueString(),
	}

	// Only enforce the password if it is non-empty
	if data.Password.ValueString() != "" {
		cred.Password = data.Password.ValueString()
	}

	err := cm.Update(ctx, data.Domain.ValueString(), data.Name.ValueString(), &cred)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Resource",
			"An unexpected error occurred while attempting to update the resource. "+
				"Please retry the operation or report this issue to the provider developers.\n\n"+
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
func (r *credentialUsernameResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data credentialUsernameResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	cm := r.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	err := cm.Delete(ctx, data.Domain.ValueString(), data.Name.ValueString())
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

// ImportState is called when performing import operations of existing resources.
func (r *credentialUsernameResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splitID := strings.Split(req.ID, "/")
	if len(splitID) < 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: \"[<folder>/]<domain>/<name>\". Got: %q", req.ID),
		)
		return
	}

	name := splitID[len(splitID)-1]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)

	domain := splitID[len(splitID)-2]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), domain)...)

	folder := strings.Trim(strings.Join(splitID[0:len(splitID)-2], "/"), "/")
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("folder"), folder)...)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), generateCredentialID(folder, name))...)
}
