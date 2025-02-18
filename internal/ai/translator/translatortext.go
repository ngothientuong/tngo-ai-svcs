package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

type TranslatorClient struct {
	Endpoint   string
	Key        string
	Region     string
	APIVersion string
}

func NewTranslatorClient(endpoint, key, region, apiVersion string) *TranslatorClient {
	return &TranslatorClient{
		Endpoint:   endpoint,
		Key:        key,
		Region:     region,
		APIVersion: apiVersion,
	}
}

func (c *TranslatorClient) GetLanguages(scope string) ([]byte, error) {
	if scope == "" {
		scope = "translation,transliteration,dictionary"
	}
	url := fmt.Sprintf("%s/languages?api-version=%s&scope=%s", c.Endpoint, c.APIVersion, scope)
	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key":    c.Key,
		"Ocp-Apim-Subscription-Region": c.Region,
		"Content-Type":                 "application/json",
	}

	resp, err := client.Get(url, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(json.RawMessage(bodyBytes), "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to format JSON: %v", err)
	}

	return prettyJSON, nil
}

func (c *TranslatorClient) Translate(texts []string, to []string) ([]byte, error) {
	urlStr := fmt.Sprintf("%s/translate?api-version=%s", c.Endpoint, c.APIVersion)
	params := url.Values{}
	for _, lang := range to {
		params.Add("to", lang)
	}
	urlStr = fmt.Sprintf("%s&%s", urlStr, params.Encode())

	var body []map[string]interface{}
	for _, text := range texts {
		body = append(body, map[string]interface{}{
			"Text": text,
		})
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key":    c.Key,
		"Ocp-Apim-Subscription-Region": c.Region,
		"Content-Type":                 "application/json; charset=UTF-8",
	}

	resp, err := client.Post(urlStr, bytes.NewBuffer(bodyBytes), headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed: %s", string(bodyBytes))
	}

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(json.RawMessage(bodyBytes), "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to format JSON: %v", err)
	}

	return prettyJSON, nil
}

func (c *TranslatorClient) Transliterate(texts []string, language, fromScript, toScript string) ([]byte, error) {
	urlStr := fmt.Sprintf("%s/transliterate?api-version=%s&language=%s&fromScript=%s&toScript=%s", c.Endpoint, c.APIVersion, language, fromScript, toScript)

	var body []map[string]interface{}
	for _, text := range texts {
		body = append(body, map[string]interface{}{
			"Text": text,
		})
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key":    c.Key,
		"Ocp-Apim-Subscription-Region": c.Region,
		"Content-Type":                 "application/json; charset=UTF-8",
	}

	resp, err := client.Post(urlStr, bytes.NewBuffer(bodyBytes), headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed: %s", string(bodyBytes))
	}

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(json.RawMessage(bodyBytes), "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to format JSON: %v", err)
	}

	return prettyJSON, nil
}

func (c *TranslatorClient) Detect(texts []string) ([]byte, error) {
	url := fmt.Sprintf("%s/detect?api-version=%s", c.Endpoint, c.APIVersion)

	var body []map[string]interface{}
	for _, text := range texts {
		body = append(body, map[string]interface{}{
			"Text": text,
		})
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key":    c.Key,
		"Ocp-Apim-Subscription-Region": c.Region,
		"Content-Type":                 "application/json; charset=UTF-8",
	}

	resp, err := client.Post(url, bytes.NewBuffer(bodyBytes), headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed: %s", string(bodyBytes))
	}

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(json.RawMessage(bodyBytes), "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to format JSON: %v", err)
	}

	return prettyJSON, nil
}

func (c *TranslatorClient) DictionaryLookup(texts []string, from, to string) ([]byte, error) {
	urlStr := fmt.Sprintf("%s/dictionary/lookup?api-version=%s&from=%s&to=%s", c.Endpoint, c.APIVersion, from, to)

	var body []map[string]interface{}
	for _, text := range texts {
		body = append(body, map[string]interface{}{
			"Text": text,
		})
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key":    c.Key,
		"Ocp-Apim-Subscription-Region": c.Region,
		"Content-Type":                 "application/json; charset=UTF-8",
	}

	resp, err := client.Post(urlStr, bytes.NewBuffer(bodyBytes), headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed: %s", string(bodyBytes))
	}

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(json.RawMessage(bodyBytes), "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to format JSON: %v", err)
	}

	return prettyJSON, nil
}

func (c *TranslatorClient) DictionaryExamples(texts []map[string]string, from, to string) ([]byte, error) {
	urlStr := fmt.Sprintf("%s/dictionary/examples?api-version=%s&from=%s&to=%s", c.Endpoint, c.APIVersion, from, to)

	bodyBytes, err := json.Marshal(texts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key":    c.Key,
		"Ocp-Apim-Subscription-Region": c.Region,
		"Content-Type":                 "application/json; charset=UTF-8",
	}

	resp, err := client.Post(urlStr, bytes.NewBuffer(bodyBytes), headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed: %s", string(bodyBytes))
	}

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(json.RawMessage(bodyBytes), "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to format JSON: %v", err)
	}

	return prettyJSON, nil
}
