package jenkins

import (
	"fmt"
	"testing"

	jenkins "github.com/bndr/gojenkins"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccJenkinsJob_basic(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJenkinsJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccJenkinsJobConfig(randString),
			},
		},
	})
}

func testAccCheckJenkinsJobDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*jenkins.Jenkins)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jenkins_job" {
			continue
		}

		_, err := client.GetJob(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Job %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccJenkinsJobConfig(randString string) string {
	return fmt.Sprintf(`
resource jenkins_job foo {
  name = "tf-acc-test-%s"
  template = file("resource_jenkins_job_test.xml")

  parameters = {
	  description = "Acceptance testing Jenkins provider"
  }

}
`, randString)
}
