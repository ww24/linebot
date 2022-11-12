# linebot GSA
resource "google_service_account" "linebot" {
  account_id   = var.name
  display_name = "${var.name} Service Account"
}

resource "google_project_iam_member" "linebot-firestore" {
  project = var.project
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_project_iam_member" "linebot-cloudtrace" {
  project = var.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_project_iam_member" "linebot-cloudprofiler" {
  project = var.project
  role    = "roles/cloudprofiler.agent"
  member  = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_project_iam_member" "linebot-cloudtasks" {
  project = var.project
  for_each = toset([
    "roles/cloudtasks.viewer",
    "roles/cloudtasks.enqueuer",
    "roles/cloudtasks.taskDeleter",
  ])
  role   = each.value
  member = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_project_iam_member" "linebot-service-account-user" {
  project = var.project
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_storage_bucket_iam_member" "linebot-storage" {
  bucket = google_storage_bucket.image.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_secret_manager_secret_iam_member" "linebot-secret" {
  for_each = toset([
    google_secret_manager_secret.line-channel-secret.id,
    google_secret_manager_secret.line-channel-access-token.id,
  ])
  secret_id = each.value
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.linebot.email}"
}

# invoker GSA
resource "google_service_account" "invoker" {
  account_id   = "${var.name}-invoker"
  display_name = "${var.name}-invoker Service Account"
}

resource "google_project_iam_member" "invoker-cloudrun" {
  project = var.project
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.invoker.email}"
}

# screenshot GSA
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

resource "google_storage_bucket_iam_member" "screenshot-storage" {
  bucket = google_storage_bucket.image.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.screenshot.email}"
}
