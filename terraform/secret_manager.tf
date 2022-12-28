resource "google_secret_manager_secret" "line-channel-secret" {
  secret_id = "line-channel-secret"

  labels = {
    service = local.name
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
    service = local.name
  }

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "line-channel-access-token" {
  secret      = google_secret_manager_secret.line-channel-access-token.id
  secret_data = var.line_channel_access_token
}

resource "google_secret_manager_secret" "maxmind-license-key" {
  secret_id = "maxmind-license-key"

  labels = {
    service = local.name
  }

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "maxmind-license-key" {
  secret      = google_secret_manager_secret.maxmind-license-key.id
  secret_data = var.maxmind_license_key
}
