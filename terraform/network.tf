resource "google_compute_network" "main" {
  name                    = "main"
  description             = "main vpc network"
  routing_mode            = "GLOBAL"
  auto_create_subnetworks = false
}
