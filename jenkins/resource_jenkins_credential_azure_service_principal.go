package jenkins

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AzureServicePrincipalCredentials struct representing credential for storing Azure service credentials
type AzureServicePrincipalCredentials struct {
	XMLName     xml.Name                             `xml:"com.microsoft.azure.util.AzureCredentials"`
	ID          string                               `xml:"id"`
	Scope       string                               `xml:"scope"`
	Description string                               `xml:"description"`
	Data        AzureServicePrincipalCredentialsData `xml:"data"`
}

type AzureServicePrincipalCredentialsData struct {
	SubscriptionId          string `xml:"subscriptionId"`
	ClientId                string `xml:"clientId"`
	ClientSecret            string `xml:"clientSecret"`
	CertificateId           string `xml:"certificateId"`
	Tenant                  string `xml:"tenant"`
	AzureEnvironmentName    string `xml:"azureEnvironmentName"`
	ServiceManagementURL    string `xml:"serviceManagementURL"`
	AuthenticationEndpoint  string `xml:"authenticationEndpoint"`
	ResourceManagerEndpoint string `xml:"resourceManagerEndpoint"`
	GraphEndpoint           string `xml:"graphEndpoint"`
}

type credentialAzureServicePrincipalResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Folder                  types.String `tfsdk:"folder"`
	Description             types.String `tfsdk:"description"`
	Domain                  types.String `tfsdk:"domain"`
	Scope                   types.String `tfsdk:"scope"`
	SubscriptionId          types.String `tfsdk:"subscription_id"`
	ClientId                types.String `tfsdk:"client_id"`
	ClientSecret            types.String `tfsdk:"client_secret"`
	CertificateId           types.String `tfsdk:"certificate_id"`
	Tenant                  types.String `tfsdk:"tenant"`
	AzureEnvironmentName    types.String `tfsdk:"azure_environment_name"`
	ServiceManagementURL    types.String `tfsdk:"service_management_url"`
	AuthenticationEndpoint  types.String `tfsdk:"authentication_endpoint"`
	ResourceManagerEndpoint types.String `tfsdk:"resource_manager_endpoint"`
	GraphEndpoint           types.String `tfsdk:"graph_endpoint"`
}

type credentialAzureServicePrincipalResource struct {
	*resourceHelper
}

// Ensure the implementation satisfies the desired interfaces.
var _ resource.ResourceWithConfigure = &credentialAzureServicePrincipalResource{}

func newCredentialAzureServicePrincipalResource() resource.Resource {
	return &credentialAzureServicePrincipalResource{
		resourceHelper: newResourceHelper(),
	}
}

// Metadata should return the full name of the resource.
func (r *credentialAzureServicePrincipalResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_azure_service_principal"
}

// Schema should return the schema for this resource.
func (r *credentialAzureServicePrincipalResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Manages an Azure Service Principal credential within Jenkins. This credential may then be referenced within jobs that are created.

~> The "client_secret" property may leave plain-text secret id in your state file. If using the property to manage the secret id in Terraform, ensure that your state file is properly secured and encrypted at rest.

~> The Jenkins installation that uses this resource is expected to have the [Azure Credentials Plugin](https://plugins.jenkins.io/azure-credentials/) installed in their system.`,
		Attributes: r.schemaCredential(map[string]schema.Attribute{
			"subscription_id": schema.StringAttribute{
				MarkdownDescription: "The Azure subscription id mapped to the Azure Service Principal.",
				Required:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The client id (application id) of the Azure Service Principal.",
				Required:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "The client secret of the Azure Service Principal. Cannot be used with `certificate_id`. Has to be specified, if `certificate_id` is not specified.",
				Sensitive:           true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("client_secret"), path.MatchRoot("certificate_id")),
				},
			},
			"certificate_id": schema.StringAttribute{
				MarkdownDescription: "The certificate reference of the Azure Service Principal, pointing to a Jenkins certificate credential. Cannot be used with `client_secret`. Has to be specified, if `client_secret` is not specified.",
				Sensitive:           true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("client_secret"), path.MatchRoot("certificate_id")),
				},
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The Azure Tenant ID of the Azure Service Principal.",
				Required:            true,
			},
			"azure_environment_name": schema.StringAttribute{
				MarkdownDescription: `The Azure Cloud enviroment name. Allowed values are "Azure", "Azure China", "Azure Germany", "Azure US Government".`,
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Azure"),
				Validators: []validator.String{
					stringvalidator.OneOf("Azure", "Azure China", "Azure Germany", "Azure US Government"),
				},
			},
			"service_management_url": schema.StringAttribute{
				MarkdownDescription: "Override the Azure management endpoint URL for the selected Azure environment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"authentication_endpoint": schema.StringAttribute{
				MarkdownDescription: "Override the Azure Active Directory endpoint for the selected Azure environment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"resource_manager_endpoint": schema.StringAttribute{
				MarkdownDescription: "Override the Azure resource manager endpoint URL for the selected Azure environment.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"graph_endpoint": schema.StringAttribute{
				MarkdownDescription: "Override the Azure graph endpoint URL for the selected Azure environment.",
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
func (r *credentialAzureServicePrincipalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data credentialAzureServicePrincipalResourceModel

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

	credData := AzureServicePrincipalCredentialsData{
		SubscriptionId:          data.SubscriptionId.ValueString(),
		ClientId:                data.ClientId.ValueString(),
		ClientSecret:            data.ClientSecret.ValueString(),
		CertificateId:           data.CertificateId.ValueString(),
		Tenant:                  data.Tenant.ValueString(),
		AzureEnvironmentName:    data.AzureEnvironmentName.ValueString(),
		ServiceManagementURL:    data.ServiceManagementURL.ValueString(),
		AuthenticationEndpoint:  data.AuthenticationEndpoint.ValueString(),
		ResourceManagerEndpoint: data.ResourceManagerEndpoint.ValueString(),
		GraphEndpoint:           data.GraphEndpoint.ValueString(),
	}

	cred := AzureServicePrincipalCredentials{
		ID:          data.Name.ValueString(),
		Scope:       data.Scope.ValueString(),
		Description: data.Description.ValueString(),
		Data:        credData,
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
func (r *credentialAzureServicePrincipalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data credentialAzureServicePrincipalResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cm := r.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	cred := AzureServicePrincipalCredentials{}
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

	// NOTE: We are NOT setting the password here, as the password returned by GetSingle is garbage
	// Password only applies to Create/Update operations if the "password" property is non-empty

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update is called to update the state of the resource. Config, planned
// state, and prior state values should be read from the
// UpdateRequest and new state values set on the UpdateResponse.
func (r *credentialAzureServicePrincipalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data credentialAzureServicePrincipalResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cm := r.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	credData := AzureServicePrincipalCredentialsData{
		SubscriptionId:          data.SubscriptionId.ValueString(),
		ClientId:                data.ClientId.ValueString(),
		Tenant:                  data.Tenant.ValueString(),
		AzureEnvironmentName:    data.AzureEnvironmentName.ValueString(),
		ServiceManagementURL:    data.ServiceManagementURL.ValueString(),
		AuthenticationEndpoint:  data.AuthenticationEndpoint.ValueString(),
		ResourceManagerEndpoint: data.ResourceManagerEndpoint.ValueString(),
		GraphEndpoint:           data.GraphEndpoint.ValueString(),
	}

	cred := AzureServicePrincipalCredentials{
		ID:          data.Name.ValueString(),
		Scope:       data.Scope.ValueString(),
		Description: data.Description.ValueString(),
		Data:        credData,
	}

	// Only enforce the password if it is non-empty
	if data.ClientSecret.ValueString() != "" {
		cred.Data.ClientSecret = data.ClientSecret.ValueString()
	}

	if data.CertificateId.ValueString() != "" {
		cred.Data.ClientId = data.CertificateId.ValueString()
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
func (r *credentialAzureServicePrincipalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data credentialAzureServicePrincipalResourceModel

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
