// Package main defines the Jenkins Terraform Provider entrypoint.
//
// This file and the folder structure within the `jenkins/` subfolder conform to the Terraform provider expectations and
// best practices at https://www.terraform.io/docs/extend/. Please see the generated documentation at
// https://registry.terraform.io/providers/taiidani/jenkins for how to use the provider within Terraform itself.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/taiidani/terraform-provider-jenkins/jenkins"
)

const providerAddress = "registry.terraform.io/taiidani/jenkins"

func main() {
	ctx := context.Background()

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// Load the legacy sdkv2 provider using the v6 protocol.
	// This should be removed after the migration is complete.
	upgradedSdkProvider, err := tf5to6server.UpgradeServer(ctx, jenkins.Provider().GRPCProvider) //nolint:staticcheck
	if err != nil {
		log.Fatal(err)
	}

	// Mux both the legacy sdkv2 provider and the new Framework provider.
	// This will be unnecessary after the migration is complete.
	muxServer, err := tf6muxserver.NewMuxServer(ctx,
		providerserver.NewProtocol6(jenkins.New()),
		func() tfprotov6.ProviderServer {
			return upgradedSdkProvider
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}
	err = tf6server.Serve(
		providerAddress,
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
