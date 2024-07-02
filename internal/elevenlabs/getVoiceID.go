package elevenlabs

import "fmt"

func (c *Client) GetVoiceID(youtubeID string) (string, error) {
	if youtubeID == "" {
		return "", fmt.Errorf("youtubeID is empty")
	}
	
	return c.cloneVoice(youtubeID)

	// return "4srV5pKnTwmwqQLucA8p", nil
}
