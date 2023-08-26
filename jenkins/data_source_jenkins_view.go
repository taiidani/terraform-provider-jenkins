package jenkins

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJenkinsView() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJenkinsViewRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The unique name of the Jenkins view.",
				Required:    true,
			},
			"folder": {
				Type:             schema.TypeString,
				Description:      "The folder namespace that the job exists in.",
				Optional:         true,
				ValidateDiagFunc: validateFolderName,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description for the view.",
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url for the view.",
				Computed:    true,
			},
		},
	}
}

func dataSourceJenkinsViewRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	folderName := d.Get("folder").(string)
	d.SetId(formatFolderName(folderName + "/" + name))

	return resourceJenkinsViewRead(ctx, d, meta)
}
