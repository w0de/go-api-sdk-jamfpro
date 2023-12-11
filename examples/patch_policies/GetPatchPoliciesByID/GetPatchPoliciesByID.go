package main

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client" // Import http_client for logging
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

func main() {
	// Define the path to the JSON configuration file
	configFilePath := "/Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-jamfpro/clientauth.json"

	// Load the client OAuth credentials from the configuration file
	authConfig, err := jamfpro.LoadClientAuthConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load client OAuth configuration: %v", err)
	}

	// Instantiate the default logger and set the desired log level
	logger := http_client.NewDefaultLogger()
	logLevel := http_client.LogLevelDebug // LogLevelNone // LogLevelWarning // LogLevelInfo  // LogLevelDebug

	// Configuration for the jamfpro
	config := jamfpro.Config{
		InstanceName:       authConfig.InstanceName,
		OverrideBaseDomain: authConfig.OverrideBaseDomain,
		LogLevel:           logLevel,
		Logger:             logger,
		ClientID:           authConfig.ClientID,
		ClientSecret:       authConfig.ClientSecret,
	}

	// Create a new jamfpro client instance
	client, err := jamfpro.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Jamf Pro client: %v", err)
	}

	// The ID of the patch policy you want to retrieve
	patchPolicyID := 1 // Replace with the actual ID you want to retrieve

	// Call the GetPatchPoliciesByID function
	patchPolicy, err := client.GetPatchPoliciesByID(patchPolicyID)
	if err != nil {
		log.Fatalf("Error fetching patch policy by ID: %v", err)
	}

	// Convert the response into pretty XML for printing
	output, err := xml.MarshalIndent(patchPolicy, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling patch policy to XML: %v", err)
	}

	// Print the pretty XML
	fmt.Printf("Patch Policy (ID: %d):\n%s\n", patchPolicyID, string(output))
}