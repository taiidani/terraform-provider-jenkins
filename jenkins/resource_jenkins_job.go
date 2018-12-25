package jenkins

import (
	"log"

	jenkins "github.com/bndr/gojenkins"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceJenkinsJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceJenkinsJobCreate,
		Read:   resourceJenkinsJobRead,
		Update: resourceJenkinsJobUpdate,
		Delete: resourceJenkinsJobDelete,
		Exists: resourceJenkinsJobExists,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The unique name of the JenkinsCI job. Folders may be specified as foldername/name.",
				Required:    true,
				ForceNew:    true,
			},
			"template": &schema.Schema{
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

func resourceJenkinsJobExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*jenkins.Jenkins)
	name := d.Id()

	log.Printf("[DEBUG] jenkins::exists - Checking if job %q exists", name)

	_, err := client.GetJob(name)
	if err != nil {
		log.Printf("[DEBUG] jenkins::exists - Job %q does not exist: %v", name, err)
		d.SetId("")
		return false, nil
	}

	log.Printf("[DEBUG] jenkins::exists - Job %q exists", name)
	return true, nil
}

func resourceJenkinsJobCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*jenkins.Jenkins)
	name := formatJobName(d.Get("name").(string))
	baseName, folders := parseJobName(d.Get("name").(string))

	xml, err := renderTemplate(d.Get("template").(string), d)
	if err != nil {
		log.Printf("[ERROR] jenkins::create - Error binding config.xml template to %q: %v", name, err)
		return err
	}

	_, err = client.CreateJobInFolder(xml, baseName, folders...)
	if err != nil {
		log.Printf("[ERROR] jenkins::create - Error creating job for %q: %v", name, err)
		return err
	}

	log.Printf("[DEBUG] jenkins::create - job %q created", name)
	d.SetId(name)

	return resourceJenkinsJobRead(d, meta)
}

func resourceJenkinsJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*jenkins.Jenkins)
	name := d.Id()

	log.Printf("[DEBUG] jenkins::read - Looking for job %q", name)

	job, err := client.GetJob(name)
	if err != nil {
		log.Printf("[DEBUG] jenkins::read - Job %q does not exist: %v", name, err)
		return err
	}

	config, err := job.GetConfig()
	if err != nil {
		log.Printf("[DEBUG] jenkins::read - Job %q could not extract configuration: %v", name, err)
		return err
	}

	log.Printf("[DEBUG] jenkins::read - Job %q exists", name)
	d.Set("template", config)

	return nil
}

func resourceJenkinsJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*jenkins.Jenkins)
	name := d.Id()

	// grab job by current name
	job, err := client.GetJob(name)

	xml, err := renderTemplate(d.Get("template").(string), d)
	if err != nil {
		log.Printf("[ERROR] jenkins::update - Error binding config.xml template to %q: %v", name, err)
		return err
	}

	err = job.UpdateConfig(xml)
	if err != nil {
		log.Printf("[ERROR] jenkins::update - Error updating job %q configuration: %v", name, err)
		return err
	}

	return resourceJenkinsJobRead(d, meta)
}

func resourceJenkinsJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*jenkins.Jenkins)
	name := d.Id()

	log.Printf("[DEBUG] jenkins::delete - Removing %q", name)

	ok, err := client.DeleteJob(name)

	log.Printf("[DEBUG] jenkins::delete - %q removed: %t", name, ok)
	return err
}
