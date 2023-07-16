# Terraform Provider

![Unit Tests](https://github.com/taiidani/terraform-provider-jenkins/workflows/Unit%20Tests/badge.svg)
![Acceptance Tests](https://github.com/taiidani/terraform-provider-jenkins/workflows/Acceptance%20Tests/badge.svg)

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

This is a community provider and is not supported by Hashicorp.

## Installation

This provider has been published to the Terraform Registry at https://registry.terraform.io/providers/taiidani/jenkins. Please visit the registry for documentation and installation instructions.

## Contributors

This repository is actively maintained as free time allows for me to work on it. I am accepting Issues and Pull Requests and if you are interested in helping as a maintainer I am open to the possibility. Please allow for 2-3 weeks of turnaround time on your contributions.

The scope of the provider covers the entire (extendable) Jenkins API provided that the https://github.com/bndr/gojenkins client library supports it. I accept submissions for functionality outside of Jenkins Core but expect that the plugin(s) required are _clearly stated_ in the documentation. See [jenkins_credential_vault_approle](https://registry.terraform.io/providers/taiidani/jenkins/latest/docs/resources/credential_vault_approle) for an example of this. I can only support these extensions as much as my own ability to test them allows -- Your Mileage May Vary.

## Developing the Provider

Working on this provider requires the following:

* [Terraform](https://www.terraform.io/downloads.html) 1.0+
* [Go](http://www.golang.org) (version requirements documented in the `go.mod` file)
* [Docker Engine](https://docs.docker.com/engine/install/) 20.10+ (for running acceptance tests)

You will also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `${GOPATH}/bin` to your `$PATH`.

To compile the provider, run `make`. This will install the provider into your GOPATH and print instructions on registering it into your system.

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`. These tests require Docker to be installed on the machine that runs them, and do not create any remote resources.

```sh
$ make testacc
```

When changing a data source or resource, you may need to update the documentation. This documentation is automatically rendered by https://github.com/hashicorp/terraform-plugin-docs. To trigger a render, execute:

```sh
$ make generate
```

## Attribution

This provider design was originally inspired from the work at [dihedron/terraform-provider-jenkins](https://github.com/dihedron/terraform-provider-jenkins).
