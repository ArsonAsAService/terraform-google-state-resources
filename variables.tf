// Provider variables
variable "gcp_region" {
  type        = string
  description = "The GCP region to create the resources in."
}

variable "gcp_project" {
  type        = string
  description = "The GCP project to create the resources in."
}

// Bucket variables
variable "state_bucket_name" {
  type        = string
  description = "The name of the GCS bucket that gets created to store Terraform state."
  default     = "terraform-state"
}

variable "log_bucket_name" {
  type        = string
  description = "The name of the GCS bucket that gets created to store logs of the state bucket."
  default     = "terraform-state-logs"
}

variable "bucket_location" {
  type        = string
  description = "The GCS location to create the storage bucket in."
  default     = "US"
}

variable "days_to_retain_versions" {
  type        = number
  description = "The number of days to retains previous versions of state files. Files older than this will be deleted by a lifecycle rule."
  default     = 90
}

variable "days_to_retain_logs" {
  type        = number
  description = "The number of days to retain logs for the state bucket."
  default     = 90
}

// KMS variables
variable "state_kms_key_ring_name" {
  type        = string
  description = "The name of the KMS Key Ring used to store keys for encrypting the state bucket."
  default     = "terraform-state-key-ring"
}

variable "state_kms_key_name" {
  type        = string
  description = "The name of the KMS key used to encrypt the state bucket."
  default     = "terraform-state-key"
}

variable "state_kms_key_rotation_period" {
  type        = string
  description = "How frequently to rotate the state encryption key. The default is 7 days."
  default     = "604800s"
}

variable "logs_kms_key_ring_name" {
  type        = string
  description = "The name of the KMS Key Ring used to store keys for encrypting the log bucket."
  default     = "terraform-log-key-ring"
}

variable "logs_kms_key_name" {
  type        = string
  description = "The name of the KMS key used to encrypt the log bucket."
  default     = "terraform-log-key"
}

variable "logs_kms_key_rotation_period" {
  type        = string
  description = "How frequently to rotate the log encryption key. The default is 7 days."
  default     = "604800s"
}

// Service Account variables
variable "create_service_account" {
  type        = bool
  description = "Whether or not to create a service account for Terraform to use."
  default     = false
}

variable "service_account_name" {
  type        = string
  description = "The name of the service account that will be created for Terraform to use."
  default     = "terraform"
}

variable "labels" {
  type        = map
  description = "The labels to add to the resources created by this module"
  default     = {}
}
