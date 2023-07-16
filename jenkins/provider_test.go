package jenkins

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// Deprecated: Use testAcc6Provider
	testAccProvider *schema.Provider

	testAcc6Provider provider.Provider
	testAccProviders map[string]func() (tfprotov6.ProviderServer, error)
	testAccClient    *jenkinsAdapter
)

func init() {
	testAccProvider = Provider()
	upgradedSdkProvider, err := tf5to6server.UpgradeServer(context.Background(), testAccProvider.GRPCProvider) //nolint:staticcheck
	if err != nil {
		log.Fatal(err)
	}

	testAcc6Provider = New()
	testAccProviders = map[string]func() (tfprotov6.ProviderServer, error){
		"jenkins": func() (tfprotov6.ProviderServer, error) {
			return tf6muxserver.NewMuxServer(context.Background(),
				providerserver.NewProtocol6(testAcc6Provider),
				func() tfprotov6.ProviderServer {
					return upgradedSdkProvider
				},
			)
		},
	}

	config := Config{
		ServerURL: os.Getenv("JENKINS_URL"),
		Username:  os.Getenv("JENKINS_USERNAME"),
		Password:  os.Getenv("JENKINS_PASSWORD"),
	}
	testAccClient = newJenkinsClient(&config)
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
