package web

import (
	"log"
	"net/http"
	"os"

	"github.com/henrik392/youtube-voice-go/cmd/web/components"
	"github.com/henrik392/youtube-voice-go/internal/zonos"
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

	w.Header().Set("Content-Type", "text/html")

	if videoID == "" {
		component := components.ProcessingError("Invalid YouTube URL")
		component.Render(r.Context(), w)
		return
	}

	log.Printf("Processing video: %s (ID: %s)", videoURL, videoID)

	// Download audio
	ytProcessor := youtube.NewProcessor("./downloads")
	audioFile, err := ytProcessor.DownloadAudio(videoURL, videoID)
	if err != nil {
		log.Printf("Failed to download audio: %v", err)
		component := components.ProcessingError("Failed to download audio: " + err.Error())
		component.Render(r.Context(), w)
		return
	}

	log.Printf("Downloaded audio file: %s", audioFile)

	// Create Dia TTS client
	diaClient := zonos.NewClient(os.Getenv("FAL_KEY"))

	// Crop, upload audio and extract reference text
	log.Printf("Cropping and uploading audio...")
	croppedFilePath, err := diaClient.CropAndCompressAudio(audioFile, 30)
	if err != nil {
		log.Printf("Failed to crop audio: %v", err)
		component := components.ProcessingError("Failed to crop audio: " + err.Error())
		component.Render(r.Context(), w)
		return
	}
	defer os.Remove(croppedFilePath) // Clean up temp file

	// Upload to S3
	audioURL, err := diaClient.UploadToS3(croppedFilePath)
	if err != nil {
		log.Printf("Failed to upload audio: %v", err)
		component := components.ProcessingError("Failed to upload audio: " + err.Error())
		component.Render(r.Context(), w)
		return
	}

	log.Printf("Video processing complete for %s", videoID)

	component := components.ProcessingComplete(videoID, audioURL)
	component.Render(r.Context(), w)
}