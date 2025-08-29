package jenkins

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccJenkinscredentialAws_basic(t *testing.T) {
	var cred credentialAws

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccCheckJenkinscredentialAwsDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
				resource jenkins_credential_aws foo {
				  name = "test-aws-cred"
				  access_key = "foo"
				  secret_key = "bar"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_aws.foo", "id", "/test-aws-cred"),
					testAccCheckJenkinscredentialAwsExists("jenkins_credential_aws.foo", &cred),
				),
			},
			{
				// Update by adding description
				Config: `
				resource jenkins_credential_aws foo {
				  name = "test-aws-cred"
				  description = "new-description"
				  access_key = "foo"
				  secret_key = "bar"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinscredentialAwsExists("jenkins_credential_aws.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_aws.foo", "description", "new-description"),
				),
			},
		},
	})
}

func TestAccJenkinscredentialAws_basic_with_iam_role_arn(t *testing.T) {
	var cred credentialAws

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccCheckJenkinscredentialAwsDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
				resource jenkins_credential_aws foo {
				  name = "test-aws-cred"
				  access_key = "foo"
				  secret_key = "bar"
				  iam_role_arn = "my-role-arn"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_aws.foo", "id", "/test-aws-cred"),
					testAccCheckJenkinscredentialAwsExists("jenkins_credential_aws.foo", &cred),
				),
			},
			{
				// Update by adding description
				Config: `
				resource jenkins_credential_aws foo {
				  name = "test-aws-cred"
				  description = "new-description"
				  iam_role_arn = "my-role-arn"
				  access_key = "foo"
				  secret_key = "bar"
				}`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinscredentialAwsExists("jenkins_credential_aws.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_aws.foo", "description", "new-description"),
				),
			},
		},
	})
}

func TestAccJenkinscredentialAws_folder_with_iam_role_arn(t *testing.T) {
	var cred credentialAws
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckJenkinscredentialAwsDestroy,
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

				resource jenkins_credential_aws foo {
				  name = "test-aws-cred"
				  folder = jenkins_folder.foo_sub.id
				  access_key = "foo"
				  secret_key = "bar"
				  iam_role_arn = "my-role-arn"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_aws.foo", "id", "/job/tf-acc-test-"+randString+"/job/subfolder/test-aws-cred"),
					testAccCheckJenkinscredentialAwsExists("jenkins_credential_aws.foo", &cred),
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

				resource jenkins_credential_aws foo {
				  name = "test-aws-cred"
				  folder = jenkins_folder.foo_sub.id
				  description = "new-description"
				  access_key = "foo"
				  secret_key = "bar"
				  iam_role_arn = "my-role-arn"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinscredentialAwsExists("jenkins_credential_aws.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_aws.foo", "description", "new-description"),
					resource.TestCheckResourceAttr("jenkins_credential_aws.foo", "iam_role_arn", "my-role-arn"),
				),
			},
		},
	})
}

func TestAccJenkinscredentialAws_folder(t *testing.T) {
	var cred credentialAws
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckJenkinscredentialAwsDestroy,
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

				resource jenkins_credential_aws foo {
				  name = "test-aws-cred"
				  folder = jenkins_folder.foo_sub.id
				  access_key = "foo"
				  secret_key = "bar"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_credential_aws.foo", "id", "/job/tf-acc-test-"+randString+"/job/subfolder/test-aws-cred"),
					testAccCheckJenkinscredentialAwsExists("jenkins_credential_aws.foo", &cred),
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

				resource jenkins_credential_aws foo {
				  name = "test-aws-cred"
				  folder = jenkins_folder.foo_sub.id
				  description = "new-description"
				  access_key = "foo"
				  secret_key = "bar"
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJenkinscredentialAwsExists("jenkins_credential_aws.foo", &cred),
					resource.TestCheckResourceAttr("jenkins_credential_aws.foo", "description", "new-description"),
				),
			},
		},
	})
}

func testAccCheckJenkinscredentialAwsExists(resourceName string, cred *credentialAws) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return errors.New(resourceName + " not found")
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

func testAccCheckJenkinscredentialAwsDestroy(s *terraform.State) error {
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jenkins_credential_aws" {
			continue
		} else if _, ok := rs.Primary.Meta["name"]; !ok {
			continue
		}

		cred := credentialAws{}
		manager := testAccClient.Credentials()
		manager.Folder = formatFolderName(rs.Primary.Meta["folder"].(string))
		err := manager.GetSingle(ctx, rs.Primary.Meta["domain"].(string), rs.Primary.Meta["name"].(string), &cred)
		if err == nil {
			return fmt.Errorf("Credentials still exists: %s - %s", rs.Primary.Attributes["folder"], rs.Primary.Attributes["name"])
		}
	}

	return nil
}
