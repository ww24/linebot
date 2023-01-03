resource "google_compute_network" "main" {
  name                    = "main"
  description             = "main vpc network"
  routing_mode            = "GLOBAL"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  for_each      = local.vpc
  name          = "main-${each.key}"
  ip_cidr_range = each.value.ip_cidr_range
  region        = each.value.region
  network       = google_compute_network.main.id
}

resource "google_compute_router" "main" {
  name    = "main"
  network = google_compute_network.main.id
  region  = google_compute_subnetwork.subnet["serverless"].region
}

resource "google_compute_address" "nat-ips" {
  count  = 2
  name   = "main-router-nat-ips-${count.index}"
  region = google_compute_router.main.region
}

resource "google_compute_router_nat" "main" {
  name                               = "main"
  router                             = google_compute_router.main.name
  region                             = google_compute_router.main.region
  nat_ip_allocate_option             = "MANUAL_ONLY"
  nat_ips                            = google_compute_address.nat-ips.*.self_link
  source_subnetwork_ip_ranges_to_nat = "LIST_OF_SUBNETWORKS"
  subnetwork {
    name                    = google_compute_subnetwork.subnet["serverless"].id
    source_ip_ranges_to_nat = ["ALL_IP_RANGES"]
  }
}

resource "google_vpc_access_connector" "main" {
  name           = "main"
  machine_type   = "f1-micro"
  min_instances  = 2
  max_instances  = 3
  min_throughput = 200
  max_throughput = 300
  region         = google_compute_subnetwork.subnet["serverless"].region
  subnet {
    name = google_compute_subnetwork.subnet["serverless"].name
  }
}
