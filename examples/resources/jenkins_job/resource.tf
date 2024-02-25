resource "jenkins_folder" "example" {
  name = "folder-name"
}

resource "jenkins_job" "example" {
  name   = "example"
  folder = jenkins_folder.example.id
  template = templatefile("${path.module}/job.xml", {
    description = "An example job created from Terraform"
  })
}
