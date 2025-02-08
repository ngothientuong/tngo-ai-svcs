package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ngothientuong/tngo-ai-svcs/internal/ai/language"
	"github.com/ngothientuong/tngo-ai-svcs/pkg/helper"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("/home/tngo/ngo/projects/tngo-ai-svcs/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the endpoint, key, and API version from environment variables
	endpoint := os.Getenv("AZURE_AI_MULTISERVICE_ENDPOINT")
	key := os.Getenv("AZURE_AI_MULTISERVICE_KEY")
	apiVersion := os.Getenv("LANGUAGE_API_VERSION")

	// Create a new AnalyzeTextClient
	client := language.NewAnalyzeTextClient(endpoint, key, apiVersion)

	// Define the documents to be analyzed with longer texts
	documents := []language.Document{
		{ID: "1", Language: "en", Text: "Microsoft was founded by Bill Gates and Paul Allen. Microsoft is a technology company that develops, manufactures, licenses, supports, and sells a range of software products and services. Its best-known software products are the Microsoft Windows line of operating systems, the Microsoft Office suite, and the Internet Explorer and Edge web browsers. Its flagship hardware products are the Xbox video game consoles and the Microsoft Surface lineup of touchscreen personal computers. Microsoft ranked No. 21 in the 2020 Fortune 500 rankings of the largest United States corporations by total revenue."},
		{ID: "2", Language: "en", Text: "Pike Place Market is a public market overlooking the Elliott Bay waterfront in Seattle, Washington, United States. The Market opened on August 17, 1907, and is one of the oldest continuously operated public farmers' markets in the United States. It is a place of business for many small farmers, craftspeople, and merchants. Named after the central street, Pike Place runs northwest from Pike Street to Virginia Street. The Market is one of Seattle's most popular tourist destinations and is home to the original Starbucks coffee shop."},
	}

	// Define the parameters for the analysis
	parameters := map[string]interface{}{
		"modelVersion": "latest",
	}

	// Call the AnalyzeText method for different kinds
	kinds := []string{
		"EntityLinking",
		"EntityRecognition",
		"KeyPhraseExtraction",
		"PiiEntityRecognition",
		"SentimentAnalysis",
	}

	for _, kind := range kinds {
		response, err := client.AnalyzeText(kind, documents, parameters)
		if err != nil {
			log.Printf("Error analyzing text for %s: %v", kind, err)
			fmt.Println()
			fmt.Println()
			continue
		}

		// Print the response
		helper.PrintResponse(kind, response)
	}

	// Additional sample for LanguageDetection
	languageDetectionDocuments := []language.Document{
		{ID: "1", Text: "Hello world"},
		{ID: "2", Text: "Bonjour tout le monde"},
		{ID: "3", Text: "Hola mundo"},
		{ID: "4", Text: "Tumhara naam kya hai?"},
	}

	languageDetectionParameters := map[string]interface{}{
		"modelVersion": "latest",
	}

	languageDetectionResponse, err := client.AnalyzeText("LanguageDetection", languageDetectionDocuments, languageDetectionParameters)
	if err != nil {
		log.Fatalf("Error analyzing text for LanguageDetection: %v", err)
	}

	// Print the response for LanguageDetection
	helper.PrintResponse("LanguageDetection", languageDetectionResponse)

	// Additional samples for EntityRecognition with different parameters

	// Exclusion List
	entityRecognitionExclusionParameters := map[string]interface{}{
		"modelVersion":  "latest",
		"exclusionList": []string{"Numeric"},
		"overlapPolicy": map[string]string{"policyKind": "allowOverlap"},
	}

	entityRecognitionExclusionDocuments := []language.Document{
		{ID: "2", Language: "en", Text: "When I was 5 years old I had $90.00 dollars to my name."},
		{ID: "3", Language: "en", Text: "When we flew from LAX it seemed like we were moving at 10 meters per second. I was lucky to see Amsterdam, Eiffel Tower, and the Nile."},
	}

	entityRecognitionExclusionResponse, err := client.AnalyzeText("EntityRecognition", entityRecognitionExclusionDocuments, entityRecognitionExclusionParameters)
	if err != nil {
		log.Fatalf("Error analyzing text for EntityRecognition with exclusion list: %v", err)
	}

	// Print the response for EntityRecognition with exclusion list
	helper.PrintResponse("EntityRecognition with exclusion list", entityRecognitionExclusionResponse)

	// Inclusion List
	entityRecognitionInclusionParameters := map[string]interface{}{
		"modelVersion":  "latest",
		"inclusionList": []string{"Location"},
	}

	entityRecognitionInclusionDocuments := []language.Document{
		{ID: "2", Language: "en", Text: "When I was 5 years old I had $90.00 dollars to my name."},
		{ID: "3", Language: "en", Text: "When we flew from LAX it seemed like we were moving at 10 meters per second. I was lucky to see Amsterdam, Eiffel Tower, and the Nile."},
	}

	entityRecognitionInclusionResponse, err := client.AnalyzeText("EntityRecognition", entityRecognitionInclusionDocuments, entityRecognitionInclusionParameters)
	if err != nil {
		log.Fatalf("Error analyzing text for EntityRecognition with inclusion list: %v", err)
	}

	// Print the response for EntityRecognition with inclusion list
	helper.PrintResponse("EntityRecognition with inclusion list", entityRecognitionInclusionResponse)

	// Inference Options
	entityRecognitionInferenceParameters := map[string]interface{}{
		"modelVersion":     "latest",
		"inferenceOptions": map[string]bool{"excludeNormalizedValues": true},
	}

	entityRecognitionInferenceDocuments := []language.Document{
		{ID: "1", Language: "en", Text: "When I was 5 years old I had $90.00 dollars to my name."},
	}

	entityRecognitionInferenceResponse, err := client.AnalyzeText("EntityRecognition", entityRecognitionInferenceDocuments, entityRecognitionInferenceParameters)
	if err != nil {
		log.Fatalf("Error analyzing text for EntityRecognition with inference options: %v", err)
	}

	// Print the response for EntityRecognition with inference options
	helper.PrintResponse("EntityRecognition with inference options", entityRecognitionInferenceResponse)

	// Overlap Policy
	entityRecognitionOverlapParameters := map[string]interface{}{
		"modelVersion":  "latest",
		"overlapPolicy": map[string]string{"policyKind": "matchLongest"},
	}

	entityRecognitionOverlapDocuments := []language.Document{
		{ID: "4", Language: "en", Text: "25th April Meeting was an interesting one. At least we got to experience the WorldCup"},
	}

	entityRecognitionOverlapResponse, err := client.AnalyzeText("EntityRecognition", entityRecognitionOverlapDocuments, entityRecognitionOverlapParameters)
	if err != nil {
		log.Fatalf("Error analyzing text for EntityRecognition with overlap policy: %v", err)
	}

	// Print the response for EntityRecognition with overlap policy
	helper.PrintResponse("EntityRecognition with overlap policy", entityRecognitionOverlapResponse)
}
