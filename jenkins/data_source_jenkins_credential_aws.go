package jenkins

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type credentialAwsDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Folder             types.String `tfsdk:"folder"`
	Description        types.String `tfsdk:"description"`
	Domain             types.String `tfsdk:"domain"`
	Scope              types.String `tfsdk:"scope"`
	AccessKey          types.String `tfsdk:"access_key"`
	IamRoleArn         types.String `tfsdk:"iam_role_arn"`
	IamMfaSerialNumber types.String `tfsdk:"iam_mfa_serial_number"`
}

type credentialAwsDataSource struct {
	*dataSourceHelper
}

// Ensure the implementation satisfies the desired interfaces.
var _ datasource.DataSourceWithConfigure = &credentialAwsDataSource{}

func newCredentialAwsDataSource() datasource.DataSource {
	return &credentialAwsDataSource{
		dataSourceHelper: newDataSourceHelper(),
	}
}

func (d *credentialAwsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_aws"
}

func (d *credentialAwsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get the attributes of an AWS credential within Jenkins.",
		Attributes: d.schemaCredential(map[string]schema.Attribute{
			"access_key": schema.StringAttribute{
				MarkdownDescription: "An AWS access key ID. This is the public part of the key pair used to authenticate with AWS services.",
				Sensitive:           true,
				Computed:            true,
			},
			"iam_role_arn": schema.StringAttribute{
				MarkdownDescription: "An ARN specifying the IAM role to assume. The format should be something like: \"arn:aws:iam::123456789012:role/MyIAMRoleName\".",
				Computed:            true,
			},
			"iam_mfa_serial_number": schema.StringAttribute{
				MarkdownDescription: "The identifier for an MFA device. Either a serial number for hardware MFA devices, or an ARN for virtual devices.\n This is only required if the trust policy of the role being assumed includes a condition that requires MFA authentication.",
				Computed:            true,
			},
		}),
	}
}

func (d *credentialAwsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data credentialAwsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cm := d.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	if data.Domain.IsNull() {
		data.Domain = basetypes.NewStringValue(defaultCredentialDomain)
	}

	cred := credentialAws{}
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
	data.AccessKey = types.StringValue(cred.AccessKey)
	data.IamRoleArn = types.StringValue(cred.IamRoleArn)
	data.IamMfaSerialNumber = types.StringValue(cred.IamMfaSerialNumber)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
