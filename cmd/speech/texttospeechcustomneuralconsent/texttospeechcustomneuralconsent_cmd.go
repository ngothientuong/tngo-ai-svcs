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

    // Create a new TextToSpeechCustomConsentClient
    client := speech.NewTextToSpeechCustomConsentClient(endpoint, key, region, apiVersion)

    // Define the consent details
    consentID := "jessica-consent"
    description := "Consent for Jessica voice"
    projectID := "jessica-project"
    voiceTalentName := "Jessica Smith"
    companyName := "Contoso"
    audioURL := "https://contoso.blob.core.windows.net/public/jessica-consent.wav?mySasToken"
    locale := "en-US"

    // Check if the consent exists
    exists, err := client.CheckConsentExists(consentID)
    if err != nil {
        log.Fatalf("Error checking consent existence: %v", err)
    }

    if !exists {
        // Create the consent if it doesn't exist
        consent, err := client.CreateConsent(consentID, description, projectID, voiceTalentName, companyName, audioURL, locale)
        if err != nil {
            log.Fatalf("Error creating consent: %v", err)
        }
        fmt.Printf("Consent created: %+v\n", consent)
    } else {
        // List consents if it exists
        consents, err := client.ListConsents()
        if err != nil {
            log.Fatalf("Error listing consents: %v", err)
        }
        fmt.Println("Existing consents:")
        for _, consent := range consents {
            fmt.Printf("ID: %s, Description: %s, ProjectID: %s, VoiceTalentName: %s, CompanyName: %s, AudioURL: %s, Locale: %s, Status: %s\n",
                consent.ID, consent.Description, consent.ProjectID, consent.VoiceTalentName, consent.CompanyName, consent.AudioURL, consent.Locale, consent.Status)
        }

        // Get the consent if it exists
        consent, err := client.GetConsent(consentID)
        if err != nil {
            log.Fatalf("Error getting consent: %v", err)
        }
        fmt.Printf("Consent already exists: %+v\n", consent)
    }
}