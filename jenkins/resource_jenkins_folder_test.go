package jenkins

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJenkinsFolder_basic(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJenkinsFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJenkinsFolderConfig(randString),
			},
		},
	})
}

func TestAccJenkinsFolder_nested(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJenkinsFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_folder foo {
					name = "tf-acc-test-%s"
					description = "Terraform acceptance tests %s"
				}

				resource jenkins_folder sub {
					name = "${jenkins_folder.foo.name}/subfolder"
					description = "Terraform acceptance tests ${jenkins_folder.foo.name}"
				}`, randString, randString),
			},
		},
	})
}

func testAccCheckJenkinsFolderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(jenkinsClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jenkins_folder" {
			continue
		}

		_, err := client.GetJob(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Folder %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccJenkinsFolderConfig(randString string) string {
	return fmt.Sprintf(`
resource jenkins_folder foo {
  name = "tf-acc-test-%s"
  description = "Terraform acceptance tests %s"
}
`, randString, randString)
}
