resource "google_cloud_run_v2_job" "screenshot" {
  name     = local.name_screenshot
  location = "asia-northeast1"

  template {
    parallelism = 0
    task_count  = 1

    template {
      service_account = google_service_account.screenshot.email
      timeout         = "120s"
      max_retries     = 2

      containers {
        image = "${local.location}-docker.pkg.dev/${var.project}/${local.gar_repository}/${local.name_screenshot}:${local.image_tag}"

        resources {
          limits = {
            cpu    = "1000m" # minimum
            memory = "512Mi"
          }
        }

        env {
          name  = "SCREENSHOT_BROWSER_TIMEOUT"
          value = local.screenshot_browser_timeout
        }

        env {
          name  = "OTEL_SAMPLING_RATE"
          value = local.linebot_otel_sampling_rate
        }

        env {
          name  = "SCREENSHOT_TARGET_URL"
          value = var.screenshot_target_url
        }

        env {
          name  = "SCREENSHOT_TARGET_SELECTOR"
          value = var.screenshot_target_selector
        }

        env {
          name  = "STORAGE_IMAGE_BUCKET"
          value = google_storage_bucket.image.name
        }
      }
    }
  }
}
