resource "jenkins_credential_ssh" "global" {
  name     = "some-id"
  username = "example-username"
  privatekey = file("/some/path/id_rsa")
}

resource "jenkins_credential_ssh" "folder" {
  name       = "some-id"
  folder     = jenkins_folder.example.id
  username   = "example-username"
  privatekey = file("/some/path/id_rsa")
}
