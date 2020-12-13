# jenkins_folder Resource

Manages a folder within Jenkins.

~> The Jenkins installation that uses this resource is expected to have the [Cloudbees Folders Plugin](https://plugins.jenkins.io/cloudbees-folder) installed in their system.

## Example Usage

```hcl
resource jenkins_folder example {
  name = "folder-name"
}

resource jenkins_folder example_child {
  name   = "child-name"
  folder = jenkins_folder.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the folder being created.
* `folder` - (Optional) The folder namespace to store the subfolder in. If creating in a nested folder structure you may separate folder names with `/`, such as `parent/child`. This name cannot be changed once the folder has been created, and all parent folders must be created in advance.
* `description` - (Optional) A block of text describing the folder's purpose.
* `template` - (Optional) A Jenkins-compatible XML template to describe the folder. You can retrieve an existing folder's XML by appending `/config.xml` to its URL and viewing the source in your browser. The `template` property is rendered using a Golang template that takes the other resource arguments as variables. Do not include the XML prolog in the definition. If `template` is not provided this will default to a "best-guess" folder definition.
* `permissions` - (Optional) A list of strings containing Jenkins permissions assigments to users and groups for the folder. For example:

```hcl
  permissions = [
    "hudson.model.Item.Build:username",
    "hudson.model.Item.Cancel:username",
    "hudson.model.Item.Configure:username",
    "hudson.model.Item.Create:username",
    "hudson.model.Item.Delete:username",
  ]
```

## Attribute Reference

All arguments above are exported.
