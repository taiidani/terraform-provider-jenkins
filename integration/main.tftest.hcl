provider "docker" {}
provider "random" {}

run "setup" {
  module {
    source = "./setup"
  }

  providers = {
    docker = docker
    random = random
  }
}

run "jobs" {
  module {
    source = "./jobs"
  }

  variables {
    port = run.setup.port
  }

  providers = {
    random = random
  }

  assert {
    condition     = chomp(coalesce(data.jenkins_job.pipeline_scm.template, "")) == chomp(local.pipeline_scm_template)
    error_message = "${data.jenkins_job.pipeline_scm.name} produced inconsistent XML"
  }

  assert {
    condition     = chomp(data.jenkins_job.pipeline_inline.template) == chomp(local.pipeline_inline_template)
    error_message = "${data.jenkins_job.pipeline_inline.name} produced inconsistent XML"
  }

  assert {
    condition     = chomp(data.jenkins_job.freestyle.template) == chomp(local.freestyle_template)
    error_message = "${data.jenkins_job.freestyle.name} produced inconsistent XML"
  }
}


run "credentials" {
  module {
    source = "./credentials"
  }

  variables {
    port = run.setup.port
  }

  providers = {
    random = random
  }
}
