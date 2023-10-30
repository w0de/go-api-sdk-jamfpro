package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

func loadScriptFromFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func main() {
	// Define the path to the JSON configuration file inside the main function
	configFilePath := "/Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-jamfpro/clientauth.json"

	// Load the client OAuth credentials from the configuration file
	authConfig, err := jamfpro.LoadClientAuthConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load client OAuth configuration: %v", err)
	}

	// Configuration for the jamfpro
	config := jamfpro.Config{
		InstanceName: authConfig.InstanceName,
		DebugMode:    true,
		Logger:       jamfpro.NewDefaultLogger(),
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
	}

	// Create a new jamfpro client instance
	client, err := jamfpro.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Jamf Pro client: %v", err)
	}

	// Load the script from a file
	scriptPath := "/Users/dafyddwatkins/GitHub/deploymenttheory/go-api-sdk-jamfpro/examples/support_files/computer_extensioin_attribute.sh"
	scriptContent, err := loadScriptFromFile(scriptPath)
	if err != nil {
		log.Fatalf("Failed to load script from file: %v", err)
	}

	// Define the new computer extension attribute
	attribute := &jamfpro.ResponseComputerExtensionAttribute{
		Name:        "Computer Extension Attribute Script Test",
		Description: "Computer Extension Attribute SCript Test",
		DataType:    "String", // String / Integer / Date (YYYY-MM-DD hh:mm:ss)
		InputType: jamfpro.ComputerExtensionAttributeInputType{
			Type:     "Script",
			Script:   scriptContent,
			Platform: "Mac", // Set this to the desired platform: "Mac" or "Windows".
		},
		InventoryDisplay: "General", // General / Hardware / Operating System / User and Location / Purchasing / Extension Attribute
		ReconDisplay:     "Extension Attributes",
	}

	// Call CreateComputerExtensionAttribute function
	createdAttribute, err := client.CreateComputerExtensionAttribute(attribute)
	if err != nil {
		log.Fatalf("Error creating Computer Extension Attribute: %v", err)
	}

	// Pretty print the created attribute in XML
	createdAttributeXML, err := xml.MarshalIndent(createdAttribute, "", "    ") // Indent with 4 spaces
	if err != nil {
		log.Fatalf("Error marshaling created Computer Extension Attribute data: %v", err)
	}
	fmt.Println("Created Computer Extension Attribute:\n", string(createdAttributeXML))
}
