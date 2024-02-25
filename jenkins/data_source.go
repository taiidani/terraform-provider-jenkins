package jenkins

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type (
	// dataSourceHelper provides assistive snippets of logic to help reduce duplication in
	// each data source definition.
	dataSourceHelper struct {
		client *jenkinsAdapter
	}
)

func newDataSourceHelper() *dataSourceHelper {
	return &dataSourceHelper{}
}

// Configure should register the client for the resource.
func (d *dataSourceHelper) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dataSourceHelper) schema(s map[string]schema.Attribute) map[string]schema.Attribute {
	if _, ok := s["id"]; !ok {
		s["id"] = schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The full canonical job path, e.g. `/job/job-name`",
		}
	}
	if _, ok := s["name"]; !ok {
		s["name"] = schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The name of the resource being read.",
		}
	}
	if _, ok := s["folder"]; !ok {
		s["folder"] = schema.StringAttribute{
			MarkdownDescription: "The folder namespace containing this resource.",
			Optional:            true,
		}
	}

	return s
}

func (d *dataSourceHelper) schemaCredential(s map[string]schema.Attribute) map[string]schema.Attribute {
	// Pull in the base schema
	s = d.schema(s)

	// Add credential-specific attributes
	if _, ok := s["description"]; !ok {
		s["description"] = schema.StringAttribute{
			MarkdownDescription: "A human readable description of the credentials being stored.",
			Computed:            true,
		}
	}
	if _, ok := s["domain"]; !ok {
		s["domain"] = schema.StringAttribute{
			MarkdownDescription: "The domain store containing this resource.",
			Optional:            true,
		}
	}
	if _, ok := s["scope"]; !ok {
		s["scope"] = schema.StringAttribute{
			MarkdownDescription: `The visibility of the credentials to Jenkins agents. This will be either "GLOBAL" or "SYSTEM".`,
			Computed:            true,
		}
	}

	return s
}
