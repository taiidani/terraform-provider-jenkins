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

func TestAccJenkinsCredentialSSH_basic(t *testing.T) {
	var cred jenkins.SSHCredentials

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccCheckJenkinsCredentialSSHDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
				resource jenkins_credential_ssh foo {
				  name = "test-ssh"
				  username = "test-ssh-user"
				  privatekey = "Some fake private key"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_ssh.foo", "id", "/test-ssh"),
					testAccCheckJenkinsCredentialSSHExists("jenkins_credential_ssh.foo", &cred),
				),
			},
			{
				// Update by changing privatekey
				Config: `
				resource jenkins_credential_ssh foo {
				  name = "test-ssh"
                                  username = "test-ssh-user"
                                  privatekey = "Some other fake private key"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialSSHExists("jenkins_credential_ssh.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_ssh.foo", "privatekey", "Some other fake private key"),
				),
			},
			{
				// Update by changing adding passphrase
				Config: `
                                resource jenkins_credential_ssh foo {
                                  name = "test-ssh"
                                  username = "test-ssh-user"
                                  privatekey = "Some other fake private key"
				  passphrase = "SuperSecret"
                                }`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialSSHExists("jenkins_credential_ssh.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_ssh.foo", "passphrase", "SuperSecret"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialSSH_folder(t *testing.T) {
	var cred jenkins.SSHCredentials
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckJenkinsCredentialSSHDestroy,
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

				resource jenkins_credential_ssh foo {
				  name = "test-ssh"
				  folder = jenkins_folder.foo_sub.id
                                  username = "test-ssh-user"
                                  privatekey = "Some fake private key"

				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_ssh.foo", "id", "/job/tf-acc-test-"+randString+"/job/subfolder/test-ssh"),
					testAccCheckJenkinsCredentialSSHExists("jenkins_credential_ssh.foo", &cred),
				),
			},
			{
				// Update by changing privatekey
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

				resource jenkins_credential_ssh foo {
				  name = "test-ssh"
				  folder = jenkins_folder.foo_sub.id
                                  username = "test-ssh-user"
                                  privatekey = "Some other fake private key"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialSSHExists("jenkins_credential_ssh.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_ssh.foo", "privatekey", "Some other fake private key"),
				),
			},
			{
				// Update by changing privatekey
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

                                resource jenkins_credential_ssh foo {
                                  name = "test-ssh"
                                  folder = jenkins_folder.foo_sub.id
                                  username = "test-ssh-user"
                                  privatekey = "Some other fake private key"
				  passphrase = "SuperSecret"
                                }`, randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialSSHExists("jenkins_credential_ssh.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_ssh.foo", "passphrase", "SuperSecret"),
				),
			},
		},
	})
}

func testAccCheckJenkinsCredentialSSHExists(resourceName string, cred *jenkins.SSHCredentials) resource.TestCheckFunc {
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

func testAccCheckJenkinsCredentialSSHDestroy(s *terraform.State) error {
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jenkins_credential_ssh" {
			continue
		} else if _, ok := rs.Primary.Meta["name"]; !ok {
			continue
		}

		cred := jenkins.SSHCredentials{}
		manager := testAccClient.Credentials()
		manager.Folder = formatFolderName(rs.Primary.Meta["folder"].(string))
		err := manager.GetSingle(ctx, rs.Primary.Meta["domain"].(string), rs.Primary.Meta["name"].(string), &cred)
		if err == nil {
			return fmt.Errorf("Credentials still exists: %s - %s", rs.Primary.Attributes["folder"], rs.Primary.Attributes["name"])
		}
	}

	return nil
}
