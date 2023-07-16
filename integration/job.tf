resource "jenkins_job" "pipeline" {
  name     = "pipeline"
  folder   = jenkins_folder.example.id
  template = templatefile("${path.module}/pipeline.xml", {
    description = "An example pipeline job"
  })
}

resource "jenkins_job" "freestyle" {
  name     = "freestyle"
  folder   = jenkins_folder.example.id
  template = templatefile("${path.module}/freestyle.xml", {
    description = "An example freestyle job"
  })
}
