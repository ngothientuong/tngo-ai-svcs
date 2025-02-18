package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ngothientuong/tngo-ai-svcs/internal/ai/translator"
)

func main() {
	// Load environment variables
	err := godotenv.Load("/home/tngo/ngo/projects/tngo-ai-svcs/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get API credentials
	endpoint := os.Getenv("TRANSLATOR_DOCUMENT_ENDPOINT")
	key := os.Getenv("TRANSLATOR_DOCUMENT_KEY")
	apiVersion := os.Getenv("TRANSLATOR_DOCUMENT_API_VERSION")
	region := os.Getenv("TRANSLATOR_DOCUMENT_REGION")

	// Validate input
	if endpoint == "" || key == "" || apiVersion == "" || region == "" {
		log.Fatalf("Missing required credentials")
	}

	// Create a new DocumentTranslationClient
	client := translator.NewDocumentTranslationClient(endpoint, key, region, apiVersion)

	// Define command-line flags
	action := flag.String("action", "", "Action to perform (translate-document, start-batch-translation, get-translation-status)")
	documentPath := flag.String("document", "", "Path to the document to translate")
	targetLanguage := flag.String("target-language", "", "Target language for translation")
	sourceLanguage := flag.String("source-language", "", "Source language for translation (optional)")
	glossaryPath := flag.String("glossary", "", "Path to the glossary file (optional)")
	outputFilePath := flag.String("output", "", "Path to the output file")
	allowFallback := flag.Bool("allow-fallback", true, "Allow fallback to general system if custom system doesn't exist")
	category := flag.String("category", "generalnn", "Category for translation (optional)")
	sourceContainer := flag.String("source-container", "", "Source container for batch translation")
	targetContainer := flag.String("target-container", "", "Target container for batch translation")
	batchID := flag.String("batch-id", "", "Batch ID for checking translation status")
	flag.Parse()

	// Perform the requested action
	switch *action {
	case "translate-document":
		if *documentPath == "" || *targetLanguage == "" || *outputFilePath == "" {
			log.Fatalf("Missing required parameters for translate-document action. Example usage: go run translatordocument_cmd.go -action translate-document -document /path/to/document -target-language fr -output /path/to/output")
		}
		err = client.TranslateDocument(*documentPath, *targetLanguage, *sourceLanguage, *glossaryPath, *outputFilePath, *allowFallback, *category)
	case "start-batch-translation":
		if *sourceContainer == "" || *targetContainer == "" || *targetLanguage == "" {
			log.Fatalf("Missing required parameters for start-batch-translation action. Example usage: go run translatordocument_cmd.go -action start-batch-translation -source-container source-container-url -target-container target-container-url -target-language fr")
		}
		result, err := client.StartBatchTranslation(*sourceContainer, *targetContainer, *targetLanguage)
		if err != nil {
			log.Fatalf("Error performing action: %v", err)
		}
		fmt.Printf("Result: %s\n", result)
	case "get-translation-status":
		if *batchID == "" {
			log.Fatalf("Missing required parameters for get-translation-status action. Example usage: go run translatordocument_cmd.go -action get-translation-status -batch-id batch-id")
		}
		result, err := client.GetTranslationStatus(*batchID)
		if err != nil {
			log.Fatalf("Error performing action: %v", err)
		}
		fmt.Printf("Result: %s\n", result)
	default:
		log.Fatalf("Invalid action specified. Example usage: go run translatordocument_cmd.go -action translate-document -document /path/to/document -target-language fr -output /path/to/output")
	}

	if err != nil {
		log.Fatalf("Error performing action: %v", err)
	}
}
