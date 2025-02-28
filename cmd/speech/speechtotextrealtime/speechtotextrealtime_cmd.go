package main

import (
	"flag"
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

	// Get API credentials
	subscriptionKey := os.Getenv("SPEECH_KEY")
	region := os.Getenv("SPEECH_REGION")

	// Parse command-line arguments
	streamURL := flag.String("url", "", "Live stream video URL (YouTube, RTSP, HLS)")
	format := flag.String("format", "mp4", "Video format (mp4, mkv, youtube, rtsp, hls)")
	flag.Parse()

	// Validate input
	if subscriptionKey == "" || region == "" || *streamURL == "" {
		log.Fatalf("Missing required credentials or stream URL")
	}

	// Start transcription from live video stream
	err = speech.TranscribeFromLiveStream(subscriptionKey, region, *streamURL, *format)
	if err != nil {
		log.Fatalf("Speech recognition failed: %v", err)
	}
}
