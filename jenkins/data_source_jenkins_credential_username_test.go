package jenkins

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJenkinsCredentialUsernameDataSource_basic(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_credential_username foo {
				  name = "tf-acc-test-%s"
				  description = "Terraform acceptance tests %s"
				  username = "foo"
				  password = "bar"
				}

				data jenkins_credential_username foo {
					name   = jenkins_credential_username.foo.name
					domain = "_"
				}`, randString, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_username.foo", "id", "/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_username.foo", "id", "/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_username.foo", "name", "tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_username.foo", "description", "Terraform acceptance tests "+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_username.foo", "username", "foo"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialUsernameDataSource_nested(t *testing.T) {
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

				resource jenkins_credential_username sub {
					name = "subfolder"
					folder = jenkins_folder.foo.id
					description = "Terraform acceptance tests %s"
					username = "foo"
					password = "bar"
				}

				data jenkins_credential_username sub {
					name   = jenkins_credential_username.sub.name
					domain = "_"
					folder = jenkins_credential_username.sub.folder
				}`, randString, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_folder.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_credential_username.sub", "id", "/job/tf-acc-test-"+randString+"/subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_credential_username.sub", "name", "subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_credential_username.sub", "folder", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_username.sub", "description", "Terraform acceptance tests "+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_username.sub", "username", "foo"),
				),
			},
		},
	})
}
