package jenkins

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type credentialAws struct {
	XMLName            xml.Name `xml:"com.cloudbees.jenkins.plugins.awscredentials.AWSCredentialsImpl"`
	ID                 string   `xml:"id"`
	Scope              string   `xml:"scope"`
	Description        string   `xml:"description"`
	AccessKey          string   `xml:"accessKey"`
	SecretKey          string   `xml:"secretKey"`
	IamRoleArn         string   `xml:"iamRoleArn"`
	IamMfaSerialNumber string   `xml:"iamMfaSerialNumber"`
}

type credentialAwsResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Folder             types.String `tfsdk:"folder"`
	Description        types.String `tfsdk:"description"`
	Domain             types.String `tfsdk:"domain"`
	Scope              types.String `tfsdk:"scope"`
	AccessKey          types.String `tfsdk:"access_key"`
	SecretKey          types.String `tfsdk:"secret_key"`
	IamRoleArn         types.String `tfsdk:"iam_role_arn"`
	IamMfaSerialNumber types.String `tfsdk:"iam_mfa_serial_number"`
}

type credentialAwsResource struct {
	*resourceHelper
}

// Ensure the implementation satisfies the desired interfaces.
var _ resource.ResourceWithConfigure = &credentialAwsResource{}

func newcredentialAwsResource() resource.Resource {
	return &credentialAwsResource{
		resourceHelper: newResourceHelper(),
	}
}

// Metadata should return the full name of the resource.
func (r *credentialAwsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_aws"
}

// Schema should return the schema for this resource.
func (r *credentialAwsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Manages an AWS credential within Jenkins.

~> The "secret_key" property may leave plain-text secret id in your state file. If using the property to manage the secret id in Terraform, ensure that your state file is properly secured and encrypted at rest.

~> The Jenkins installation that uses this resource is expected to have the [AWS Credentials Plugin](https://plugins.jenkins.io/aws-credentials/) installed in their system.`,
		Attributes: r.schemaCredential(map[string]schema.Attribute{
			"access_key": schema.StringAttribute{
				MarkdownDescription: "An AWS access key ID. This is the public part of the key pair used to authenticate with AWS services.",
				Optional:            true,
				Sensitive:           true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"secret_key": schema.StringAttribute{
				MarkdownDescription: "An AWS secret access key. This is the private part of the key pair used to authenticate with AWS services.",
				Optional:            true,
				Sensitive:           true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"iam_role_arn": schema.StringAttribute{
				MarkdownDescription: "An ARN specifying the IAM role to assume. The format should be something like: \"arn:aws:iam::123456789012:role/MyIAMRoleName\".",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"iam_mfa_serial_number": schema.StringAttribute{
				MarkdownDescription: "The identifier for an MFA device. Either a serial number for hardware MFA devices, or an ARN for virtual devices.\n This is only required if the trust policy of the role being assumed includes a condition that requires MFA authentication.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
		}),
	}
}

// Create is called when the provider must create a new resource. Config
// and planned state values should be read from the
// CreateRequest and new state values set on the CreateResponse.
func (r *credentialAwsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data credentialAwsResourceModel

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

	cred := credentialAws{
		ID:                 data.Name.ValueString(),
		Scope:              data.Scope.ValueString(),
		Description:        data.Description.ValueString(),
		AccessKey:          data.AccessKey.ValueString(),
		SecretKey:          data.SecretKey.ValueString(),
		IamRoleArn:         data.IamRoleArn.ValueString(),
		IamMfaSerialNumber: data.IamMfaSerialNumber.ValueString(),
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
func (r *credentialAwsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data credentialAwsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cm := r.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	cred := credentialAws{}
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
	data.AccessKey = types.StringValue(cred.AccessKey)
	data.IamRoleArn = types.StringValue(cred.IamRoleArn)
	data.IamMfaSerialNumber = types.StringValue(cred.IamMfaSerialNumber)
	// NOTE: We are NOT setting the secret here, as the secret returned by GetSingle is garbage
	// Secret only applies to Create/Update operations if the "secret_id" property is non-empty

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update is called to update the state of the resource. Config, planned
// state, and prior state values should be read from the
// UpdateRequest and new state values set on the UpdateResponse.
func (r *credentialAwsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data credentialAwsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cm := r.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	cred := credentialAws{
		ID:                 data.Name.ValueString(),
		Scope:              data.Scope.ValueString(),
		Description:        data.Description.ValueString(),
		AccessKey:          data.AccessKey.ValueString(),
		IamRoleArn:         data.IamRoleArn.ValueString(),
		IamMfaSerialNumber: data.IamMfaSerialNumber.ValueString(),
	}

	// Only enforce the password if it is non-empty
	if data.SecretKey.ValueString() != "" {
		cred.SecretKey = data.SecretKey.ValueString()
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
func (r *credentialAwsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data credentialAwsResourceModel

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
