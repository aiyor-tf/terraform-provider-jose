terraform {
  required_providers {
    jose = {
      source = "registry.terraform.io/aiyor-tf/jose"
    }
  }
}

provider "jose" {}
