resource "google_secret_manager_secret" "line-channel-secret" {
  secret_id = "line-channel-secret"

  labels = {
    service = var.name
  }

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "line-channel-secret" {
  secret      = google_secret_manager_secret.line-channel-secret.id
  secret_data = var.line_channel_secret
}

resource "google_secret_manager_secret" "line-channel-access-token" {
  secret_id = "line-channel-access-token"

  labels = {
    service = var.name
  }

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "line-channel-access-token" {
  secret      = google_secret_manager_secret.line-channel-access-token.id
  secret_data = var.line_channel_access_token
}
