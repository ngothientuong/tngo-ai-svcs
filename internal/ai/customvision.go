package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

func CreateProject(url, training_key, project_name string) (*Project, error) {
	client := httpclient.NewClient(training_key)
	resp, err := client.Post(url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Project creation response: %s\n", responseBody)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create project: %s", resp.Status)
	}

	var project Project
	err = json.NewDecoder(bytes.NewBuffer(responseBody)).Decode(&project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func CreateTag(url, training_key string, projectID uuid.UUID, tagName string) (*Tag, error) {
	client := httpclient.NewClient(training_key)
	resp, err := client.Post(url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tag Tag
	err = json.NewDecoder(resp.Body).Decode(&tag)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func UploadImages(url, training_key string, projectID uuid.UUID, imageFiles [][]byte, tagID uuid.UUID) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for i, imageFile := range imageFiles {
		part, _ := writer.CreateFormFile("imageData", fmt.Sprintf("image%d.jpg", i))
		part.Write(imageFile)
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}

	client := httpclient.NewClient(training_key)
	resp, err := client.Post(url, body, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusMultiStatus {
		return fmt.Errorf("failed to upload images: %s", resp.Status)
	}

	return nil
}

func TrainProject(url, training_key string, projectID uuid.UUID) (*Iteration, error) {
	client := httpclient.NewClient(training_key)
	resp, err := client.Post(url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var iteration Iteration
	err = json.NewDecoder(resp.Body).Decode(&iteration)
	if err != nil {
		return nil, err
	}

	return &iteration, nil
}

func GetIteration(url, training_key string, projectID uuid.UUID, iterationID uuid.UUID) (*Iteration, error) {
	client := httpclient.NewClient(training_key)
	resp, err := client.Get(url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var iteration Iteration
	err = json.NewDecoder(resp.Body).Decode(&iteration)
	if err != nil {
		return nil, err
	}

	return &iteration, nil
}

func PublishIteration(url, training_key, iteration_publish_name, prediction_resource_id string, projectID uuid.UUID, iterationID uuid.UUID) error {
	client := httpclient.NewClient(training_key)
	resp, err := client.Post(url, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to publish iteration: %s", resp.Status)
	}

	return nil
}

func QuickTestImage(url, training_key, sampleDataDirectory, imagePath string, projectID uuid.UUID) (*PredictionResults, error) {
	imageFile, _ := os.ReadFile(path.Join(sampleDataDirectory, imagePath))
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("imageData", "image.jpg")
	part.Write(imageFile)
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}

	client := httpclient.NewClient(training_key)
	resp, err := client.Post(url, body, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results PredictionResults
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

type Project struct {
	ID   uuid.UUID `json:"Id"`
	Name string    `json:"Name"`
}

type Domain struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

type Tag struct {
	ID   uuid.UUID `json:"Id"`
	Name string    `json:"Name"`
}

type Iteration struct {
	ID     uuid.UUID `json:"Id"`
	Status string    `json:"Status"`
}

type PredictionResults struct {
	Predictions []Prediction `json:"Predictions"`
}

type Prediction struct {
	TagName     string  `json:"TagName"`
	Probability float64 `json:"Probability"`
}
