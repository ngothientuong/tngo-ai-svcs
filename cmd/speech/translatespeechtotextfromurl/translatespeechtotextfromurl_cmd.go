package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ngothientuong/tngo-ai-svcs/internal/ai/speech"
)

func main() {
	// Load environment variables
	err := godotenv.Load("/home/tngo/ngo/projects/tngo-ai-svcs/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Load API keys
	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")

	// Parse command-line arguments
	url := flag.String("url", "", "Live stream or video file URL")
	format := flag.String("format", "mp4", "Format (mp4, mkv, youtube, rtsp, hls)")
	fromLang := flag.String("from", "en-US", "Source language")
	toLang := flag.String("to", "es", "Target language for translation")
	flag.Parse()

	// Validate input
	if speechKey == "" || speechRegion == "" || *url == "" {
		log.Fatalf("Missing required credentials or URL. Example run: go run translatespeechtotextfromurl_cmd.go -url \"https://www.youtube.com/watch?v=Db7sLhE3MnI\" -format youtube -from en-US -to ja")
	}

	// Start speech translation
	fmt.Printf("🎤 Translating speech from %s (%s → %s)...\n", *url, *fromLang, *toLang)
	err = speech.TranslateSpeechFromURL(speechKey, speechRegion, *url, *format, *fromLang, *toLang)
	if err != nil {
		log.Fatalf("Translation failed: %v", err)
	}
}
