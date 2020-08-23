# jenkins_credential_username Resource

Manages a username credential within Jenkins. This username may then be referenced within jobs that are created.

~> The "password" property may leave plain-text passwords in your state file. If using the property to manage the password in Terraform, ensure that your state file is properly secured and encrypted at rest.

~> When using this resource within a folder context it can conflict with the [folder resource](folder) template. When using these in combination you may need to add a lifecycle `ignore_changes` rule to the folder's `template` property.

## Example Usage

```hcl
resource jenkins_credential_username example {
  name     = "example-username"
  username = "example"
  password = "super-secret"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the credentials being created. This maps to the ID property within Jenkins, and cannot be changed once set.
* `domain` - (Optional) The domain store to place the credentials into. If not set will default to the global credentials store.
* `folder` - (Optional) The folder namespace to store the credentials in. If not set will default to global Jenkins credentials.
* `scope` - (Optional) The visibility of the credentials to Jenkins agents. This must be set to either "GLOBAL" or "SYSTEM". If not set will default to "GLOBAL".
* `description` - (Optional) A human readable description of the credentials being stored.
* `username` - (Required) The username to be associated with the credentials.
* `password` - (Optional) The password to be associated with the credentials. If empty then the password property will become unmanaged and expected to be set manually within Jenkins. If set then the password will be updated only upon changes -- if the password is set manually within Jenkins then it will not reconcile this drift until the next time the password property is changed.

## Attribute Reference

All arguments above are exported.
