package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ngothientuong/tngo-ai-svcs/internal/ai/contentsafety"
	"github.com/ngothientuong/tngo-ai-svcs/internal/config"
)

func main() {
	config.LoadEnv()

	endpoint := os.Getenv("CONTENT_SAFETY_ENDPOINT")
	subscriptionKey := os.Getenv("CONTENT_SAFETY_KEY")
	apiVersion := os.Getenv("CONTENT_SAFETY_API_VERSION")

	if endpoint == "" || subscriptionKey == "" || apiVersion == "" {
		log.Println("One or more environment variables are not set.")
		return
	}

	client := contentsafety.NewContentSafetyClient(endpoint, subscriptionKey, apiVersion)

	// Example for text content
	textContent := "Russian kills every single Ukrainian in sight with a gun, a knife, or just his bare hands. He is a monster."
	mediaType := contentsafety.MediaTypeText
	blocklists := []string{}

	detectionResult, err := client.Detect(mediaType, textContent, blocklists)
	if err != nil {
		log.Fatalf("Error detecting content safety: %v", err)
	}

	rejectThresholds := map[contentsafety.Category]int{
		contentsafety.CategoryHate:     4,
		contentsafety.CategorySelfHarm: 4,
		contentsafety.CategorySexual:   4,
		contentsafety.CategoryViolence: 4,
	}

	decisionResult, err := client.MakeDecision(detectionResult, rejectThresholds)
	if err != nil {
		log.Fatalf("Error making decision: %v", err)
	}

	fmt.Printf("Decision: %+v\n", decisionResult)
}
