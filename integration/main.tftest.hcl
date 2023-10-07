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
