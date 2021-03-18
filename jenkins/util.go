package jenkins

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// formatFolderName will format a folder name in the way that Jenkins expects, with "name/job/name" separators.
// Deduplication will be performed so that it is safe to pass an already-formatted job into this function.
func formatFolderName(name string) string {
	split := strings.Split(name, "/")

	ret := []string{}
	for _, segment := range split {
		if segment == "" || segment == "job" {
			continue
		}
		ret = append(ret, segment)
	}
	return strings.Join(ret, "/job/")
}

// formatFolderID will format a set of folders in the way that Jenkins expects for the "folder" property, with "/job/name/job/name" separators.
func formatFolderID(folders []string) string {
	if len(folders) == 0 {
		return ""
	}
	return "/job/" + formatFolderName(strings.Join(folders, "/"))
}

// extractFolders prepares a job name for some folder-aware client library calls.
// These calls are different from other calls in that they expect the folders to be specified
// as a series of parameters with no "/job/" separators.
//
// This func will strip out the "/job/" separators from the given string and only return
// the apparent "path" to the folder.
func extractFolders(folder string) (folders []string) {
	for _, item := range strings.Split(folder, "/") {
		if item == "" || item == "job" {
			continue
		}
		folders = append(folders, item)
	}

	return
}

// parseCanonicalJobID will take a canonical Jenkins ID and extract out the base name of the job
// as well as the folder segments that are part of it.
func parseCanonicalJobID(id string) (name string, folders []string) {
	if id == "" {
		return
	}

	folders = extractFolders(id)
	return folders[len(folders)-1], folders[0 : len(folders)-1]
}

// folderExists will validate that a given folder name exists
func folderExists(client jenkinsClient, name string) error {
	folders := extractFolders(name)
	if len(folders) > 0 {
		folderName, parentFolders := parseCanonicalJobID(name)
		_, err := client.GetFolder(folderName, parentFolders...)
		if err != nil {
			return err
		}
	}

	return nil
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

func generateCredentialID(folder, name string) string {
	return fmt.Sprintf("%s/%s", folder, name)
}
