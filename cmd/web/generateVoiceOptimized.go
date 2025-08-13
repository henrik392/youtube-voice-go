package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/henrik392/youtube-voice-go/cmd/web/components"
	"github.com/henrik392/youtube-voice-go/internal/diatts"
)

func GenerateVoiceOptimizedHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	audioURL := r.FormValue("audio_url")
	refText := r.FormValue("ref_text")
	videoID := r.FormValue("video_id")

	log.Printf("Generating voice with pre-processed audio")
	log.Printf("Text: %s", text)
	log.Printf("Audio URL: %s", audioURL)
	log.Printf("Reference Text: %s", refText)
	log.Printf("Video ID: %s", videoID)

	if text == "" || audioURL == "" || refText == "" {
		serveError(w, r, "Missing required parameters: text, audio_url, or ref_text")
		return
	}

	// Create Dia TTS client
	diaClient := diatts.NewClient(os.Getenv("FAL_KEY"))

	// Format the target text for Dia TTS (needs [S1] format)
	formattedText := fmt.Sprintf("[S1] %s", text)
	log.Printf("Formatted text: %s", formattedText)

	// Generate speech using pre-processed audio URL and reference text
	log.Printf("Starting voice cloning with Dia TTS using pre-processed data...")
	audioData, err := diaClient.VoiceCloneWithURL(formattedText, audioURL, refText)
	if err != nil {
		log.Printf("Failed to generate speech: %v", err)
		serveError(w, r, "Failed to generate speech: "+err.Error())
		return
	}

	log.Printf("Generated speech successfully! Audio data size: %d bytes", len(audioData))

	// Save the generated speech
	uuid := uuid.New()
	speechFilePath := filepath.Join("./downloads", fmt.Sprintf("%s_speech_%s.wav", videoID, uuid.String()))
	err = diaClient.SaveAudioFile(audioData, speechFilePath)
	if err != nil {
		serveError(w, r, "Failed to save speech to file: "+err.Error())
		return
	}

	log.Printf("Saved speech to file: %s", speechFilePath)

	audioURL = fmt.Sprintf("/serve-audio?path=%s", speechFilePath)
	audioPlayer := components.AudioPlayer(audioURL, "")
	err = audioPlayer.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}