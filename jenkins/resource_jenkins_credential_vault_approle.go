package jenkins

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// VaultAppRoleCredentials struct representing credential for storing Vault AppRole role id and secret id
type VaultAppRoleCredentials struct {
	XMLName     xml.Name `xml:"com.datapipe.jenkins.vault.credentials.VaultAppRoleCredential"`
	ID          string   `xml:"id"`
	Scope       string   `xml:"scope"`
	Description string   `xml:"description"`
	Namespace   string   `xml:"namespace"`
	Path        string   `xml:"path"`
	RoleID      string   `xml:"roleId"`
	SecretID    string   `xml:"secretId"`
}

func resourceJenkinsCredentialVaultAppRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJenkinsCredentialVaultAppRoleCreate,
		ReadContext:   resourceJenkinsCredentialVaultAppRoleRead,
		UpdateContext: resourceJenkinsCredentialVaultAppRoleUpdate,
		DeleteContext: resourceJenkinsCredentialVaultAppRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceJenkinsCredentialVaultAppRoleImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The identifier assigned to the credentials.",
				Required:    true,
				ForceNew:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "The domain namespace that the credentials will be added to.",
				Optional:    true,
				Default:     "_",
				// In-place updates should be possible, but gojenkins does not support move operations
				ForceNew: true,
			},
			"folder": {
				Type:        schema.TypeString,
				Description: "The folder namespace that the credentials will be added to.",
				Optional:    true,
				ForceNew:    true,
			},
			"scope": {
				Type:             schema.TypeString,
				Description:      "The Jenkins scope assigned to the credentials.",
				Optional:         true,
				Default:          "GLOBAL",
				ValidateDiagFunc: validateCredentialScope,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The credentials descriptive text.",
				Optional:    true,
				Default:     "Managed by Terraform",
			},
			"namespace": {
				Type:        schema.TypeString,
				Description: "Namespace of the roles approle backend.",
				Optional:    true,
			},
			"path": {
				Type:        schema.TypeString,
				Description: "Path of the roles approle backend.",
				Optional:    true,
				Default:     "approle",
			},
			"role_id": {
				Type:        schema.TypeString,
				Description: "The roles role_id.",
				Required:    true,
			},
			"secret_id": {
				Type:        schema.TypeString,
				Description: "The roles secret_id. If left empty will be unmanaged.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceJenkinsCredentialVaultAppRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	cm := client.Credentials()
	cm.Folder = formatFolderName(d.Get("folder").(string))
	// return diag.FromErr(fmt.Errorf("invalid folder name '%s', '%s'", cm.Folder, d.Get("folder").(string)))
	// Validate that the folder exists
	if err := folderExists(ctx, client, cm.Folder); err != nil {
		return diag.FromErr(fmt.Errorf("invalid folder name '%s' specified: %w", cm.Folder, err))
	}

	cred := VaultAppRoleCredentials{
		ID:          d.Get("name").(string),
		Scope:       d.Get("scope").(string),
		Description: d.Get("description").(string),
		Namespace:   d.Get("namespace").(string),
		Path:        d.Get("path").(string),
		RoleID:      d.Get("role_id").(string),
		SecretID:    d.Get("secret_id").(string),
	}

	domain := d.Get("domain").(string)
	err := cm.Add(ctx, domain, cred)
	if err != nil {
		return diag.Errorf("Could not create vault approle credentials: %s", err)
	}

	d.SetId(generateCredentialID(d.Get("folder").(string), cred.ID))
	return resourceJenkinsCredentialVaultAppRoleRead(ctx, d, meta)
}

func resourceJenkinsCredentialVaultAppRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cm := meta.(jenkinsClient).Credentials()
	cm.Folder = formatFolderName(d.Get("folder").(string))

	cred := VaultAppRoleCredentials{}
	err := cm.GetSingle(
		ctx,
		d.Get("domain").(string),
		d.Get("name").(string),
		&cred,
	)

	if err != nil {
		if strings.HasSuffix(err.Error(), "404") {
			// Job does not exist
			d.SetId("")
			return nil
		}

		return diag.Errorf("Could not read vault approle credentials: %s", err)
	}

	d.SetId(generateCredentialID(d.Get("folder").(string), cred.ID))
	_ = d.Set("scope", cred.Scope)
	_ = d.Set("description", cred.Description)
	_ = d.Set("namespace", cred.Namespace)
	_ = d.Set("path", cred.Path)
	_ = d.Set("role_id", cred.RoleID)
	// NOTE: We are NOT setting the password here, as the password returned by GetSingle is garbage
	// Password only applies to Create/Update operations if the "password" property is non-empty

	return nil
}

func resourceJenkinsCredentialVaultAppRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cm := meta.(jenkinsClient).Credentials()
	cm.Folder = formatFolderName(d.Get("folder").(string))

	domain := d.Get("domain").(string)
	cred := VaultAppRoleCredentials{
		ID:          d.Get("name").(string),
		Scope:       d.Get("scope").(string),
		Description: d.Get("description").(string),
		Namespace:   d.Get("namespace").(string),
		Path:        d.Get("path").(string),
		RoleID:      d.Get("role_id").(string),
	}

	// Only enforce the password if it is non-empty
	if d.Get("secret_id").(string) != "" {
		cred.SecretID = d.Get("secret_id").(string)
	}

	err := cm.Update(ctx, domain, d.Get("name").(string), &cred)
	if err != nil {
		return diag.Errorf("Could not update vault approle credentials: %s", err)
	}

	d.SetId(generateCredentialID(d.Get("folder").(string), cred.ID))
	return resourceJenkinsCredentialVaultAppRoleRead(ctx, d, meta)
}

func resourceJenkinsCredentialVaultAppRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cm := meta.(jenkinsClient).Credentials()
	cm.Folder = formatFolderName(d.Get("folder").(string))

	err := cm.Delete(
		ctx,
		d.Get("domain").(string),
		d.Get("name").(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceJenkinsCredentialVaultAppRoleImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{d}

	splitID := strings.Split(d.Id(), "/")
	if len(splitID) < 2 {
		return ret, fmt.Errorf("import ID was improperly formatted. Imports need to be in the format \"[<folder>/]<domain>/<name>\"")
	}

	name := splitID[len(splitID)-1]
	_ = d.Set("name", name)

	domain := splitID[len(splitID)-2]
	_ = d.Set("domain", domain)

	folder := strings.Trim(strings.Join(splitID[0:len(splitID)-2], "/"), "/")
	_ = d.Set("folder", folder)

	d.SetId(generateCredentialID(folder, name))
	return ret, nil
}
