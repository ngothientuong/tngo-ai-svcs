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
	"time"

	"github.com/google/uuid"
	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

type Project struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Settings     Settings  `json:"settings"`
	Created      string    `json:"created"`
	LastModified string    `json:"lastModified"`
	ThumbnailUri string    `json:"thumbnailUri"`
}

type Settings struct {
	DomainId string `json:"domainId"`
}

type Tag struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	ImageCount  int       `json:"imageCount"`
}

type Image struct {
	ID               uuid.UUID `json:"id"`
	Created          string    `json:"created"`
	Width            int       `json:"width"`
	Height           int       `json:"height"`
	ResizedImageUri  string    `json:"resizedImageUri"`
	OriginalImageUri string    `json:"originalImageUri"`
	ThumbnailUri     string    `json:"thumbnailUri"`
	Tags             []Tag     `json:"tags"`
}

type ImageCreateSummary struct {
	IsBatchSuccessful bool    `json:"isBatchSuccessful"`
	Images            []Image `json:"images"`
}

type Iteration struct {
	ID                        uuid.UUID            `json:"id"`
	Name                      string               `json:"name"`
	Status                    string               `json:"status"`
	Created                   string               `json:"created"`
	LastModified              string               `json:"lastModified"`
	ProjectId                 uuid.UUID            `json:"projectId"`
	Exportable                bool                 `json:"exportable"`
	DomainId                  *uuid.UUID           `json:"domainId,omitempty"`
	ExportableTo              []string             `json:"exportableTo"`
	TrainingType              string               `json:"trainingType"`
	ReservedBudgetInHours     int                  `json:"reservedBudgetInHours"`
	PublishName               string               `json:"publishName"`
	ClassificationType        string               `json:"classificationType"`
	CustomBaseModelInfo       *CustomBaseModelInfo `json:"customBaseModelInfo,omitempty"`
	OriginalPublishResourceId string               `json:"originalPublishResourceId,omitempty"`
	TrainingErrorDetails      string               `json:"trainingErrorDetails,omitempty"`
	TrainingTimeInMinutes     int                  `json:"trainingTimeInMinutes"`
	TrainedAt                 string               `json:"trainedAt"`
}

type CustomBaseModelInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type IterationPerformance struct {
	PerTagPerformance     []TagPerformance `json:"perTagPerformance"`
	Precision             float64          `json:"precision"`
	PrecisionStdDeviation float64          `json:"precisionStdDeviation"`
	Recall                float64          `json:"recall"`
	RecallStdDeviation    float64          `json:"recallStdDeviation"`
}

type TagPerformance struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Precision             float64   `json:"precision"`
	PrecisionStdDeviation float64   `json:"precisionStdDeviation"`
	Recall                float64   `json:"recall"`
	RecallStdDeviation    float64   `json:"recallStdDeviation"`
}

type ImagePrediction struct {
	ID          uuid.UUID    `json:"id"`
	Project     uuid.UUID    `json:"project"`
	Iteration   uuid.UUID    `json:"iteration"`
	Created     string       `json:"created"`
	Predictions []Prediction `json:"predictions"`
}

type Prediction struct {
	TagId       uuid.UUID `json:"tagId"`
	TagName     string    `json:"tagName"`
	Probability float64   `json:"probability"`
}

type PredictionResults struct {
	Predictions []Prediction `json:"Predictions"`
}

type ImageFileCreateEntry struct {
	Name     string `json:"name"`
	Contents string `json:"contents"`
}

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

	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Tag creation response: %s\n", responseBody)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create tag: %s", resp.Status)
	}

	var tag Tag
	err = json.NewDecoder(bytes.NewBuffer(responseBody)).Decode(&tag)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func GetTags(url, training_key, projectId string, iterationId *string) ([]Tag, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	params := map[string]interface{}{}
	if iterationId != nil {
		params["iterationId"] = *iterationId
	}
	resp, err := client.Get(url, headers, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get tags: status code %d", resp.StatusCode)
	}

	var tags []Tag
	err = json.NewDecoder(resp.Body).Decode(&tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func GetTagById(url, training_key, projectId, tagId string, iterationId *string) (*Tag, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	fullURL := fmt.Sprintf("%s/%s/tags/%s", url, projectId, tagId)
	params := map[string]interface{}{}
	if iterationId != nil {
		params["iterationId"] = *iterationId
	}
	resp, err := client.Get(fullURL, headers, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get tag: status code %d", resp.StatusCode)
	}

	var tag Tag
	err = json.NewDecoder(resp.Body).Decode(&tag)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func GetTagByName(url, training_key, projectId, tagName string, iterationId *string) (*Tag, error) {
	tags, err := GetTags(url, training_key, projectId, iterationId)
	if err != nil {
		return nil, err
	}
	fmt.Println("Looping through tags")
	for _, tag := range tags {
		if tag.Name == tagName {
			return &tag, nil
		}
	}

	return nil, fmt.Errorf("tag with name %s not found", tagName)
}

func GetProjects(url, training_key string) ([]Project, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	resp, err := client.Get(url, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get projects: status code %d", resp.StatusCode)
	}

	var projects []Project
	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func GetProjectById(url, training_key, projectId string) (*Project, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	fullURL := fmt.Sprintf("%s/%s", url, projectId)
	resp, err := client.Get(fullURL, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get project: status code %d", resp.StatusCode)
	}

	var project Project
	err = json.NewDecoder(resp.Body).Decode(&project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func GetProjectByName(url, training_key, projectName string) (*Project, error) {
	projects, err := GetProjects(url, training_key)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if project.Name == projectName {
			return &project, nil
		}
	}

	return nil, fmt.Errorf("project with name %s not found", projectName)
}

func GetImages(url, training_key string, params map[string]interface{}) ([]Image, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	resp, err := client.Get(url, headers, params)
	if err != nil {
		fmt.Println("Error getting images:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get images: status code %d", resp.StatusCode)
	}

	var images []Image
	err = json.NewDecoder(resp.Body).Decode(&images)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func GetImageByIds(url, training_key, projectId string, imageIds []string) ([]Image, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	fullURL := fmt.Sprintf("%s/%s/images", url, projectId)
	params := map[string]interface{}{
		"imageIds": imageIds,
	}
	resp, err := client.Get(fullURL, headers, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get images: status code %d", resp.StatusCode)
	}

	var images []Image
	err = json.NewDecoder(resp.Body).Decode(&images)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func CreateImagesFromData(url, training_key string, imageFiles [][]byte, tagID uuid.UUID) (*ImageCreateSummary, error) {
	// Append tagIds to the URL
	url = fmt.Sprintf("%s?tagIds=%s", url, tagID.String())

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for i, imageFile := range imageFiles {
		part, err := writer.CreateFormFile("imageData", fmt.Sprintf("image%d.jpg", i))
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %v", err)
		}
		_, err = part.Write(imageFile)
		if err != nil {
			return nil, fmt.Errorf("failed to write image file to form: %v", err)
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
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusMultiStatus {
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response body: %s\n", responseBody)
		return nil, fmt.Errorf("failed to upload images: %s", resp.Status)
	}

	var summary ImageCreateSummary
	err = json.NewDecoder(resp.Body).Decode(&summary)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func DeleteImages(url, training_key string, imageIds []string) error {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	fullURL := fmt.Sprintf("%s/images", url)
	params := map[string]interface{}{
		"imageIds": imageIds,
	}

	resp, err := client.Delete(fullURL, headers, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to delete images: %s", resp.Status)
	}

	return nil
}

func CreateImagesFromFiles(url, training_key string, images []ImageFileCreateEntry, tagIds []string) (*ImageCreateSummary, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	fullURL := fmt.Sprintf("%s/images/files", url)
	body := map[string]interface{}{
		"images": images,
		"tagIds": tagIds,
	}
	resp, err := client.Post(fullURL, body, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusMultiStatus {
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response body: %s\n", responseBody)
		return nil, fmt.Errorf("failed to create images from files: %s", resp.Status)
	}

	var summary ImageCreateSummary
	err = json.NewDecoder(resp.Body).Decode(&summary)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func TrainProject(url, training_key string, params interface{}) (*Iteration, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	resp, err := client.Post(url, params, headers, nil)
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

func GetIterations(url, training_key string) ([]Iteration, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	resp, err := client.Get(url, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get iterations: status code %d", resp.StatusCode)
	}

	var iterations []Iteration
	err = json.NewDecoder(resp.Body).Decode(&iterations)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return iterations, nil
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

func GetIterationLatest(url, training_key string) (*Iteration, error) {
	iterations, err := GetIterations(url, training_key)
	if err != nil {
		return nil, fmt.Errorf("failed to get iterations: %v", err)
	}

	if len(iterations) == 0 {
		return nil, fmt.Errorf("no iterations found for url %s", url)
	}

	latestIteration := iterations[0]
	latestTime, err := ParseTime(latestIteration.Created, latestIteration.LastModified)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}

	for _, iteration := range iterations[1:] {
		createdTime, err := ParseTime(iteration.Created, iteration.LastModified)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time: %v", err)
		}
		if createdTime.After(latestTime) {
			latestIteration = iteration
			latestTime = createdTime
		}
	}

	return &latestIteration, nil
}

func GetIterationPerformance(url, training_key string, params map[string]string) (*IterationPerformance, error) {
	client := httpclient.NewClient()
	headers := map[string]string{
		"Training-Key": training_key,
	}
	resp, err := client.Get(url, headers, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get iteration performance: status code %d", resp.StatusCode)
	}

	var performance IterationPerformance
	err = json.NewDecoder(resp.Body).Decode(&performance)
	if err != nil {
		return nil, err
	}

	return &performance, nil
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

func QuickTestImage(url, training_key, sampleDataDirectory, imagePath string, params map[string]string) (*ImagePrediction, error) {
	imageFile, err := os.ReadFile(path.Join(sampleDataDirectory, imagePath))
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %v", err)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("imageData", "image.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = part.Write(imageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to write image file to form: %v", err)
	}
	writer.Close()

	client := httpclient.NewClient()
	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"Training-Key": training_key,
	}
	resp, err := client.Post(url, body, headers, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Debugging output
	fmt.Printf("Uploading image to URL: %s\n", url)
	fmt.Printf("Headers: %v\n", headers)

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to classify image: status code %d: %s - %s", resp.StatusCode, http.StatusText(resp.StatusCode), responseBody)
	}

	var results ImagePrediction
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

func QuickTestImageUrl(url, training_key string, imageUrl string, params interface{}) (*ImagePrediction, error) {
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

	var results ImagePrediction
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

/*
HELPER FUNCTIONS
*/
// func AreImageSetsDifferent(previousImages []Image, currentImages [][]byte) bool {
// 	if len(previousImages) != len(currentImages) {
// 		return true
// 	}

// 	// Sort previous images by their URIs
// 	sort.Slice(previousImages, func(i, j int) bool {
// 		return previousImages[i].OriginalImageUri < previousImages[j].OriginalImageUri
// 	})

// 	// Sort current images by their hashes
// 	sort.Slice(currentImages, func(i, j int) bool {
// 		return string(HashImage(currentImages[i])) < string(HashImage(currentImages[j]))
// 	})

// 	for i, image := range previousImages {
// 		imageData, err := DownloadImageData(image.OriginalImageUri)
// 		if err != nil {
// 			return true
// 		}
// 		if !bytes.Equal(HashImage(imageData), HashImage(currentImages[i])) {
// 			return true
// 		}
// 	}

// 	return false
// }

// func DownloadImageData(imageUri string) ([]byte, error) {
// 	resp, err := http.Get(imageUri)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
// 	}

// 	imageData, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return imageData, nil
// }

// func HashImage(imageData []byte) []byte {
// 	hash := sha256.Sum256(imageData)
// 	return hash[:]
// }

func ParseTime(created, lastModified string) (time.Time, error) {
	createdTime, err := time.Parse(time.RFC3339, created)
	if err != nil {
		return time.Time{}, err
	}
	lastModifiedTime, err := time.Parse(time.RFC3339, lastModified)
	if err != nil {
		return time.Time{}, err
	}
	if lastModifiedTime.After(createdTime) {
		return lastModifiedTime, nil
	}
	return createdTime, nil
}
