resource "google_service_account" "linebot" {
  account_id   = var.name
  display_name = "${var.name} Service Account"
}

resource "google_project_iam_member" "firestore" {
  project = var.project
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_project_iam_member" "cloudtrace" {
  project = var.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_project_iam_member" "cloudprofiler" {
  project = var.project
  role    = "roles/cloudprofiler.agent"
  member  = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_project_iam_member" "cloudtasks-viewer" {
  project = var.project
  role    = "roles/cloudtasks.viewer"
  member  = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_project_iam_member" "cloudtasks-enqueuer" {
  project = var.project
  role    = "roles/cloudtasks.enqueuer"
  member  = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_project_iam_member" "cloudtasks-deleter" {
  project = var.project
  role    = "roles/cloudtasks.taskDeleter"
  member  = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_storage_bucket_iam_member" "image" {
  bucket = google_storage_bucket.image.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_secret_manager_secret_iam_member" "line-channel-secret" {
  secret_id = google_secret_manager_secret.line-channel-secret.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_secret_manager_secret_iam_member" "line-channel-access-token" {
  secret_id = google_secret_manager_secret.line-channel-access-token.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_service_account" "invoker" {
  account_id   = "${var.name}-invoker"
  display_name = "${var.name}-invoker Service Account"
}

resource "google_project_iam_member" "invoker-service-account-user" {
  project = var.project
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.invoker.email}"
}

resource "google_service_account" "screenshot" {
  account_id   = var.name_screenshot
  display_name = "screenshot Service Account"
}

resource "google_project_iam_member" "screenshot-cloudtrace" {
  project = var.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.screenshot.email}"
}

resource "google_project_iam_member" "screenshot-cloudprofiler" {
  project = var.project
  role    = "roles/cloudprofiler.agent"
  member  = "serviceAccount:${google_service_account.screenshot.email}"
}
