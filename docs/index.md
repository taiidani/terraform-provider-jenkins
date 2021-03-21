# Jenkins Provider

The Jenkins provider is used to interact with the Jenkins API. The provider needs to be configured with the proper credentials before it can be used.

## Example Usage

```hcl
# Configure the Jenkins Provider
provider "jenkins" {
  server_url = "https://jenkins.url" # Or use JENKINS_URL env var
  username   = "username"            # Or use JENKINS_USERNAME env var
  password   = "password"            # Or use JENKINS_PASSWORD env var
  ca_cert = ""                       # Or use JENKINS_CA_CERT env var
}

# Create a Jenkins job
resource "jenkins_job" "example" {
  # ...
}
```

## Authentication

Jenkins uses a user/password challenge for authentication. It requires a username & password for determining identity and permissions. This method also supports Jenkins' various authentication plugins, such as GitHub OAuth (through the use of Personal Access Tokens).

### Static credentials ###

Static credentials can be provided by adding a `username` and `password` in-line in the Jenkins provider block:

Usage:

```hcl
provider "jenkins" {
  server_url = "https://jenkins.url"
  username   = "username"
  password   = "password"
}
```

### Environment variables

You can provide your credentials via the `JENKINS_USERNAME` and `JENKINS_PASSWORD`, environment variables. `JENKINS_URL` is also available which will assign the `server_url` property.

```hcl
provider "jenkins" {}
```

Usage:

```sh
$ export JENKINS_SERVER_URL="https://jenkins.url"
$ export JENKINS_USERNAME="username"
$ export JENKINS_PASSWORD="password"
$ terraform plan
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html) (e.g. `alias` and `version`), the following arguments are supported in the Jenkins `provider` block:

* `server_url` - (Required) This is the Jenkins server URL. It should be fully qualified (e.g. `https://...`) and point to the root of the Jenkins server location.

* `username` - (Required) This is Jenkins username for authentication.

* `password` - (Required) This is the Jenkins password for authentication. If you are using the GitHub OAuth authentication method, enter your Personal Access Token here.

* `ca_cert` - (Optional) This is the path to the self-signed certificate that may be required in order to authenticate to your Jenkins instance.
