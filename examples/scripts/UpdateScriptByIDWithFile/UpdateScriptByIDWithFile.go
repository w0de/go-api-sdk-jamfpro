package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

const (
	maxConcurrentRequestsAllowed = 5 // Maximum allowed concurrent requests.
	defaultTokenLifespan         = 30 * time.Minute
	defaultBufferPeriod          = 5 * time.Minute
	scriptFilePath               = "/Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-jamfpro/examples/support_files/scriptfile.sh" // Path to the script file
)

func main() {
	// Define the path to the JSON configuration file inside the main function
	configFilePath := "/Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-jamfpro/clientauth.json"

	// Load the client OAuth credentials from the configuration file
	authConfig, err := http_client.LoadClientAuthConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load client OAuth configuration: %v", err)
	}

	// Configuration for the jamfpro
	config := jamfpro.Config{
		InstanceName:          authConfig.InstanceName,
		DebugMode:             true,
		Logger:                jamfpro.NewDefaultLogger(),
		MaxConcurrentRequests: maxConcurrentRequestsAllowed,
		TokenLifespan:         defaultTokenLifespan,
		BufferPeriod:          defaultBufferPeriod,
		ClientID:              authConfig.ClientID,
		ClientSecret:          authConfig.ClientSecret,
	}

	// Create a new jamfpro client instanceclient,
	client, err := jamfpro.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Jamf Pro client: %v", err)
	}

	// Load script contents from a file
	file, err := os.Open(scriptFilePath)
	if err != nil {
		log.Fatalf("Failed to open script file: %v", err)
	}
	defer file.Close()

	scriptContents, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read script file: %v", err)
	}

	// Define a sample script for testing
	updatedScript := &jamfpro.ResponseScript{
		ID:             194, // Assuming ID 194 for this example
		Name:           "Updated Sample Script",
		Category:       "None",
		Filename:       "string",
		Info:           "Updated Script information",
		Notes:          "Updated Sample Script",
		Priority:       "After",
		OSRequirements: "string",
		ScriptContents: string(scriptContents),
	}

	// Call UpdateScriptByID function
	resultScript, err := client.UpdateScriptByID(updatedScript)
	if err != nil {
		log.Fatalf("Error updating script: %v", err)
	}

	// Pretty print the updated script details in XML
	resultScriptXML, err := xml.MarshalIndent(resultScript, "", "    ") // Indent with 4 spaces
	if err != nil {
		log.Fatalf("Error marshaling updated script data: %v", err)
	}
	fmt.Println("Updated Script Details:\n", string(resultScriptXML))
}
