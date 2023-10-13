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

  assert {
    condition     = output.azure_service_principal.client_id == "123"
    error_message = "${nonsensitive(output.azure_service_principal.client_id)} did not contain expected \"123\" value"
  }
  assert {
    condition     = output.secret_file.filename == "secret-file.txt"
    error_message = "${nonsensitive(output.secret_file.filename)} did not contain expected \"secret-file.txt\" value"
  }
  assert {
    condition     = output.secret_text.secret == "barsoom"
    error_message = "${nonsensitive(output.secret_text.secret)} did not contain expected \"barsoom\" value"
  }
  assert {
    condition     = output.username.username == jenkins_credential_username.global.username
    error_message = "${output.username.username} data value did not match resource value"
  }
  assert {
    condition     = output.vault_approle.role_id == jenkins_credential_vault_approle.global.role_id
    error_message = "${output.vault_approle.role_id} data value did not match resource value"
  }
  assert {
    condition     = output.vault_approle.namespace == jenkins_credential_vault_approle.global.namespace
    error_message = "${output.vault_approle.namespace} data value did not match resource value"
  }
}
