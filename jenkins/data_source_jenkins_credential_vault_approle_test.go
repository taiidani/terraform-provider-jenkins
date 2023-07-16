package jenkins

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJenkinsCredentialVaultAppRoleDataSource_basic(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_credential_vault_approle foo {
					name = "tf-acc-test-%s"
					description = "Terraform acceptance tests %s"
					role_id = "foo"
					secret_id = "bar"
				}

				data jenkins_credential_vault_approle foo {
					name   = jenkins_credential_vault_approle.foo.name
					domain = "_"
				}`, randString, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_vault_approle.foo", "id", "/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.foo", "id", "/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.foo", "name", "tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.foo", "description", "Terraform acceptance tests "+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.foo", "role_id", "foo"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialVaultAppRoleDataSource_nested(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_folder foo {
					name = "tf-acc-test-%s"
				}

				resource jenkins_credential_vault_approle sub {
					name = "subfolder"
					folder = jenkins_folder.foo.id
					description = "Terraform acceptance tests %s"
					role_id = "foo"
					secret_id = "bar"
				}

				data jenkins_credential_vault_approle sub {
					name   = jenkins_credential_vault_approle.sub.name
					domain = "_"
					folder = jenkins_credential_vault_approle.sub.folder
				}`, randString, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_folder.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_credential_vault_approle.sub", "id", "/job/tf-acc-test-"+randString+"/subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.sub", "name", "subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.sub", "folder", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.sub", "description", "Terraform acceptance tests "+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.sub", "role_id", "foo"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialVaultAppRoleDataSource_basic_namespaced(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_credential_vault_approle foo {
					name = "tf-acc-test-%s"
					description = "Terraform acceptance tests %s"
					role_id = "foo"
					secret_id = "bar"
					namespace = "my-namespace"
				}

				data jenkins_credential_vault_approle foo {
					name   = jenkins_credential_vault_approle.foo.name
					domain = "_"
				}`, randString, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_vault_approle.foo", "id", "/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.foo", "id", "/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.foo", "name", "tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.foo", "description", "Terraform acceptance tests "+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.foo", "role_id", "foo"),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.foo", "namespace", "my-namespace"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialVaultAppRoleDataSource_nested_namespaced(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_folder foo {
					name = "tf-acc-test-%s"
				}

				resource jenkins_credential_vault_approle sub {
					name = "subfolder"
					folder = jenkins_folder.foo.id
					description = "Terraform acceptance tests %s"
					role_id = "foo"
					secret_id = "bar"
					namespace = "my-namespace"
				}

				data jenkins_credential_vault_approle sub {
					name   = jenkins_credential_vault_approle.sub.name
					domain = "_"
					folder = jenkins_credential_vault_approle.sub.folder
				}`, randString, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_folder.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_credential_vault_approle.sub", "id", "/job/tf-acc-test-"+randString+"/subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.sub", "name", "subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.sub", "folder", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.sub", "description", "Terraform acceptance tests "+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.sub", "role_id", "foo"),
					resource.TestCheckResourceAttr("data.jenkins_credential_vault_approle.sub", "namespace", "my-namespace"),
				),
			},
		},
	})
}
