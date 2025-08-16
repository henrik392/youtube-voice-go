package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/henrik392/youtube-voice-go/cmd/web/components"
	"github.com/henrik392/youtube-voice-go/internal/zonos"
	"github.com/henrik392/youtube-voice-go/internal/youtube"
)

func GenerateVoiceEnhancedHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	audioMode := r.FormValue("audio-mode")

	log.Printf("Text: %s", text)
	log.Printf("Audio mode: %s", audioMode)

	if text == "" {
		serveError(w, r, "Please provide text to generate speech")
		return
	}

	if len(text) > 500 {
		serveError(w, r, "Text must be 500 characters or less")
		return
	}

	var audioFile string
	var err error

	// Handle different input modes
	switch audioMode {
	case "url":
		audioFile, err = handleURLInput(r)
	case "file":
		audioFile, err = handleFileInput(r)
	case "microphone":
		audioFile, err = handleMicrophoneInput(r)
	default:
		serveError(w, r, "Invalid audio input mode")
		return
	}

	if err != nil {
		serveError(w, r, fmt.Sprintf("Failed to process audio input: %v", err))
		return
	}

	log.Printf("Processing audio file: %s", audioFile)

	// Create Zonos client
	log.Printf("Creating Zonos client...")
	diaClient := zonos.NewClient(os.Getenv("FAL_KEY"))

	// Generate speech using Zonos voice cloning
	log.Printf("Starting voice cloning with Zonos...")
	audioData, err := diaClient.VoiceClone(text, audioFile)
	if err != nil {
		log.Printf("Failed to generate speech: %v", err)
		serveError(w, r, "Failed to generate speech: "+err.Error())
		return
	}

	log.Printf("Generated speech successfully! Audio data size: %d bytes", len(audioData))

	// Generate unique filename for the output
	uuid := uuid.New()
	speechFilePath := filepath.Join("./downloads", fmt.Sprintf("speech_%s.wav", uuid.String()))
	err = diaClient.SaveAudioFile(audioData, speechFilePath)
	if err != nil {
		serveError(w, r, "Failed to save speech to file: "+err.Error())
		return
	}

	log.Printf("Saved speech to file: %s", speechFilePath)

	audioURL := fmt.Sprintf("/serve-audio?path=%s", url.QueryEscape(speechFilePath))
	audioPlayer := components.AudioPlayer(audioURL, "")

	err = audioPlayer.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleURLInput(r *http.Request) (string, error) {
	videoURL := r.FormValue("url")
	if videoURL == "" {
		return "", fmt.Errorf("please provide a video URL")
	}

	videoID := youtube.ExtractVideoID(videoURL)
	if videoID == "" {
		return "", fmt.Errorf("invalid video URL")
	}

	log.Printf("Processing video URL: %s (ID: %s)", videoURL, videoID)

	ytProcessor := youtube.NewProcessor("./downloads")
	audioFile, err := ytProcessor.DownloadAudio(videoURL, videoID)
	if err != nil {
		return "", fmt.Errorf("failed to download audio: %v", err)
	}

	return audioFile, nil
}

func handleFileInput(r *http.Request) (string, error) {
	// The file should have been uploaded already via the upload endpoint
	// Look for fileId in form data
	fileId := r.FormValue("file-id")
	if fileId == "" {
		// Try to handle direct file upload in this request
		return handleDirectFileUpload(r)
	}

	// Construct file path from fileId
	// We need to find the file in downloads directory
	files, err := filepath.Glob("downloads/" + fileId + ".*")
	if err != nil || len(files) == 0 {
		return "", fmt.Errorf("uploaded file not found")
	}

	return files[0], nil
}

func handleDirectFileUpload(r *http.Request) (string, error) {
	err := r.ParseMultipartForm(50 << 20) // 50MB
	if err != nil {
		return "", fmt.Errorf("file too large")
	}

	file, header, err := r.FormFile("audio-file")
	if err != nil {
		return "", fmt.Errorf("no file provided")
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	validTypes := []string{
		"audio/mpeg", "audio/mp3", "audio/wav", "audio/wave",
		"audio/mp4", "audio/m4a", "audio/ogg", "audio/flac",
	}
	
	isValid := false
	for _, validType := range validTypes {
		if strings.Contains(contentType, validType) || strings.Contains(strings.ToLower(header.Filename), strings.TrimPrefix(validType, "audio/")) {
			isValid = true
			break
		}
	}

	if !isValid {
		return "", fmt.Errorf("invalid file type. Please upload MP3, WAV, M4A, OGG, or FLAC files")
	}

	// Save file
	fileID := uuid.New().String()
	ext := filepath.Ext(header.Filename)
	filename := fileID + ext
	filePath := filepath.Join("downloads", filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to save file")
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("failed to save file")
	}

	return filePath, nil
}

func handleMicrophoneInput(r *http.Request) (string, error) {
	// The recording should have been saved already via the save-recording endpoint
	// Look for recordingId in form data
	recordingId := r.FormValue("recording-id")
	if recordingId == "" {
		// Try to handle direct recording upload in this request
		return handleDirectRecordingUpload(r)
	}

	// Construct file path from recordingId
	recordingFile := "downloads/" + recordingId + ".webm"
	if _, err := os.Stat(recordingFile); os.IsNotExist(err) {
		return "", fmt.Errorf("recorded audio not found")
	}

	return recordingFile, nil
}

func handleDirectRecordingUpload(r *http.Request) (string, error) {
	// Read the recording data from request body or form
	var recordingData []byte
	var err error

	// Check if it's multipart form data
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		err = r.ParseMultipartForm(50 << 20)
		if err != nil {
			return "", fmt.Errorf("failed to parse form data")
		}

		file, _, err := r.FormFile("recording")
		if err != nil {
			return "", fmt.Errorf("no recording data provided")
		}
		defer file.Close()

		recordingData, err = io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("failed to read recording data")
		}
	} else {
		// Handle JSON payload
		var payload struct {
			RecordingData string `json:"recordingData"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			return "", fmt.Errorf("invalid recording data")
		}

		if payload.RecordingData == "" {
			return "", fmt.Errorf("no recording data provided")
		}

		// The recording data might be base64 encoded blob data
		recordingData = []byte(payload.RecordingData)
	}

	if len(recordingData) == 0 {
		return "", fmt.Errorf("no recording data provided")
	}

	// Save recording to file
	recordingID := uuid.New().String()
	filename := recordingID + ".webm"
	filePath := filepath.Join("downloads", filename)

	err = os.WriteFile(filePath, recordingData, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save recording")
	}

	return filePath, nil
}