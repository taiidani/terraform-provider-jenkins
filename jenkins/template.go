package jenkins

import (
	"bytes"
	"log"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Job contains all the data pertaining to a Jenkins job, in a format that is
// easy to use with Golang text/templates
type job struct {
	Name        string
	Description string
	Permissions []string
	Parameters  map[string]string
}

func renderTemplate(data string, d *schema.ResourceData) (string, error) {
	log.Printf("[DEBUG] jenkins::xml - Binding template:\n%s", data)

	// create and parse the config.xml template
	tpl, err := template.New("template").Parse(data)
	if err != nil {
		log.Printf("[ERROR] jenkins::xml - Error parsing template: %v", err)
		return "", err
	}

	// now copy the input parameters into a data structure that is compatible
	// with the config.xml template
	j := &job{
		Name:       d.Get("name").(string),
		Parameters: map[string]string{},
	}
	if value, ok := d.GetOk("description"); ok {
		j.Description = value.(string)
	}
	if value, ok := d.GetOk("permissions"); ok {
		value := value.(*schema.Set)
		for _, v := range value.List() {
			j.Permissions = append(j.Permissions, v.(string))
		}
	}
	if value, ok := d.GetOk("parameters"); ok {
		value := value.(map[string]interface{})
		for k, v := range value {
			j.Parameters[k] = v.(string)
		}
	}

	// apply the job object to the template
	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, j)
	if err != nil {
		log.Printf("[ERROR] jenkis::xml - Error executing template: %v", err)
		return "", err
	}

	xml := buffer.String()
	log.Printf("[DEBUG] jenkins::xml - Bound template:\n%s", xml)
	return xml, nil
}
