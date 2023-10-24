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
