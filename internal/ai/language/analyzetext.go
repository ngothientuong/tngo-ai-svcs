package language

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

type AnalyzeTextClient struct {
	Endpoint   string
	Key        string
	APIVersion string
}

type AnalyzeTextRequest struct {
	Kind          string                 `json:"kind"`
	AnalysisInput MultiLanguageInput     `json:"analysisInput"`
	Parameters    map[string]interface{} `json:"parameters,omitempty"`
}

type MultiLanguageInput struct {
	Documents []Document `json:"documents"`
}

type Document struct {
	ID       string `json:"id"`
	Language string `json:"language,omitempty"`
	Text     string `json:"text"`
}

type AnalyzeTextResponse struct {
	Results interface{} `json:"results"`
}

func NewAnalyzeTextClient(endpoint, key, apiVersion string) *AnalyzeTextClient {
	return &AnalyzeTextClient{
		Endpoint:   endpoint,
		Key:        key,
		APIVersion: apiVersion,
	}
}

func (c *AnalyzeTextClient) AnalyzeText(kind string, documents []Document, parameters map[string]interface{}) (*AnalyzeTextResponse, error) {
	url := fmt.Sprintf("%s/language/:analyze-text?api-version=%s", c.Endpoint, c.APIVersion)
	requestBody, err := json.Marshal(AnalyzeTextRequest{
		Kind: kind,
		AnalysisInput: MultiLanguageInput{
			Documents: documents,
		},
		Parameters: parameters,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": c.Key,
		"Content-Type":              "application/json",
	}

	resp, err := client.Post(url, bytes.NewBuffer(requestBody), headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to analyze text: %s", body)
	}

	var response AnalyzeTextResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}
