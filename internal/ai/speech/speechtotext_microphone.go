package speech

import (
	"fmt"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	speechsdk "github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// TranscribeFromMicrophone transcribes speech from a microphone input
func TranscribeFromMicrophone(subscriptionKey, region string) error {
	// Create speech config
	speechConfig, err := speechsdk.NewSpeechConfigFromSubscription(subscriptionKey, region)
	if err != nil {
		return fmt.Errorf("failed to create speech config: %v", err)
	}
	defer speechConfig.Close()

	// Create audio config from microphone input
	audioConfig, err := audio.NewAudioConfigFromDefaultMicrophoneInput()
	if err != nil {
		return fmt.Errorf("failed to create audio config: %v", err)
	}
	defer audioConfig.Close()

	// Create speech recognizer
	speechRecognizer, err := speechsdk.NewSpeechRecognizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		return fmt.Errorf("failed to create speech recognizer: %v", err)
	}
	defer speechRecognizer.Close()

	// Continuous transcription
	fmt.Println("🎤 Transcribing from microphone...")

	done := make(chan struct{})

	// Handle intermediate results
	speechRecognizer.Recognizing(func(event speechsdk.SpeechRecognitionEventArgs) {
		fmt.Printf("🔵 Partial: %s\n", event.Result.Text)
	})

	// Handle final results
	speechRecognizer.Recognized(func(event speechsdk.SpeechRecognitionEventArgs) {
		fmt.Printf("✅ Final: %s\n", event.Result.Text)
	})

	// Handle errors
	speechRecognizer.Canceled(func(event speechsdk.SpeechRecognitionCanceledEventArgs) {
		fmt.Printf("❌ Canceled: %v\n", event.ErrorDetails)
		close(done)
	})

	// Handle session stop
	speechRecognizer.SessionStopped(func(event speechsdk.SessionEventArgs) {
		fmt.Println("🔴 Session stopped")
		close(done)
	})

	// Start continuous recognition
	errChan := speechRecognizer.StartContinuousRecognitionAsync()
	if err := <-errChan; err != nil {
		return fmt.Errorf("failed to start continuous recognition: %v", err)
	}

	<-done
	return nil
}
