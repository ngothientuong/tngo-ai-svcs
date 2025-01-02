package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
				query.Add(key, fmt.Sprintf("%v", value))
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
				query.Add(key, fmt.Sprintf("%v", value))
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
