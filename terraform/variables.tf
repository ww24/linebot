variable "project" {
  type        = string
  description = "GCP Project ID"
}

variable "google_credentials" {
  type        = string
  description = "GCP Service Account (credential json value)"
}

variable "image_tag" {
  type    = string
  default = "latest"
}

variable "image_bucket" {
  type        = string
  description = "image bucket name"
}

variable "geolite2_bucket" {
  type        = string
  description = "MaxMind GeoLite2 bucket name"
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

variable "service_endpoint" {
  type        = string
  description = "Cloud Run Service Endpoint (https://*.a.run.app)"
}

variable "maxmind_license_key" {
  type        = string
  description = "MaxMind License Key"
}

locals {
  # GCP location
  location = "asia-northeast1"

  # Application name
  name            = "linebot"
  name_screenshot = "screenshot"

  # Google Artifact Registry
  gar_repository = "ww24"

  # OpenTelemetry sampling rate
  linebot_otel_sampling_rate = "1"
}
