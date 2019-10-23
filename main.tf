/**
 * # GCP Terraform State Resources
 *
 * ## Description
 *
 * This module describes the resources required to securely store Terraform
 * state in GCP. It creates:
 *
 * * A Cloud Storage bucket for storing state
 * * A Cloud KMS cryptoKey (and Keyring) used to encrypt state
 * * A Cloud Storage bucket for storing logs from the state bucket
 * * A Cloud KMS cryptoKey (and Keyring) used to encrypt logs
 * * A Service Account that has permission to edit objects in the GCS bucket and
 *   use the cryptoKeys to encrypt/decrypt state
 *
 * Creating a service account for Terraform to use to interact with GCP is
 * beyond the scope of this module as the required permissions are highly
 * specific to the org/project.
 */

locals {
  project_services = [
    "storage-component.googleapis.com",
    "storage-api.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com",
    "cloudkms.googleapis.com",
  ]

  default_labels = {
    "terraform_managed" = "true"
    "terraform"         = "true"
  }
}

locals {
  labels = merge(var.labels, local.default_labels)
}

data "google_project" "proj" {
  project_id = var.gcp_project
}

resource "google_project_service" "services" {
  count   = length(local.project_services)
  project = data.google_project.proj.project_id

  service            = element(local.project_services, count.index)
  disable_on_destroy = false
}

resource "google_kms_key_ring" "state" {
  name     = var.state_kms_key_ring_name
  location = "global"

  depends_on = [
    google_project_service.services,
  ]
}

resource "google_kms_crypto_key" "state" {
  name            = var.state_kms_key_name
  key_ring        = google_kms_key_ring.state.self_link
  rotation_period = var.state_kms_key_rotation_period

  depends_on = [
    google_project_service.services,
  ]
}

resource "google_kms_key_ring" "logs" {
  name     = var.logs_kms_key_ring_name
  location = "global"

  depends_on = [
    google_project_service.services,
  ]
}

resource "google_kms_crypto_key" "logs" {
  name            = var.logs_kms_key_name
  key_ring        = google_kms_key_ring.logs.self_link
  rotation_period = var.logs_kms_key_rotation_period

  depends_on = [
    google_project_service.services,
  ]
}

resource "google_storage_bucket" "logs" {
  project = data.google_project.proj.project_id

  name               = var.log_bucket_name
  location           = var.bucket_location
  bucket_policy_only = true
  labels             = local.labels

  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      age = var.days_to_retain_logs
    }
  }

  encryption {
    default_kms_key_name = google_kms_crypto_key.logs.self_link
  }

  depends_on = [
    google_project_service.services,
  ]
}


resource "google_storage_bucket" "state" {
  project = data.google_project.proj.project_id

  name               = var.state_bucket_name
  location           = var.bucket_location
  bucket_policy_only = true
  labels             = local.labels

  versioning {
    enabled = true
  }

  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      age        = var.days_to_retain_versions
      with_state = "ARCHIVED"
    }
  }

  logging {
    log_bucket = google_storage_bucket.logs.name
  }

  encryption {
    default_kms_key_name = google_kms_crypto_key.state.self_link
  }

  depends_on = [
    google_project_service.services,
  ]
}

resource "google_service_account" "terraform" {
  count   = "${var.create_service_account ? 1 : 0}"
  project = data.google_project.proj.project_id

  account_id   = var.service_account_name
  display_name = "Service Account for Terraform"

  depends_on = [
    google_project_service.services,
  ]
}

resource "google_kms_key_ring_iam_member" "terraform" {
  count = "${var.create_service_account ? 1 : 0}"

  key_ring_id = google_kms_key_ring.state.self_link
  role        = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member      = "serviceAccount:${google_service_account.terraform[0].email}"

  depends_on = [
    google_project_service.services,
  ]
}

resource "google_storage_bucket_iam_member" "terraform" {
  count = "${var.create_service_account ? 1 : 0}"

  bucket = google_storage_bucket.state.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.terraform[0].email}"

  depends_on = [
    google_project_service.services,
  ]
}
