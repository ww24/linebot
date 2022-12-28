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
