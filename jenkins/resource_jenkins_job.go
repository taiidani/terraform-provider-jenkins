package jenkins

import (
	"context"
	"fmt"
	"log"
	"strings"

	jenkins "github.com/bndr/gojenkins"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceJenkinsJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJenkinsJobCreate,
		ReadContext:   resourceJenkinsJobRead,
		UpdateContext: resourceJenkinsJobUpdate,
		DeleteContext: resourceJenkinsJobDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The unique name of the JenkinsCI job. Folders may be specified as foldername/name.",
				Required:    true,
				ForceNew:    true,
			},
			"template": {
				Type:             schema.TypeString,
				Description:      "The configuration file template, used to communicate with Jenkins.",
				Required:         true,
				DiffSuppressFunc: templateDiff,
			},
			"parameters": {
				Type:        schema.TypeMap,
				Description: "The set of parameters to be rendered in the template when generating a valid config.xml file.",
				Optional:    true,
				Elem:        schema.TypeString,
			},
		},
	}
}

func resourceJenkinsJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jenkins.Jenkins)
	name := formatJobName(d.Get("name").(string))
	baseName, folders := parseJobName(d.Get("name").(string))

	xml, err := renderTemplate(d.Get("template").(string), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::create - Error binding config.xml template to %q: %w", name, err))
	}

	_, err = client.CreateJobInFolder(xml, baseName, folders...)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::create - Error creating job for %q: %w", name, err))
	}

	log.Printf("[DEBUG] jenkins::create - job %q created", name)
	d.SetId(name)

	return resourceJenkinsJobRead(ctx, d, meta)
}

func resourceJenkinsJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	name := d.Id()

	log.Printf("[DEBUG] jenkins::read - Looking for job %q", name)

	job, err := client.GetJob(name)
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			// Job does not exist
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("jenkins::read - Job %q does not exist: %w", name, err))
	}

	config, err := job.GetConfig()
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::read - Job %q could not extract configuration: %v", name, err))
	}

	log.Printf("[DEBUG] jenkins::read - Job %q exists", name)
	if err := d.Set("template", config); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceJenkinsJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jenkins.Jenkins)
	name := d.Id()

	// grab job by current name
	job, err := client.GetJob(name)

	xml, err := renderTemplate(d.Get("template").(string), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::update - Error binding config.xml template to %q: %w", name, err))
	}

	err = job.UpdateConfig(xml)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::update - Error updating job %q configuration: %w", name, err))
	}

	return resourceJenkinsJobRead(ctx, d, meta)
}

func resourceJenkinsJobDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	name := d.Id()

	log.Printf("[DEBUG] jenkins::delete - Removing %q", name)

	ok, err := client.DeleteJob(name)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] jenkins::delete - %q removed: %t", name, ok)
	return nil
}
