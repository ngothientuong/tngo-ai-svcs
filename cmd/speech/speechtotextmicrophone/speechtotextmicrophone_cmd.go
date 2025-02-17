package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ngothientuong/tngo-ai-svcs/internal/ai/speech"
)

func main() {
	// Load environment variables
	err := godotenv.Load("/home/tngo/ngo/projects/tngo-ai-svcs/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get API credentials
	subscriptionKey := os.Getenv("SPEECH_KEY")
	region := os.Getenv("SPEECH_REGION")

	// Validate input
	if subscriptionKey == "" || region == "" {
		log.Fatalf("Missing required credentials")
	}

	// Start transcription from microphone
	err = speech.TranscribeFromMicrophone(subscriptionKey, region)
	if err != nil {
		log.Fatalf("Speech recognition failed: %v", err)
	}
}
