# jenkins_credential_vault_approle Resource

Manages a Vault AppRole credential within Jenkins. This credential may then be referenced within jobs that are created.

~> The "secret_id" property may leave plain-text secret id in your state file. If using the property to manage the secret id in Terraform, ensure that your state file is properly secured and encrypted at rest.

~> The Jenkins installation that uses this resource is expected to have the [Hashicorp Vault Plugin](https://plugins.jenkins.io/hashicorp-vault-plugin/) installed in their system.

## Example Usage

```hcl
resource "jenkins_credential_vault_approle" "example" {
  name     = "example-approle"
  role_id = "example"
  secret_id = "super-secret"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the credentials being created. This maps to the ID property within Jenkins, and cannot be changed once set.
* `domain` - (Optional) The domain store to place the credentials into. If not set will default to the global credentials store.
* `folder` - (Optional) The folder namespace to store the credentials in. If not set will default to global Jenkins credentials.
* `scope` - (Optional) The visibility of the credentials to Jenkins agents. This must be set to either "GLOBAL" or "SYSTEM". If not set will default to "GLOBAL".
* `description` - (Optional) A human readable description of the credentials being stored.
* `namespace` - (Optional) The Vault namespace of the approle credential.
* `path` - (Optional) The unique name of the approle auth backend. Defaults to `approle`.
* `role_id` - (Required) The role_id to be associated with the credentials.
* `secret_id` - (Optional) The secret_id to be associated with the credentials. If empty then the secret_id property will become unmanaged and expected to be set manually within Jenkins. If set then the secret_id will be updated only upon changes -- if the secret_id is set manually within Jenkins then it will not reconcile this drift until the next time the secret_id property is changed.

## Attribute Reference

All arguments above are exported.
