terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = ">= 0.5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.5.0"
    }
  }
}
