package computervision

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

type ImageAnalysisClient struct {
	Endpoint string
	Key      string
}

type ImageAnalysisRequest struct {
	URL string `json:"url"`
}

type VisualFeature string

const (
	VisualFeatureTags          VisualFeature = "tags"
	VisualFeatureObjects       VisualFeature = "objects"
	VisualFeatureCaption       VisualFeature = "caption"
	VisualFeatureDenseCaptions VisualFeature = "denseCaptions"
	VisualFeatureRead          VisualFeature = "read"
	VisualFeatureSmartCrops    VisualFeature = "smartCrops"
	VisualFeaturePeople        VisualFeature = "people"
)

type ImageAnalysisResponse struct {
	Tags       TagsResult    `json:"tagsResult"`
	Objects    ObjectsResult `json:"objectsResult"`
	Caption    CaptionResult `json:"captionResult"`
	ReadResult ReadResult    `json:"readResult"`
}

type TagsResult struct {
	Values []Tag `json:"values"`
}

type ObjectsResult struct {
	Values []DetectedObject `json:"values"`
}

type CaptionResult struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
}

type ReadResult struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Lines []Line `json:"lines"`
}

type Line struct {
	Text            string        `json:"text"`
	BoundingPolygon []BoundingBox `json:"boundingPolygon"`
	Words           []Word        `json:"words"`
}

type Word struct {
	Text            string        `json:"text"`
	BoundingPolygon []BoundingBox `json:"boundingPolygon"`
	Confidence      float64       `json:"confidence"`
}

type BoundingBox struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Tag struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

type DetectedObject struct {
	BoundingBox Rectangle `json:"boundingBox"`
	Tags        []Tag     `json:"tags"`
}

type Rectangle struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

func NewImageAnalysisClient(endpoint, key string) *ImageAnalysisClient {
	return &ImageAnalysisClient{
		Endpoint: endpoint,
		Key:      key,
	}
}

func (c *ImageAnalysisClient) AnalyzeImage(imageURL string, visualFeatures []VisualFeature) (*ImageAnalysisResponse, error) {
	features := ""
	for i, feature := range visualFeatures {
		if i > 0 {
			features += ","
		}
		features += string(feature)
	}
	url := fmt.Sprintf("%s/computervision/imageanalysis:analyze?api-version=2024-02-01&features=%s&model-version=latest&language=en", c.Endpoint, features)
	requestBody, err := json.Marshal(ImageAnalysisRequest{URL: imageURL})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": c.Key,
		"Content-Type":              "application/json",
	}

	fmt.Printf("Sending request to URL: %s\n", url)
	fmt.Printf("Request body: %s\n", requestBody)
	resp, err := client.Post(url, bytes.NewBuffer(requestBody), headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to analyze image: %s", body)
	}
	fmt.Printf("Response status code: %d\n", resp.StatusCode)
	fmt.Printf("Response before decoding: %v\n", resp)
	var response ImageAnalysisResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	fmt.Printf("Response: %+v\n", response)

	return &response, nil
}
