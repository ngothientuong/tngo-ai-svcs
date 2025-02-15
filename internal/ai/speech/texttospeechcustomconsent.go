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

type TextToSpeechCustomConsentClient struct {
	Endpoint   string
	Key        string
	Region     string
	APIVersion string
}

type Consent struct {
	ID                 string `json:"id"`
	Description        string `json:"description"`
	ProjectID          string `json:"projectId"`
	VoiceTalentName    string `json:"voiceTalentName"`
	CompanyName        string `json:"companyName"`
	AudioURL           string `json:"audioUrl"`
	Locale             string `json:"locale"`
	Status             string `json:"status"`
	CreatedDateTime    string `json:"createdDateTime"`
	LastActionDateTime string `json:"lastActionDateTime"`
}

type ListConsentsResponse struct {
	Value []Consent `json:"value"`
}

func NewTextToSpeechCustomConsentClient(endpoint, key, region, apiVersion string) *TextToSpeechCustomConsentClient {
	return &TextToSpeechCustomConsentClient{
		Endpoint:   endpoint,
		Key:        key,
		Region:     region,
		APIVersion: apiVersion,
	}
}

func (c *TextToSpeechCustomConsentClient) CheckConsentExists(consentID string) (bool, error) {
	url := fmt.Sprintf("%s/customvoice/consents/%s?api-version=%s", c.Endpoint, consentID, c.APIVersion)

	token, err := aitoken.GetToken(c.Endpoint, c.Key)
	if err != nil {
		return false, fmt.Errorf("failed to get token: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	resp, err := client.Get(url, headers, nil)
	if err != nil {
		return false, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to check consent existence: %s", body)
	}
}

func (c *TextToSpeechCustomConsentClient) GetConsent(consentID string) (*Consent, error) {
	url := fmt.Sprintf("%s/customvoice/consents/%s?api-version=%s", c.Endpoint, consentID, c.APIVersion)

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
		return nil, fmt.Errorf("failed to get consent: %s", body)
	}

	var consent Consent
	err = json.NewDecoder(resp.Body).Decode(&consent)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &consent, nil
}

func (c *TextToSpeechCustomConsentClient) CreateConsent(consentID, description, projectID, voiceTalentName, companyName, audioURL, locale string) (*Consent, error) {
	url := fmt.Sprintf("%s/customvoice/consents/%s?api-version=%s", c.Endpoint, consentID, c.APIVersion)

	token, err := aitoken.GetToken(c.Endpoint, c.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	consent := Consent{
		ID:              consentID,
		Description:     description,
		ProjectID:       projectID,
		VoiceTalentName: voiceTalentName,
		CompanyName:     companyName,
		AudioURL:        audioURL,
		Locale:          locale,
	}
	requestBody, err := json.Marshal(consent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	resp, err := client.Put(url, bytes.NewBuffer(requestBody), headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create consent: %s", body)
	}

	var createdConsent Consent
	err = json.NewDecoder(resp.Body).Decode(&createdConsent)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &createdConsent, nil
}

func (c *TextToSpeechCustomConsentClient) DeleteConsent(consentID string) error {
	url := fmt.Sprintf("%s/customvoice/consents/%s?api-version=%s", c.Endpoint, consentID, c.APIVersion)

	token, err := aitoken.GetToken(c.Endpoint, c.Key)
	if err != nil {
		return fmt.Errorf("failed to get token: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	resp, err := client.Delete(url, headers, nil)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete consent: %s", body)
	}

	return nil
}

func (c *TextToSpeechCustomConsentClient) ListConsents() ([]Consent, error) {
	url := fmt.Sprintf("%s/customvoice/consents?api-version=%s", c.Endpoint, c.APIVersion)

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
		return nil, fmt.Errorf("failed to list consents: %s", body)
	}

	var response ListConsentsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return response.Value, nil
}

func (c *TextToSpeechCustomConsentClient) PostConsent(consentID, description, projectID, voiceTalentName, companyName, audioURL, locale string) (*Consent, error) {
	url := fmt.Sprintf("%s/customvoice/consents/%s?api-version=%s", c.Endpoint, consentID, c.APIVersion)

	token, err := aitoken.GetToken(c.Endpoint, c.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	consent := Consent{
		ID:              consentID,
		Description:     description,
		ProjectID:       projectID,
		VoiceTalentName: voiceTalentName,
		CompanyName:     companyName,
		AudioURL:        audioURL,
		Locale:          locale,
	}
	requestBody, err := json.Marshal(consent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	resp, err := client.Post(url, bytes.NewBuffer(requestBody), headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to post consent: %s", body)
	}

	var createdConsent Consent
	err = json.NewDecoder(resp.Body).Decode(&createdConsent)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &createdConsent, nil
}
