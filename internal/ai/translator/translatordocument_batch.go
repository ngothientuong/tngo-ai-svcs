package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

func (c *DocumentTranslationClient) StartBatchTranslation(sourceContainer, targetContainer, targetLanguage, prefix, suffix, sourceLanguage, storageType string, glossaries []map[string]string) ([]byte, error) {
	urlStr := fmt.Sprintf("%s/translator/document/batches?api-version=%s", c.Endpoint, c.APIVersion)

	source := map[string]interface{}{
		"sourceUrl": sourceContainer,
	}
	if prefix != "" || suffix != "" {
		source["filter"] = map[string]string{
			"prefix": prefix,
			"suffix": suffix,
		}
	}
	if sourceLanguage != "" {
		source["language"] = sourceLanguage
	}
	source["storageSource"] = "AzureBlob"

	target := map[string]interface{}{
		"targetUrl":     targetContainer,
		"category":      "general",
		"language":      targetLanguage,
		"storageSource": "AzureBlob",
	}
	if len(glossaries) > 0 {
		target["glossaries"] = glossaries
	}

	body := map[string]interface{}{
		"inputs": []map[string]interface{}{
			{
				"source": source,
				"targets": []map[string]interface{}{
					target,
				},
				"storageType": storageType,
			},
		},
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	fmt.Printf("Request Body: %s\n", string(bodyBytes))

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key":    c.Key,
		"Ocp-Apim-Subscription-Region": c.Region,
		"Content-Type":                 "application/json",
	}

	resp, err := client.Post(urlStr, bytes.NewBuffer(bodyBytes), headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed: %s", string(bodyBytes))
	}

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Printf("Response Body: %s\n", string(bodyBytes))

	return bodyBytes, nil
}
