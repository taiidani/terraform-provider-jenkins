package jenkins

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const resourceJenkinsFolderTmpl = `<com.cloudbees.hudson.plugins.folder.Folder plugin="cloudbees-folder@6.6">
  <actions/>
  <description>{{ .Description }}</description>
  <properties>
    <com.cloudbees.hudson.plugins.folder.properties.AuthorizationMatrixProperty>
      <inheritanceStrategy class="org.jenkinsci.plugins.matrixauth.inheritance.InheritParentStrategy"/>
      {{ range $value := .Permissions }}
      <permission>{{ $value }}</permission>
      {{ end }}
    </com.cloudbees.hudson.plugins.folder.properties.AuthorizationMatrixProperty>
  </properties>
  <icon class="com.cloudbees.hudson.plugins.folder.icons.StockFolderIcon"/>
</com.cloudbees.hudson.plugins.folder.Folder>
`

func resourceJenkinsFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceJenkinsJobCreate,
		Read:   resourceJenkinsJobRead,
		Update: resourceJenkinsJobUpdate,
		Delete: resourceJenkinsJobDelete,
		Exists: resourceJenkinsJobExists,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The unique name of the JenkinsCI folder. Subfolders may be specified as foldername/name.",
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of this folder's purpose.",
				Optional:    true,
			},
			"permissions": {
				Type:        schema.TypeSet,
				Description: "The Jenkins permissions sets that provide additional access to this folder.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"template": &schema.Schema{
				Type:             schema.TypeString,
				Description:      "The configuration file template, used to communicate with Jenkins.",
				Optional:         true,
				Default:          resourceJenkinsFolderTmpl,
				DiffSuppressFunc: templateDiff,
			},
		},
	}
}
