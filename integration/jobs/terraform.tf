terraform {
  required_providers {
    jenkins = {
      source  = "taiidani/jenkins"
      version = ">= 0.5.0"
    }
  }
}

variable "port" {
  description = "The port that the Jenkins setup has been published on"
}

provider "jenkins" {
  server_url = "http://localhost:${var.port}"
  username   = "admin"
  password   = "admin"
  ca_cert    = ""
}
