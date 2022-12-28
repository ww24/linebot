provider "google" {
  credentials = var.google_credentials
  project     = var.project
  region      = local.location
}

provider "google-beta" {
  credentials = var.google_credentials
  project     = var.project
  region      = local.location
}
