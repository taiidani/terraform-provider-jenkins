resource "jenkins_folder" "example" {
  name        = "folder-name"
  description = "A sample folder"

  security {
    permissions = [
      "com.cloudbees.plugins.credentials.CredentialsProvider.Create:anonymous",
      "com.cloudbees.plugins.credentials.CredentialsProvider.Delete:authenticated",
      "hudson.model.Item.Cancel:authenticated",
      "hudson.model.Item.Discover:anonymous",
    ]
  }
}

resource "jenkins_folder" "example_subfolder" {
  name        = "subfolder"
  folder      = jenkins_folder.example.id
  description = "A sample subfolder"
}
