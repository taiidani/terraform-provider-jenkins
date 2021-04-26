# jenkins_credential_secret_file Resource

Manages a secret file credential within Jenkins. This secret file may then be referenced within jobs that are created.

## Example Usage

```hcl
resource "jenkins_credential_secret_file" "example" {
  name        = "example-secret-file"
  filename    = "secret-file.txt"
  secretbytes = base64encode("My secret file very secret content.")
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the credentials being created. This maps to the ID property within Jenkins, and cannot be changed once set.
* `domain` - (Optional) The domain store to place the credentials into. If not set will default to the global credentials store.
* `folder` - (Optional) The folder namespace to store the credentials in. If not set will default to global Jenkins credentials.
* `scope` - (Optional) The visibility of the credentials to Jenkins agents. This must be set to either "GLOBAL" or "SYSTEM". If not set will default to "GLOBAL".
* `description` - (Optional) A human readable description of the credentials being stored.
* `filename` - (Required) The secret file filename on jenkins server side.
* `secretbytes` - (Required) The secret file, base64 encoded content. It can be sourced directly from local file with filebase64(path) TF function or given directly.


## Attribute Reference

All arguments above are exported.
