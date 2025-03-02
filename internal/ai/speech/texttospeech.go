package speech

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ngothientuong/tngo-ai-svcs/pkg/aitoken"
	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

type TextToSpeechClient struct {
	Endpoint string
	Key      string
	Region   string
}

type SynthesisRequest struct {
	Text   string `json:"text"`
	Voice  string `json:"voice"`
	Format string `json:"format"`
}

type Voice struct {
	Name         string   `json:"name"`
	ShortName    string   `json:"shortName"`
	Gender       string   `json:"gender"`
	Locale       string   `json:"locale"`
	LocaleName   string   `json:"localeName"`
	SampleRate   string   `json:"sampleRateHertz"`
	VoiceType    string   `json:"voiceType"`
	Status       string   `json:"status"`
	StyleList    []string `json:"styleList"`
	RolePlayList []string `json:"rolePlayList"`
}

func NewTextToSpeechClient(endpoint, key, region string) *TextToSpeechClient {
	return &TextToSpeechClient{
		Endpoint: endpoint,
		Key:      key,
		Region:   region,
	}
}

func (c *TextToSpeechClient) SynthesizeSpeech(text, voice, format string) ([]byte, error) {
	url := fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/v1", c.Region)
	requestBody := fmt.Sprintf("<speak version='1.0' xml:lang='en-US'><voice xml:lang='en-US' xml:gender='Female' name='%s'>%s</voice></speak>", voice, text)

	token, err := aitoken.GetToken(c.Endpoint, c.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Authorization":            "Bearer " + token,
		"Content-Type":             "application/ssml+xml",
		"X-Microsoft-OutputFormat": format,
		"User-Agent":               "curl",
	}

	resp, err := client.Post(url, bytes.NewBuffer([]byte(requestBody)), headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to synthesize speech: %s", body)
	}

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return audioData, nil
}

func (c *TextToSpeechClient) GetVoices() ([]Voice, error) {
	url := fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/voices/list", c.Region)

	token, err := aitoken.GetToken(c.Endpoint, c.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	resp, err := client.Get(url, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get voices: %s", body)
	}

	var voices []Voice
	err = json.NewDecoder(resp.Body).Decode(&voices)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return voices, nil
}
