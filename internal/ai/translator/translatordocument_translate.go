package translator

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

type DocumentTranslationClient struct {
	Endpoint   string
	Key        string
	Region     string
	APIVersion string
}

func NewDocumentTranslationClient(endpoint, key, region, apiVersion string) *DocumentTranslationClient {
	return &DocumentTranslationClient{
		Endpoint:   endpoint,
		Key:        key,
		Region:     region,
		APIVersion: apiVersion,
	}
}

/*
Reference: https://learn.microsoft.com/en-us/azure/ai-services/translator/document-translation/overview
🛠 Supported Formats for Synchronous Translation
✅ Word (.docx)
✅ PowerPoint (.pptx)
✅ Excel (.xlsx)
✅ Text (.txt, .csv, .tsv, .html, .xml, .xlf, .xliff)
*/
func (c *DocumentTranslationClient) TranslateDocument(documentPath, targetLanguage, sourceLanguage, glossaryPath, outputFilePath string, allowFallback bool, category string) error {
	urlStr := fmt.Sprintf("%s/translator/document:translate?targetLanguage=%s&api-version=%s", c.Endpoint, targetLanguage, c.APIVersion)
	if sourceLanguage != "" {
		urlStr += fmt.Sprintf("&sourceLanguage=%s", sourceLanguage)
	}
	if category != "" {
		urlStr += fmt.Sprintf("&category=%s", category)
	}
	urlStr += fmt.Sprintf("&allowFallback=%t", allowFallback)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add document file
	documentFile, err := os.Open(documentPath)
	if err != nil {
		return fmt.Errorf("failed to open document file: %v", err)
	}
	defer documentFile.Close()

	part, err := writer.CreatePart(map[string][]string{
		"Content-Disposition": {fmt.Sprintf(`form-data; name="document"; filename="%s"`, filepath.Base(documentPath))},
		"Content-Type":        {http.DetectContentType([]byte(filepath.Ext(documentPath)))},
	})
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, documentFile)
	if err != nil {
		return fmt.Errorf("failed to copy document file: %v", err)
	}

	// Add glossary file if provided
	if glossaryPath != "" {
		glossaryFile, err := os.Open(glossaryPath)
		if err != nil {
			return fmt.Errorf("failed to open glossary file: %v", err)
		}
		defer glossaryFile.Close()

		part, err := writer.CreatePart(map[string][]string{
			"Content-Disposition": {fmt.Sprintf(`form-data; name="glossary"; filename="%s"`, filepath.Base(glossaryPath))},
			"Content-Type":        {http.DetectContentType([]byte(filepath.Ext(glossaryPath)))},
		})
		if err != nil {
			return fmt.Errorf("failed to create form file: %v", err)
		}
		_, err = io.Copy(part, glossaryFile)
		if err != nil {
			return fmt.Errorf("failed to copy glossary file: %v", err)
		}
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": c.Key,
		"Content-Type":              writer.FormDataContentType(),
	}

	resp, err := client.Post(urlStr, body, headers, nil)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed: %s", string(bodyBytes))
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to output file: %v", err)
	}

	return nil
}
