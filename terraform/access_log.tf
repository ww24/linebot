resource "google_pubsub_schema" "access_log_schema_v1" {
  name       = "${var.name}-access-log-v1"
  type       = "AVRO"
  definition = file("access_log_schema/v1.avsc")
}

resource "google_pubsub_topic" "access_log" {
  name = "${var.name}-access-log"

  schema_settings {
    schema   = google_pubsub_schema.access_log_schema_v1.id
    encoding = "BINARY"
  }
}

resource "google_pubsub_subscription" "access_log_bq" {
  name                       = "${var.name}-access-log-bq"
  topic                      = google_pubsub_topic.access_log.name
  ack_deadline_seconds       = 10
  message_retention_duration = "604800s"

  bigquery_config {
    table               = "${var.project}:${google_bigquery_table.access_log.dataset_id}.${google_bigquery_table.access_log.table_id}"
    use_topic_schema    = true
    write_metadata      = false
    drop_unknown_fields = true
  }

  depends_on = [google_bigquery_table_iam_member.pubsub_sa_bigquery]
}

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

resource "google_bigquery_dataset" "access_log" {
  dataset_id    = "${var.name}_access_log"
  friendly_name = "${var.name} access log"
  description   = "${var.name} access log dataset"
  location      = "US"
}

resource "google_bigquery_table" "access_log" {
  dataset_id = google_bigquery_dataset.access_log.dataset_id
  table_id   = "${var.name}_access_log"
  clustering = ["timestamp"]
  schema     = file("access_log_schema/v1.json")

  time_partitioning {
    expiration_ms            = 31536000000 # 1 year
    type                     = "DAY"
    field                    = "timestamp"
    require_partition_filter = true
  }
}

resource "google_storage_bucket" "geolite2" {
  project                     = var.project
  name                        = var.geolite2_bucket
  location                    = "US-CENTRAL1"
  storage_class               = "STANDARD"
  uniform_bucket_level_access = true

  labels = {
    service = var.name
  }
}
