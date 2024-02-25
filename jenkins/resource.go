package jenkins

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type (
	// resourceHelper provides assistive snippets of logic to help reduce duplication in
	// each resource definition.
	resourceHelper struct {
		client *jenkinsAdapter
	}
)

func newResourceHelper() *resourceHelper {
	return &resourceHelper{}
}

// Configure should register the client for the resource.
func (r *resourceHelper) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ImportState is called when performing import operations of existing resources.
func (r *resourceHelper) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *resourceHelper) schema(s map[string]schema.Attribute) map[string]schema.Attribute {
	if _, ok := s["id"]; !ok {
		s["id"] = schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The full canonical job path, e.g. `/job/job-name`",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
	}
	if _, ok := s["name"]; !ok {
		s["name"] = schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The name of the resource being created. This maps to the ID property within Jenkins, and cannot be changed once set.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		}
	}
	if _, ok := s["folder"]; !ok {
		s["folder"] = schema.StringAttribute{
			MarkdownDescription: "The folder namespace to store the resource in. If not set will default to global Jenkins.",
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		}
	}

	return s
}
func (r *resourceHelper) schemaCredential(s map[string]schema.Attribute) map[string]schema.Attribute {
	// Pull in the base schema
	s = r.schema(s)

	// Add credential-specific attributes
	if _, ok := s["description"]; !ok {
		s["description"] = schema.StringAttribute{
			MarkdownDescription: "A human readable description of the credentials being stored.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("Managed by Terraform"),
		}
	}
	if _, ok := s["domain"]; !ok {
		s["domain"] = schema.StringAttribute{
			MarkdownDescription: "The domain store to place the credentials into. If not set will default to the global credentials store.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(defaultCredentialDomain),
			PlanModifiers: []planmodifier.String{
				// In-place updates should be possible, but gojenkins does not support move operations
				stringplanmodifier.RequiresReplace(),
			},
		}
	}
	if _, ok := s["scope"]; !ok {
		s["scope"] = schema.StringAttribute{
			MarkdownDescription: `The visibility of the credentials to Jenkins agents. This must be set to either "GLOBAL" or "SYSTEM". If not set will default to "GLOBAL".`,
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("GLOBAL"),
			Validators: []validator.String{
				stringvalidator.OneOf(supportedCredentialScopes...),
			},
		}
	}

	return s
}
