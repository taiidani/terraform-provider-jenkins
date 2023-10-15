package jenkins

import (
	"context"
	"fmt"
	"strings"

	"github.com/bndr/gojenkins"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ViewResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Folder           types.String `tfsdk:"folder"`
	Description      types.String `tfsdk:"description"`
	AssignedProjects types.List   `tfsdk:"assigned_projects"`
	URL              types.String `tfsdk:"url"`
}

type ViewResource struct {
	*resourceHelper
}

// Ensure the implementation satisfies the desired interfaces.
var _ resource.ResourceWithConfigure = &ViewResource{}

func newViewResource() resource.Resource {
	return &ViewResource{
		resourceHelper: newResourceHelper(),
	}
}

// Metadata should return the full name of the resource.
func (r *ViewResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_view"
}

// Schema should return the schema for this resource.
func (r *ViewResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Manages a view within Jenkins.

~> Due to API client limitations, updates to some attributes may be restricted.`,
		Attributes: r.schema(map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique name of this view.",
			},
			"assigned_projects": schema.ListAttribute{
				MarkdownDescription: "The list of projects assigned to the view. For example, the name of a folder.",
				Optional:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description for the view.",
				Computed:            true, // No way to update or set description with the gojenkins client at the moment.
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The url for the view.",
				Computed:            true,
			},
		}),
	}
}

// Create is called when the provider must create a new resource. Config
// and planned state values should be read from the
// CreateRequest and new state values set on the CreateResponse.
func (r *ViewResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ViewResourceModel

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

	view, err := cm.J.CreateView(ctx, data.Name.ValueString(), gojenkins.LIST_VIEW)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while creating the resource. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)

		return
	}

	assigedProjects := data.AssignedProjects.Elements()
	for _, project := range assigedProjects {
		projectName := strings.Trim(project.String(), "\"")
		_, err := view.AddJob(ctx, projectName)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Assign View Projects",
				fmt.Sprintf("Error adding %q to Jenkins view %q: %s", projectName, data.Name.ValueString(), err),
			)

			_, err := cm.J.Requester.Post(ctx, "/view/"+data.Name.ValueString()+"/doDelete", nil, nil, nil)
			if err != nil {
				resp.Diagnostics.AddError(
					"Unable to Delete Resource",
					"An unexpected error occurred while deleting the resource. "+
						"Please report this issue to the provider developers.\n\n"+
						"Error: "+err.Error(),
				)
			}

			return
		}
	}

	// Convert from the API data model to the Terraform data model
	// and set any unknown attribute values.
	data.ID = types.StringValue(view.GetName())
	data.Name = types.StringValue(view.GetName())
	data.Description = types.StringValue(view.GetDescription())
	data.URL = types.StringValue(view.GetUrl())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read is called when the provider must read resource values in order
// to update state. Planned state values should be read from the
// ReadRequest and new state values set on the ReadResponse.
func (r *ViewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ViewResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cm := r.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	view, err := cm.J.GetView(ctx, data.ID.ValueString())
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

	data.ID = types.StringValue(view.GetName())
	data.Name = types.StringValue(view.GetName())
	data.Description = types.StringValue(view.GetDescription())
	data.URL = types.StringValue(view.GetUrl())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update is called to update the state of the resource. Config, planned
// state, and prior state values should be read from the
// UpdateRequest and new state values set on the UpdateResponse.
func (r *ViewResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ViewResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Skip any updating of the resource.
	// No update-functionality in gojenkins.

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete is called when the provider must delete the resource. Config
// values may be read from the DeleteRequest.
//
// If execution completes without error, the framework will automatically
// call DeleteResponse.State.RemoveResource(), so it can be omitted
// from provider logic.
func (r *ViewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ViewResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	cm := r.client.Credentials()
	cm.Folder = formatFolderName(data.Folder.ValueString())

	_, err := cm.J.Requester.Post(ctx, "/view/"+data.Name.ValueString()+"/doDelete", nil, nil, nil)
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
