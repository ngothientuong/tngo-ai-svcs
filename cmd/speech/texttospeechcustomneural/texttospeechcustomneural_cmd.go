package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ngothientuong/tngo-ai-svcs/internal/ai/speech"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("/home/tngo/ngo/projects/tngo-ai-svcs/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the endpoint, key, region, and API version from environment variables
	endpoint := os.Getenv("SPEECH_ENDPOINT")
	key := os.Getenv("SPEECH_KEY")
	region := os.Getenv("SPEECH_REGION")
	apiVersion := os.Getenv("SPEECH_TO_TEXT_API_VERSION")

	// Create a new TextToSpeechCustomNeuralClient
	client := speech.NewTextToSpeechCustomNeuralClient(endpoint, key, region, apiVersion)

	// Define the project details
	projectID := "tuong-custom-voice-project"
	kind := "ProfessionalVoice"
	description := "This is a custom voice project for creating a professional voice."

	// Check if the project exists
	exists, err := client.CheckProjectExists(projectID)
	if err != nil {
		log.Fatalf("Error checking project existence: %v", err)
	}

	if !exists {
		// Create the project if it doesn't exist
		project, err := client.CreateProject(projectID, kind, description)
		if err != nil {
			log.Fatalf("Error creating project: %v", err)
		}
		fmt.Printf("Project created: %+v\n", project)
	} else {
		// List projects if it exists
		projects, err := client.ListProjects()
		if err != nil {
			log.Fatalf("Error listing projects: %v", err)
		}
		fmt.Println("Existing projects:")
		for _, project := range projects {
			fmt.Printf("ID: %s, Kind: %s, Description: %s\n", project.ID, project.Kind, project.Description)
		}

		// Get the project ID if it exists
		projectID, err := client.GetProjectID(kind)
		if err != nil {
			log.Fatalf("Error getting project ID: %v", err)
		}
		fmt.Printf("Project already exists with ID: %s\n", projectID)
	}
}
