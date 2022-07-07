package jenkins

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJenkinsPluginsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: "data jenkins_plugins foo {}",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.jenkins_plugins.foo", "id", "jenkins-data-source-plugins-id"),
				),
			},
		},
	})
}
