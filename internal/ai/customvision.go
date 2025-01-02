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

func CreateProject(url, training_key string, params interface{}) (*Project, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	resp, err := client.Post(url, nil, headers, params)
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

func CreateTag(url, training_key string, params interface{}) (*Tag, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	resp, err := client.Post(url, nil, headers, params)
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

func UploadImages(url, training_key string, imageFiles [][]byte, tagID uuid.UUID) error {
	// Append tagIds to the URL
	url = fmt.Sprintf("%s?tagIds=%s", url, tagID.String())

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for i, imageFile := range imageFiles {
		part, err := writer.CreateFormFile("imageData", fmt.Sprintf("image%d.jpg", i))
		if err != nil {
			return fmt.Errorf("failed to create form file: %v", err)
		}
		_, err = part.Write(imageFile)
		if err != nil {
			return fmt.Errorf("failed to write image file to form: %v", err)
		}
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"Training-Key": training_key,
	}

	// Debugging output
	fmt.Printf("Uploading images to URL: %s\n", url)
	fmt.Printf("Headers: %v\n", headers)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusMultiStatus {
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response body: %s\n", responseBody)
		return fmt.Errorf("failed to upload images: %s", resp.Status)
	}

	return nil
}

func TrainProject(url, training_key string, params interface{}) (*Iteration, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	resp, err := client.Post(url, nil, headers, params)
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

func GetIteration(url, training_key string, params interface{}) (*Iteration, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	resp, err := client.Get(url, headers, params)
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

func PublishIteration(url, training_key string, params interface{}) error {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	resp, err := client.Post(url, nil, headers, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to publish iteration: %s", resp.Status)
	}

	return nil
}

func QuickTestImage(url, training_key, sampleDataDirectory, imagePath string, params interface{}) (*PredictionResults, error) {
	imageFile, err := os.ReadFile(path.Join(sampleDataDirectory, imagePath))
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %v", err)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("imageData", "image.jpg")
	part.Write(imageFile)
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"Training-Key": training_key,
	}

	client := httpclient.NewClient()
	resp, err := client.Post(url, body, headers, params)
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

func QuickTestImageUrl(url, training_key string, imageUrl string, params interface{}) (*PredictionResults, error) {
	payload := map[string]string{"url": imageUrl}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(payload)

	headers := map[string]string{
		"Content-Type": "application/json",
		"Training-Key": training_key,
	}

	client := httpclient.NewClient()
	resp, err := client.Post(url, body, headers, params)
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
