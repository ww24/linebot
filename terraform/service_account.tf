resource "google_service_account" "linebot" {
  account_id   = "linebot"
  display_name = "linebot Service Account"
}

resource "google_project_iam_member" "firestore" {
  role   = "roles/datastore.user"
  member = "serviceAccount:${google_service_account.linebot.email}"
}

resource "google_project_iam_member" "cloudtrace" {
  role   = "roles/cloudtrace.agent"
  member = "serviceAccount:${google_service_account.linebot.email}"
}


resource "google_project_iam_member" "cloudprofiler" {
  role   = "roles/cloudprofiler.agent"
  member = "serviceAccount:${google_service_account.linebot.email}"
}
