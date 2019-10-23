provider "google" {
  version = "~> 2.17"
  region  = var.gcp_region
  project = var.gcp_project
}
