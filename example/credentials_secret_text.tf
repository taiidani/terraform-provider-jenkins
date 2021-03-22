resource "jenkins_credential_secret_text" "global" {
  name   = "global-username"
  secret = "barsoom"
}

resource "jenkins_credential_secret_text" "folder" {
  name   = "folder-username"
  folder = jenkins_folder.example.id
  secret = "barsoom"
}
