# Terraform Provider

![Unit Tests](https://github.com/taiidani/terraform-provider-jenkins/workflows/Unit%20Tests/badge.svg)
![Acceptance Tests](https://github.com/taiidani/terraform-provider-jenkins/workflows/Acceptance%20Tests/badge.svg)
[![codecov](https://codecov.io/gh/taiidani/terraform-provider-jenkins/branch/master/graph/badge.svg)](https://codecov.io/gh/taiidani/terraform-provider-jenkins)

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

This is a community provider and is not supported by Hashicorp.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.11+ (to build the provider plugin)

## Installation

Install the provider with:

```bash
go install github.com/taiidani/terraform-provider-jenkins
```

Then copy or link the resulting binary to your [terraform.d plugins folder](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins). On macOS this might look like:

```bash
ln -s "$(go env GOPATH)/bin/terraform-provider-jenkins" "$HOME/.terraform.d/plugins/terraform-provider-jenkins"
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

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
