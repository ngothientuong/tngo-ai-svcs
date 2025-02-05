package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
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

	imageURL := "https://images-prod.dazeddigital.com/900/azure/dazed-prod/1390/6/1396041.jpg"
	content, err := imageURLToBase64(imageURL)
	if err != nil {
		log.Fatalf("failed to convert image URL to base64: %v", err)
	}

	mediaType := contentsafety.MediaTypeImage
	blocklists := []string{}

	detectionResult, err := client.Detect(mediaType, content, blocklists)
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

func imageURLToBase64(imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %v", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)
	return base64Image, nil
}
