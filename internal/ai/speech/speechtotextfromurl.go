package speech

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	speechsdk "github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

// TranscribeFromLiveStream transcribes speech from a live video stream
func TranscribeFromLiveStream(subscriptionKey, region, streamURL, format string) error {
	// Create speech config
	speechConfig, err := speechsdk.NewSpeechConfigFromSubscription(subscriptionKey, region)
	if err != nil {
		return fmt.Errorf("failed to create speech config: %v", err)
	}
	defer speechConfig.Close()

	// Get direct audio URL if YouTube
	if format == "youtube" {
		out, err := exec.Command("yt-dlp", "-g", streamURL).Output()
		if err != nil {
			return fmt.Errorf("failed to get direct YouTube audio URL: %v", err)
		}
		streamURL = string(out)
	}

	// FFmpeg command to extract live audio
	ffmpegCmd := exec.Command("ffmpeg", "-i", streamURL, "-vn", "-ac", "1", "-ar", "16000", "-f", "wav", "pipe:1")
	stdout, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	if err := ffmpegCmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %v", err)
	}
	defer ffmpegCmd.Wait() // Ensure FFmpeg exits cleanly

	// Create a PushAudioInputStream
	audioFormat, err := audio.GetDefaultInputFormat()
	if err != nil {
		return fmt.Errorf("failed to get default audio format: %v", err)
	}
	defer audioFormat.Close()

	audioStream, err := audio.CreatePushAudioInputStreamFromFormat(audioFormat)
	if err != nil {
		return fmt.Errorf("failed to create push audio input stream: %v", err)
	}
	defer audioStream.Close()

	// Read the audio stream and send to Azure Speech SDK
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

	// Initialize recognizer with streaming audio input
	audioConfig, err := audio.NewAudioConfigFromStreamInput(audioStream)
	if err != nil {
		return fmt.Errorf("failed to create audio config: %v", err)
	}
	defer audioConfig.Close()

	speechRecognizer, err := speechsdk.NewSpeechRecognizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		return fmt.Errorf("failed to create speech recognizer: %v", err)
	}
	defer speechRecognizer.Close()

	// Continuous transcription
	fmt.Println("🎤 Transcribing live stream...")

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
