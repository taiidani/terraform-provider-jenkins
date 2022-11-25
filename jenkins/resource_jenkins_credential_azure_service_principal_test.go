package jenkins

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJenkinsCredentialAzureServicePrincipal_basic(t *testing.T) {
	var cred AzureServicePrincipalCredentials

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckJenkinsCredentialAzureServicePrincipalDestroy,
			testAccCheckJenkinsFolderDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: `
				resource "jenkins_folder" "example" {
					name        = "azure-service-principal-test-folder"
					description = "A sample folder"
				
					security {
						permissions = [
							"com.cloudbees.plugins.credentials.CredentialsProvider.Create:anonymous",
							"com.cloudbees.plugins.credentials.CredentialsProvider.Delete:authenticated",
							"hudson.model.Item.Cancel:authenticated",
							"hudson.model.Item.Discover:anonymous",
						]
					}
				}
				  
				resource jenkins_credential_azure_service_principal foo {
					name = "bla"
					folder = jenkins_folder.example.id
					description = "blabla"
					subscription_id = "123"
					client_id = "123"
					client_secret = "super-secret"
					tenant = "456"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialAzureServicePrincipalExists("jenkins_credential_azure_service_principal.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_azure_service_principal.foo", "id", "/job/azure-service-principal-test-folder/bla"),
				),
			},
		},
	})
}

func testAccCheckJenkinsCredentialAzureServicePrincipalExists(resourceName string, cred *AzureServicePrincipalCredentials) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(jenkinsClient)
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf(resourceName + " not found")
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID is not set")
		}

		manager := client.Credentials()
		manager.Folder = formatFolderName(rs.Primary.Attributes["folder"])
		err := manager.GetSingle(ctx, rs.Primary.Attributes["domain"], rs.Primary.Attributes["name"], cred)
		if err != nil {
			return fmt.Errorf("Unable to retrieve credentials for %s - %s: %w", rs.Primary.Attributes["folder"], rs.Primary.Attributes["name"], err)
		}

		return nil
	}
}

func testAccCheckJenkinsCredentialAzureServicePrincipalDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(jenkinsClient)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jenkins_credential_azure_service_principal" {
			continue
		} else if _, ok := rs.Primary.Meta["name"]; !ok {
			continue
		}

		cred := AzureServicePrincipalCredentials{}
		manager := client.Credentials()
		manager.Folder = formatFolderName(rs.Primary.Meta["folder"].(string))
		err := manager.GetSingle(ctx, rs.Primary.Meta["domain"].(string), rs.Primary.Meta["name"].(string), &cred)
		if err == nil {
			return fmt.Errorf("Credentials still exists: %s - %s", rs.Primary.Attributes["folder"], rs.Primary.Attributes["name"])
		}
	}

	return nil
}
