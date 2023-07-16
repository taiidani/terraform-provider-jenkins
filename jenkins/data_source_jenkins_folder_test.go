package jenkins

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJenkinsFolderDataSource_basic(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_folder foo {
				  name = "tf-acc-test-%s"
				  display_name = "TF Acceptance Test %s"
				  description = "Terraform acceptance tests %s"
				}

				data jenkins_folder foo {
					name = jenkins_folder.foo.name
				}`, randString, randString, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_folder.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_folder.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_folder.foo", "name", "tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_folder.foo", "description", "Terraform acceptance tests "+randString),
					resource.TestCheckResourceAttr("data.jenkins_folder.foo", "display_name", "TF Acceptance Test "+randString),
				),
			},
		},
	})
}

func TestAccJenkinsFolderDataSource_nested(t *testing.T) {
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

				resource jenkins_folder sub {
					name = "subfolder"
					display_name = "TF Acceptance Test %s"
					folder = jenkins_folder.foo.id
					description = "Terraform acceptance tests %s"
				}

				data jenkins_folder sub {
					name = jenkins_folder.sub.name
					folder = jenkins_folder.sub.folder
				}`, randString, randString, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_folder.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_folder.sub", "id", "/job/tf-acc-test-"+randString+"/job/subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_folder.sub", "name", "subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_folder.sub", "folder", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_folder.sub", "description", "Terraform acceptance tests "+randString),
					resource.TestCheckResourceAttr("data.jenkins_folder.sub", "display_name", "TF Acceptance Test "+randString),
				),
			},
		},
	})
}
