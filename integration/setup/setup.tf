resource "random_pet" "name" {
  prefix = "jenkins"
}

resource "docker_image" "jenkins" {
  name = "jenkins"

  build {
    context = "."
    tag     = ["jenkins:terraformtest"]
  }
}

resource "docker_container" "jenkins" {
  name  = random_pet.name.id
  image = docker_image.jenkins.image_id
  wait  = true

  env = [
    "JAVA_OPTS=-Djenkins.install.runSetupWizard=false",
  ]

  ports {
    internal = 8080
    ip       = "127.0.0.1"
  }
}

output "name" {
  description = "The name of the docker container spun up"
  value       = docker_container.jenkins.name
}

output "port" {
  description = "The port that the Jenkins setup has been published on"
  value       = docker_container.jenkins.ports[0].external
}
