package jenkins

import (
	"context"
	"fmt"
	"strings"

	jenkins "github.com/bndr/gojenkins"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type credentialSSHResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Folder      types.String `tfsdk:"folder"`
	Description types.String `tfsdk:"description"`
	Domain      types.String `tfsdk:"domain"`
	Scope       types.String `tfsdk:"scope"`
	Username    types.String `tfsdk:"username"`
	PrivateKey  types.String `tfsdk:"privatekey"`
	Passphrase  types.String `tfsdk:"passphrase"`
}

type credentialSSHResource struct {
	*resourceHelper
}

// Ensure the implementation satisfies the desired interfaces.
var _ resource.ResourceWithConfigure = &credentialSSHResource{}

func newCredentialSSHResource() resource.Resource {
	return &credentialSSHResource{
		resourceHelper: newResourceHelper(),
	}
}

// Metadata should return the full name of the resource.
func (r *credentialSSHResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_ssh"
}

// Schema should return the schema for this resource.
func (r *credentialSSHResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Manages a SSH credential within Jenkins. This SSH credential may then be referenced within jobs that are created.

~> The "passphrase" and "privatekey" properties may leave plain-text values in your state file. Ensure that your state file is properly secured and encrypted at rest.`,
		Attributes: r.schemaCredential(map[string]schema.Attribute{
			"username": schema.StringAttribute{
				MarkdownDescription: "Username",
				Required:            true,
			},
			"privatekey": schema.StringAttribute{
				MarkdownDescription: "Private SSH key, can be given as string or read from file with 'file()' terraform function.",
				Required:            true,
				Sensitive:           true,
			},
			"passphrase": schema.StringAttribute{
				MarkdownDescription: "Passphrase for privatekey. This has to be skipped if private key was created without passphrase.",
				Optional:            true,
				Sensitive:           true,
			},
		}),
	}
}

// Create is called when the provider must create a new resource. Config
// and planned state values should be read from the
// CreateRequest and new state values set on the CreateResponse.
func (r *credentialSSHResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data credentialSSHResourceModel

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

	cred := jenkins.SSHCredentials{
		ID:          data.Name.ValueString(),
		Scope:       data.Scope.ValueString(),
		Description: data.Description.ValueString(),
		Username:    data.Username.ValueString(),
		PrivateKeySource: &jenkins.PrivateKey{
			Class: jenkins.KeySourceDirectEntryType,
			Value: data.PrivateKey.ValueString(),
		},
	}

	passphrase := data.Passphrase.ValueString()
	if len(passphrase) > 0 {
		cred.Passphrase = passphrase
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
func (r *credentialSSHResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data credentialSSHResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cm := r.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	cred := jenkins.SSHCredentials{}
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

	// NOTE: We are NOT setting the secrets here, as the secrets returned by GetSingle are garbage
	// Secret only applies to Create/Update operations if the "privatekey" or "passphrase" properties are non-empty

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update is called to update the state of the resource. Config, planned
// state, and prior state values should be read from the
// UpdateRequest and new state values set on the UpdateResponse.
func (r *credentialSSHResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data credentialSSHResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cm := r.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	cred := jenkins.SSHCredentials{
		ID:          data.Name.ValueString(),
		Scope:       data.Scope.ValueString(),
		Description: data.Description.ValueString(),
		Username:    data.Username.ValueString(),
		PrivateKeySource: &jenkins.PrivateKey{
			Class: jenkins.KeySourceDirectEntryType,
			Value: data.PrivateKey.ValueString(),
		},
	}

	// Only enforce the secrets if non-empty
	if data.Passphrase.ValueString() != "" {
		cred.Passphrase = data.Passphrase.ValueString()
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
func (r *credentialSSHResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data credentialSSHResourceModel

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
