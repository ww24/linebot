data "google_cloud_run_service" "screenshot" {
  name     = var.name_screenshot
  location = var.location
}

locals {
  ss_current_image = data.google_cloud_run_service.screenshot.template != null ? data.google_cloud_run_service.screenshot.template[0].spec[0].containers[0].image : null
  ss_new_image     = "${var.location}-docker.pkg.dev/${var.project}/${var.gar_repository}/${var.image_name_screenshot}:${var.image_tag}"
  ss_image         = (local.ss_current_image != null && var.image_tag == "latest") ? local.ss_current_image : local.ss_new_image
}

resource "google_cloud_run_service" "screenshot" {
  name     = var.name_screenshot
  location = var.location
  project  = var.project

  template {
    spec {
      service_account_name = google_service_account.screenshot.email

      timeout_seconds = 60
      containers {
        image = local.ss_image

        resources {
          limits = {
            cpu    = "1000m"
            memory = "500Mi"
          }
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1"
        "autoscaling.knative.dev/minScale" = "1"
      }

      labels = {
        service = var.name
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true
}

resource "google_cloud_run_service_iam_policy" "invoker" {
  location    = google_cloud_run_service.screenshot.location
  project     = google_cloud_run_service.screenshot.project
  service     = google_cloud_run_service.screenshot.name
  policy_data = data.google_iam_policy.invoker.policy_data
}

data "google_iam_policy" "invoker" {
  binding {
    role = "roles/run.invoker"
    members = [
      "serviceAccount:${google_service_account.linebot.email}",
    ]
  }
}
