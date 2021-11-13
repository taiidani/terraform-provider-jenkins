resource "jenkins_credential_vault_approle" "global" {
  name      = "global-approle"
  role_id   = "foo"
  secret_id = "bar"
  namespace = "my-namespace"
}

resource "jenkins_credential_vault_approle" "folder" {
  name      = "folder-approle"
  folder    = jenkins_folder.example.id
  role_id   = "foo"
  secret_id = "bar"
  namespace = "my-namespace"
}