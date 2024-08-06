package elevenlabs

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	APIKey  string
	BaseURL string
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://api.elevenlabs.io/v1",
	}
}

func (c *Client) postJSON(endpoint string, payload []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	return c.doRequest(req)
}

func (c *Client) postFormData(endpoint string, formData *bytes.Buffer, contentType string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
	req, err := http.NewRequest("POST", url, formData)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	return c.doRequest(req)
}

func (c *Client) getRequest(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	return c.doRequest(req)
}

func (c *Client) deleteRequest(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	return c.doRequest(req)
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("xi-api-key", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// ready body

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	return body, nil
}
