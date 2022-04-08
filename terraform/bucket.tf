resource "google_storage_bucket" "image" {
  project                     = var.project
  name                        = var.image_bucket
  location                    = "ASIA-NORTHEAST1"
  storage_class               = "STANDARD"
  uniform_bucket_level_access = true

  labels = {
    service = var.name
  }
}
