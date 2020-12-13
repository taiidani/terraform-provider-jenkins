package jenkins

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		CreateContext: resourceJenkinsJobCreate,
		ReadContext:   resourceJenkinsJobRead,
		UpdateContext: resourceJenkinsJobUpdate,
		DeleteContext: resourceJenkinsJobDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The unique name of the JenkinsCI folder.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateJobName,
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
			"permissions": {
				Type:        schema.TypeSet,
				Description: "The Jenkins permissions sets that provide additional access to this folder.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"template": {
				Type:             schema.TypeString,
				Description:      "The configuration file template, used to communicate with Jenkins.",
				Optional:         true,
				Default:          resourceJenkinsFolderTmpl,
				DiffSuppressFunc: templateDiff,
			},
		},
	}
}
