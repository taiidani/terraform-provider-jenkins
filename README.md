# Jenkins Terraform Provider

[![test](https://github.com/taiidani/terraform-provider-jenkins/actions/workflows/test.yml/badge.svg)](https://github.com/taiidani/terraform-provider-jenkins/actions/workflows/test.yml)

- Website: https://developer.hashicorp.com/terraform
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)

This is a community provider and is not supported by HashiCorp or IBM.

> [!WARNING]
> This project is seeking new ownership, and is only supported as free time allows. I will continue to process Dependabot PRs but am no longer accepting Issues and Pull Requests at this time. See https://github.com/taiidani/terraform-provider-jenkins/issues/274 for discussion on new ownership. If a blessed fork is not established by January 2027 I will be archiving this repository without one.

## Installation

This provider has been published to the Terraform Registry at https://registry.terraform.io/providers/taiidani/jenkins. Please visit the registry for documentation and installation instructions.

## Contributors

The scope of the provider covers the entire (extendable) Jenkins API provided that the https://github.com/bndr/gojenkins client library supports it. I accept submissions for functionality outside of Jenkins Core but expect that the plugin(s) required are _clearly stated_ in the documentation. See [jenkins_credential_vault_approle](https://registry.terraform.io/providers/taiidani/jenkins/latest/docs/resources/credential_vault_approle) for an example of this. I can only support these extensions as much as my own ability to test them allows -- Your Mileage May Vary.

## Developing the Provider

Working on this provider requires the following:

* [Mise](https://mise.jdx.dev/) (for managing tools and tasks)
* [Docker Engine](https://docs.docker.com/engine/install/) 20.10+ (for running acceptance tests)

You will also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `${GOPATH}/bin` to your `$PATH`.

To compile the provider, run `make`. This will install the provider into your GOPATH and print instructions on registering it into your system.

There are many tests available for the provider. Run `mise run` and examine the options under the "test" namespace to select one.

```sh
$ mise run
```

In order to run the full suite of Acceptance tests, run `mise run test:acceptance`. These tests require Docker to be installed on the machine that runs them, and do not create any remote resources.

```sh
$ mise run test:acceptance
```

In order to run the integration tests, navigate to the tests folder and run [terraform test](https://developer.hashicorp.com/terraform/language/tests) within it. These tests require Docker to be installed on the machine that runs them, and do not create any remote resources.

```sh
$ cd integration
$ terraform init
$ terraform test
```

When changing a data source or resource, you may need to update the documentation. This documentation is automatically rendered by https://github.com/hashicorp/terraform-plugin-docs. To trigger a render, execute:

```sh
$ mise run generate
```

## Attribution

This provider design was originally inspired from the work at [dihedron/terraform-provider-jenkins](https://github.com/dihedron/terraform-provider-jenkins).
