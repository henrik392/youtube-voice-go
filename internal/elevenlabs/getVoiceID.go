package elevenlabs

import "fmt"

func GetVoiceId(youtubeID string) (string, error) {
	if youtubeID == "" {
		return "", fmt.Errorf("youtubeID is empty")
	}

	return "4srV5pKnTwmwqQLucA8p", nil
}
