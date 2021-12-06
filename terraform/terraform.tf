terraform {
  required_version = "~> 1.0.7"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.2.1"
    }
  }

  backend "remote" {
    organization = "ww24"

    workspaces {
      name = "chatbot"
    }
  }
}
