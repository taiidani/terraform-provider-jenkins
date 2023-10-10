resource "jenkins_credential_username" "global" {
  name     = "global-username"
  username = "foo"
  # Passwords may be unmanaged
  # password = "barsoom"
}

data "jenkins_credential_username" "global" {
  depends_on = [jenkins_credential_username.global]
  name       = "global-username"
}

output "username" {
  value = data.jenkins_credential_username.global
}

resource "jenkins_credential_username" "folder" {
  name     = "folder-username"
  folder   = jenkins_folder.example.id
  username = "folder-foo"
  password = "barsoom"
}
