data "jenkins_credential_username" "example" {
  name   = "name"
  folder = jenkins_folder.example.id
}
