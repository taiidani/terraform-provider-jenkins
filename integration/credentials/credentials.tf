resource "jenkins_credential_username" "global" {
  name     = "global-username"
  username = "foo"
  # Passwords may be unmanaged
  # password = "barsoom"
}

resource "jenkins_credential_username" "folder" {
  name     = "folder-username"
  folder   = jenkins_folder.example.id
  username = "folder-foo"
  password = "barsoom"
}
