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

resource "google_service_account" "invoker" {
  account_id   = "${var.name}-invoker"
  display_name = "${var.name}-invoker Service Account"
}
