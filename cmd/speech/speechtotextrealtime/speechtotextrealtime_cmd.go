package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/ngothientuong/tngo-ai-svcs/internal/ai/speech"
)

var (
	mu       sync.Mutex
	clients  = make(map[*websocket.Conn]bool)
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	mode     = "both" // Default mode: "both", can be "text-only" or "speech-only"
)

// ModeChangeRequest handles mode switching via WebSocket
type ModeChangeRequest struct {
	Mode string `json:"mode"`
}

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Load API keys separately
	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")
	translationKey := os.Getenv("TRANSLATION_KEY")
	translationEndpoint := os.Getenv("TRANSLATION_ENDPOINT")

	// Serve Web Interface
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// WebSocket handler for real-time updates
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return
		}
		defer conn.Close()

		mu.Lock()
		clients[conn] = true
		mu.Unlock()

		// Keep connection alive
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				mu.Lock()
				delete(clients, conn)
				mu.Unlock()
				break
			}

			// Handle mode change requests from frontend
			var request ModeChangeRequest
			if err := json.Unmarshal(message, &request); err == nil && request.Mode != "" {
				mu.Lock()
				mode = request.Mode
				mu.Unlock()
				log.Println("Mode changed to:", mode)
			}
		}
	})

	// Start speech-to-speech translation
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		go func() {
			err := speech.TranslateSpeechFromMicrophone(speechKey, speechRegion, translationKey, translationEndpoint, "vi", mode)
			if err != nil {
				log.Printf("Translation failed: %v", err)
			}
		}()

		fmt.Fprintln(w, "Translation session started")
	})

	// Stop speech translation
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		err := speech.StopSpeechRecognition()
		if err != nil {
			http.Error(w, "Failed to stop recognition", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Translation session stopped")
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
