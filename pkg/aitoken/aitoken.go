package aitoken

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ngothientuong/tngo-ai-svcs/pkg/httpclient"
)

var cachedToken string
var tokenExpiry time.Time

func GetToken(endpoint, key string) (string, error) {
	if cachedToken != "" && time.Now().Before(tokenExpiry) {
		return cachedToken, nil
	}

	url := fmt.Sprintf("%s/sts/v1.0/issueToken", endpoint)
	client := httpclient.NewClient()
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": key,
		"Content-Type":              "application/x-www-form-urlencoded",
	}

	resp, err := client.Post(url, strings.NewReader(""), headers, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get token: %s", body)
	}

	token, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token: %v", err)
	}

	cachedToken = string(token)
	tokenExpiry = time.Now().Add(9 * time.Minute) // Token is valid for 10 minutes

	return cachedToken, nil
}
