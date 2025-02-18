package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ngothientuong/tngo-ai-svcs/internal/ai/translator"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	// Load environment variables
	err := godotenv.Load("/home/tngo/ngo/projects/tngo-ai-svcs/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get API credentials
	endpoint := os.Getenv("TRANSLATOR_ENDPOINT")
	key := os.Getenv("TRANSLATOR_KEY")
	apiVersion := os.Getenv("TRANSLATOR_API_VERSION")
	region := os.Getenv("TRANSLATOR_REGION")

	// Validate input
	if endpoint == "" || key == "" || apiVersion == "" || region == "" {
		log.Fatalf("Missing required credentials")
	}

	// Create a new TranslatorClient
	client := translator.NewTranslatorClient(endpoint, key, region, apiVersion)

	// Define command-line flags
	var texts arrayFlags
	var translations arrayFlags
	var from arrayFlags
	var to arrayFlags
	action := flag.String("action", "", "Action to perform (languages, translate, transliterate, detect, dictionary-lookup, dictionary-examples)")
	flag.Var(&texts, "text", "Text to process (can be specified multiple times)")
	flag.Var(&translations, "translation", "Translation to process (can be specified multiple times)")
	flag.Var(&from, "from", "Source language or script (can be specified multiple times)")
	flag.Var(&to, "to", "Target language or script (can be specified multiple times)")
	language := flag.String("language", "", "Language for transliterate action")
	scope := flag.String("scope", "translation,transliteration,dictionary", "Scope for GetLanguages action")
	flag.Parse()

	// Perform the requested action
	var result []byte
	switch *action {
	case "languages":
		result, err = client.GetLanguages(*scope)
	case "translate":
		if len(texts) == 0 || len(to) == 0 {
			log.Fatalf("Missing required parameters for translate action. Example usage: go run translatortext_cmd.go -action translate -text \"Hello, world!\" -to fr")
		}
		result, err = client.Translate(texts, to)
	case "transliterate":
		if len(texts) == 0 || len(from) == 0 || len(to) == 0 || *language == "" {
			log.Fatalf("Missing required parameters for transliterate action. Example usage: go run translatortext_cmd.go -action transliterate -text \"こんにちは\" -from jpan -to latn -language ja")
		}
		result, err = client.Transliterate(texts, *language, from[0], to[0])
	case "detect":
		if len(texts) == 0 {
			log.Fatalf("Missing required parameters for detect action. Example usage: go run translatortext_cmd.go -action detect -text \"Bonjour tout le monde\"")
		}
		result, err = client.Detect(texts)
	case "dictionary-lookup":
		if len(texts) == 0 || len(from) == 0 || len(to) == 0 {
			log.Fatalf("Missing required parameters for dictionary-lookup action. Example usage: go run translatortext_cmd.go -action dictionary-lookup -text \"example\" -from en -to fr")
		}
		result, err = client.DictionaryLookup(texts, from[0], to[0])
	case "dictionary-examples":
		if len(texts) == 0 || len(from) == 0 || len(to) == 0 {
			log.Fatalf("Missing required parameters for dictionary-examples action. Example usage: go run translatortext_cmd.go -action dictionary-examples -text \"example\" -from en -to fr")
		}

		// Step 1: Call Dictionary Lookup to Get Translations
		fmt.Println("Fetching dictionary lookup translations...")
		lookupResult, err := client.DictionaryLookup(texts, from[0], to[0])
		if err != nil {
			log.Fatalf("Error fetching dictionary lookup: %v", err)
		}

		// Step 2: Extract translations from response
		var lookupResponse []map[string]interface{}
		err = json.Unmarshal(lookupResult, &lookupResponse)
		if err != nil {
			log.Fatalf("Failed to parse lookup response: %v", err)
		}

		// Step 3: Prepare Dictionary Examples Requests
		var exampleRequests []map[string]string
		for _, entry := range lookupResponse {
			if translations, exists := entry["translations"].([]interface{}); exists {
				for _, t := range translations {
					if translation, ok := t.(map[string]interface{}); ok {
						text := entry["normalizedSource"].(string)
						translatedText := translation["normalizedTarget"].(string)
						exampleRequests = append(exampleRequests, map[string]string{
							"Text":        text,
							"Translation": translatedText,
						})
					}
				}
			}
		}

		// Step 4: Call Dictionary Examples for Each Translation
		fmt.Println("Fetching dictionary examples for translations...")
		result, err = client.DictionaryExamples(exampleRequests, from[0], to[0])
		if err != nil {
			log.Fatalf("Error fetching dictionary examples: %v", err)
		}
	default:
		log.Fatalf("Invalid action specified. Example usage: go run translatortext_cmd.go -action translate -text \"Hello, world!\" -to fr")
	}

	if err != nil {
		log.Fatalf("Error performing action: %v", err)
	}

	fmt.Printf("Result: %s\n", result)
}
