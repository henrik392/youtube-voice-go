package elevenlabs

import (
	"encoding/json"
	"fmt"
)

const MAX_VOICES int = 10

func (c *Client) GetVoiceID(youtubeID string) (string, error) {
	if youtubeID == "" {
		return "", fmt.Errorf("youtubeID is empty")
	}

	voiceID, err := c.getSavedVoiceID(youtubeID)
	if err != nil {
		return "", fmt.Errorf("failed to get saved voice ID: %v", err)
	}

	if voiceID != "" {
		return voiceID, nil
	}

	err = c.removeVoiceIfMaxReached()

	return c.cloneVoice(youtubeID)

	// return "4srV5pKnTwmwqQLucA8p", nil
}

type VoicesResponse struct {
	Voices []struct {
		VoiceID  string `json:"voice_id"`
		Name     string `json:"name"`
		Category string `json:"category"`
	} `json:"voices"`
}

type Voice struct {
	VoiceID string
	Name    string
}

func (c *Client) getVoices() ([]Voice, error) {
	endpoint := "voices"
	body, err := c.getRequest(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get voices: %v", err)
	}

	var response VoicesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	voices := make([]Voice, 0, len(response.Voices))
	for _, voice := range response.Voices {
		if voice.Category == "cloned" {
			voices = append(voices, Voice{
				VoiceID: voice.VoiceID,
				Name:    voice.Name,
			})
		}
	}

	return voices, nil
}

func (c *Client) getSavedVoiceID(youtubeID string) (string, error) {
	voices, err := c.getVoices()
	if err != nil {
		return "", fmt.Errorf("failed to get voices: %v", err)
	}

	for _, voice := range voices {
		if voice.Name == youtubeID {
			return voice.VoiceID, nil
		}
	}

	return "", nil
}

func (c *Client) removeVoiceIfMaxReached() error {
	voices, err := c.getVoices()

	if err != nil {
		return fmt.Errorf("failed to get voices: %v", err)
	}

	if len(voices) >= MAX_VOICES {
		voiceID := voices[0].VoiceID
		err = c.removeVoice(voiceID)
		if err != nil {
			return fmt.Errorf("failed to remove voice: %v", err)
		}
	}

	return nil
}

func (c *Client) removeVoice(voiceID string) error {
	endpoint := fmt.Sprintf("voices/%s", voiceID)
	_, err := c.deleteRequest(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete voice: %v", err)
	}

	fmt.Printf("Voice %s removed\n", voiceID)
	return nil
}
