data "jenkins_credential_vault_approle" "example" {
  name   = "name"
  folder = jenkins_folder.example.id
}
