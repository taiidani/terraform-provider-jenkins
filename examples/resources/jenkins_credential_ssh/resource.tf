resource "jenkins_credential_ssh" "example" {
  name       = "example-id"
  username   = "example-username"
  privatekey = file("/some/path/id_rsa")
  passphrase = "Super_Secret_Pass"
}
