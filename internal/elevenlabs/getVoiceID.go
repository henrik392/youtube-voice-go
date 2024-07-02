package elevenlabs

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetVoiceID(youtubeID string) (string, error) {
	if youtubeID == "" {
		return "", fmt.Errorf("youtubeID is empty")
	}

	voiceID, err := c.getSavedVoiceID(youtubeID)
	if err != nil {
		return "", fmt.Errorf("failed to get saved voice ID: %v", err)
	}

	fmt.Println("Saved Voice ID:", voiceID)

	if voiceID != "" {
		return voiceID, nil
	}

	return c.cloneVoice(youtubeID)

	// return "4srV5pKnTwmwqQLucA8p", nil
}

type VoicesResponse struct {
	Voices []struct {
		VoiceID string `json:"voice_id"`
		Name    string `json:"name"`
	} `json:"voices"`
}

func (c *Client) getSavedVoiceID(youtubeID string) (string, error) {
	endpoint := "voices"
	body, err := c.getRequest(endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to get voices: %v", err)
	}

	// Unmarshal the JSON response into the VoicesResponse struct
	var response VoicesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Loop through the voices and return the voice ID if the name matches the youtube ID
	for _, voice := range response.Voices {
		fmt.Println("Voice Name:", voice.Name, "Youtube ID:", youtubeID)
		if voice.Name == youtubeID {
			fmt.Println("Found Voice ID:", voice.VoiceID)
			return voice.VoiceID, nil
		}
	}

	return "", nil
}
