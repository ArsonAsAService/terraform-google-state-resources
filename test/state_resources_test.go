package test

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/storage"

	"github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestNoServiceAccount(t *testing.T) {
	t.Parallel()

	projectID, regionID, testName := GetProjectAndRegionAndTestName(t, "no-sa")
	// Copy everything to a temporary location
	testFolder := SetupTestFolder(t, "..", ".")

	// Build some expectations
	expectedStateBucketName := fmt.Sprintf("%s-state", testName)
	expectedLogBucketName := fmt.Sprintf("%s-logs", testName)
	expectedStateKeyRingName := fmt.Sprintf("%s-state-key-ring", testName)
	expectedStateKeyName := fmt.Sprintf("%s-state-key", testName)
	expectedLogKeyRingName := fmt.Sprintf("%s-logs-key-ring", testName)
	expectedLogKeyName := fmt.Sprintf("%s-logs-key", testName)
	expectedLabels := map[string]string{
		"terraform":         "true",
		"terraform_managed": "true",
		"foo":               "bar",
	}
	expectedLifecycle := storage.Lifecycle{
		Rules: []storage.LifecycleRule{
			storage.LifecycleRule{
				Action: storage.LifecycleAction{
					Type: "Delete",
				},
				Condition: storage.LifecycleCondition{
					AgeInDays: 90,
					Liveness:  2,
				},
			},
		},
	}
	expectedStateEncryption := storage.BucketEncryption{
		DefaultKMSKeyName: GetCryptoKeyID(projectID, expectedStateKeyRingName, expectedStateKeyName),
	}
	expectedStateLogging := storage.BucketLogging{
		LogBucket:       expectedLogBucketName,
		LogObjectPrefix: expectedStateBucketName,
	}
	expectedLogEncryption := storage.BucketEncryption{
		DefaultKMSKeyName: GetCryptoKeyID(projectID, expectedLogKeyRingName, expectedLogKeyName),
	}

	// Configure options supplied to terraform
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: testFolder,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"gcp_region":              regionID,
			"gcp_project":             projectID,
			"state_bucket_name":       expectedStateBucketName,
			"log_bucket_name":         expectedLogBucketName,
			"state_kms_key_ring_name": expectedStateKeyRingName,
			"state_kms_key_name":      expectedStateKeyName,
			"logs_kms_key_ring_name":  expectedLogKeyRingName,
			"logs_kms_key_name":       expectedLogKeyName,
			"service_account_name":    fmt.Sprintf("%s-terraform", testName),
			"labels": map[string]string{
				"foo": "bar",
			},
		},

		EnvVars: GetTerraformEnvVars(t),
	}

	// At the end of the test, run `terraform destroy` to clean up any resources
	// that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if
	// there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Reach into the GCS client library since terratest does not have support for
	// introspecting most of the properties of GCS buckets.
	ctx := context.Background()
	client, _ := storage.NewClient(ctx)

	// State bucket assertions
	stateBucket := client.Bucket(expectedStateBucketName)
	stateBucketAttrs, err := stateBucket.Attrs(ctx)
	if err != nil {
		t.Fatal(err)
	}

	gcp.AssertStorageBucketExists(t, expectedStateBucketName)
	assert.Equal(t, "US", stateBucketAttrs.Location, "State bucket should be in the US location.")
	assert.Equal(t, true, stateBucketAttrs.VersioningEnabled, true, "State bucket should have versioning enabled.")
	assert.Equal(t, expectedLifecycle, stateBucketAttrs.Lifecycle, "State bucket should have lifecycle rules")
	assert.Equal(t, &expectedStateEncryption, stateBucketAttrs.Encryption, "State bucket should have encryption enabled")
	assert.Equal(t, expectedLabels, stateBucketAttrs.Labels, "State bucket should have terraform labels.")
	assert.Equal(t, &expectedStateLogging, stateBucketAttrs.Logging, "State bucket should have logging to the log bucket configured.")

	// Log bucket assertions
	logBucket := client.Bucket(expectedLogBucketName)
	logBucketAttrs, err := logBucket.Attrs(ctx)
	if err != nil {
		t.Fatal(err)
	}

	gcp.AssertStorageBucketExists(t, expectedLogBucketName)
	assert.Equal(t, "US", logBucketAttrs.Location, "Log bucket should be in the US location.")
	assert.Equal(t, false, logBucketAttrs.VersioningEnabled, true, "Log bucket should not have versioning enabled.")
	assert.Equal(t, expectedLifecycle, logBucketAttrs.Lifecycle, "Log bucket should have lifecycle rules")
	assert.Equal(t, &expectedLogEncryption, logBucketAttrs.Encryption, "Log bucket should have encryption enabled")
	assert.Equal(t, expectedLabels, logBucketAttrs.Labels, "Log bucket should have terraform labels.")
	assert.Empty(t, logBucketAttrs.Logging, "Log bucket should not have logging configured.")
}
