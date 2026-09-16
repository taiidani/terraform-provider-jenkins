package jenkins

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	jenkins "github.com/bndr/gojenkins"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	//go:embed "resource_jenkins_job_test.xml"
	testXML []byte

	//go:embed "resource_jenkins_job_test_want.xml"
	testXMLWant string
)

func TestAccJenkinsJob_basic(t *testing.T) {
	testDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(testDir, "test.xml"), testXML, 0644)
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccCheckJenkinsJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource jenkins_job foo {
	name = "tf-acc-test-%s"
	template = templatefile("%s/test.xml", {
		description = "Acceptance testing Jenkins provider"
	})
}`, randString, testDir),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_job.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_job.foo", "name", "tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_job.foo", "template", strings.TrimSpace(testXMLWant)),
				),
			},
		},
	})
}

// TestAccJenkinsJob_utf8 is a regression test for a non-ASCII character (an
// em-dash, U+2014) in a job's config.xml.
//
// The provider updates config.xml through gojenkins' Job.UpdateConfig
// ("POST /job/<name>/config.xml"). Jenkins reads that body with getReader(),
// which defaults to ISO-8859-1 and ignores the document's encoding
// declaration. gojenkins before v1.2.0 sent no charset, so the em-dash's UTF-8
// bytes decoded to an illegal XML character and the update failed with a 500
// ("Failed to persist config.xml"). v1.2.0 sends "application/xml;charset=utf-8"
// and the description round-trips.
//
// The bug is on update, not create: createItem decodes as UTF-8 regardless of
// the charset header, so step one passes even against the old client. Step two
// rewrites the description to the em-dash, exercising UpdateConfig.
func TestAccJenkinsJob_utf8(t *testing.T) {
	testDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(testDir, "test.xml"), testXML, 0644)
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	jobConfig := func(description string) string {
		return fmt.Sprintf(`
resource jenkins_job foo {
	name = "tf-acc-test-%s"
	template = templatefile("%s/test.xml", {
		description = %q
	})
}`, randString, testDir, description)
	}

	const utf8Description = "Acceptance testing Jenkins provider — UTF-8 round-trip"
	wantUTF8Template := strings.Replace(
		strings.TrimSpace(testXMLWant),
		"Acceptance testing Jenkins provider",
		utf8Description,
		1,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccCheckJenkinsJobDestroy,
		Steps: []resource.TestStep{
			{
				// Create with an ASCII description.
				Config: jobConfig("Acceptance testing Jenkins provider"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_job.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_job.foo", "template", strings.TrimSpace(testXMLWant)),
				),
			},
			{
				// Update the description to include a non-ASCII character. This
				// drives Job.UpdateConfig, the path that 500s without the
				// gojenkins charset fix.
				Config: jobConfig(utf8Description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_job.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_job.foo", "template", wantUTF8Template),
				),
			},
		},
	})
}

func TestAccJenkinsJob_nested(t *testing.T) {
	testDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(testDir, "test.xml"), testXML, 0644)
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccCheckJenkinsFolderDestroy,
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
	template = templatefile("%s/test.xml", {
		description = "Acceptance testing Jenkins provider"
	})
}`, randString, randString, testDir),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jenkins_folder.foo", "id", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_folder.foo", "name", "tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_job.sub", "id", "/job/tf-acc-test-"+randString+"/job/subfolder"),
					resource.TestCheckResourceAttr("jenkins_job.sub", "name", "subfolder"),
					resource.TestCheckResourceAttr("jenkins_job.sub", "folder", "/job/tf-acc-test-"+randString),
					resource.TestCheckResourceAttr("jenkins_job.sub", "template", strings.TrimSpace(testXMLWant)),
				),
			},
		},
	})
}

func testAccCheckJenkinsJobDestroy(s *terraform.State) error {
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jenkins_job" {
			continue
		}

		_, err := testAccClient.GetJob(ctx, rs.Primary.ID)
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
					mockDeleteJobInFolder: func(ctx context.Context, name string, parentIDs ...string) (bool, error) {
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
					mockDeleteJobInFolder: func(ctx context.Context, name string, parentIDs ...string) (bool, error) {
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
					mockGetJob: func(ctx context.Context, id string, parentIDs ...string) (*jenkins.Job, error) {
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
					mockGetJob: func(ctx context.Context, id string, parentIDs ...string) (*jenkins.Job, error) {
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
