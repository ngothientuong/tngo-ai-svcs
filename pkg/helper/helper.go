package helper

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ngothientuong/tngo-ai-svcs/internal/ai/language"
)

func PrintResponse(kind string, response *language.AnalyzeTextResponse) {
	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal response for %s: %v", kind, err)
	}
	fmt.Printf("AnalyzeText response for %s:\n%s\n", kind, string(responseJSON))
}
