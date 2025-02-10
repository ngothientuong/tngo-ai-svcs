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

	// Get the endpoint, key, and region from environment variables
	endpoint := os.Getenv("SPEECH_ENDPOINT")
	key := os.Getenv("SPEECH_KEY")
	region := os.Getenv("SPEECH_REGION")

	// Create a new TextToSpeechClient
	client := speech.NewTextToSpeechClient(endpoint, key, region)

	// Define the text, voice, and format for the synthesis
	text := "Hello, world!"
	voice := "en-US-AvaMultilingualNeural"
	format := "audio-16khz-128kbitrate-mono-mp3"

	// Call the SynthesizeSpeech method
	audioData, err := client.SynthesizeSpeech(text, voice, format)
	if err != nil {
		log.Fatalf("Error synthesizing speech: %v", err)
	}

	// Save the audio data to a file
	outputFile := "output.mp3"
	err = os.WriteFile(outputFile, audioData, 0644)
	if err != nil {
		log.Fatalf("Error saving audio file: %v", err)
	}

	fmt.Printf("Audio saved to %s\n", outputFile)
}
