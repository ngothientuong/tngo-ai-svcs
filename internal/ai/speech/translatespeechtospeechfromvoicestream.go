package speech

import (
	"fmt"
	"sync"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	speechsdk "github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"github.com/gorilla/websocket"
)

var (
	clients               = make(map[*websocket.Conn]bool)
	mu                    sync.Mutex
	translationRecognizer *speechsdk.TranslationRecognizer
)

// TranslateSpeechFromMicrophone starts speech translation
func TranslateSpeechFromMicrophone(speechKey, speechRegion, translationKey, translationEndpoint, targetLanguage, mode string) error {
	// Create speech translation config
	speechConfig, err := speechsdk.NewSpeechTranslationConfigFromSubscription(speechKey, speechRegion)
	if err != nil {
		return fmt.Errorf("failed to create speech translation config: %v", err)
	}
	defer speechConfig.Close()

	err = speechConfig.AddTargetLanguage(targetLanguage)
	if err != nil {
		return fmt.Errorf("failed to add target language: %v", err)
	}

	// Create microphone audio config
	audioConfig, err := audio.NewAudioConfigFromDefaultMicrophoneInput()
	if err != nil {
		return fmt.Errorf("failed to create audio config: %v", err)
	}
	defer audioConfig.Close()

	// Create translation recognizer
	translationRecognizer, err = speechsdk.NewTranslationRecognizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		return fmt.Errorf("failed to create translation recognizer: %v", err)
	}

	// Start continuous recognition
	errChan := translationRecognizer.StartContinuousRecognitionAsync()
	if err := <-errChan; err != nil {
		return fmt.Errorf("failed to start continuous recognition: %v", err)
	}

	return nil
}

// StopSpeechRecognition stops translation
func StopSpeechRecognition() error {
	if translationRecognizer == nil {
		return fmt.Errorf("Recognizer not running")
	}

	errChan := translationRecognizer.StopContinuousRecognitionAsync()
	if err := <-errChan; err != nil {
		return fmt.Errorf("Failed to stop recognition: %v", err)
	}

	translationRecognizer.Close()
	translationRecognizer = nil
	return nil
}
