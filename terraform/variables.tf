variable "location" {
  type    = string
  default = "asia-northeast1"
}

variable "project" {
  type        = string
  description = "GCP Project ID"
}

variable "google_credentials" {
  type        = string
  description = "GCP Service Account (credential json value)"
}

variable "name" {
  type    = string
  default = "linebot"
}

variable "gar_repository" {
  type    = string
  default = "ww24"
}

variable "image_name" {
  type    = string
  default = "linebot"
}

variable "image_tag" {
  type    = string
  default = "latest"
}

// application environments
variable "line_channel_secret" {
  type        = string
  description = "LINE Channel Secret"
}

variable "line_channel_access_token" {
  type        = string
  description = "LINE Channel Access Token"
}

variable "allow_conv_ids" {
  type        = string
  description = "Allowed list, conversation ids"
}

variable "cloud_tasks_queue" {
  type    = string
  default = "linebot"
}

variable "service_endpoint" {
  type        = string
  description = "Cloud Run Service Endpoint (https://*.a.run.app)"
}
