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

data "jenkins_folder" "example" {
  depends_on = [jenkins_folder.example]
  name       = jenkins_folder.example.name
}

resource "jenkins_folder" "example_subfolder" {
  name        = "subfolder"
  folder      = jenkins_folder.example.id
  description = "A sample subfolder"
}

data "jenkins_folder" "example_subfolder" {
  depends_on = [jenkins_folder.example_subfolder]
  name       = jenkins_folder.example_subfolder.name
  folder     = jenkins_folder.example_subfolder.folder
}
