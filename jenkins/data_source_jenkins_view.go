package jenkins

import (
	"context"
	"fmt"
	"log"
	"strings"

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

func resourceJenkinsViewRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	name, _ := parseCanonicalJobID(d.Id())

	log.Printf("[DEBUG] jenkins::read - Looking for view %q", name)

	view, err := client.GetView(ctx, name)
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			// View does not exist
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("jenkins::read - View %q does not exist: %w", name, err))
	}

	description := view.GetDescription()
	err = d.Set("description", description)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::read - Description could not be set for View %q, %w", name, err))
	}

	url := view.GetUrl()
	err = d.Set("url", url)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::read - Url could not be set for View %q, %w", name, err))
	}

	name = view.GetName()
	err = d.Set("name", name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::read - Name could not be set for View %q, %w", name, err))
	}

	return nil
}
