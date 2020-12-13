# Terraform Provider

![Unit Tests](https://github.com/taiidani/terraform-provider-jenkins/workflows/Unit%20Tests/badge.svg)
![Acceptance Tests](https://github.com/taiidani/terraform-provider-jenkins/workflows/Acceptance%20Tests/badge.svg)
[![codecov](https://codecov.io/gh/taiidani/terraform-provider-jenkins/branch/master/graph/badge.svg)](https://codecov.io/gh/taiidani/terraform-provider-jenkins)

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

This is a community provider and is not supported by Hashicorp.

## Installation

This provider has been published to the Terraform Registry at https://registry.terraform.io/providers/taiidani/jenkins. Please visit the registry for documentation and installation instructions.

## Developing the Provider

Working on this provider requires the following:

* [Terraform](https://www.terraform.io/downloads.html) 0.14+
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

## Attribution

This provider design was originally inspired from the work at [dihedron/terraform-provider-jenkins](https://github.com/dihedron/terraform-provider-jenkins).
