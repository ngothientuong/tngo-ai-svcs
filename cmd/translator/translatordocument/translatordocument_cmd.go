package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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
	action := flag.String("action", "", "Action to perform (translate-document, start-batch-translation, get-translation-status, get-documents-status, get-document-status)")
	documentPath := flag.String("document", "", "Path to the document to translate")
	targetLanguage := flag.String("target-language", "", "Target language for translation")
	sourceLanguage := flag.String("source-language", "", "Source language for translation (optional)")
	glossaryPath := flag.String("glossary", "", "Path to the glossary file (optional)")
	outputFilePath := flag.String("output", "", "Path to the output file")
	allowFallback := flag.Bool("allow-fallback", true, "Allow fallback to general system if custom system doesn't exist")
	category := flag.String("category", "generalnn", "Category for translation (optional)")
	sourceContainer := flag.String("source-container", "", "Source container for batch translation")
	targetContainer := flag.String("target-container", "", "Target container for batch translation")
	prefix := flag.String("prefix", "", "Prefix filter for source documents")
	suffix := flag.String("suffix", "", "Suffix filter for source documents")
	storageType := flag.String("storage-type", "Folder", "Storage type for input documents (Folder or File)")
	glossaries := flag.String("glossaries", "", "Comma-separated list of glossaries in the format glossaryUrl,format,version")
	batchID := flag.String("batch-id", "", "Batch ID for checking translation status")
	documentID := flag.String("document-id", "", "Document ID for checking document status")
	flag.Parse()

	// Parse glossaries
	var glossaryList []map[string]string
	if *glossaries != "" {
		glossaryItems := strings.Split(*glossaries, ",")
		for _, item := range glossaryItems {
			parts := strings.Split(item, ",")
			if len(parts) == 3 {
				glossaryList = append(glossaryList, map[string]string{
					"glossaryUrl":   parts[0],
					"format":        parts[1],
					"version":       parts[2],
					"storageSource": "AzureBlob",
				})
			}
		}
	}

	// Perform the requested action
	switch *action {
	case "translate-document":
		if *documentPath == "" || *targetLanguage == "" || *outputFilePath == "" {
			log.Fatalf("Missing required parameters for translate-document action. Example usage: go run translatordocument_cmd.go -action translate-document -document /path/to/document -target-language fr -output /path/to/output")
		}
		err = client.TranslateDocument(*documentPath, *targetLanguage, *sourceLanguage, *glossaryPath, *outputFilePath, *allowFallback, *category)
	case "start-batch-translation":
		if *sourceContainer == "" || *targetContainer == "" || *targetLanguage == "" {
			log.Println(`Requires: 1. enablement of system-managed identity of the AI Services; 2. Role assignment to Storage Blob Data Contributor to the system-managed identity; 3. Storage account with a container for the source documents having folder name "{folderName}" for --prefix flag and files in the folder with name ending with ".pdf" for --suffix flag and a target container for the translated documents.`)
			log.Fatalf(`Missing required parameters for start-batch-translation action. Example usage: go run translatordocument_cmd.go -action start-batch-translation -source-container https://tngodemo1storageaccount.blob.core.windows.net/tngodemo1translator -target-container https://tngodemo1storageaccount.blob.core.windows.net/tngodemo1translateddocuments  -target-language vi --prefix 'demotranslator/' --suffix ".pdf" --storage-type "Folder"`)
		}
		result, err := client.StartBatchTranslation(*sourceContainer, *targetContainer, *targetLanguage, *prefix, *suffix, *sourceLanguage, *storageType, glossaryList)
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
	case "get-documents-status":
		if *batchID == "" {
			log.Fatalf("Missing required parameters for get-documents-status action. Example usage: go run translatordocument_cmd.go -action get-documents-status -batch-id batch-id")
		}
		result, err := client.GetDocumentsStatus(*batchID)
		if err != nil {
			log.Fatalf("Error performing action: %v", err)
		}
		fmt.Printf("Result: %s\n", result)
	case "get-document-status":
		if *batchID == "" || *documentID == "" {
			log.Fatalf("Missing required parameters for get-document-status action. Example usage: go run translatordocument_cmd.go -action get-document-status -batch-id batch-id -document-id document-id")
		}
		result, err := client.GetDocumentStatus(*batchID, *documentID)
		if err != nil {
			log.Fatalf("Error performing action: %v", err)
		}
		fmt.Printf("Result: %s\n", result)
	default:
		log.Fatalf("Invalid action specified. Example usage: go run translatordocument_cmd.go -action get-translation-status -batch-id batch-id")
	}

	if err != nil {
		log.Fatalf("Error performing action: %v", err)
	}
}
