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

func TestAccJenkinsCredentialSecretFile_basic(t *testing.T) {
	var cred jenkins.FileCredentials

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccCheckJenkinsCredentialSecretFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
				resource jenkins_credential_secret_file foo {
				  name = "test-secret-file"
				  filename = "secret.txt"
				  secretbytes = base64encode("This is a test.")
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_secret_file.foo", "id", "/test-secret-file"),
					testAccCheckJenkinsCredentialSecretFileExists("jenkins_credential_secret_file.foo", &cred),
				),
			},
			{
				// Update by changing secretbytes
				Config: `
				resource jenkins_credential_secret_file foo {
				  name = "test-secret-file"
				  filename = "secret.txt"
				  secretbytes = base64encode("This is a new secret content.")
				}`,
				// In comparison below I use already base64 encoded value of: VGhpcyBpcyBhIG5ldyBzZWNyZXQgY29udGVudC4=
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialSecretFileExists("jenkins_credential_secret_file.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_secret_file.foo", "secretbytes", "VGhpcyBpcyBhIG5ldyBzZWNyZXQgY29udGVudC4="),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialSecretFile_folder(t *testing.T) {
	var cred jenkins.FileCredentials
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckJenkinsCredentialSecretFileDestroy,
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

				resource jenkins_credential_secret_file foo {
					name = "test-secret-file"
					folder = jenkins_folder.foo_sub.id
					filename = "secret.txt"
					secretbytes = "VGhpcyBpcyBhIHRlc3Qu"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_secret_file.foo", "id", "/job/tf-acc-test-"+randString+"/job/subfolder/test-secret-file"),
					testAccCheckJenkinsCredentialSecretFileExists("jenkins_credential_secret_file.foo", &cred),
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

				resource jenkins_credential_secret_file foo {
					name = "test-secret-file"
					folder = jenkins_folder.foo_sub.id
					description = "new-description"
                                        filename = "secret.txt"
                                        secretbytes = "VGhpcyBpcyBhIHRlc3Qu"

				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialSecretFileExists("jenkins_credential_secret_file.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_secret_file.foo", "description", "new-description"),
				),
			},
		},
	})
}

func testAccCheckJenkinsCredentialSecretFileExists(resourceName string, cred *jenkins.FileCredentials) resource.TestCheckFunc {
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

func testAccCheckJenkinsCredentialSecretFileDestroy(s *terraform.State) error {
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jenkins_credential_secret_file" {
			continue
		} else if _, ok := rs.Primary.Meta["name"]; !ok {
			continue
		}

		cred := jenkins.FileCredentials{}
		manager := testAccClient.Credentials()
		manager.Folder = formatFolderName(rs.Primary.Meta["folder"].(string))
		err := manager.GetSingle(ctx, rs.Primary.Meta["domain"].(string), rs.Primary.Meta["name"].(string), &cred)
		if err == nil {
			return fmt.Errorf("Credentials still exists: %s - %s", rs.Primary.Attributes["folder"], rs.Primary.Attributes["name"])
		}
	}

	return nil
}
