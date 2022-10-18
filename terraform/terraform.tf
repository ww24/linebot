terraform {
  required_version = "~> 1.1.2"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.41.0"
    }
  }

  backend "remote" {
    organization = "ww24"

    workspaces {
      name = "chatbot"
    }
  }
}
