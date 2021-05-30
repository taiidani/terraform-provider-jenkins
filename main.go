// Package main defines the Jenkins Terraform Provider entrypoint.
//
// This file and the folder structure within the `jenkins/` subfolder conform to the Terraform provider expectations and
// best practices at https://www.terraform.io/docs/extend/. Please see the generated documentation at
// https://registry.terraform.io/providers/taiidani/jenkins for how to use the provider within Terraform itself.
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/taiidani/terraform-provider-jenkins/jenkins"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return jenkins.Provider()
		},
	})
}
