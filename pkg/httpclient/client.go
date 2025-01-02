package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ngothientuong/tngo-ai-svcs/pkg/errorcustom"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Post(urlStr string, payload interface{}, headers map[string]string, params interface{}) (*http.Response, error) {
	var body *bytes.Buffer
	if payload != nil {
		body = new(bytes.Buffer)
		json.NewEncoder(body).Encode(payload)
	} else {
		body = new(bytes.Buffer) // Initialize body as an empty buffer
	}

	// Parse the URL and add parameters if params is not nil
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if params != nil {
		query := parsedURL.Query()
		switch p := params.(type) {
		case map[string]string:
			for key, value := range p {
				query.Add(key, value)
			}
		case map[string]interface{}:
			for key, value := range p {
				switch v := value.(type) {
				case string:
					query.Add(key, v)
				case *string:
					if v != nil {
						query.Add(key, *v)
					}
				case int, int8, int16, int32, int64:
					query.Add(key, fmt.Sprintf("%d", v))
				case uint, uint8, uint16, uint32, uint64:
					query.Add(key, fmt.Sprintf("%d", v))
				case float32, float64:
					query.Add(key, fmt.Sprintf("%f", v))
				case bool:
					query.Add(key, fmt.Sprintf("%t", v))
				case fmt.Stringer:
					query.Add(key, v.String())
				case []string:
					query.Add(key, fmt.Sprintf("[%s]", strings.Join(v, ",")))
				default:
					return nil, fmt.Errorf("unsupported param type for key %s: %T", key, value)
				}
			}
		default:
			return nil, fmt.Errorf("unsupported params type: %T", params)
		}
		parsedURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequest("POST", parsedURL.String(), body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	fmt.Println("Request URL: ", parsedURL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errorName, exists := errorcustom.ErrorNames[resp.StatusCode]
		if !exists {
			errorName = "UnknownError"
		}
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, &errorcustom.CustomVisionError{
			StatusCode: resp.StatusCode,
			ErrorName:  errorName,
			Message:    string(responseBody),
		}
	}

	return resp, nil
}

func (c *Client) Get(urlStr string, headers map[string]string, params interface{}) (*http.Response, error) {
	// Parse the URL and add parameters if params is not nil
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if params != nil {
		query := parsedURL.Query()
		switch p := params.(type) {
		case map[string]string:
			for key, value := range p {
				query.Add(key, value)
			}
		case map[string]interface{}:
			for key, value := range p {
				switch v := value.(type) {
				case string:
					query.Add(key, v)
				case *string:
					if v != nil {
						query.Add(key, *v)
					}
				case int, int8, int16, int32, int64:
					query.Add(key, fmt.Sprintf("%d", v))
				case uint, uint8, uint16, uint32, uint64:
					query.Add(key, fmt.Sprintf("%d", v))
				case float32, float64:
					query.Add(key, fmt.Sprintf("%f", v))
				case bool:
					query.Add(key, fmt.Sprintf("%t", v))
				case fmt.Stringer:
					query.Add(key, v.String())
				case []string:
					query.Add(key, fmt.Sprintf("[%s]", strings.Join(v, ",")))
				default:
					return nil, fmt.Errorf("unsupported param type for key %s: %T", key, value)
				}
			}
		default:
			return nil, fmt.Errorf("unsupported params type: %T", params)
		}
		parsedURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	fmt.Println("Request URL: ", parsedURL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errorName, exists := errorcustom.ErrorNames[resp.StatusCode]
		if !exists {
			errorName = "UnknownError"
		}
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, &errorcustom.CustomVisionError{
			StatusCode: resp.StatusCode,
			ErrorName:  errorName,
			Message:    string(responseBody),
		}
	}

	return resp, nil
}

func (c *Client) Delete(urlStr string, headers map[string]string, params interface{}) (*http.Response, error) {
	// Parse the URL and add parameters if params is not nil
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if params != nil {
		query := parsedURL.Query()
		switch p := params.(type) {
		case map[string]string:
			for key, value := range p {
				query.Add(key, value)
			}
		case map[string]interface{}:
			for key, value := range p {
				query.Add(key, fmt.Sprintf("%v", value))
			}
		default:
			return nil, fmt.Errorf("unsupported params type: %T", params)
		}
		parsedURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequest("DELETE", parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	fmt.Println("Request URL: ", parsedURL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("failed to make request: status code %d", resp.StatusCode)
	}

	return resp, nil
}

func (c *Client) Put(urlStr string, payload interface{}, headers map[string]string, params interface{}) (*http.Response, error) {
	var body *bytes.Buffer
	if payload != nil {
		body = new(bytes.Buffer)
		json.NewEncoder(body).Encode(payload)
	} else {
		body = new(bytes.Buffer) // Initialize body as an empty buffer
	}

	// Parse the URL and add parameters if params is not nil
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if params != nil {
		query := parsedURL.Query()
		switch p := params.(type) {
		case map[string]string:
			for key, value := range p {
				query.Add(key, value)
			}
		case map[string]interface{}:
			for key, value := range p {
				query.Add(key, fmt.Sprintf("%v", value))
			}
		default:
			return nil, fmt.Errorf("unsupported params type: %T", params)
		}
		parsedURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequest("PUT", parsedURL.String(), body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	fmt.Println("Request URL: ", parsedURL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to make request: status code %d", resp.StatusCode)
	}

	return resp, nil
}
