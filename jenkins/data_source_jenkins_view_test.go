package jenkins

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJenkinsViewDataSource_basic(t *testing.T) {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource jenkins_view foo {
				  name = "tf-acc-test-%s"
				}

				data jenkins_view foo {
					name = jenkins_view.foo.name
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_view.foo", "id", "tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_view.foo", "id", "tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("data.jenkins_view.foo", "name", "tf-acc-test-"+randString),
				),
			},
		},
	})
}
