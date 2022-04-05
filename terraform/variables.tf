variable "location" {
  type    = string
  default = "asia-northeast1"
}

variable "project" {
  type = string
}

// credentials json value
variable "google_credentials" {
  type = string
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
  type = string
}

variable "line_channel_access_token" {
  type = string
}

variable "allow_conv_ids" {
  type = string
}

variable "cloud_tasks_queue" {
  type    = string
  default = "linebot"
}

variable "service_endpoint" {
  type = string
}
