package elevenlabs

import (
	"encoding/json"
	"fmt"
	"os"
)

func (c *Client) TextToSpeech(voiceID, text string) ([]byte, error) {
	endpoint := fmt.Sprintf("text-to-speech/%s", voiceID)

	payload := map[string]interface{}{
		"text":     text,
		"model_id": "eleven_monolingual_v1",
		"voice_settings": map[string]float64{
			"stability":        0.5,
			"similarity_boost": 0.5,
		},
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %w", err)
	}

	return c.postJSON(endpoint, jsonPayload)
}

func (c *Client) SaveAudioFile(audio []byte, filename string) error {
	return os.WriteFile(filename, audio, 0644)
}
