package speech

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

type TextToSpeechCustomNeuralClient struct {
	Endpoint   string
	Key        string
	Region     string
	APIVersion string
}

type Project struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Description string `json:"description,omitempty"`
}

type ListProjectsResponse struct {
	Value []Project `json:"value"`
}

func NewTextToSpeechCustomNeuralClient(endpoint, key, region, apiVersion string) *TextToSpeechCustomNeuralClient {
	return &TextToSpeechCustomNeuralClient{
		Endpoint:   endpoint,
		Key:        key,
		Region:     region,
		APIVersion: apiVersion,
	}
}

func (c *TextToSpeechCustomNeuralClient) CheckProjectExists(projectID string) (bool, error) {
	url := fmt.Sprintf("%s/customvoice/projects/%s?api-version=%s", c.Endpoint, projectID, c.APIVersion)

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": c.Key,
		"Content-Type":              "application/json",
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
		return false, fmt.Errorf("failed to check project existence: %s", body)
	}
}

func (c *TextToSpeechCustomNeuralClient) GetProjectID(projectName string) (string, error) {
	url := fmt.Sprintf("%s/customvoice/projects?api-version=%s", c.Endpoint, c.APIVersion)

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": c.Key,
		"Content-Type":              "application/json",
	}

	resp, err := client.Get(url, headers, nil)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get projects: %s", body)
	}

	var response ListProjectsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	for _, project := range response.Value {
		if project.Kind == projectName {
			return project.ID, nil
		}
	}

	return "", fmt.Errorf("project with name %s not found", projectName)
}

func (c *TextToSpeechCustomNeuralClient) CreateProject(projectID, kind, description string) (*Project, error) {
	url := fmt.Sprintf("%s/customvoice/projects/%s?api-version=%s", c.Endpoint, projectID, c.APIVersion)

	project := Project{
		ID:          projectID,
		Kind:        kind,
		Description: description,
	}
	requestBody, err := json.Marshal(project)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": c.Key,
		"Content-Type":              "application/json",
	}

	resp, err := client.Put(url, bytes.NewBuffer(requestBody), headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create project: %s", body)
	}

	var createdProject Project
	err = json.NewDecoder(resp.Body).Decode(&createdProject)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &createdProject, nil
}

func (c *TextToSpeechCustomNeuralClient) DeleteProject(projectID string) error {
	url := fmt.Sprintf("%s/customvoice/projects/%s?api-version=%s", c.Endpoint, projectID, c.APIVersion)

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": c.Key,
		"Content-Type":              "application/json",
	}

	resp, err := client.Delete(url, headers, nil)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete project: %s", body)
	}

	return nil
}

func (c *TextToSpeechCustomNeuralClient) ListProjects() ([]Project, error) {
	url := fmt.Sprintf("%s/customvoice/projects?api-version=%s", c.Endpoint, c.APIVersion)

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": c.Key,
		"Content-Type":              "application/json",
	}

	resp, err := client.Get(url, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list projects: %s", body)
	}

	var response ListProjectsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return response.Value, nil
}
