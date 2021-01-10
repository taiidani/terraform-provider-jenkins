package jenkins

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJenkinsJob() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJenkinsJobRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The unique name of the JenkinsCI job.",
				Required:         true,
				ValidateDiagFunc: validateJobName,
			},
			"folder": {
				Type:             schema.TypeString,
				Description:      "The folder namespace that the job exists in.",
				Optional:         true,
				ValidateDiagFunc: validateFolderName,
			},
			"template": {
				Type:        schema.TypeString,
				Description: "The configuration file template, used to communicate with Jenkins.",
				Computed:    true,
			},
		},
	}
}

func dataSourceJenkinsJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	folderName := d.Get("folder").(string)
	d.SetId(formatFolderName(folderName + "/" + name))

	return resourceJenkinsJobRead(ctx, d, meta)
}
