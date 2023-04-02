# jenkins_credential_secret_text Resource

Manages a secret text credential within Jenkins. This secret text may then be referenced within jobs that are created.

## Example Usage

```hcl
resource "jenkins_credential_secret_text" "example" {
  name     = "example-username"
  secret   = "super-secret"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the credentials being created. This maps to the ID property within Jenkins, and cannot be changed once set.
* `domain` - (Optional) The domain store to place the credentials into. If not set will default to the global credentials store.
* `folder` - (Optional) The folder namespace to store the credentials in. If not set will default to global Jenkins credentials.
* `scope` - (Optional) The visibility of the credentials to Jenkins agents. This must be set to either "GLOBAL" or "SYSTEM". If not set will default to "GLOBAL".
* `description` - (Optional) A human readable description of the credentials being stored. If not set will default to "Managed by Terraform".
* `secret` - (Required) The secret text to be associated with the credentials.

## Attribute Reference

All arguments above are exported.

## Import

Secret text credential may be imported by their canonical name, e.g.

```sh
$ terraform import jenkins_credential_secret_text.example "[<folder>/]<domain>/<name>"
```
