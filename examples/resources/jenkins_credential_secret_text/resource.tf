resource "jenkins_credential_secret_text" "example" {
  name   = "example-username"
  secret = "super-secret"
}
