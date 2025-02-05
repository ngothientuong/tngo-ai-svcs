package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ngothientuong/tngo-ai-svcs/internal/ai/computervision"
	"github.com/ngothientuong/tngo-ai-svcs/internal/config"
)

func main() {
	config.LoadEnv()

	subscriptionKey := os.Getenv("COMPUTER_VISION_KEY")
	endpoint := os.Getenv("COMPUTER_VISION_ENDPOINT")

	if subscriptionKey == "" || endpoint == "" {
		log.Println("One or more environment variables are not set.")
		return
	}

	log.Printf("Using endpoint: %s", endpoint)
	log.Printf("Using key: %s", subscriptionKey)

	client := computervision.NewImageAnalysisClient(endpoint, subscriptionKey)

	imageURL1 := "https://aka.ms/azsdk/image-analysis/sample.jpg" // Replace with your image URL
	imageURL2 := "https://encrypted-tbn2.gstatic.com/shopping?q=tbn:ANd9GcR1kz5pVfNQf8xvsfVsO50VVwK7sgonr1IyfUs5p-5p5wHj4FNfGVaBSlnS-yHi6Ab1y2WBHJIMlEDHSvycMh1vv5GEmMw67HQVJCK7Ao2-ZZ1CkUziJJuH388"

	// Analyze the first image for caption, tags, and objects
	visualFeatures1 := []computervision.VisualFeature{
		computervision.VisualFeatureCaption,
		computervision.VisualFeatureTags,
		computervision.VisualFeatureObjects,
	}
	log.Printf("Analyzing image: %s", imageURL1)
	result1, err := client.AnalyzeImage(imageURL1, visualFeatures1)
	if err != nil {
		log.Fatalf("Error analyzing image: %v", err)
	}

	if result1.Caption.Text != "" {
		fmt.Printf("Caption: %s\n", result1.Caption.Text)
	}
	if len(result1.Tags.Values) > 0 {
		fmt.Printf("Tags: %v\n", result1.Tags.Values)
	}
	if len(result1.Objects.Values) > 0 {
		fmt.Printf("Objects: %v\n", result1.Objects.Values)
	}

	// Analyze the second image for text (OCR)
	visualFeatures2 := []computervision.VisualFeature{
		computervision.VisualFeatureRead,
	}
	log.Printf("Analyzing image: %s", imageURL2)
	result2, err := client.AnalyzeImage(imageURL2, visualFeatures2)
	if err != nil {
		log.Fatalf("Error analyzing image: %v", err)
	}

	if len(result2.ReadResult.Blocks) > 0 {
		fmt.Printf("Read - OCR: %v\n", result2.ReadResult.Blocks)
	}
}
