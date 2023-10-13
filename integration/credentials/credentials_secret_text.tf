resource "jenkins_credential_secret_text" "global" {
  name   = "global-secret-text"
  secret = "barsoom"
}

output "secret_text" {
  value     = jenkins_credential_secret_text.global
  sensitive = true
}

resource "jenkins_credential_secret_text" "folder" {
  name   = "folder-secret-text"
  folder = jenkins_folder.example.id
  secret = "barsoom"
}
