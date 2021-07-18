package jenkins

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJenkinsFolder() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJenkinsFolderRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The unique name of the JenkinsCI folder.",
				Required:         true,
				ValidateDiagFunc: validateJobName,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The name of the folder to display in the UI.",
				Computed:    true,
			},
			"folder": {
				Type:             schema.TypeString,
				Description:      "The folder namespace that the folder exists in.",
				Optional:         true,
				ValidateDiagFunc: validateFolderName,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of this folder's purpose.",
				Computed:    true,
			},
			"security": {
				Type:        schema.TypeSet,
				Description: "The Jenkins project-based security configuration.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"inheritance_strategy": {
							Type:        schema.TypeString,
							Description: "The strategy for applying these permissions sets to existing inherited permissions.",
							Computed:    true,
						},
						"permissions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The Jenkins permissions sets that provide access to this folder.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"template": {
				Type:        schema.TypeString,
				Description: "The configuration file template, used to communicate with Jenkins.",
				Computed:    true,
			},
		},
	}
}

func dataSourceJenkinsFolderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	folderName := d.Get("folder").(string)
	d.SetId(formatFolderName(folderName + "/" + name))

	return resourceJenkinsFolderRead(ctx, d, meta)
}
