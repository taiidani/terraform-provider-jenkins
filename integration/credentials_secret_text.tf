resource "jenkins_credential_secret_text" "global" {
  name   = "global-secret-text"
  secret = "barsoom"
}

resource "jenkins_credential_secret_text" "folder" {
  name   = "folder-secret-text"
  folder = jenkins_folder.example.id
  secret = "barsoom"
}
