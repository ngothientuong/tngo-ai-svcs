package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

func (c *DocumentTranslationClient) StartBatchTranslation(sourceContainer, targetContainer, targetLanguage string) ([]byte, error) {
	urlStr := fmt.Sprintf("%s/translator/document/batches?api-version=%s", c.Endpoint, c.APIVersion)

	body := map[string]interface{}{
		"inputs": []map[string]interface{}{
			{
				"source": map[string]string{
					"sourceUrl": sourceContainer,
				},
				"targets": []map[string]string{
					{
						"targetUrl": targetContainer,
						"language":  targetLanguage,
					},
				},
			},
		},
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

	return bodyBytes, nil
}
