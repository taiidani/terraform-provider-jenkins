package jenkins

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bndr/gojenkins"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceJenkinsView() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJenkinsViewCreate,
		ReadContext:   resourceJenkinsViewRead,
		UpdateContext: resourceJenkinsViewUpdate,
		DeleteContext: resourceJenkinsViewDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"assigned_projects": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Computed:    true, // No way to update or set description with the gojenkins client at the moment.
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url for the view.",
				Computed:    true,
			},
		},
	}
}

func resourceJenkinsViewCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	cm := client.Credentials()
	name := d.Get("name").(string)

	view, err := cm.J.CreateView(ctx, name, gojenkins.LIST_VIEW)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error creating the Jenkins View: %s", err))
	}

	assigedProjects := d.Get("assigned_projects").([]interface{})
	for _, project := range assigedProjects {
		_, err := view.AddJob(ctx, project.(string))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error adding %s to Jenkins view %s: %s", project.(string), name, err))
		}
	}
	d.SetId(view.GetName())
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

func resourceJenkinsViewUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil // No update-functionality in gojenkins.
}

func resourceJenkinsViewDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	cm := client.Credentials()
	name := d.Get("name").(string)

	_, err := cm.J.Requester.Post(ctx, "/view/"+name+"/doDelete", nil, nil, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting the Jenkins view: %s", err))
	}

	return nil
}
