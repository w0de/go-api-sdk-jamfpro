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

	// New directory binding data
	newBinding := &jamfpro.ResponseDirectoryBinding{
		Name:       "New Binding",
		Priority:   1,
		Domain:     "example.com",
		Username:   "user@example.com",
		Password:   "password",
		ComputerOU: "CN=Computers,DC=example,DC=com",
		Type:       "Active Directory",
	}

	// Create new directory binding
	createdBinding, err := client.CreateDirectoryBinding(newBinding)
	if err != nil {
		fmt.Println("Error creating directory binding:", err)
		return
	}

	// Pretty print the created directory binding in xml
	createdBindingXML, err := xml.MarshalIndent(createdBinding, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling created binding data: %v", err)
	}
	fmt.Printf("Created Directory Binding:\n%s\n", string(createdBindingXML))
}
