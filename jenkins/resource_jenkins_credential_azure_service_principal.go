package jenkins

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// VaultAppRoleCredentials struct representing credential for storing Vault AppRole role id and secret id
type AzureServicePrincipalCredentials struct {
	XMLName     xml.Name                             `xml:"com.microsoft.azure.util.AzureCredentials"`
	Id          string                               `xml:"id"`
	Scope       string                               `xml:"scope"`
	Description string                               `xml:"description"`
	Data        AzureServicePrincipalCredentialsData `xml:"data"`
}

type AzureServicePrincipalCredentialsData struct {
	SubscriptionId          string `xml:"subscriptionId"`
	ClientId                string `xml:"clientId"`
	ClientSecret            string `xml:"clientSecret"`
	CertificateId           string `xml:"certificateId"`
	Tenant                  string `xml:"tenant"`
	AzureEnvironmentName    string `xml:"azureEnvironmentName"`
	ServiceManagementURL    string `xml:"serviceManagementURL"`
	AuthenticationEndpoint  string `xml:"authenticationEndpoint"`
	ResourceManagerEndpoint string `xml:"resourceManagerEndpoint"`
	GraphEndpoint           string `xml:"graphEndpoint"`
}

func resourceJenkinsCredentialAzureServicePrincipal() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJenkinsCredentialAzureServicePrincipalCreate,
		ReadContext:   resourceJenkinsCredentialAzureServicePrincipalRead,
		UpdateContext: resourceJenkinsCredentialAzureServicePrincipalUpdate,
		DeleteContext: resourceJenkinsCredentialAzureServicePrincipalDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceJenkinsCredentialAzureServicePrincipalImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The credential id of the Azure serivce principal credential created in Jenkins.",
				Required:    true,
				ForceNew:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "The Jenkins domain that the credentials will be added to.",
				Optional:    true,
				Default:     "_",
				// In-place updates should be possible, but gojenkins does not support move operations
				ForceNew: true,
			},
			"folder": {
				Type:        schema.TypeString,
				Description: "The Jenkins folder that the credentials will be added to.",
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
				Description: "An optional description to help tell similar credentials apart.",
				Optional:    true,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Description: "The Azure subscription id.",
				Required:    true,
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "The client id (application id) of the Azure Service Principal.",
				Required:    true,
			},
			"client_secret": {
				Type:         schema.TypeString,
				Description:  "The client secret of the Azure Service Principal. Cannot be used with certificate_id.",
				Sensitive:    true,
				Optional:     true,
				ExactlyOneOf: []string{"client_secret", "certificate_id"},
			},
			"certificate_id": {
				Type:         schema.TypeString,
				Description:  "The certificate reference of the Azure Service Principal, pointing to a Jenkins certificate credential. Cannot be used with client_secret.",
				Sensitive:    true,
				Optional:     true,
				ExactlyOneOf: []string{"client_secret", "certificate_id"},
			},
			"tenant": {
				Type:        schema.TypeString,
				Description: "The Azure Tenant ID of the Azure Service Principal.",
				Required:    true,
			},
			"azure_environment_name": {
				Type:         schema.TypeString,
				Description:  `The Azure Cloud enviroment name. Allowed values are "Azure", "Azure China", "Azure Germany", "Azure US Government".`,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Azure", "Azure China", "Azure Germany", "Azure US Government"}, false),
				Default:      "Azure",
			},
			"service_management_url": {
				Type:        schema.TypeString,
				Description: "Override the Azure management endpoint URL for the selected Azure environment.",
				Optional:    true,
				Default:     "",
			},
			"authentication_endpoint": {
				Type:        schema.TypeString,
				Description: "Override the Azure Active Directory endpoint for the selected Azure environment.",
				Optional:    true,
				Default:     "",
			},
			"resource_manager_endpoint": {
				Type:        schema.TypeString,
				Description: "Override the Azure resource manager endpoint URL for the selected Azure environment.",
				Optional:    true,
				Default:     "",
			},
			"graph_endpoint": {
				Type:        schema.TypeString,
				Description: "Override the Azure graph endpoint URL for the selected Azure environment.",
				Optional:    true,
				Default:     "",
			},
		},
	}
}

func resourceJenkinsCredentialAzureServicePrincipalCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(jenkinsClient)
	cm := client.Credentials()
	cm.Folder = formatFolderName(d.Get("folder").(string))
	// return diag.FromErr(fmt.Errorf("invalid folder name '%s', '%s'", cm.Folder, d.Get("folder").(string)))
	// Validate that the folder exists
	if err := folderExists(ctx, client, cm.Folder); err != nil {
		return diag.FromErr(fmt.Errorf("invalid folder name '%s' specified: %w", cm.Folder, err))
	}

	credData := AzureServicePrincipalCredentialsData{
		SubscriptionId:          d.Get("subscription_id").(string),
		ClientId:                d.Get("client_id").(string),
		ClientSecret:            d.Get("client_secret").(string),
		CertificateId:           d.Get("certificate_id").(string),
		Tenant:                  d.Get("tenant").(string),
		AzureEnvironmentName:    d.Get("azure_environment_name").(string),
		ServiceManagementURL:    d.Get("service_management_url").(string),
		AuthenticationEndpoint:  d.Get("authentication_endpoint").(string),
		ResourceManagerEndpoint: d.Get("resource_manager_endpoint").(string),
		GraphEndpoint:           d.Get("graph_endpoint").(string),
	}

	cred := AzureServicePrincipalCredentials{
		Id:          d.Get("name").(string),
		Scope:       d.Get("scope").(string),
		Description: d.Get("description").(string),
		Data:        credData,
	}

	domain := d.Get("domain").(string)
	err := cm.Add(ctx, domain, cred)
	if err != nil {
		return diag.Errorf("Could not create Azure service principal credentials: %s", err)
	}

	d.SetId(generateCredentialID(d.Get("folder").(string), cred.Id))
	return resourceJenkinsCredentialAzureServicePrincipalRead(ctx, d, meta)
}

func resourceJenkinsCredentialAzureServicePrincipalRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cm := meta.(jenkinsClient).Credentials()
	cm.Folder = formatFolderName(d.Get("folder").(string))

	cred := AzureServicePrincipalCredentials{}
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

		return diag.Errorf("Could not read Azure service principal credentials: %s", err)
	}

	d.SetId(generateCredentialID(d.Get("folder").(string), cred.Id))
	_ = d.Set("description", cred.Description)
	_ = d.Set("scope", cred.Scope)

	// NOTE: We are NOT setting the password here, as the password returned by GetSingle is garbage
	// Password only applies to Create/Update operations if the "password" property is non-empty

	return nil
}

func resourceJenkinsCredentialAzureServicePrincipalUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cm := meta.(jenkinsClient).Credentials()
	cm.Folder = formatFolderName(d.Get("folder").(string))

	domain := d.Get("domain").(string)

	credData := AzureServicePrincipalCredentialsData{
		SubscriptionId:          d.Get("subscription_id").(string),
		ClientId:                d.Get("client_id").(string),
		Tenant:                  d.Get("tenant").(string),
		AzureEnvironmentName:    d.Get("azure_environment_name").(string),
		ServiceManagementURL:    d.Get("service_management_url").(string),
		AuthenticationEndpoint:  d.Get("authentication_endpoint").(string),
		ResourceManagerEndpoint: d.Get("resource_manager_endpoint").(string),
		GraphEndpoint:           d.Get("graph_endpoint").(string),
	}

	cred := AzureServicePrincipalCredentials{
		Id:          d.Get("name").(string),
		Scope:       d.Get("scope").(string),
		Description: d.Get("description").(string),
		Data:        credData,
	}

	// Only enforce the password if it is non-empty
	if d.Get("client_secret").(string) != "" {
		cred.Data.ClientSecret = d.Get("client_secret").(string)
	}

	if d.Get("certificate_id").(string) != "" {
		cred.Data.ClientId = d.Get("certificate_id").(string)
	}

	err := cm.Update(ctx, domain, d.Get("name").(string), &cred)
	if err != nil {
		return diag.Errorf("Could not update Azure Service Principal credentials: %s", err)
	}

	d.SetId(generateCredentialID(d.Get("folder").(string), cred.Id))
	return resourceJenkinsCredentialAzureServicePrincipalRead(ctx, d, meta)
}

func resourceJenkinsCredentialAzureServicePrincipalDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceJenkinsCredentialAzureServicePrincipalImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	ret := []*schema.ResourceData{d}

	splitID := strings.Split(d.Id(), "/")
	if len(splitID) < 2 {
		return ret, fmt.Errorf("Import ID was improperly formatted. Imports need to be in the format \"[<folder>/]<domain>/<name>\"")
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
