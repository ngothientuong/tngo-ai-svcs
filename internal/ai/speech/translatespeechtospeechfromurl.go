package speech

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	speechsdk "github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

var (
	ttsMutex sync.Mutex
)

// SpeechToSpeechTranslation translates speech from a video/livestream source and outputs spoken translation
func SpeechToSpeechTranslation(subscriptionKey, region, url, format, fromLang, toLang string, saveAudio bool) error {
	// Validate input languages
	if !isValidSourceLanguage(fromLang) || !isValidTargetLanguage(toLang) {
		return fmt.Errorf("invalid language: %s → %s. Check Azure supported languages", fromLang, toLang)
	}

	// Create speech translation config
	config, err := speechsdk.NewSpeechTranslationConfigFromSubscription(subscriptionKey, region)
	if err != nil {
		return fmt.Errorf("failed to create speech translation config: %v", err)
	}
	defer config.Close()

	// Set source and target languages
	err = config.SetSpeechRecognitionLanguage(fromLang)
	if err != nil {
		return fmt.Errorf("failed to set speech recognition language: %v", err)
	}
	err = config.AddTargetLanguage(toLang)
	if err != nil {
		return fmt.Errorf("failed to add target language: %v", err)
	}

	// Create audio stream from URL
	audioStream, err := NewPushAudioInputStreamFromURL(url, format)
	if err != nil {
		return fmt.Errorf("failed to create audio input stream: %v", err)
	}
	defer audioStream.Close()

	// Create audio config
	audioConfig, err := audio.NewAudioConfigFromStreamInput(audioStream)
	if err != nil {
		return fmt.Errorf("failed to create audio config: %v", err)
	}
	defer audioConfig.Close()

	// Create translation recognizer
	recognizer, err := speechsdk.NewTranslationRecognizerFromConfig(config, audioConfig)
	if err != nil {
		return fmt.Errorf("failed to create translation recognizer: %v", err)
	}
	defer recognizer.Close()

	// Open an audio file for saving (if enabled)
	var audioFile *os.File
	if saveAudio {
		audioFile, err = os.Create("translated_audio.wav")
		if err != nil {
			return fmt.Errorf("failed to create audio file: %v", err)
		}
		defer audioFile.Close()
	}

	// Handle translation results
	recognizer.Recognized(func(event speechsdk.TranslationRecognitionEventArgs) {
		translations := event.Result.GetTranslations()
		if translatedText, ok := translations[toLang]; ok {
			fmt.Printf("✅ Final Translation: %s\n", translatedText)
			go textToSpeech(subscriptionKey, region, translatedText, toLang, saveAudio, audioFile)
		}
	})

	recognizer.Canceled(func(event speechsdk.TranslationRecognitionCanceledEventArgs) {
		fmt.Printf("❌ Canceled: %v\n", event.ErrorDetails)
		closeDoneChan() // Use existing helper
	})

	recognizer.SessionStopped(func(event speechsdk.SessionEventArgs) {
		fmt.Println("🔴 Translation session stopped")
		closeDoneChan() // Use existing helper
	})

	// Start continuous recognition
	errChan := recognizer.StartContinuousRecognitionAsync()
	if err := <-errChan; err != nil {
		return fmt.Errorf("failed to start continuous recognition: %v", err)
	}

	// Wait for recognition to complete (reuse existing helper)
	waitForDoneChan()
	return nil
}

// textToSpeech converts translated text to spoken audio and plays it
func textToSpeech(subscriptionKey, region, text, lang string, saveAudio bool, audioFile *os.File) {
	ttsMutex.Lock()
	defer ttsMutex.Unlock()

	// Create speech config
	config, err := speechsdk.NewSpeechConfigFromSubscription(subscriptionKey, region)
	if err != nil {
		fmt.Printf("❌ TTS config error: %v\n", err)
		return
	}
	defer config.Close()

	// Set language
	err = config.SetSpeechSynthesisLanguage(lang)
	if err != nil {
		fmt.Printf("❌ Failed to set TTS language: %v\n", err)
		return
	}

	// Create an audio output stream
	audioStream, err := audio.CreatePullAudioOutputStream()
	if err != nil {
		fmt.Printf("❌ Failed to create audio output stream: %v\n", err)
		return
	}
	defer audioStream.Close()

	// Create synthesizer
	audioConfig, err := audio.NewAudioConfigFromStreamOutput(audioStream)
	if err != nil {
		fmt.Printf("❌ Failed to create TTS audio config: %v\n", err)
		return
	}
	defer audioConfig.Close()

	synthesizer, err := speechsdk.NewSpeechSynthesizerFromConfig(config, audioConfig)
	if err != nil {
		fmt.Printf("❌ Failed to create TTS synthesizer: %v\n", err)
		return
	}
	defer synthesizer.Close()

	// Perform speech synthesis
	resultChan := synthesizer.SpeakTextAsync(text)
	result := <-resultChan

	// Handle errors
	if result.Error != nil {
		fmt.Printf("❌ TTS error: %v\n", result.Error)
		return
	}

	// Read synthesized speech audio from the stream
	audioBuffer, err := audioStream.Read(4096)
	if err != nil && err.Error() != "EOF" {
		fmt.Printf("❌ Audio stream read error: %v\n", err)
		return
	}

	// Save audio to file if enabled
	if saveAudio && audioFile != nil {
		_, err := audioFile.Write(audioBuffer)
		if err != nil {
			fmt.Printf("❌ Failed to save synthesized speech: %v\n", err)
		}
	}

	// Play audio in real-time (Linux: `aplay`, Mac: `afplay`, Windows: `powershell`)
	cmd := exec.Command("aplay", "/dev/stdin")
	stdin, _ := cmd.StdinPipe()
	cmd.Start()
	stdin.Write(audioBuffer)
	stdin.Close()
	cmd.Wait()
}
