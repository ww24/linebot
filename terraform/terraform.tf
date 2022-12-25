terraform {
  required_version = "~> 1.3.4"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.47.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 4.47.0"
    }
  }

  backend "remote" {
    organization = "ww24"

    workspaces {
      name = "chatbot"
    }
  }
}
