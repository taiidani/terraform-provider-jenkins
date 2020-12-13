package jenkins

import (
	"context"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	jenkins "github.com/bndr/gojenkins"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJenkinsJob_basic(t *testing.T) {
	xml, _ := ioutil.ReadFile("resource_jenkins_job_test.xml")
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJenkinsJobDestroy,
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
				}`, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_job.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_job.foo", "name", "tf-acc-test-"+randString),
				),
			},
		},
	})
}

func TestAccJenkinsJob_nested(t *testing.T) {
	xml, _ := ioutil.ReadFile("resource_jenkins_job_test.xml")
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

				resource jenkins_job sub {
					name = "subfolder"
					folder = jenkins_folder.foo.id
					template = <<EOT
				  `+string(xml)+`
				  EOT

					parameters = {
						description = "Acceptance testing Jenkins provider"
					}
				}`, randString, randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_folder.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_folder.foo", "name", "tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_job.sub", "id", "/job/tf-acc-test-"+randString+"/job/subfolder"),
					resource.TestCheckResourceAttr("jenkins_job.sub", "name", "subfolder"),
					resource.TestCheckResourceAttr("jenkins_job.sub", "folder", "/job/tf-acc-test-"+randString),
				),
			},
		},
	})
}

func testAccCheckJenkinsJobDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(jenkinsClient)

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

func Test_resourceJenkinsJobDelete(t *testing.T) {
	type args struct {
		ctx  context.Context
		d    *schema.ResourceData
		meta jenkinsClient
	}
	tests := []struct {
		name string
		args args
		want diag.Diagnostics
	}{
		{
			name: "success",
			args: args{
				meta: &mockJenkinsClient{
					mockDeleteJobInFolder: func(name string, parentIDs ...string) (bool, error) {
						return true, nil
					},
				},
				d: schema.TestResourceDataRaw(t, resourceJenkinsJob().Schema, map[string]interface{}{}),
			},
		},
		{
			name: "error",
			args: args{
				meta: &mockJenkinsClient{
					mockDeleteJobInFolder: func(name string, parentIDs ...string) (bool, error) {
						return false, fmt.Errorf("omg")
					},
				},
				d: schema.TestResourceDataRaw(t, resourceJenkinsJob().Schema, map[string]interface{}{}),
			},
			want: diag.Diagnostics{
				diag.Diagnostic{Summary: "omg"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceJenkinsJobDelete(tt.args.ctx, tt.args.d, tt.args.meta); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceJenkinsJobDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceJenkinsJobRead(t *testing.T) {
	type args struct {
		ctx  context.Context
		d    *schema.ResourceData
		meta jenkinsClient
	}
	tests := []struct {
		name string
		args args
		want diag.Diagnostics
	}{
		{
			name: "missing-job",
			args: args{
				meta: &mockJenkinsClient{
					mockGetJob: func(id string, parentIDs ...string) (*jenkins.Job, error) {
						return nil, fmt.Errorf("404")
					},
				},
				d: schema.TestResourceDataRaw(t, resourceJenkinsJob().Schema, map[string]interface{}{}),
			},
		},
		{
			name: "error-job",
			args: args{
				meta: &mockJenkinsClient{
					mockGetJob: func(id string, parentIDs ...string) (*jenkins.Job, error) {
						return nil, fmt.Errorf("500")
					},
				},
				d: schema.TestResourceDataRaw(t, resourceJenkinsJob().Schema, map[string]interface{}{}),
			},
			want: diag.Diagnostics{
				diag.Diagnostic{Summary: "jenkins::read - Job \"\" does not exist: 500"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceJenkinsJobRead(tt.args.ctx, tt.args.d, tt.args.meta); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceJenkinsJobRead() = %v, want %v", got, tt.want)
			}
		})
	}
}
