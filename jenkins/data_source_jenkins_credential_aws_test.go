package jenkins

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJenkinsCredentialAwsDataSource_basic(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_credential_aws foo {
					name = "tf-acc-test-%s"
					description = "Terraform acceptance tests %s"
					access_key = "foo"
				}

				data jenkins_credential_aws foo {
					name   = jenkins_credential_aws.foo.name
					domain = "`+defaultCredentialDomain+`"
				}`, randString, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_aws.foo", "id", "/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.foo", "id", "/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.foo", "name", "tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.foo", "description", "Terraform acceptance tests "+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.foo", "access_key", "foo"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialAwsDataSource_nested(t *testing.T) {
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

				resource jenkins_credential_aws sub {
					name = "subfolder"
					folder = jenkins_folder.foo.id
					description = "Terraform acceptance tests %s"
					access_key = "foo"
				}

				data jenkins_credential_aws sub {
					name   = jenkins_credential_aws.sub.name
					domain = "`+defaultCredentialDomain+`"
					folder = jenkins_credential_aws.sub.folder
				}`, randString, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_folder.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_credential_aws.sub", "id", "/job/tf-acc-test-"+randString+"/subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.sub", "name", "subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.sub", "folder", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.sub", "description", "Terraform acceptance tests "+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.sub", "access_key", "foo"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialAwsDataSource_basic_with_iam_role_arn(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_credential_aws foo {
					name = "tf-acc-test-%s"
					description = "Terraform acceptance tests %s"
					access_key = "foo"
					iam_role_arn = "my-role-arn"
				}

				data jenkins_credential_aws foo {
					name   = jenkins_credential_aws.foo.name
					domain = "`+defaultCredentialDomain+`"
				}`, randString, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_aws.foo", "id", "/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.foo", "id", "/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.foo", "name", "tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.foo", "description", "Terraform acceptance tests "+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.foo", "access_key", "foo"),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.foo", "iam_role_arn", "my-role-arn"),
				),
			},
		},
	})
}

func TestAccJenkinsCredentialAwsDataSource_nested__with_iam_role_arn(t *testing.T) {
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

				resource jenkins_credential_aws sub {
					name = "subfolder"
					folder = jenkins_folder.foo.id
					description = "Terraform acceptance tests %s"
					access_key = "foo"
					iam_role_arn = "my-role-arn"
				}

				data jenkins_credential_aws sub {
					name   = jenkins_credential_aws.sub.name
					domain = "`+defaultCredentialDomain+`"
					folder = jenkins_credential_aws.sub.folder
				}`, randString, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_folder.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_credential_aws.sub", "id", "/job/tf-acc-test-"+randString+"/subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.sub", "name", "subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.sub", "folder", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.sub", "description", "Terraform acceptance tests "+randString),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.sub", "access_key", "foo"),
					resource.TestCheckResourceAttr("data.jenkins_credential_aws.sub", "iam_role_arn", "my-role-arn"),
				),
			},
		},
	})
}
