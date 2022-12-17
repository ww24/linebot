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

# MaxMind GeoLite2
resource "google_bigquery_dataset" "geolite2" {
  dataset_id    = "geolite2"
  friendly_name = "GeoLite2"
  description   = "geolite2 dataset"
  location      = "US"
}

resource "google_bigquery_table" "geolite2-city-blocks" {
  dataset_id = google_bigquery_dataset.geolite2.dataset_id
  table_id   = "GeoLite2-City-Blocks"
  schema     = file("geolite2/geolite2_city_blocks_schema.json")
}

resource "google_bigquery_table" "geolite2-city-locations" {
  dataset_id = google_bigquery_dataset.geolite2.dataset_id
  table_id   = "GeoLite2-City-Locations"
  schema     = file("geolite2/geolite2_city_locations_schema.json")
}

# Data Transfer Service (GCS => BigQuery)
resource "google_bigquery_data_transfer_config" "load-geolite2-city-blocks" {
  display_name           = "Load geolite2.GeoLite2-City-Blocks"
  location               = "US"
  data_source_id         = "google_cloud_storage"
  schedule               = "every day 17:00" # 02:00 JST
  destination_dataset_id = google_bigquery_dataset.geolite2.dataset_id
  service_account_name   = google_service_account.access-log.email
  params = {
    destination_table_name_template = google_bigquery_table.geolite2-city-blocks.table_id
    write_disposition               = "MIRROR"
    data_path_template              = "gs://${var.geolite2_bucket}/GeoLite2-City-Blocks-IPv*.csv"
    file_format                     = "CSV"
    max_bad_records                 = "0"
    skip_leading_rows               = "1"
    ignore_unknown_values           = "true"
  }

  depends_on = [google_bigquery_table.geolite2-city-blocks]
}

resource "google_bigquery_data_transfer_config" "load-geoLite2-city-locations" {
  display_name           = "Load geolite2.GeoLite2-City-Locations"
  location               = "US"
  data_source_id         = "google_cloud_storage"
  schedule               = "every day 17:00" # 02:00 JST
  destination_dataset_id = google_bigquery_dataset.geolite2.dataset_id
  service_account_name   = google_service_account.access-log.email
  params = {
    destination_table_name_template = google_bigquery_table.geolite2-city-locations.table_id
    write_disposition               = "MIRROR"
    data_path_template              = "gs://${var.geolite2_bucket}/GeoLite2-City-Locations-en.csv"
    file_format                     = "CSV"
    max_bad_records                 = "0"
    skip_leading_rows               = "1"
    ignore_unknown_values           = "true"
  }

  depends_on = [google_bigquery_table.geolite2-city-locations]
}

# Scheduled Queries
resource "google_bigquery_data_transfer_config" "transform-geolite2-city" {
  display_name           = "Transform geolite2.GeoLite2-City"
  location               = "US"
  data_source_id         = "scheduled_query"
  schedule               = "every day 18:00" # 03:00 JST
  destination_dataset_id = google_bigquery_dataset.geolite2.dataset_id
  service_account_name   = google_service_account.access-log.email
  params = {
    destination_table_name_template = "GeoLite2-City"
    write_disposition               = "WRITE_TRUNCATE"
    query                           = file("geolite2/transform_geolite2_city.sql")
  }
}

resource "google_bigquery_data_transfer_config" "snapshot-geolite2-city" {
  display_name         = "Snapshot geolite2.GeoLite2-City"
  location             = "US"
  data_source_id       = "scheduled_query"
  schedule             = "every day 19:00" # 04:00 JST
  service_account_name = google_service_account.access-log.email
  params = {
    query = file("geolite2/snapshot_geolite2_city.sql")
  }
}
