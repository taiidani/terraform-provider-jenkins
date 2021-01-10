package jenkins

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJenkinsJobDataSource_basic(t *testing.T) {
	xml, _ := ioutil.ReadFile("resource_jenkins_job_test.xml")
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_job foo {
					name = "tf-acc-test-%s"
					template = <<EOT
				  `+string(xml)+`
				  EOT

					parameters = {
						description = "Acceptance testing Jenkins provider"
					}
				}

				data jenkins_job foo {
					name = jenkins_job.foo.name
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_job.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_job.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_job.foo", "name", "tf-acc-test-"+randString),
				),
			},
		},
	})
}

func TestAccJenkinsJobDataSource_nested(t *testing.T) {
	xml, _ := ioutil.ReadFile("resource_jenkins_job_test.xml")
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_folder foo {
					name = "tf-acc-test-%s"
				}

				resource jenkins_job sub {
					name = "subfolder"
					folder = jenkins_folder.foo.id
					template = <<EOT
				  `+string(xml)+`
				  EOT

					parameters = {
						description = "Acceptance testing Jenkins provider"
					}
				}

				data jenkins_job sub {
					name = jenkins_job.sub.name
					folder = jenkins_job.sub.folder
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_folder.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_job.sub", "id", "/job/tf-acc-test-"+randString+"/job/subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_job.sub", "name", "subfolder"),
					resource.TestCheckResourceAttr("data.jenkins_job.sub", "folder", "/job/tf-acc-test-"+randString),
				),
			},
		},
	})
}
