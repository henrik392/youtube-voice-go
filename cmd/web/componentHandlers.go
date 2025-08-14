package web

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/henrik392/youtube-voice-go/cmd/web/components"
	"github.com/google/uuid"
)

// FileUploadHandler serves the file upload component
func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	component := components.FileUpload()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error rendering FileUpload component: %v", err)
	}
}

// MicrophoneHandler serves the microphone recording component
func MicrophoneHandler(w http.ResponseWriter, r *http.Request) {
	component := components.MicrophoneRecorder()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error rendering MicrophoneRecorder component: %v", err)
	}
}

// UploadAudioHandler handles audio file uploads
func UploadAudioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form with max memory of 50MB
	err := r.ParseMultipartForm(50 << 20) // 50MB
	if err != nil {
		http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("audio-file")
	if err != nil {
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
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
		http.Error(w, "Invalid file type. Please upload MP3, WAV, M4A, OGG, or FLAC files", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	fileID := uuid.New().String()
	ext := filepath.Ext(header.Filename)
	filename := fileID + ext

	// Ensure downloads directory exists
	downloadsDir := "downloads"
	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Failed to create downloads directory: %v", err)
		return
	}

	// Create file
	filePath := filepath.Join(downloadsDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Failed to create file: %v", err)
		return
	}
	defer dst.Close()

	// Copy file data
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Failed to save file: %v", err)
		return
	}

	// Return success response with file ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "fileId": "` + fileID + `", "filename": "` + filename + `"}`))
}

// SaveRecordingHandler handles microphone recordings
func SaveRecordingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the recording data from request body
	recordingData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read recording data", http.StatusBadRequest)
		return
	}

	if len(recordingData) == 0 {
		http.Error(w, "No recording data provided", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	recordingID := uuid.New().String()
	filename := recordingID + ".webm"

	// Ensure downloads directory exists
	downloadsDir := "downloads"
	if err := os.MkdirAll(downloadsDir, 0755); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Failed to create downloads directory: %v", err)
		return
	}

	// Save recording to file
	filePath := filepath.Join(downloadsDir, filename)
	err = os.WriteFile(filePath, recordingData, 0644)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Failed to save recording: %v", err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "recordingId": "` + recordingID + `", "filename": "` + filename + `"}`))
}