package jenkins

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccJenkinsCredentialVaultAppRole_basic(t *testing.T) {
	var cred VaultAppRoleCredentials

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccCheckJenkinsCredentialVaultAppRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
				resource jenkins_credential_vault_approle foo {
				  name = "test-approle"
				  role_id = "foo"
				  secret_id = "bar"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_vault_approle.foo", "id", "/test-approle"),
					testAccCheckJenkinsCredentialVaultAppRoleExists("jenkins_credential_vault_approle.foo", &cred),
				),
			},
			{
				// Update by adding description
				Config: `
				resource jenkins_credential_vault_approle foo {
				  name = "test-approle"
				  description = "new-description"
				  role_id = "foo"
				  secret_id = "bar"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialVaultAppRoleExists("jenkins_credential_vault_approle.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_vault_approle.foo", "description", "new-description"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialVaultAppRole_basic_namespaced(t *testing.T) {
	var cred VaultAppRoleCredentials

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccCheckJenkinsCredentialVaultAppRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
				resource jenkins_credential_vault_approle foo {
				  name = "test-approle"
				  role_id = "foo"
				  secret_id = "bar"
				  namespace = "my-namespace"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_vault_approle.foo", "id", "/test-approle"),
					testAccCheckJenkinsCredentialVaultAppRoleExists("jenkins_credential_vault_approle.foo", &cred),
				),
			},
			{
				// Update by adding description
				Config: `
				resource jenkins_credential_vault_approle foo {
				  name = "test-approle"
				  description = "new-description"
				  namespace = "my-namespace"
				  role_id = "foo"
				  secret_id = "bar"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialVaultAppRoleExists("jenkins_credential_vault_approle.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_vault_approle.foo", "description", "new-description"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialVaultAppRole_folder_namespaced(t *testing.T) {
	var cred VaultAppRoleCredentials
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckJenkinsCredentialVaultAppRoleDestroy,
			testAccCheckJenkinsFolderDestroy,
		),
		Steps: []resource.TestStep{
			{
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

				resource jenkins_credential_vault_approle foo {
				  name = "test-approle"
				  folder = jenkins_folder.foo_sub.id
				  role_id = "foo"
				  secret_id = "bar"
				  namespace = "my-namespace"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_vault_approle.foo", "id", "/job/tf-acc-test-"+randString+"/job/subfolder/test-approle"),
					testAccCheckJenkinsCredentialVaultAppRoleExists("jenkins_credential_vault_approle.foo", &cred),
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

				resource jenkins_credential_vault_approle foo {
				  name = "test-approle"
				  folder = jenkins_folder.foo_sub.id
				  description = "new-description"
				  role_id = "foo"
				  secret_id = "bar"
				  namespace = "my-namespace"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialVaultAppRoleExists("jenkins_credential_vault_approle.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_vault_approle.foo", "description", "new-description"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialVaultAppRole_folder(t *testing.T) {
	var cred VaultAppRoleCredentials
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckJenkinsCredentialVaultAppRoleDestroy,
			testAccCheckJenkinsFolderDestroy,
		),
		Steps: []resource.TestStep{
			{
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

				resource jenkins_credential_vault_approle foo {
				  name = "test-approle"
				  folder = jenkins_folder.foo_sub.id
				  role_id = "foo"
				  secret_id = "bar"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_vault_approle.foo", "id", "/job/tf-acc-test-"+randString+"/job/subfolder/test-approle"),
					testAccCheckJenkinsCredentialVaultAppRoleExists("jenkins_credential_vault_approle.foo", &cred),
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

				resource jenkins_credential_vault_approle foo {
				  name = "test-approle"
				  folder = jenkins_folder.foo_sub.id
				  description = "new-description"
				  role_id = "foo"
				  secret_id = "bar"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinsCredentialVaultAppRoleExists("jenkins_credential_vault_approle.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_vault_approle.foo", "description", "new-description"),
				),
			},
		},
	})
}

func testAccCheckJenkinsCredentialVaultAppRoleExists(resourceName string, cred *VaultAppRoleCredentials) resource.TestCheckFunc {
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

func testAccCheckJenkinsCredentialVaultAppRoleDestroy(s *terraform.State) error {
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jenkins_credential_vault_approle" {
			continue
		} else if _, ok := rs.Primary.Meta["name"]; !ok {
			continue
		}

		cred := VaultAppRoleCredentials{}
		manager := testAccClient.Credentials()
		manager.Folder = formatFolderName(rs.Primary.Meta["folder"].(string))
		err := manager.GetSingle(ctx, rs.Primary.Meta["domain"].(string), rs.Primary.Meta["name"].(string), &cred)
		if err == nil {
			return fmt.Errorf("Credentials still exists: %s - %s", rs.Primary.Attributes["folder"], rs.Primary.Attributes["name"])
		}
	}

	return nil
}
