package jenkins

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJenkinsCredentialVaultAppRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJenkinsCredentialVaultAppRoleRead,
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
			"namespace": {
				Type:        schema.TypeString,
				Description: "Namespace of the roles approle backend.",
				Optional:    true,
			},
			"path": {
				Type:        schema.TypeString,
				Description: "Path of the roles approle backend.",
				Computed:    true,
			},
			"role_id": {
				Type:        schema.TypeString,
				Description: "The roles role_id.",
				Computed:    true,
			},
		},
	}
}

func dataSourceJenkinsCredentialVaultAppRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	folderName := d.Get("folder").(string)
	d.SetId(formatFolderName(folderName + "/" + name))

	return resourceJenkinsCredentialVaultAppRoleRead(ctx, d, meta)
}
