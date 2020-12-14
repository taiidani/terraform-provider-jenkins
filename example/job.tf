resource "jenkins_folder" "example" {
  name        = "folder-name"
  description = "A sample folder"
}

resource "jenkins_job" "pipeline" {
  name     = "pipeline"
  folder   = jenkins_folder.example.id
  template = file("${path.module}/pipeline.xml")

  parameters = {
    description = "An example pipeline job"
  }
}

resource "jenkins_job" "freestyle" {
  name     = "freestyle"
  folder   = jenkins_folder.example.id
  template = file("${path.module}/freestyle.xml")

  parameters = {
    description = "An example freestyle job"
  }
}
