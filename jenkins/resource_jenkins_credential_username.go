package jenkins

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	jenkins "github.com/bndr/gojenkins"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var supportedCredentialScopes = []string{"SYSTEM", "GLOBAL"}

func resourceJenkinsCredentialUsername() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJenkinsCredentialUsernameCreate,
		ReadContext:   resourceJenkinsCredentialUsernameRead,
		UpdateContext: resourceJenkinsCredentialUsernameUpdate,
		DeleteContext: resourceJenkinsCredentialUsernameDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceJenkinsCredentialUsernameImport,
		},
		Schema: map[string]*schema.Schema{
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
			"username": {
				Type:        schema.TypeString,
				Description: "The credentials user username.",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The credentials user password. If left empty will be unmanaged.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func validateCredentialScope(v interface{}, p cty.Path) diag.Diagnostics {
	for _, supported := range supportedCredentialScopes {
		if v == supported {
			return nil
		}
	}
	return diag.Errorf("Invalid scope: %s. Supported scopes are: %s", v, strings.Join(supportedCredentialScopes, ", "))
}

func resourceJenkinsCredentialUsernameCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cm := meta.(jenkinsClient).Credentials()
	cm.Folder = d.Get("folder").(string)

	// Generate a unique ID for the new credentials
	id := fmt.Sprintf("terraform-%s-%d", d.Get("username").(string), time.Now().Unix())

	cred := jenkins.UsernameCredentials{
		ID:          id,
		Scope:       d.Get("scope").(string),
		Description: d.Get("description").(string),
		Username:    d.Get("username").(string),
		Password:    d.Get("password").(string),
	}

	domain := d.Get("domain").(string)
	err := cm.Add(domain, cred)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::create - Error creating username credentials for %q: %w", id, err))
	}

	log.Printf("[DEBUG] jenkins::create - username credentials %q created", id)
	d.SetId(id)

	return resourceJenkinsCredentialUsernameRead(ctx, d, meta)
}

func resourceJenkinsCredentialUsernameRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] jenkins::read - Looking for job %q", d.Id())

	cm := meta.(jenkinsClient).Credentials()
	cm.Folder = d.Get("folder").(string)

	cred := jenkins.UsernameCredentials{}
	err := cm.GetSingle(
		d.Get("domain").(string),
		d.Id(),
		&cred,
	)

	if err != nil {
		if strings.HasSuffix(err.Error(), "404") {
			// Job does not exist
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("jenkins::read - Job %q request error: %w", d.Id(), err))
	}

	d.SetId(cred.ID)
	d.Set("scope", cred.Scope)
	d.Set("description", cred.Description)
	d.Set("username", cred.Username)

	// NOTE: We are NOT setting the password here, as the password returned by GetSingle is garbage
	// Password only applies to Create operations
	return nil
}

func resourceJenkinsCredentialUsernameUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cm := meta.(jenkinsClient).Credentials()
	cm.Folder = d.Get("folder").(string)

	domain := d.Get("domain").(string)
	cred := jenkins.UsernameCredentials{
		ID:          d.Id(),
		Scope:       d.Get("scope").(string),
		Description: d.Get("description").(string),
		Username:    d.Get("username").(string),
	}

	// Only enforce the password if it is non-empty
	if d.Get("password").(string) != "" {
		cred.Password = d.Get("password").(string)
	}

	err := cm.Update(domain, d.Id(), &cred)
	if err != nil {
		return diag.FromErr(fmt.Errorf("jenkins::update - Error updating username credentials for %q: %w", d.Id(), err))
	}

	log.Printf("[DEBUG] jenkins::update - username credentials %q updated", d.Id())
	return resourceJenkinsCredentialUsernameRead(ctx, d, meta)
}

func resourceJenkinsCredentialUsernameDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cm := meta.(jenkinsClient).Credentials()
	cm.Folder = d.Get("folder").(string)

	log.Printf("[DEBUG] jenkins::delete - Removing %q", d.Id())

	err := cm.Delete(
		d.Get("domain").(string),
		d.Id(),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] jenkins::delete - %q removed", d.Id())
	return nil
}

func resourceJenkinsCredentialUsernameImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{d}

	splitID := strings.Split(d.Id(), "/")
	if len(splitID) < 2 {
		return ret, fmt.Errorf("Import ID was improperly formatted. Imports need to be in the format \"[<folder>/]<domain>/<id>\"")
	}

	d.SetId(splitID[len(splitID)-1])
	d.Set("domain", splitID[len(splitID)-2])
	d.Set("folder", strings.Trim(strings.Join(splitID[0:len(splitID)-2], "/"), "/"))
	return ret, nil
}
