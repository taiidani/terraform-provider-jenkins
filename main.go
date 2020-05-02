package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/taiidani/terraform-provider-jenkins/jenkins"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: jenkins.Provider,
	})
}
