resource "jenkins_credential_vault_approle" "foo" {
    name = "global-approle"
    role_id = "foo"
    secret_id = "bar"
    namespace = "my-namespace"
}

resource "jenkins_credential_vault_approle" "foo" {
    name = "global-approle-folder"
    folder = jenkins_folder.example.id
    role_id = "foo"
    secret_id = "bar"
    namespace = "my-namespace"
}