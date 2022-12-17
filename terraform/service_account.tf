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

resource "google_pubsub_topic_iam_member" "linebot-access-log-publisher" {
  topic  = google_pubsub_topic.access_log.name
  role   = "roles/pubsub.publisher"
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

# allows Cloud Pub/Sub Service Account to push BigQuery Dataset
# https://cloud.google.com/iam/docs/service-agents
resource "google_project_service_identity" "pubsub" {
  provider = google-beta
  service  = "pubsub.googleapis.com"
}

resource "google_bigquery_table_iam_member" "pubsub_sa_bigquery" {
  dataset_id = google_bigquery_table.access_log.dataset_id
  table_id   = google_bigquery_table.access_log.table_id
  for_each   = toset(["roles/bigquery.metadataViewer", "roles/bigquery.dataEditor"])
  role       = each.key
  member     = "serviceAccount:${google_project_service_identity.pubsub.email}"
}

# access-log GSA
resource "google_service_account" "access-log" {
  account_id   = "${var.name}-access-log"
  display_name = "${var.name}-access-log Service Account"
}

resource "google_bigquery_dataset_iam_member" "access-log" {
  dataset_id = google_bigquery_dataset.geolite2.dataset_id
  role       = "roles/bigquery.admin"
  member     = "serviceAccount:${google_service_account.access-log.email}"
}

resource "google_project_iam_member" "access-log" {
  project = var.project
  role    = "roles/bigquery.jobUser"
  member  = "serviceAccount:${google_service_account.access-log.email}"
}

# BigQuery Connection
resource "google_bigquery_connection" "geolite2" {
  connection_id = "geolite2"
  location      = "US"
  cloud_resource {}
}

resource "google_storage_bucket_iam_member" "geolite2" {
  bucket = google_storage_bucket.geolite2.name
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${google_bigquery_connection.geolite2.cloud_resource[0].service_account_id}"
}
