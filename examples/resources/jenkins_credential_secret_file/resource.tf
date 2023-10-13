resource "jenkins_credential_secret_file" "example" {
  name        = "example-secret-file"
  filename    = "secret-file.txt"
  secretbytes = base64encode("My secret file very secret content.")
}
