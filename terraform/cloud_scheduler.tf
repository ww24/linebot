resource "google_cloud_scheduler_job" "scheduler" {
  name             = local.name
  description      = "${local.name} scheduler"
  schedule         = "0 * * * *"
  time_zone        = "Asia/Tokyo"
  attempt_deadline = "180s"

  http_target {
    http_method = "POST"
    uri         = "${google_cloud_run_service.linebot.status[0].url}/scheduler"
    headers = {
      content-type = "application/json"
    }

    oidc_token {
      service_account_email = google_service_account.invoker.email
      audience              = google_cloud_run_service.linebot.status[0].url
    }
  }

  retry_config {
    retry_count          = 3
    min_backoff_duration = "1s"
    max_backoff_duration = "10s"
    max_doublings        = 2
  }
}

locals {
  screenshot_job = {
    location = google_cloud_run_v2_job.screenshot.location
    name     = google_cloud_run_v2_job.screenshot.name
  }
}

resource "google_cloud_scheduler_job" "screenshot" {
  name             = local.name_screenshot
  description      = "${local.name_screenshot} scheduler"
  schedule         = "0 * * * *"
  time_zone        = "Asia/Tokyo"
  attempt_deadline = "180s"
  region           = local.screenshot_job.location

  retry_config {
    max_backoff_duration = "3600s"
    min_backoff_duration = "5s"
    max_retry_duration   = "0s"
    max_doublings        = 5
    retry_count          = 0
  }

  http_target {
    http_method = "POST"
    uri         = "https://${local.screenshot_job.location}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${var.project}/jobs/${local.screenshot_job.name}:run"

    oauth_token {
      service_account_email = google_service_account.invoker.email
    }
  }
}
