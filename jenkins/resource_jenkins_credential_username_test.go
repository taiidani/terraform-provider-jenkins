package jenkins

import (
	"fmt"
	"testing"

	jenkins "github.com/bndr/gojenkins"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJenkinsCredentialUsername_basic(t *testing.T) {
	var cred jenkins.UsernameCredentials
	// randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJenkinsCredentialUsernameDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_credential_username foo {
				  name = "test-username"
				  username = "foo"
				  password = "bar"
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_username.foo", "id", "/test-username"),
					testAccCheckJenkinsCredentialUsernameExists("jenkins_credential_username.foo", &cred),
				),
			},
			{
				// Update by adding description
				Config: fmt.Sprintf(`
				resource jenkins_credential_username foo {
				  name = "test-username"
				  description = "new-description"
				  username = "foo"
				  password = "bar"
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_username.foo", "description", "new-description"),
					testAccCheckJenkinsCredentialUsernameExists("jenkins_credential_username.foo", &cred),
					testAccCheckJenkinsCredentialUsernameDescriptionUpdated(&cred, "new-description"),
				),
			},
		},
	})
}

func testAccCheckJenkinsCredentialUsernameExists(resourceName string, cred *jenkins.UsernameCredentials) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(jenkinsClient)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf(resourceName + " not found")
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID is not set")
		}

		err := client.Credentials().GetSingle(rs.Primary.Attributes["domain"], rs.Primary.Attributes["name"], cred)
		if err != nil {
			return fmt.Errorf("Unable to retrieve credentials for %s: %w", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccCheckJenkinsCredentialUsernameDescriptionUpdated(cred *jenkins.UsernameCredentials, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cred.Description != description {
			return fmt.Errorf("Description was not set")
		}

		return nil
	}
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
			return fmt.Errorf("Credentials %s still exists", rs.Primary.ID)
		}
	}

	return nil
}
