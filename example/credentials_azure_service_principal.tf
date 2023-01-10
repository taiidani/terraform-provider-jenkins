   resource "jenkins_credential_azure_service_principal" "azure_service_principal_test_credential" {
    name = "bla"
    folder = jenkins_folder.example.id
    description = "blabla"
    subscription_id = "123"
    client_id = "123"
    client_secret = "super-secret"
    tenant = "456"
}
