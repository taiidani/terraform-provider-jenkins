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

check "pipeline_scm_xml" {
  data "jenkins_job" "pipeline_scm" {
    name   = "pipeline-scm"
    folder = jenkins_folder.example.id
  }

  assert {
    condition     = trimspace(data.jenkins_job.pipeline_scm.template) == trimspace(local.pipeline_scm_template)
    error_message = "${data.jenkins_job.pipeline_scm.name} produced inconsistent XML"
  }
}

resource "jenkins_job" "pipeline_inline" {
  name     = "pipeline-inline"
  folder   = jenkins_folder.example.id
  template = local.pipeline_inline_template
}

check "pipeline_inline_xml" {
  data "jenkins_job" "pipeline_inline" {
    name   = "pipeline-inline"
    folder = jenkins_folder.example.id
  }

  assert {
    condition     = trimspace(data.jenkins_job.pipeline_inline.template) == trimspace(local.pipeline_inline_template)
    error_message = "${data.jenkins_job.pipeline_inline.name} produced inconsistent XML"
  }
}

resource "jenkins_job" "freestyle" {
  name     = "freestyle"
  folder   = jenkins_folder.example.id
  template = local.freestyle_template
}

check "freestyle_xml" {
  data "jenkins_job" "freestyle" {
    name   = "freestyle"
    folder = jenkins_folder.example.id
  }

  assert {
    condition     = trimspace(data.jenkins_job.freestyle.template) == trimspace(local.freestyle_template)
    error_message = "${data.jenkins_job.freestyle.name} produced inconsistent XML"
  }
}
