package jenkins

import (
	"context"
	"fmt"
	"testing"

	jenkins "github.com/bndr/gojenkins"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccJenkinsCredentialSecretText_basic(t *testing.T) {
	var cred jenkins.StringCredentials

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccCheckJenkinsCredentialSecretTextDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
				resource jenkins_credential_secret_text foo {
				  name = "test-secret-text"
				  secret = "very-secret"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_secret_text.foo", "id", "/test-secret-text"),
					testAccCheckJenkinsCredentialSecretTextExists("jenkins_credential_secret_text.foo", &cred),
				),
			},
			{
				// Update by adding description
				Config: `
				resource jenkins_credential_secret_text foo {
				  name = "test-secret-text"
				  description = "new-description"
				  secret = "very-secret"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialSecretTextExists("jenkins_credential_secret_text.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_secret_text.foo", "description", "new-description"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialSecretText_folder(t *testing.T) {
	var cred jenkins.StringCredentials
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckJenkinsCredentialSecretTextDestroy,
			testAccCheckJenkinsFolderDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_folder foo {
					name = "tf-acc-test-%s"
					description = "Terraform acceptance testing"
				}

				resource jenkins_folder foo_sub {
					name = "subfolder"
					folder = jenkins_folder.foo.id
					description = "Terraform acceptance testing"
				}

				resource jenkins_credential_secret_text foo {
				  name = "test-secret-text"
				  folder = jenkins_folder.foo_sub.id
				  secret = "very-secret"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_secret_text.foo", "id", "/job/tf-acc-test-"+randString+"/job/subfolder/test-secret-text"),
					testAccCheckJenkinsCredentialSecretTextExists("jenkins_credential_secret_text.foo", &cred),
				),
			},
			{
				// Update by adding description
				Config: fmt.Sprintf(`
				resource jenkins_folder foo {
					name = "tf-acc-test-%s"
					description = "Terraform acceptance testing"

					lifecycle {
						ignore_changes = [template]
					}
				}

				resource jenkins_folder foo_sub {
					name = "subfolder"
					folder = jenkins_folder.foo.id
					description = "Terraform acceptance testing"

					lifecycle {
						ignore_changes = [template]
					}
				}

				resource jenkins_credential_secret_text foo {
				  name = "test-secret-text"
				  folder = jenkins_folder.foo_sub.id
				  description = "new-description"
				  secret = "very-secret"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialSecretTextExists("jenkins_credential_secret_text.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_secret_text.foo", "description", "new-description"),
				),
			},
		},
	})
}

func testAccCheckJenkinsCredentialSecretTextExists(resourceName string, cred *jenkins.StringCredentials) resource.TestCheckFunc {
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

func testAccCheckJenkinsCredentialSecretTextDestroy(s *terraform.State) error {
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jenkins_credential_secret_text" {
			continue
		} else if _, ok := rs.Primary.Meta["name"]; !ok {
			continue
		}

		cred := jenkins.StringCredentials{}
		manager := testAccClient.Credentials()
		manager.Folder = formatFolderName(rs.Primary.Meta["folder"].(string))
		err := manager.GetSingle(ctx, rs.Primary.Meta["domain"].(string), rs.Primary.Meta["name"].(string), &cred)
		if err == nil {
			return fmt.Errorf("Credentials still exists: %s - %s", rs.Primary.Attributes["folder"], rs.Primary.Attributes["name"])
		}
	}

	return nil
}
