package jenkins

import (
	"fmt"
	"testing"

	jenkins "github.com/bndr/gojenkins"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJenkinsCredentialUsername_basic(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJenkinsCredentialUsernameDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_credential_username foo {
				  name = "test-username-%s"
				  username = "foo"
				  password = "bar"
				}`, randString),
			},
		},
	})
}

func testAccCheckJenkinsCredentialUsernameDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(jenkinsClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jenkins_credential_username" {
			continue
		}

		cred := jenkins.UsernameCredentials{}
		err := client.Credentials().GetSingle(rs.Primary.Meta["domain"].(string), rs.Primary.Meta["name"].(string), &cred)
		if err == nil {
			return fmt.Errorf("Job %s still exists", rs.Primary.ID)
		}
	}

	return nil
}
