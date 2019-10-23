# GCP Terraform State Resources

## Description

This module describes the resources required to securely store Terraform
state in GCP. It creates:

* A Cloud Storage bucket for storing state
* A Cloud KMS cryptoKey (and Keyring) used to encrypt state
* A Cloud Storage bucket for storing logs from the state bucket
* A Cloud KMS cryptoKey (and Keyring) used to encrypt logs
* A Service Account that has permission to edit objects in the GCS bucket and
  use the cryptoKeys to encrypt/decrypt state

Creating a service account for Terraform to use to interact with GCP is
beyond the scope of this module as the required permissions are highly
specific to the org/project.
