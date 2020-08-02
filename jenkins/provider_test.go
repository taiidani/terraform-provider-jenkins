package jenkins

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"jenkins": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := testAccProvider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("JENKINS_URL"); v == "" {
		t.Fatal("JENKINS_URL must be set for acceptance tests")
	}
	if v := os.Getenv("JENKINS_USERNAME"); v == "" {
		t.Fatal("JENKINS_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("JENKINS_PASSWORD"); v == "" {
		t.Fatal("JENKINS_PASSWORD must be set for acceptance tests")
	}
}
