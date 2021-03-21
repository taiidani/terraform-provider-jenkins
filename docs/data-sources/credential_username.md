# jenkins_credential_username Data Source

Get the attributes of a username credential within Jenkins.

## Example Usage

```hcl
data "jenkins_credential_username" "example" {
  name        = "job-name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource being read.
* `domain` - (Optional) The domain store to place the credentials into. If not set will default to the global credentials store.
* `folder` - (Optional) The folder namespace containing this resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The full canonical job path, E.G. `/job/job-name`.
* `description` - A human readable description of the credentials being stored.
* `scope` - The visibility of the credentials to Jenkins agents. This must be set to either "GLOBAL" or "SYSTEM".
* `username` - The username to be associated with the credentials.
