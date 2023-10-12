resource "jenkins_credential_azure_service_principal" "foo" {
  name            = "example-secret"
  subscription_id = "01234567-89ab-cdef-0123-456789abcdef"
  client_id       = "abcdef01-2345-6789-0123-456789abcdef"
  client_secret   = "super-secret"
  tenant          = "01234567-89ab-cdef-abcd-456789abcdef"
}
