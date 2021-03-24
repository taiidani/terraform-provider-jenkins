# Connects to the instance launched through "docker-compose up -d"
# Once done with testing, clean up the instance with "docker-compose down --volumes"
provider "jenkins" {
  server_url = "http://localhost:8080" # Or use JENKINS_URL env var
  username   = "admin"                 # Or use JENKINS_USERNAME env var
  password   = "admin"                 # Or use JENKINS_PASSWORD env var
  ca_cert    = ""                      # Or use JENKINS_CA_CERT env var
}

terraform {
  required_providers {
    jenkins = {
      source  = "taiidani/jenkins"
      version = ">= 0.5.0"
    }
  }
}
