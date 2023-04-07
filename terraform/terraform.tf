terraform {
  required_version = "~> 1.3.6"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.60.2"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 4.57.0"
    }
  }

  backend "remote" {
    organization = "ww24"

    workspaces {
      name = "chatbot"
    }
  }
}
