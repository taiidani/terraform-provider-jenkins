# jenkins_folder Data Source

Get the attributes of a folder within Jenkins.

~> The Jenkins installation that uses this resource is expected to have the [Cloudbees Folders Plugin](https://plugins.jenkins.io/cloudbees-folder) installed in their system.

## Example Usage

```hcl
data jenkins_folder example {
  name        = "folder-name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the folder being read.
* `folder` - (Optional) The folder namespace containing this folder.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - A block of text describing the folder's purpose.
* `template` - A Jenkins-compatible XML template to describe the folder.
