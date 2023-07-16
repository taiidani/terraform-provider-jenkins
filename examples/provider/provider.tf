# Configure the Jenkins Provider
provider "jenkins" {
  server_url = "http://localhost:8080" # Or JENKINS_URL env var
  username   = "admin"                 # Or JENKINS_USERNAME env var
  password   = "admin"                 # Or JENKINS_PASSWORD env var
  ca_cert    = ""                      # Or JENKINS_CA_CERT env var
}
