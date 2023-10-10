resource "jenkins_credential_vault_approle" "example" {
  name      = "example-approle"
  role_id   = "example"
  secret_id = "super-secret"
}
