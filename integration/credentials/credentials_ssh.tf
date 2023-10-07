resource "jenkins_credential_ssh" "global" {
  name       = "some-id"
  username   = "example-username"
  privatekey = file("${path.module}/id_ed25519")
}

resource "jenkins_credential_ssh" "folder" {
  name       = "some-id"
  folder     = jenkins_folder.example.id
  username   = "example-username"
  privatekey = file("${path.module}/id_ed25519")
}
