terraform {
  required_version = "~> 0.15.1"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 3.67.0"
    }
  }

  backend "remote" {
    organization = "ww24"

    workspaces {
      name = "chatbot"
    }
  }
}
