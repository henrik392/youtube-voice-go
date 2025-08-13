package diatts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
	APIKey  string
	BaseURL string
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://fal.run/fal-ai/dia-tts/voice-clone",
	}
}

type Request struct {
	Text         string `json:"text"`
	RefAudioURL  string `json:"ref_audio_url"`
	RefText      string `json:"ref_text"`
}

type Response struct {
	Audio struct {
		URL string `json:"url"`
	} `json:"audio"`
}

func (c *Client) VoiceClone(text, refAudioURL, refText string) ([]byte, error) {
	payload := Request{
		Text:        text,
		RefAudioURL: refAudioURL,
		RefText:     refText,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Key %s", c.APIKey))

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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d\nBody: %s", resp.StatusCode, body)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return c.downloadAudio(response.Audio.URL)
}

func (c *Client) downloadAudio(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading audio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download audio: status code %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) SaveAudioFile(audio []byte, filename string) error {
	return os.WriteFile(filename, audio, 0644)
}

func (c *Client) GenerateRefAudioURL(audioFilePath, baseURL string) string {
	return fmt.Sprintf("%s/serve-audio?path=%s", baseURL, audioFilePath)
}

func (c *Client) ExtractReferenceText(audioFilePath string) (string, error) {
	// TODO: Implement actual speech-to-text extraction
	// For now, return a placeholder that matches the Dia TTS expected format
	return "[S1] This is sample reference text extracted from the audio.", nil
}