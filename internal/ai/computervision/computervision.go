package computervision

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

type Index struct {
	Name                 string         `json:"name"`
	MetadataSchema       MetadataSchema `json:"metadataSchema"`
	Features             []Feature      `json:"features"`
	UserData             interface{}    `json:"userData"`
	ETag                 string         `json:"eTag"`
	CreatedDateTime      string         `json:"createdDateTime"`
	LastModifiedDateTime string         `json:"lastModifiedDateTime"`
}

type MetadataSchema struct {
	Language string          `json:"language"`
	Fields   []MetadataField `json:"fields"`
}

type MetadataField struct {
	Name       string `json:"name"`
	Searchable bool   `json:"searchable"`
	Filterable bool   `json:"filterable"`
	Type       string `json:"type"`
}

type Feature struct {
	Name         string `json:"name"`
	ModelVersion string `json:"modelVersion"`
	Domain       string `json:"domain"`
}

type Ingestion struct {
	Name                 string `json:"name"`
	State                string `json:"state"`
	CreatedDateTime      string `json:"createdDateTime"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
}

type IngestionRequest struct {
	Videos []IngestionDocument `json:"videos"`
}

type IngestionDocument struct {
	Mode        string            `json:"mode"`
	DocumentId  string            `json:"documentId"`
	DocumentUrl string            `json:"documentUrl"`
	Metadata    map[string]string `json:"metadata"`
}

type SearchQuery struct {
	QueryText string        `json:"queryText"`
	Filters   SearchFilters `json:"filters"`
}

type SearchFilters struct {
	StringFilters  []StringFilter `json:"stringFilters"`
	FeatureFilters []string       `json:"featureFilters"`
}

type StringFilter struct {
	FieldName string   `json:"fieldName"`
	Values    []string `json:"values"`
}

type SearchResponse struct {
	Value []SearchResult `json:"value"`
}

type SearchResult struct {
	DocumentId   string  `json:"documentId"`
	DocumentKind string  `json:"documentKind"`
	Start        string  `json:"start"`
	End          string  `json:"end"`
	Best         string  `json:"best"`
	Relevance    float64 `json:"relevance"`
}

func CreateIndex(url, subscriptionKey string, params interface{}) (*Index, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": subscriptionKey,
		"Content-Type":              "application/json",
	}
	resp, err := client.Put(url, params, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create index due to status code: %d", resp.StatusCode)
	}

	var index Index
	err = json.NewDecoder(bytes.NewBuffer(responseBody)).Decode(&index)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v, response: %s", err, responseBody)
	}

	fmt.Println("Exiting CreateIndex function")
	return &index, nil
}

func DeleteIndex(url, subscriptionKey string) error {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": subscriptionKey,
	}
	resp, err := client.Delete(url, headers, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete index: %s", resp.Status)
	}

	return nil
}

func CheckIndexExists(url, subscriptionKey string) (bool, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": subscriptionKey,
	}
	resp, err := client.Get(url, headers, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		return false, fmt.Errorf("failed to check if index exists: %s", resp.Status)
	}
}

func GetIndex(url, subscriptionKey string) (*Index, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": subscriptionKey,
	}
	resp, err := client.Get(url, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get index: %s", resp.Status)
	}

	var index Index
	err = json.NewDecoder(resp.Body).Decode(&index)
	if err != nil {
		return nil, err
	}

	return &index, nil
}

func CheckIngestionExists(url, subscriptionKey string) (bool, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": subscriptionKey,
	}
	resp, err := client.Get(url, headers, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		responseBody, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to check if ingestion exists: %s, response: %s", resp.Status, responseBody)
	}
}

func CreateIngestion(url, subscriptionKey string, params interface{}) (*Ingestion, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": subscriptionKey,
		"Content-Type":              "application/json",
	}

	// Convert params to JSON
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %v", err)
	}

	resp, err := client.Put(url, bytes.NewBuffer(jsonParams), headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Ingestion creation response: %s\n", responseBody)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("failed to create ingestion: %s", resp.Status)
	}

	var ingestion Ingestion
	err = json.NewDecoder(bytes.NewBuffer(responseBody)).Decode(&ingestion)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v, response: %s", err, responseBody)
	}

	return &ingestion, nil
}

func UpdateIngestion(url, subscriptionKey string, params interface{}) (*Ingestion, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": subscriptionKey,
		"Content-Type":              "application/json",
	}

	// Convert params to JSON
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %v", err)
	}

	resp, err := client.Patch(url, bytes.NewBuffer(jsonParams), headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Ingestion update response: %s\n", responseBody)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("failed to update ingestion: %s", resp.Status)
	}

	var ingestion Ingestion
	err = json.NewDecoder(bytes.NewBuffer(responseBody)).Decode(&ingestion)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v, response: %s", err, responseBody)
	}

	return &ingestion, nil
}

func GetIngestionStatus(url, subscriptionKey string) (*Ingestion, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": subscriptionKey,
	}
	resp, err := client.Get(url, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get ingestion status: %s", resp.Status)
	}

	var ingestion Ingestion
	err = json.NewDecoder(resp.Body).Decode(&ingestion)
	if err != nil {
		return nil, err
	}

	return &ingestion, nil
}

func ListIngestions(url, subscriptionKey string) ([]Ingestion, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": subscriptionKey,
	}
	resp, err := client.Get(url, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list ingestions: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Printf("Response body: %s\n", body)

	var result struct {
		Value []Ingestion `json:"value"`
	}
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result.Value, nil
}

func ListDocuments(url, subscriptionKey string) ([]IngestionDocument, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": subscriptionKey,
	}
	resp, err := client.Get(url, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list documents: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Printf("Response body: %s\n", body)

	var result struct {
		Value []IngestionDocument `json:"value"`
	}
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result.Value, nil
}

func SearchByText(url, subscriptionKey string, params interface{}) (*SearchResponse, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": subscriptionKey,
		"Content-Type":              "application/json",
	}

	// Convert params to JSON
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %v", err)
	}

	resp, err := client.Post(url, bytes.NewBuffer(jsonParams), headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search by text: %s", resp.Status)
	}

	// Read response raw body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Printf("Response body: %s\n", body)

	var result SearchResponse
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result, nil
}
