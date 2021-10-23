data "google_cloud_run_service" "linebot" {
  name     = var.name
  location = var.location
}

locals {
  current_image = data.google_cloud_run_service.linebot.template != null ? data.google_cloud_run_service.linebot.template[0].spec[0].containers[0].image : null
  new_image     = "${var.location}-docker.pkg.dev/${var.project}/${var.gar_repository}/${var.image_name}:${var.image_tag}"
  image         = (local.current_image != null && var.image_tag == "latest") ? local.current_image : local.new_image
  image_tag     = split(":", image)[1]
}

resource "google_cloud_run_service" "linebot" {
  name     = var.name
  location = var.location
  project  = var.project

  template {
    spec {
      service_account_name = google_service_account.linebot.email

      timeout_seconds = 60
      containers {
        image = local.image

        resources {
          limits = {
            cpu    = "1000m"
            memory = "300Mi"
          }
        }

        env {
          name  = "LINE_CHANNEL_SECRET"
          value = var.line_channel_secret
        }

        env {
          name  = "LINE_CHANNEL_ACCESS_TOKEN"
          value = var.line_channel_access_token
        }

        env {
          name  = "ALLOW_CONV_IDS"
          value = var.allow_conv_ids
        }
      }
    }

    metadata {
      # The revision name must be prefixed by the name of the enclosing Service or Configuration with a trailing -
      # Resource name must use only lowercase letters, numbers and '-'. Must begin with a letter and cannot end with a '-'. Maximum length is 63 characters.
      name = local.image_tag == "latest" ? null : "${var.name}-v${replace(var.image_tag, ".", "-")}"

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

  autogenerate_revision_name = local.image_tag == "latest"
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.linebot.location
  project  = google_cloud_run_service.linebot.project
  service  = google_cloud_run_service.linebot.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}
