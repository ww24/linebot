resource "google_app_engine_application" "app" {
  project     = var.project
  location_id = var.location
}

resource "google_cloud_tasks_queue" "linebot" {
  name     = var.name
  location = var.location

  rate_limits {
    max_dispatches_per_second = 2
    max_concurrent_dispatches = 1
  }

  retry_config {
    max_attempts       = 100
    min_backoff        = "0.100s"
    max_backoff        = "10s"
    max_doublings      = 4
    max_retry_duration = "30s"
  }

  stackdriver_logging_config {
    sampling_ratio = 1
  }
}
