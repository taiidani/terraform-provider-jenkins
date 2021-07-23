# jenkins_folder Resource

Manages a folder within Jenkins.

~> The Jenkins installation that uses this resource is expected to have the [Cloudbees Folders Plugin](https://plugins.jenkins.io/cloudbees-folder) installed in their system.

## Example Usage

```hcl
resource "jenkins_folder" "example" {
  name        = "folder-name"
  description = "A top-level folder"
}

resource "jenkins_folder" "example_child" {
  name        = "child-name"
  folder      = jenkins_folder.example.id
  description = "A nested subfolder"

  security {
    permissions = [
      "com.cloudbees.plugins.credentials.CredentialsProvider.Create:anonymous",
      "com.cloudbees.plugins.credentials.CredentialsProvider.Delete:authenticated",
      "hudson.model.Item.Cancel:authenticated",
      "hudson.model.Item.Discover:anonymous",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the folder being created.
* `display_name` - (Optional) The name of the folder to be displayed in the UI.
* `folder` - (Optional) The folder namespace to store the subfolder in. If creating in a nested folder structure you may separate folder names with `/`, such as `parent/child`. This name cannot be changed once the folder has been created, and all parent folders must be created in advance.
* `description` - (Optional) A block of text describing the folder's purpose.
* `security` - (Optional) An optional block defining a project-based authorization strategy, documented below.

### security

~> This block may need the [Matrix Authorization Strategy Plugin](https://plugins.jenkins.io/matrix-auth/) installed and enabled in the system's Global Security settings in order to function properly.

* `inheritance_strategy` - The strategy for applying these permissions sets to existing inherited permissions. Defaults to "org.jenkinsci.plugins.matrixauth.inheritance.InheritParentStrategy".
* `permissions` - A list of strings containing Jenkins permissions assigments to users and groups for the folder. For example:

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

In addition to all arguments above, the following attributes are exported:

* `id` - The full canonical folder path, E.G. `/job/parent`.
* `template` - A Jenkins-compatible XML template to describe the folder. You can retrieve an existing folder's XML by appending `/config.xml` to its URL and viewing the source in your browser.

## Import

Folders may be imported by their canonical name, e.g.

```sh
$ terraform import jenkins_folder.example /job/folder-name
```
