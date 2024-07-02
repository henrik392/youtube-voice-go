package elevenlabs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// CloneVoice clones a voice by uploading an audio to the elevenlabs API, a voice id is returned.
// It takes a YouTube ID as input and returns the response from the server as a string (voice id).
// If the YouTube ID is empty or if there is an error during the process, an error is returned.
func (c *Client) cloneVoice(youtubeID string) (string, error) {
	if youtubeID == "" {
		return "", fmt.Errorf("youtubeID is empty")
	}

	// get audioFilePath of audio file in downloads folder
	audioFilePath, err := filepath.Abs(filepath.Join("./downloads", fmt.Sprintf("%s.mp3", youtubeID)))
	if err != nil {
		return "", fmt.Errorf("failed to get path or audio does not exist: %v", err)
	}

	// Perpare the multipart form data
	var formData bytes.Buffer
	writer := multipart.NewWriter(&formData)

	// Add `name` part
	err = writer.WriteField("name", youtubeID)
	if err != nil {
		return "", fmt.Errorf("failed to write 'name' field: %v", err)
	}

	// Add 'files' part
	file, err := os.Open(audioFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %v", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("files", filepath.Base(audioFilePath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file to part: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Ready to add voice

	endpoint := "voices/add"
	contentType := writer.FormDataContentType()

	response, err := c.postFormData(endpoint, &formData, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to post form data: %v", err)
	}

	// the resopnse is a JSON object with a 'voice_id' field like {"voice_id":"Yc8gJFBQEo23EEbwIFtd"}, return the voice id as a string
	var voiceResponse struct {
		VoiceID string `json:"voice_id"`
	}
	err = json.Unmarshal(response, &voiceResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return voiceResponse.VoiceID, nil
}
