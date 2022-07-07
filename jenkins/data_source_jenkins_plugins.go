package jenkins

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJenkinsPlugins() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluginsRead,
		Schema: map[string]*schema.Schema{
			"list": {
				Type:        schema.TypeList,
				Description: "The list of the Jenkins plugins.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourcePluginsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jenkinsAdapter)

	p, err := client.GetPlugins(ctx, 1)
	if err != nil {
		return diag.FromErr(err)
	}

	// e.g. ["git-server:1.11",  "workflow-api:1153.vb_912c0e47fb_a_", ...]
	var plugins []string
	for _, v := range p.Raw.Plugins {
		plugins = append(plugins, fmt.Sprintf("%s:%s", v.ShortName, v.Version))
	}

	if err := d.Set("list", plugins); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("jenkins-data-source-plugins-id")
	return nil
}
