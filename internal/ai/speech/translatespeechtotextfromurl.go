package speech

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	speechsdk "github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

var (
	doneChan chan struct{}
)

// TranslateSpeechFromURL translates speech from a given URL to the specified language
func TranslateSpeechFromURL(subscriptionKey, region, url, format, fromLang, toLang string) error {
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

	// Set source language (speech recognition language)
	err = config.SetSpeechRecognitionLanguage(fromLang)
	if err != nil {
		return fmt.Errorf("failed to set speech recognition language: %v", err)
	}

	// Set target language
	err = config.AddTargetLanguage(toLang)
	if err != nil {
		return fmt.Errorf("failed to add target language: %v", err)
	}

	// Create a PushAudioInputStream to stream audio from FFmpeg
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

	// Handle translation results
	recognizer.Recognizing(func(event speechsdk.TranslationRecognitionEventArgs) {
		translations := event.Result.GetTranslations()
		if translatedText, ok := translations[toLang]; ok {
			fmt.Printf("🔵 Partial Translation: %s\n", translatedText)
		}
	})

	recognizer.Recognized(func(event speechsdk.TranslationRecognitionEventArgs) {
		translations := event.Result.GetTranslations()
		if translatedText, ok := translations[toLang]; ok {
			fmt.Printf("✅ Final Translation: %s\n", translatedText)
		}
	})

	recognizer.Canceled(func(event speechsdk.TranslationRecognitionCanceledEventArgs) {
		fmt.Printf("❌ Canceled: %v\n", event.ErrorDetails)
		closeDoneChan()
	})

	recognizer.SessionStopped(func(event speechsdk.SessionEventArgs) {
		fmt.Println("🔴 Translation session stopped")
		closeDoneChan()
	})

	// Start continuous recognition
	errChan := recognizer.StartContinuousRecognitionAsync()
	if err := <-errChan; err != nil {
		return fmt.Errorf("failed to start continuous recognition: %v", err)
	}

	// Wait for recognition to complete
	waitForDoneChan()
	return nil
}

// Helper function to close the done channel safely
func closeDoneChan() {
	mu.Lock()
	defer mu.Unlock()
	if doneChan != nil {
		close(doneChan)
		doneChan = nil
	}
}

// Helper function to wait for translation completion
func waitForDoneChan() {
	mu.Lock()
	doneChan = make(chan struct{})
	mu.Unlock()
	<-doneChan
}

// NewPushAudioInputStreamFromURL streams audio from a URL using FFmpeg
func NewPushAudioInputStreamFromURL(url, format string) (*audio.PushAudioInputStream, error) {
	// If input is a YouTube URL, get direct audio stream
	if format == "youtube" {
		out, err := exec.Command("yt-dlp", "-g", url).Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get YouTube audio stream URL: %v", err)
		}
		url = strings.TrimSpace(string(out))
	}

	// FFmpeg command to extract live audio
	ffmpegCmd := exec.Command("ffmpeg", "-i", url, "-vn", "-ac", "1", "-ar", "16000", "-f", "wav", "pipe:1")
	stdout, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	if err := ffmpegCmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start ffmpeg: %v", err)
	}

	// Create PushAudioInputStream
	audioStream, err := audio.CreatePushAudioInputStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create push audio input stream: %v", err)
	}

	// Stream FFmpeg output to Azure SDK
	go func() {
		buf := bufio.NewReader(stdout)
		for {
			data := make([]byte, 4096)
			n, err := buf.Read(data)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("Audio read error: %v\n", err)
				break
			}
			if n > 0 {
				err = audioStream.Write(data[:n])
				if err != nil {
					fmt.Printf("Audio write error: %v\n", err)
					break
				}
			}
		}
		audioStream.Close()
	}()

	return audioStream, nil
}

// isValidSourceLanguage checks if the source language is supported by Azure Speech-to-Text
func isValidSourceLanguage(lang string) bool {
	supportedSourceLanguages := []string{"en-US", "es-ES", "fr-FR", "de-DE", "zh-CN", "ja-JP", "ar-SA", "ru-RU", "it-IT", "pt-BR"}
	for _, l := range supportedSourceLanguages {
		if l == lang {
			return true
		}
	}
	return false
}

// isValidTargetLanguage checks if the target language is supported by Azure Translation
func isValidTargetLanguage(lang string) bool {
	supportedTargetLanguages := []string{"en", "es", "fr", "de", "zh", "ja", "ar", "ru", "it", "pt"}
	for _, l := range supportedTargetLanguages {
		if l == lang {
			return true
		}
	}
	return false
}
