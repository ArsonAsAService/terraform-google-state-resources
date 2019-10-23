package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

// GetProjectAndRegionAndTestName returns the ID of the a GCP Project, a random
// GCP region, and a name to use for the test run
func GetProjectAndRegionAndTestName(t *testing.T, namePrefix string) (string, string, string) {
	projectID := gcp.GetGoogleProjectIDFromEnvVar(t)
	regionID := gcp.GetRandomRegion(t, projectID, []string{}, []string{})

	randomSuffix := strings.ToLower(random.UniqueId())
	testName := fmt.Sprintf("%s-%s", namePrefix, randomSuffix)
	return projectID, regionID, testName
}

// SetupTestFolder copies the root folder to a temporary location and then
// returns the absolute path to the supplied folder relative to the root. This
// enables tests to be run in parallel without clobbering state folders.
func SetupTestFolder(t *testing.T, rootFolder string, terraformFolderRelativeToRoot string) string {
	return test_structure.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)
}

// GetTerraformEnvVars returns environment variables that will be passed to the
// shell that terratest uses to execute Terraform.
func GetTerraformEnvVars(t *testing.T) map[string]string {
	vars := map[string]string{
		"GOOGLE_CREDENTIALS": gcp.GetGoogleCredentialsFromEnvVar(t),
		// Since this project has asdf support, ensure that asdf picks the right
		// Terraform version when terratest invokes terraform, regardless of what
		// the global version is set to.
		"ASDF_TERRAFORM_VERSION": "0.12.12",
	}
	return vars
}

// GetCryptoKeyID returns a string with the fully qualified ID of the provided
// by constructing a string using the project ID and key ring ID. It assumes
// that the keyring is in the "global" location.
func GetCryptoKeyID(projectID string, keyRingID string, keyName string) string {
	return fmt.Sprintf("projects/%s/locations/global/keyRings/%s/cryptoKeys/%s", projectID, keyRingID, keyName)
}
