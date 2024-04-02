terraform {
  required_version = "~> 1.5.7"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.11.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 5.23.0"
    }
  }

  backend "remote" {
    organization = "ww24"

    workspaces {
      name = "chatbot"
    }
  }
}
