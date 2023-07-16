package jenkins

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccJenkinsCredentialAzureServicePrincipal_basic(t *testing.T) {
	var cred AzureServicePrincipalCredentials

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckJenkinsCredentialAzureServicePrincipalDestroy,
			testAccCheckJenkinsFolderDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: `
				resource jenkins_credential_azure_service_principal foo {
					name = "bla"
					description = "blabla"
					subscription_id = "12345"
					client_id = "123"
					client_secret = "super-secret"
					tenant = "456"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialAzureServicePrincipalExists("jenkins_credential_azure_service_principal.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "id", "/bla"),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "description", "blabla"),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "subscription_id", "12345"),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "client_id", "123"),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "client_secret", "super-secret"),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "tenant", "456"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialAzureServicePrincipal_folder(t *testing.T) {
	var cred AzureServicePrincipalCredentials
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckJenkinsCredentialAzureServicePrincipalDestroy,
			testAccCheckJenkinsFolderDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "jenkins_folder" "example" {
					name        = "azure-service-principal-test-folder-%s"
					description = "A sample folder"
				}

				resource jenkins_credential_azure_service_principal foo {
					name = "bla"
					folder = jenkins_folder.example.id
					description = "blabla"
					subscription_id = "123"
					client_id = "123"
					client_secret = "super-secret"
					tenant = "456"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialAzureServicePrincipalExists("jenkins_credential_azure_service_principal.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "id", "/job/azure-service-principal-test-folder-"+randString+"/bla"),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "folder", "/job/azure-service-principal-test-folder-"+randString),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "description", "blabla"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "jenkins_folder" "example" {
					name        = "azure-service-principal-test-folder-%s"
					description = "A sample folder"
				}

				resource jenkins_credential_azure_service_principal foo {
					name = "bla"
					folder = jenkins_folder.example.id
					description = "blablablabla"
					subscription_id = "123"
					client_id = "123"
					client_secret = "super-secret"
					tenant = "456"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialAzureServicePrincipalExists("jenkins_credential_azure_service_principal.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "id", "/job/azure-service-principal-test-folder-"+randString+"/bla"),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "folder", "/job/azure-service-principal-test-folder-"+randString),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "description", "blablablabla"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialAzureServicePrincipal_folder_certificate(t *testing.T) {
	var cred AzureServicePrincipalCredentials
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckJenkinsCredentialAzureServicePrincipalDestroy,
			testAccCheckJenkinsFolderDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "jenkins_folder" "example" {
					name        = "azure-service-principal-test-folder-%s"
					description = "A sample folder"
				}

				resource jenkins_credential_azure_service_principal foo {
					name = "bla"
					folder = jenkins_folder.example.id
					description = "blabla"
					subscription_id = "123"
					client_id = "123"
					certificate_id = "my-cred-id/123"
					tenant = "456"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialAzureServicePrincipalExists("jenkins_credential_azure_service_principal.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "id", "/job/azure-service-principal-test-folder-"+randString+"/bla"),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "certificate_id", "my-cred-id/123"),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "jenkins_folder" "example" {
					name        = "azure-service-principal-test-folder-%s"
					description = "A sample folder"
				}

				resource jenkins_credential_azure_service_principal foo {
					name = "bla"
					folder = jenkins_folder.example.id
					description = "blablablabla"
					subscription_id = "123"
					client_id = "123"
					client_secret = "super-secret"
					tenant = "456"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialAzureServicePrincipalExists("jenkins_credential_azure_service_principal.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "id", "/job/azure-service-principal-test-folder-"+randString+"/bla"),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "client_secret", "super-secret"),
				),
			},
		},
	})
}

func testAccCheckJenkinsCredentialAzureServicePrincipalExists(resourceName string, cred *AzureServicePrincipalCredentials) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf(resourceName + " not found")
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID is not set")
		}

		manager := testAccClient.Credentials()
		manager.Folder = formatFolderName(rs.Primary.Attributes["folder"])
		err := manager.GetSingle(ctx, rs.Primary.Attributes["domain"], rs.Primary.Attributes["name"], cred)
		if err != nil {
			return fmt.Errorf("Unable to retrieve credentials for %s - %s: %w", rs.Primary.Attributes["folder"], rs.Primary.Attributes["name"], err)
		}

		return nil
	}
}

func testAccCheckJenkinsCredentialAzureServicePrincipalDestroy(s *terraform.State) error {
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jenkins_credential_azure_service_principal" {
			continue
		} else if _, ok := rs.Primary.Meta["name"]; !ok {
			continue
		}

		cred := AzureServicePrincipalCredentials{}
		manager := testAccClient.Credentials()
		manager.Folder = formatFolderName(rs.Primary.Meta["folder"].(string))
		err := manager.GetSingle(ctx, rs.Primary.Meta["domain"].(string), rs.Primary.Meta["name"].(string), &cred)
		if err == nil {
			return fmt.Errorf("Credentials still exists: %s - %s", rs.Primary.Attributes["folder"], rs.Primary.Attributes["name"])
		}
	}

	return nil
}
