package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/henrik392/youtube-voice-go/internal/diatts"
	"github.com/henrik392/youtube-voice-go/internal/youtube"
)

type ProcessVideoResponse struct {
	Success    bool   `json:"success"`
	VideoID    string `json:"video_id"`
	AudioURL   string `json:"audio_url,omitempty"`
	RefText    string `json:"ref_text,omitempty"`
	Error      string `json:"error,omitempty"`
}

func ProcessVideoHandler(w http.ResponseWriter, r *http.Request) {
	videoURL := r.FormValue("url")
	videoID := youtube.ExtractVideoID(videoURL)

	w.Header().Set("Content-Type", "application/json")

	if videoID == "" {
		response := ProcessVideoResponse{
			Success: false,
			Error:   "Invalid YouTube URL",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("Processing video: %s (ID: %s)", videoURL, videoID)

	// Download audio
	ytProcessor := youtube.NewProcessor("./downloads")
	audioFile, err := ytProcessor.DownloadAudio(videoURL, videoID)
	if err != nil {
		log.Printf("Failed to download audio: %v", err)
		response := ProcessVideoResponse{
			Success: false,
			VideoID: videoID,
			Error:   "Failed to download audio: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("Downloaded audio file: %s", audioFile)

	// Create Dia TTS client
	diaClient := diatts.NewClient(os.Getenv("FAL_KEY"))

	// Crop, upload audio and extract reference text
	log.Printf("Cropping and uploading audio...")
	croppedFilePath, err := diaClient.CropAndCompressAudio(audioFile, 15)
	if err != nil {
		log.Printf("Failed to crop audio: %v", err)
		response := ProcessVideoResponse{
			Success: false,
			VideoID: videoID,
			Error:   "Failed to crop audio: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer os.Remove(croppedFilePath) // Clean up temp file

	// Upload to S3
	audioURL, err := diaClient.UploadToS3(croppedFilePath)
	if err != nil {
		log.Printf("Failed to upload audio: %v", err)
		response := ProcessVideoResponse{
			Success: false,
			VideoID: videoID,
			Error:   "Failed to upload audio: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Extract reference text using speech-to-text
	log.Printf("Extracting reference text...")
	refText, err := diaClient.ExtractReferenceTextFromURL(audioURL)
	if err != nil {
		log.Printf("Failed to extract reference text: %v", err)
		response := ProcessVideoResponse{
			Success: false,
			VideoID: videoID,
			Error:   "Failed to extract reference text: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("Video processing complete for %s", videoID)

	response := ProcessVideoResponse{
		Success:  true,
		VideoID:  videoID,
		AudioURL: audioURL,
		RefText:  refText,
	}
	json.NewEncoder(w).Encode(response)
}