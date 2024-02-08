package main

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-http-client/httpclient"
	"github.com/deploymenttheory/go-api-http-client/logger"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

func main() {
	// Define the path to the JSON configuration file
	configFilePath := "/Users/dafyddwatkins/localtesting/clientauth.json"

	// Load the client OAuth credentials from the configuration file
	authConfig, err := jamfpro.LoadAuthConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load client OAuth configuration: %v", err)
	}

	// Instantiate the default logger and set the desired log level
	logLevel := logger.LogLevelWarn // LogLevelNone / LogLevelDebug / LogLevelInfo / LogLevelError

	// Configuration for the jamfpro
	config := httpclient.Config{
		InstanceName: authConfig.InstanceName,
		Auth: httpclient.AuthConfig{
			ClientID:     authConfig.ClientID,
			ClientSecret: authConfig.ClientSecret,
		},
		LogLevel: logLevel,
	}

	// Create a new jamfpro client instance
	client, err := jamfpro.BuildClient(config)
	if err != nil {
		log.Fatalf("Failed to create Jamf Pro client: %v", err)
	}

	// Define the advanced computer search details
	newSearch := &jamfpro.ResourceAdvancedComputerSearch{
		Name:   "jamf pro SDK - Advanced Search Name",
		ViewAs: "Standard Web Page",
		Criteria: jamfpro.SharedContainerCriteria{
			Size: 4,
			Criterion: []jamfpro.SharedSubsetCriteria{
				{
					Name:         "Building",
					Priority:     0,
					AndOr:        "and",
					SearchType:   "is",
					Value:        "square",
					OpeningParen: true,
					ClosingParen: false,
				},
				{
					Name:         "Model",
					Priority:     1,
					AndOr:        "and",
					SearchType:   "is",
					Value:        "macbook air",
					OpeningParen: false,
					ClosingParen: true,
				},
				{
					Name:         "Computer Name",
					Priority:     2,
					AndOr:        "or",
					SearchType:   "matches regex",
					Value:        "thing",
					OpeningParen: true,
					ClosingParen: false,
				},
				{
					Name:         "Licensed Software",
					Priority:     3,
					AndOr:        "and",
					SearchType:   "has",
					Value:        "office",
					OpeningParen: false,
					ClosingParen: true,
				},
			},
		},
		DisplayFields: []jamfpro.SharedAdvancedSearchContainerDisplayField{
			{
				DisplayField: []jamfpro.SharedAdvancedSearchSubsetDisplayField{
					{
						Name: "Activation Lock Manageable",
					},
					{
						Name: "Apple Silicon",
					},
					{
						Name: "Architecture Type",
					},
					{
						Name: "Available RAM Slots",
					},
				},
			},
		},
		Site: jamfpro.SharedResourceSite{
			ID:   -1,
			Name: "None",
		},
	}

	// Marshal the newSearch object into XML for logging
	newSearchJSON, err := xml.MarshalIndent(newSearch, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling new search to XML:", err)
		return
	}
	fmt.Printf("New Advanced Computer Search Request:\n%s\n", string(newSearchJSON))

	// Create the advanced computer search
	createdSearch, err := client.CreateAdvancedComputerSearch(newSearch)
	if err != nil {
		fmt.Println("Error creating advanced computer search:", err)
		return
	}

	// Print the created advanced computer search details
	createdSearchXML, err := xml.MarshalIndent(createdSearch, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling created search to XML:", err)
		return
	}
	fmt.Printf("Created Advanced Computer Search:\n%s\n", string(createdSearchXML))
}
