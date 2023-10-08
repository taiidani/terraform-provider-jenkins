locals {
  pipeline_scm_template = templatefile("${path.module}/pipeline_scm.xml", {
    description = "An example pipeline job"
  })

  pipeline_inline_template = templatefile("${path.module}/pipeline_inline.xml", {
    description = "An example pipeline inline script job"
  })

  freestyle_template = templatefile("${path.module}/freestyle.xml", {
    description = "An example freestyle job"
  })
}

resource "jenkins_job" "pipeline_scm" {
  name     = "pipeline-scm"
  folder   = jenkins_folder.example.id
  template = local.pipeline_scm_template
}

data "jenkins_job" "pipeline_scm" {
  depends_on = [jenkins_job.pipeline_scm]
  name       = "pipeline-scm"
  folder     = jenkins_job.pipeline_scm.folder
}

resource "jenkins_job" "pipeline_inline" {
  name     = "pipeline-inline"
  folder   = jenkins_folder.example.id
  template = local.pipeline_inline_template
}

data "jenkins_job" "pipeline_inline" {
  depends_on = [jenkins_job.pipeline_inline]
  name       = "pipeline-inline"
  folder     = jenkins_job.pipeline_inline.folder
}

resource "jenkins_job" "freestyle" {
  name     = "freestyle"
  folder   = jenkins_folder.example.id
  template = local.freestyle_template
}

data "jenkins_job" "freestyle" {
  depends_on = [jenkins_job.freestyle]
  name       = "freestyle"
  folder     = jenkins_job.freestyle.folder
}
