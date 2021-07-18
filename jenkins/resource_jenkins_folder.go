package jenkins

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceJenkinsFolder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJenkinsFolderCreate,
		ReadContext:   resourceJenkinsFolderRead,
		UpdateContext: resourceJenkinsFolderUpdate,
		DeleteContext: resourceJenkinsJobDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The unique name of the JenkinsCI folder.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateJobName,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The name of the folder to display in the UI.",
				Optional:    true,
			},
			"folder": {
				Type:             schema.TypeString,
				Description:      "The folder namespace that the folder will be added to as a subfolder.",
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateFolderName,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of this folder's purpose.",
				Optional:    true,
			},
			"security": {
				Type:        schema.TypeSet,
				Description: "The Jenkins project-based security configuration.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"inheritance_strategy": {
							Type:        schema.TypeString,
							Description: "The strategy for applying these permissions sets to existing inherited permissions.",
							Optional:    true,
							Default:     "org.jenkinsci.plugins.matrixauth.inheritance.InheritParentStrategy",
						},
						"permissions": {
							Type:        schema.TypeList,
							Required:    true,
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

func resourceJenkinsFolderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	name := d.Get("name").(string)
	folderName := d.Get("folder").(string)

	// Validate that the folder exists
	if err := folderExists(ctx, client, folderName); err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::create - Could not find folder '%s': %w", folderName, err))
	}

	f := folder{
		Description: d.Get("description").(string),
		DisplayName: d.Get("display_name").(string),
	}
	f.Properties.Security = expandSecurity(d.Get("security").(*schema.Set).List())

	xml, err := f.Render()
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::create - Error binding config.xml template to %q: %w", name, err))
	}

	folders := extractFolders(folderName)
	_, err = client.CreateJobInFolder(ctx, string(xml), name, folders...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::create - Error creating job for %q in folder %s: %w", name, folderName, err))
	}

	log.Printf("[DEBUG] jenkins::create - job %q created in folder %s", name, folderName)
	d.SetId(formatFolderName(folderName + "/" + name))

	return resourceJenkinsFolderRead(ctx, d, meta)
}

func resourceJenkinsFolderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	name, folders := parseCanonicalJobID(d.Id())

	log.Printf("[DEBUG] jenkins::read - Looking for job %q", name)

	job, err := client.GetJob(ctx, name, folders...)
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			// Job does not exist
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("jenkins::read - Job %q does not exist: %w", name, err))
	}

	// Extract the raw XML configuration
	config, err := job.GetConfig(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::read - Job %q could not extract configuration: %v", job.Base, err))
	}

	log.Printf("[DEBUG] jenkins::read - Job %q exists", job.Base)
	d.SetId(job.Base)

	if err := d.Set("template", config); err != nil {
		return diag.FromErr(err)
	}

	// Next, parse the properties from the config
	f, err := parseFolder(config)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("display_name", f.DisplayName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("folder", formatFolderID(folders)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", f.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("security", flattenSecurity(f.Properties.Security)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceJenkinsFolderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	name, folders := parseCanonicalJobID(d.Id())

	// grab job by current name
	job, err := client.GetJob(ctx, name, folders...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::update - Could not find job %q: %w", name, err))
	}

	// Extract the raw XML configuration
	config, err := job.GetConfig(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::update - Job %q could not extract configuration: %v", job.Base, err))
	}

	// Next, parse the properties from the config
	f, err := parseFolder(config)
	if err != nil {
		return diag.FromErr(err)
	}

	// Then update the values
	f.Description = d.Get("description").(string)
	f.DisplayName = d.Get("display_name").(string)
	f.Properties.Security = expandSecurity(d.Get("security").(*schema.Set).List())

	// And send it back to Jenkins
	xml, err := f.Render()
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::create - Error binding config.xml template to %q: %w", name, err))
	}

	err = job.UpdateConfig(ctx, string(xml))
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::update - Error updating job %q configuration: %w", name, err))
	}

	return resourceJenkinsFolderRead(ctx, d, meta)
}

func expandSecurity(config []interface{}) *folderSecurity {
	if len(config) == 0 {
		return nil
	}

	ret := &folderSecurity{}
	data := config[0].(map[string]interface{})
	ret.InheritanceStrategy = folderPermissionInheritanceStrategy{
		Class: data["inheritance_strategy"].(string),
	}
	ret.Permission = []string{}
	for _, permission := range data["permissions"].([]interface{}) {
		ret.Permission = append(ret.Permission, permission.(string))
	}
	return ret
}

func flattenSecurity(config *folderSecurity) []map[string]interface{} {
	ret := []map[string]interface{}{}
	if config == nil {
		return ret
	}

	d := map[string]interface{}{}
	d["inheritance_strategy"] = config.InheritanceStrategy.Class
	d["permissions"] = config.Permission

	return append(ret, d)
}
