package jenkins

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJenkinsCredentialUsername() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJenkinsCredentialUsernameRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The identifier assigned to the credentials.",
				Required:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "The domain namespace that the credentials will be added to.",
				Optional:    true,
			},
			"folder": {
				Type:        schema.TypeString,
				Description: "The folder namespace that the credentials will be added to.",
				Optional:    true,
			},
			"scope": {
				Type:        schema.TypeString,
				Description: "The Jenkins scope assigned to the credentials.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The credentials descriptive text.",
				Computed:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The credentials user username.",
				Computed:    true,
			},
		},
	}
}

func dataSourceJenkinsCredentialUsernameRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	folderName := d.Get("folder").(string)
	d.SetId(formatFolderName(folderName + "/" + name))

	return resourceJenkinsCredentialUsernameRead(ctx, d, meta)
}
