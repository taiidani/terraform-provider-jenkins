data "jenkins_credential_aws" "example" {
  name   = "name"
  folder = jenkins_folder.example.id
}
