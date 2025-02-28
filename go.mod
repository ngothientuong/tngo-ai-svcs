module github.com/ngothientuong/tngo-ai-svcs

go 1.23.1

require (
	github.com/Microsoft/cognitive-services-speech-sdk-go v1.33.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/joho/godotenv v1.5.1
)

replace github.com/Microsoft/cognitive-services-speech-sdk-go => github.com/Microsoft/cognitive-services-speech-sdk-go v0.0.0-20250225193958-1f1d5d41a9bb // Support Go translation speech-to-speech
