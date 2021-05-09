resource "google_service_account" "linebot" {
  account_id   = "linebot"
  display_name = "linebot Service Account"
}

resource "google_project_iam_member" "project" {
  role   = "roles/datastore.user"
  member = "serviceAccount:${google_service_account.linebot.email}"
}
