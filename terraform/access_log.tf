resource "google_pubsub_schema" "access_log_schema_v1" {
  name       = "${local.name}-access-log-v1"
  type       = "AVRO"
  definition = file("access_log_schema/v1.avsc")
}

resource "google_pubsub_topic" "access_log_v1" {
  name = "${local.name}-access-log-v1"

  schema_settings {
    schema   = google_pubsub_schema.access_log_schema_v1.id
    encoding = "BINARY"
  }
}

resource "google_pubsub_subscription" "access_log_v1_bq" {
  name                       = "${local.name}-access-log-v1-bq"
  topic                      = google_pubsub_topic.access_log_v1.name
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

resource "google_bigquery_dataset" "access_log" {
  dataset_id    = "${local.name}_access_log"
  friendly_name = "${local.name} access log"
  description   = "${local.name} access log dataset"
  location      = "US"
}

resource "google_bigquery_table" "access_log" {
  dataset_id               = google_bigquery_dataset.access_log.dataset_id
  table_id                 = "access_log"
  clustering               = ["timestamp"]
  schema                   = file("access_log_schema/v1.json")
  require_partition_filter = true

  time_partitioning {
    expiration_ms = 31536000000 # 1 year
    type          = "DAY"
    field         = "timestamp"
  }
}

resource "google_bigquery_routine" "with_geolocation" {
  dataset_id   = google_bigquery_dataset.access_log.dataset_id
  routine_id   = "with_geolocation"
  routine_type = "TABLE_VALUED_FUNCTION"
  language     = "SQL"
  definition_body = templatefile("geolite2/function_with_geolocation.sql", {
    project = var.project,
    dataset = google_bigquery_dataset.access_log.dataset_id,
  })
  arguments {
    name          = "since"
    argument_kind = "FIXED_TYPE"
    data_type     = jsonencode({ "typeKind" : "TIMESTAMP" })
  }
  arguments {
    name          = "until"
    argument_kind = "FIXED_TYPE"
    data_type     = jsonencode({ "typeKind" : "TIMESTAMP" })
  }
}

resource "google_storage_bucket" "geolite2" {
  project                     = var.project
  name                        = var.geolite2_bucket
  location                    = "US-CENTRAL1"
  storage_class               = "STANDARD"
  uniform_bucket_level_access = true

  labels = {
    service = local.name
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
  external_data_configuration {
    connection_id         = google_bigquery_connection.geolite2.name
    autodetect            = false
    ignore_unknown_values = true
    source_uris           = ["gs://${var.geolite2_bucket}/GeoLite2-City-Blocks-IPv*.csv"]
    source_format         = "CSV"
    csv_options {
      quote             = ""
      skip_leading_rows = 1
    }
  }

  lifecycle {
    ignore_changes = [external_data_configuration[0].connection_id]
  }
}

resource "google_bigquery_table" "geolite2-city-locations" {
  dataset_id = google_bigquery_dataset.geolite2.dataset_id
  table_id   = "GeoLite2-City-Locations"
  schema     = file("geolite2/geolite2_city_locations_schema.json")
  external_data_configuration {
    connection_id         = google_bigquery_connection.geolite2.name
    autodetect            = false
    ignore_unknown_values = true
    source_uris           = ["gs://${var.geolite2_bucket}/GeoLite2-City-Locations-en.csv"]
    source_format         = "CSV"
    csv_options {
      quote             = ""
      skip_leading_rows = 1
    }
  }

  lifecycle {
    ignore_changes = [external_data_configuration[0].connection_id]
  }
}

# Scheduled Queries
resource "google_bigquery_data_transfer_config" "transform-geolite2-city" {
  display_name           = "Transform geolite2.GeoLite2-City"
  location               = "US"
  data_source_id         = "scheduled_query"
  schedule               = "every day 17:00" # 02:00 JST
  destination_dataset_id = google_bigquery_dataset.geolite2.dataset_id
  service_account_name   = google_service_account.access-log.email
  params = {
    destination_table_name_template = "GeoLite2-City"
    write_disposition               = "WRITE_TRUNCATE"
    query                           = file("geolite2/transform_geolite2_city.sql")
  }
}

resource "google_bigquery_data_transfer_config" "snapshot-geolite2-city" {
  display_name         = "Snapshot geolite2.GeoLite2_City_YYYYMMDD"
  location             = "US"
  data_source_id       = "scheduled_query"
  schedule             = "every day 18:00" # 03:00 JST
  service_account_name = google_service_account.access-log.email
  params = {
    query = file("geolite2/snapshot_geolite2_city.sql")
  }
}

resource "google_cloud_run_v2_job" "maxmind" {
  name     = "maxmind"
  location = "us-central1"

  template {
    parallelism = 1
    task_count  = 1

    template {
      service_account = google_service_account.maxmind.email
      timeout         = "60s"
      max_retries     = 1

      containers {
        image = "us.gcr.io/google.com/cloudsdktool/google-cloud-cli:554.0.0-alpine"

        resources {
          limits = {
            cpu    = "1000m" # minimum
            memory = "1Gi" # minimum
          }
        }

        command = ["bash"]
        args = [
          "-euc",
          <<-EOT
          curl -sL "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City-CSV&license_key=$${MAXMIND_LICENSE_KEY}&suffix=zip" -o City.zip
          curl -sL "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City-CSV&license_key=$${MAXMIND_LICENSE_KEY}&suffix=zip.sha256" | awk '{print $1"  City.zip"}' > shasum.txt
          sha256sum -c shasum.txt
          unzip -p City.zip "GeoLite2-**/GeoLite2-City-Blocks-IPv4.csv" > GeoLite2-City-Blocks-IPv4.csv
          unzip -p City.zip "GeoLite2-**/GeoLite2-City-Blocks-IPv6.csv" > GeoLite2-City-Blocks-IPv6.csv
          unzip -p City.zip "GeoLite2-**/GeoLite2-City-Locations-en.csv" > GeoLite2-City-Locations-en.csv
          gsutil cp GeoLite2-City-Blocks-IPv*.csv GeoLite2-City-Locations-en.csv "$${DESTINATION_URL}"
          EOT
        ]

        env {
          name = "MAXMIND_LICENSE_KEY"
          value_source {
            secret_key_ref {
              secret  = google_secret_manager_secret.maxmind-license-key.secret_id
              version = "latest"
            }
          }
        }

        env {
          name  = "DESTINATION_URL"
          value = google_storage_bucket.geolite2.url
        }
      }
    }
  }
}

locals {
  maxmind_job = {
    location = google_cloud_run_v2_job.maxmind.location
    name     = google_cloud_run_v2_job.maxmind.name
  }
}

resource "google_cloud_scheduler_job" "maxmind" {
  name             = "maxmind"
  description      = "MaxMind Scheduler"
  schedule         = "0 1 * * *"
  time_zone        = "Asia/Tokyo"
  attempt_deadline = "30s"
  region           = local.maxmind_job.location

  http_target {
    http_method = "POST"
    uri         = "https://${local.maxmind_job.location}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${var.project}/jobs/${local.maxmind_job.name}:run"

    oauth_token {
      service_account_email = google_service_account.invoker.email
    }
  }
}
