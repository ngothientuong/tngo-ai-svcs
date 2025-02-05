package contentsafety

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

type MediaType string

const (
	MediaTypeText  MediaType = "text"
	MediaTypeImage MediaType = "image"
)

type Category string

const (
	CategoryHate     Category = "Hate"
	CategorySelfHarm Category = "SelfHarm"
	CategorySexual   Category = "Sexual"
	CategoryViolence Category = "Violence"
)

type Action string

const (
	ActionAccept Action = "Accept"
	ActionReject Action = "Reject"
)

type DetectionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *DetectionError) Error() string {
	return fmt.Sprintf("DetectionError(code=%s, message=%s)", e.Code, e.Message)
}

type Decision struct {
	SuggestedAction    Action               `json:"suggestedAction"`
	ActionByCategory   map[Category]Action  `json:"actionByCategory"`
	SeverityByCategory map[Category]float64 `json:"severityByCategory"`
}

type ContentSafetyClient struct {
	Endpoint        string
	SubscriptionKey string
	APIVersion      string
}

func NewContentSafetyClient(endpoint, subscriptionKey, apiVersion string) *ContentSafetyClient {
	return &ContentSafetyClient{
		Endpoint:        endpoint,
		SubscriptionKey: subscriptionKey,
		APIVersion:      apiVersion,
	}
}

func (c *ContentSafetyClient) buildURL(mediaType MediaType) string {
	switch mediaType {
	case MediaTypeText:
		return fmt.Sprintf("%s/contentsafety/text:analyze?api-version=%s", c.Endpoint, c.APIVersion)
	case MediaTypeImage:
		return fmt.Sprintf("%s/contentsafety/image:analyze?api-version=%s", c.Endpoint, c.APIVersion)
	default:
		panic(fmt.Sprintf("Invalid Media Type %s", mediaType))
	}
}

func (c *ContentSafetyClient) buildHeaders() map[string]string {
	return map[string]string{
		"Ocp-Apim-Subscription-Key": c.SubscriptionKey,
		"Content-Type":              "application/json",
	}
}

func (c *ContentSafetyClient) buildRequestBody(mediaType MediaType, content string, blocklists []string) map[string]interface{} {
	switch mediaType {
	case MediaTypeText:
		return map[string]interface{}{
			"text":           content,
			"blocklistNames": blocklists,
		}
	case MediaTypeImage:
		return map[string]interface{}{
			"image": map[string]string{
				"content": content,
			},
		}
	default:
		panic(fmt.Sprintf("Invalid Media Type %s", mediaType))
	}
}

func (c *ContentSafetyClient) Detect(mediaType MediaType, content string, blocklists []string) (map[string]interface{}, error) {
	url := c.buildURL(mediaType)
	headers := c.buildHeaders()
	requestBody := c.buildRequestBody(mediaType, content, blocklists)
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := httpclient.NewClient()
	resp, err := client.Post(url, bytes.NewBuffer(payload), headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var detectionError DetectionError
		if err := json.Unmarshal(body, &detectionError); err != nil {
			return nil, fmt.Errorf("failed to unmarshal error response: %v", err)
		}
		return nil, &detectionError
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return result, nil
}

func (c *ContentSafetyClient) GetDetectResultByCategory(category Category, detectResult map[string]interface{}) (map[string]interface{}, error) {
	categoriesAnalysis, ok := detectResult["categoriesAnalysis"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid categoriesAnalysis format")
	}

	for _, res := range categoriesAnalysis {
		resMap, ok := res.(map[string]interface{})
		if !ok {
			continue
		}
		if resMap["category"] == string(category) {
			return resMap, nil
		}
	}
	return nil, fmt.Errorf("category %s not found in detection result", category)
}

func (c *ContentSafetyClient) MakeDecision(detectionResult map[string]interface{}, rejectThresholds map[Category]int) (*Decision, error) {
	actionResult := make(map[Category]Action)
	severityResult := make(map[Category]float64)
	finalAction := ActionAccept

	for category, threshold := range rejectThresholds {
		if threshold != -1 && threshold != 0 && threshold != 2 && threshold != 4 && threshold != 6 {
			return nil, fmt.Errorf("RejectThreshold can only be in (-1, 0, 2, 4, 6)")
		}

		cateDetectRes, err := c.GetDetectResultByCategory(category, detectionResult)
		if err != nil {
			return nil, err
		}

		severity, ok := cateDetectRes["severity"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid severity format for category %s", category)
		}

		action := ActionAccept
		if threshold != -1 && int(severity) >= threshold {
			action = ActionReject
		}
		actionResult[category] = action
		severityResult[category] = severity
		if action == ActionReject {
			finalAction = ActionReject
		}
	}

	if blocklistsMatch, ok := detectionResult["blocklistsMatch"].([]interface{}); ok && len(blocklistsMatch) > 0 {
		finalAction = ActionReject
	}

	return &Decision{
		SuggestedAction:    finalAction,
		ActionByCategory:   actionResult,
		SeverityByCategory: severityResult,
	}, nil
}
