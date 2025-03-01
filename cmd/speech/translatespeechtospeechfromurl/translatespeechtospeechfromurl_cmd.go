package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ngothientuong/tngo-ai-svcs/internal/ai/speech"
)

// printUsage prints CLI usage examples if incorrect arguments are provided
func printUsage() {
	fmt.Println("❌ Incorrect usage! Please provide valid arguments.")
	fmt.Println("✅ Examples:")
	fmt.Println("1️⃣ Translate speech from a YouTube video to Japanese:")
	fmt.Println("   go run cmd/speech_to_speech_cmd.go -url \"https://www.youtube.com/watch?v=EXAMPLE\" -format youtube -from en-US -to ja")
	fmt.Println("2️⃣ Translate from a live RTSP stream to French:")
	fmt.Println("   go run cmd/speech_to_speech_cmd.go -url \"rtsp://your-livestream-url\" -format rtsp -from en-US -to fr")
	fmt.Println("3️⃣ Save translated speech to a file:")
	fmt.Println("   go run cmd/speech_to_speech_cmd.go -url \"https://www.youtube.com/watch?v=EXAMPLE\" -format youtube -from en-US -to es -save")
	fmt.Println("🚀 Ensure you have set your Azure Speech API keys correctly in `.env` before running the command.")
}

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
	fromLang := flag.String("from", "en-US", "Source language (spoken language in video)")
	toLang := flag.String("to", "ja", "Target language for speech synthesis output")
	saveAudio := flag.Bool("save", false, "Save translated speech audio as a file")
	flag.Parse()

	// Validate input and show examples if missing arguments
	if speechKey == "" || speechRegion == "" || *url == "" {
		printUsage()
		log.Fatalf("❌ Missing required credentials or URL. Please check your input.")
	}

	// Start speech-to-speech translation
	fmt.Printf("🎤 Translating speech from %s (%s → %s)...\n", *url, *fromLang, *toLang)
	err = speech.SpeechToSpeechTranslation(speechKey, speechRegion, *url, *format, *fromLang, *toLang, *saveAudio)
	if err != nil {
		log.Fatalf("❌ Speech-to-speech translation failed: %v", err)
	}
}
