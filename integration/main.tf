provider "random" {}

module "jobs" {
  source = "./jobs"
  port   = 8080
}

module "credentials" {
  source = "./credentials"
  port   = 8080
}
