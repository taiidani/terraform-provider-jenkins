package jenkins

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func formatJobName(name string) string {
	split := strings.Split(name, "/")
	return strings.Join(split, "/job/")
}

func parseJobName(name string) (job string, folders []string) {
	split := strings.Split(name, "/")
	return split[len(split)-1], split[0 : len(split)-1]
}

func templateDiff(k, old, new string, d *schema.ResourceData) bool {
	new, _ = renderTemplate(new, d)

	// Sanitize the XML entries to prevent inadvertent inequalities
	old = strings.Replace(old, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>", "", -1)
	old = strings.Replace(old, " ", "", -1)
	old = strings.TrimSpace(old)
	new = strings.Replace(new, "<?new version=\"1.0\" encoding=\"UTF-8\"?>", "", -1)
	new = strings.Replace(new, " ", "", -1)
	new = strings.TrimSpace(new)

	log.Printf("[DEBUG] jenkins::diff - Old: %q", old)
	log.Printf("[DEBUG] jenkins::diff - New: %q", new)
	return old == new
}
