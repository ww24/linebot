data "google_cloud_run_service" "linebot" {
  name     = local.name
  location = local.location
}

locals {
  current_image = data.google_cloud_run_service.linebot.template != null ? data.google_cloud_run_service.linebot.template[0].spec[0].containers[0].image : null
  new_image     = "${local.location}-docker.pkg.dev/${var.project}/${local.gar_repository}/${local.name}:${var.image_tag}"
  image         = (local.current_image != null && var.image_tag == "latest") ? local.current_image : local.new_image
  image_tag     = split(":", local.image)[1]
}

resource "google_cloud_run_service" "linebot" {
  name     = local.name
  location = local.location
  project  = var.project

  template {
    spec {
      service_account_name = google_service_account.linebot.email

      timeout_seconds = 120
      containers {
        image = local.image

        resources {
          limits = {
            cpu    = "1000m"
            memory = "300Mi"
          }
        }

        env {
          name = "LINE_CHANNEL_SECRET"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.line-channel-secret.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "LINE_CHANNEL_ACCESS_TOKEN"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.line-channel-access-token.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name  = "ALLOW_CONV_IDS"
          value = var.allow_conv_ids
        }

        env {
          name  = "CLOUD_TASKS_LOCATION"
          value = google_cloud_tasks_queue.linebot.location
        }

        env {
          name  = "CLOUD_TASKS_QUEUE"
          value = google_cloud_tasks_queue.linebot.name
        }

        env {
          name  = "SERVICE_ENDPOINT"
          value = var.service_endpoint
        }

        env {
          name  = "STORAGE_IMAGE_BUCKET"
          value = google_storage_bucket.image.name
        }

        env {
          name  = "INVOKER_SERVICE_ACCOUNT_ID"
          value = google_service_account.invoker.unique_id
        }

        env {
          name  = "INVOKER_SERVICE_ACCOUNT_EMAIL"
          value = google_service_account.invoker.email
        }

        env {
          name  = "OTEL_SAMPLING_RATE"
          value = local.linebot_otel_sampling_rate
        }

        env {
          name  = "ACCESS_LOG_TOPIC"
          value = google_pubsub_topic.access_log_v1.name
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1"
        "autoscaling.knative.dev/minScale" = "1"
      }

      labels = {
        service = local.name
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true
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
