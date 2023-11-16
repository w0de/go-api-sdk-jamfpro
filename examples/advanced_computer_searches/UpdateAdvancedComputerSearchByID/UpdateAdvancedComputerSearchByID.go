package main

import (
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

	updatedSearch, err := client.UpdateAdvancedComputerSearchByID(7, &jamfpro.ResponseAdvancedComputerSearch{
		Name:   "Advanced Search Name",
		ViewAs: "Standard Web Page",
		Criteria: []jamfpro.AdvancedComputerSearchesCriteria{
			{
				Size: 1,
				Criterion: jamfpro.CriterionDetail{
					Name:         "Last Inventory Update",
					Priority:     0,
					AndOr:        "and",
					SearchType:   "more than x days ago",
					Value:        "7",
					OpeningParen: false,
					ClosingParen: false,
				},
			},
		},
		DisplayFields: []jamfpro.AdvancedComputerSearchesDisplayField{
			{
				Size: 1,
				DisplayField: jamfpro.DisplayFieldDetail{
					Name: "IP Address",
				},
			},
		},
		Site: jamfpro.AdvancedComputerSearchesSiteDetail{
			ID:   -1,
			Name: "None",
		},
	})
	if err != nil {
		fmt.Println("Error updating advanced computer search by ID:", err)
		return
	}

	fmt.Println("Updated Search:", updatedSearch)
}