resource "jenkins_credential_vault_approle" "global" {
  name      = "global-approle"
  namespace = "baz"
  role_id   = "foo"
  secret_id = "bar"
}

data "jenkins_credential_vault_approle" "global" {
  depends_on = [jenkins_credential_vault_approle.global]
  name       = "global-approle"
}

output "vault_approle" {
  value = data.jenkins_credential_vault_approle.global
}

resource "jenkins_credential_vault_approle" "folder" {
  name      = "folder-approle"
  folder    = jenkins_folder.example.id
  role_id   = "foo"
  secret_id = "bar"
}

resource "jenkins_credential_vault_approle" "global-namespaced" {
  name      = "global-approle-namespaced"
  role_id   = "foo"
  secret_id = "bar"
  namespace = "my-namespace"
}

resource "jenkins_credential_vault_approle" "folder-namespaced" {
  name      = "folder-approle-namespaced"
  folder    = jenkins_folder.example.id
  role_id   = "foo"
  secret_id = "bar"
  namespace = "my-namespace"
}
