package translator

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

func (c *DocumentTranslationClient) GetTranslationStatus(batchID string) ([]byte, error) {
	urlStr := fmt.Sprintf("%s/translator/document/batches/%s?api-version=%s", c.Endpoint, batchID, c.APIVersion)

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key":    c.Key,
		"Ocp-Apim-Subscription-Region": c.Region,
		"Content-Type":                 "application/json; charset=UTF-8",
	}

	resp, err := client.Get(urlStr, headers, nil)
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

	return bodyBytes, nil
}

func (c *DocumentTranslationClient) GetDocumentsStatus(batchID string) ([]byte, error) {
	urlStr := fmt.Sprintf("%s/translator/document/batches/%s/documents?api-version=%s", c.Endpoint, batchID, c.APIVersion)

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key":    c.Key,
		"Ocp-Apim-Subscription-Region": c.Region,
		"Content-Type":                 "application/json; charset=UTF-8",
	}

	resp, err := client.Get(urlStr, headers, nil)
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

	return bodyBytes, nil
}

func (c *DocumentTranslationClient) GetDocumentStatus(batchID, documentID string) ([]byte, error) {
	urlStr := fmt.Sprintf("%s/translator/document/batches/%s/documents/%s?api-version=%s", c.Endpoint, batchID, documentID, c.APIVersion)

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key":    c.Key,
		"Ocp-Apim-Subscription-Region": c.Region,
		"Content-Type":                 "application/json; charset=UTF-8",
	}

	resp, err := client.Get(urlStr, headers, nil)
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

	return bodyBytes, nil
}
