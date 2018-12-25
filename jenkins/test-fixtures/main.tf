provider "docker" {
  version = "~> 1.1"
}

resource "docker_volume" "data" {}

resource "docker_container" "jenkins" {
  image = "jenkins-provider-acc"
  name  = "jenkins-provider-acc"

  env = [
    "JAVA_OPTS=-Djenkins.install.runSetupWizard=false",
  ]

  ports {
    internal = "8080"
    external = "8080"
    ip       = "127.0.0.1"
  }

  volumes {
    volume_name    = "${docker_volume.data.name}"
    container_path = "/var/jenkins_home"
  }

  healthcheck {
    test = ["CMD", "curl", "-f", "http://localhost:8080"]
  }
}

output "container_id" {
  value = "${docker_container.jenkins.id}"
}
