// jamfproapi_jamf_pro_sso_failover_test.go
// Jamf Pro Api - Jamf Pro SSO Failover Integration Testing
// api reference: https://developer.jamf.com/jamf-pro/reference/get_v1-sso-failover

/*
	Test Strategy:

Global Setup: The testing process begins with the initialization of the Jamf Pro HTTP client.
This global setup phase involves creating a temporary API test role and setting up a corresponding
API client for integration testing. This ensures that all tests run in a consistent
and controlled environment.

Individual Test Execution: Each integration test, managed by testing.T, is executed according
to a predefined test plan. These tests utilize the temporary API client and test role established
in the setup phase. The use of testing.T facilitates granular error reporting and isolated
testing of specific functionalities within the Jamf Pro integration.

Global Teardown: Upon completion of all tests, the suite enters the teardown phase. This
involves a systematic cleanup of all test-generated resources, including the removal of
the temporary API client and the test role. This step is crucial for ensuring that the
testing environment is reset and no residual data impacts subsequent test runs.
*/
package jamfpro_integration_test

import (
	"log"
	"os"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

var (
	client    *jamfpro.Client
	roleName  string
	apiClient *jamfpro.ApiIntegration
)

// TestMain function performs global setup of the test suite, executes
// the test suite and then performs the teardown of any resources created.
func TestMain(m *testing.M) {
	var err error
	// Initialize the Jamf Pro client
	client, err = setupJamfProSSOFailoverClient()
	if err != nil {
		log.Fatalf("Failed to setup Jamf Pro client: %v", err)
	}

	// Create an API role and get its name
	roleName, err := setupJamfProSSOFailoverTemporaryTestRole(client)
	if err != nil {
		log.Fatalf("Failed to setup temporary test role: %v", err)
	}

	// Create the Jamf API client with the role
	apiClient, err := setupJamfProSSOFailoverTemporaryTestAPIClient(client, []string{roleName})
	if err != nil {
		log.Fatalf("Failed to setup temporary API client: %v", err)
	}

	// Run the tests
	exitVal := m.Run()

	// Global Teardown: first remove the API client, then the role
	teardownApiIntegration(client, apiClient.DisplayName)
	teardownApiRole(client, roleName)

	os.Exit(exitVal)
}

// setupJamfProSSOFailoverClient initializes a Jamf Pro client for testing purposes.
// It reads the OAuth credentials from a specified configuration file and uses these credentials
// to create and return a new Jamf Pro client instance. This function is used to setup a client
// for various integration tests in this suite.
func setupJamfProSSOFailoverClient() (*jamfpro.Client, error) {

	// Define the path to the JSON configuration file
	configFilePath := "/Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-jamfpro/clientauth.json"

	// Load the client OAuth credentials from the configuration file
	authConfig, err := jamfpro.LoadClientAuthConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load client OAuth configuration: %v", err)
	}

	// Instantiate the default logger and set the desired log level
	logger := http_client.NewDefaultLogger()
	logLevel := http_client.LogLevelInfo // Set http_client logging level to debug

	// Configuration for the jamfpro
	config := jamfpro.Config{
		InstanceName: authConfig.InstanceName,
		LogLevel:     logLevel,
		Logger:       logger,
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
	}

	// Create a new jamfpro client instance
	client, err := jamfpro.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Jamf Pro client: %v", err)
	}

	return client, nil
}

// setupJamfProSSOFailoverTemporaryTestRole sets up a temporary test role in Jamf Pro for testing purposes.
// It creates a new Jamf API role with specific privileges necessary for the test scenarios.
// This function returns the name of the created role for use in subsequent tests.
func setupJamfProSSOFailoverTemporaryTestRole(client *jamfpro.Client) (string, error) {
	roleName := "go-api-sdk-jamfpro-apir-sso-failover"

	newRole := &jamfpro.APIRole{
		DisplayName: roleName,
		Privileges:  []string{"Read SSO Settings", "Update SSO Settings"},
	}

	newRole, err := client.CreateJamfApiRole(newRole)
	if err != nil {
		log.Fatalf("Error creating Jamf API role: %v", err)
	}

	// Log the creation of the new API Integration
	log.Printf("Created API Role with Display Name: %s\n", newRole.DisplayName)
	log.Printf("Created API Role has the following priviledges: %s\n", newRole.Privileges)
	return roleName, nil
}

// setupJamfProSSOFailoverTemporaryTestAPIClient creates a temporary API integration in Jamf Pro
// using specific API roles. This integration is used for testing and has privileges as defined in the roles.
func setupJamfProSSOFailoverTemporaryTestAPIClient(client *jamfpro.Client, roleNames []string) (*jamfpro.ApiIntegration, error) {
	// Define the new API Integration using the provided role names
	newApiIntegration := &jamfpro.ApiIntegration{
		AuthorizationScopes:        roleNames, // Use the role names provided
		DisplayName:                "SSOFailoverTemporaryTestAPIClient",
		Enabled:                    true,
		AccessTokenLifetimeSeconds: 1200,
	}

	// Call the function to create the new API Integration
	createdApiIntegration, err := client.CreateApiIntegration(newApiIntegration)
	if err != nil {
		log.Fatalf("Error creating API Integration: %v", err)
	}

	// Log the creation of the new API Integration
	log.Printf("Created API Integration with Display Name: %s\n", createdApiIntegration.DisplayName)

	return createdApiIntegration, nil
}

// TestJamfProIntegration_GetSSOFailoverSettings tests the GetSSOFailoverSettings functionality of the Jamf Pro client.
// It verifies that the SSO failover settings can be retrieved correctly and asserts that the
// returned settings contain expected data. This test validates the ability
// of the client to interact with the Jamf Pro API and retrieve SSO failover information.
func TestJamfProIntegration_GetSSOFailoverSettings(t *testing.T) {

	failoverSettings, err := client.GetSSOFailoverSettings()
	if err != nil {
		log.Fatalf("Failed to get SSO failover settings: %v", err)
	}

	// Assert that failover URL is not nil and not empty
	if failoverSettings.FailoverURL == "" {
		t.Errorf("Expected a failover URL, got an empty string")
	}

	// Assert that generation time is not zero (assuming it's a Unix timestamp or similar)
	if failoverSettings.GenerationTime == 0 {
		t.Errorf("Expected a non-zero generation time, got zero")
	}

	// Log the retrieved failover settings for verification
	log.Printf("Retrieved SSO Failover URL: %s", failoverSettings.FailoverURL)
	log.Printf("Retrieved Generation Time: %d", failoverSettings.GenerationTime)
}

// TestJamfProIntegration_UpdateFailoverUrl tests the UpdateFailoverUrl functionality of the Jamf Pro client.
// It verifies that the SSO failover URL can be updated correctly and asserts that the
// returned settings contain the new failover URL and a new generation time.
func TestJamfProIntegration_UpdateFailoverUrl(t *testing.T) {

	// Update the SSO failover URL
	updatedFailoverSettings, err := client.UpdateFailoverUrl()
	if err != nil {
		log.Fatalf("Error updating SSO failover URL: %v", err)
	}

	// Assert that the updated failover URL is not empty
	if updatedFailoverSettings.FailoverURL == "" {
		t.Errorf("Expected a non-empty failover URL, got an empty string")
	}

	// Assert that the generation time is updated (not zero)
	if updatedFailoverSettings.GenerationTime == 0 {
		t.Errorf("Expected a non-zero generation time, got zero")
	}

	// Log the updated failover settings for verification
	log.Printf("Updated SSO Failover URL: %s", updatedFailoverSettings.FailoverURL)
	log.Printf("New Generation Time: %d", updatedFailoverSettings.GenerationTime)
}

// Helper function to delete API integration
func teardownApiIntegration(client *jamfpro.Client, integrationName string) {
	if err := client.DeleteApiIntegrationByName(integrationName); err != nil {
		log.Fatalf("Failed to delete API integration: %v", err) // Exits the program if there's an error
	} else {
		log.Printf("API integration '%s' deleted successfully", integrationName)
	}
}

// Helper function to delete API role
func teardownApiRole(client *jamfpro.Client, roleName string) {
	if err := client.DeleteJamfApiRoleByName(roleName); err != nil {
		log.Fatalf("Failed to delete API role: %v", err) // Exits the program if there's an error
	} else {
		log.Printf("API role '%s' deleted successfully", roleName)
	}
}
