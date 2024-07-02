package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/henrik392/youtube-voice-go/internal/elevenlabs"
	"github.com/henrik392/youtube-voice-go/internal/youtube"
)

func GenerateVoiceHandler(w http.ResponseWriter, r *http.Request) {
	// expects a form with a 'text' field and a 'url' field
	// 'text' is the text to be spoken
	// 'url' is the youtube video url

	url := r.FormValue("url")
	text := r.FormValue("text")
	youtubeID := youtube.GetYoutubeId(url)

	fmt.Println("URL:", url)
	fmt.Println("Text:", text)
	fmt.Println("Youtube ID:", youtubeID)

	// Download the audio from the youtube video
	err, shouldReturn := downloadAndServeAudio(youtubeID, w, r)
	if shouldReturn {
		return
	}

	elClient := elevenlabs.NewClient(os.Getenv("ELEVENLABS_API_KEY"))

	voiceID, err := elClient.GetVoiceID(youtubeID)
	if err != nil {
		log.Printf("Failed to get voice ID: %v", err)
		http.Error(w, "Failed to get voice ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	audioData, err := elClient.TextToSpeech(voiceID, text)
	fmt.Println("Voice ID:", voiceID)
	if err != nil {
		log.Printf("Failed to generate voice: %v", err)
		http.Error(w, "Failed to generate voice: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Generated speech!")

	// Save the generated speech
	uuid := uuid.New()
	speechFilePath := filepath.Join("./downloads", fmt.Sprintf("%s_speech_%s.mp3", youtubeID, uuid.String()))
	err = elClient.SaveAudioFile(audioData, speechFilePath)
	if err != nil {
		log.Printf("Failed to save generated speech: %v", err)
		http.Error(w, "Failed to save generated speech", http.StatusInternalServerError)
		return
	}

	fmt.Println("Saved speech!")

	downloadButton := DownloadButton(url)

	err = downloadButton.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func downloadAndServeAudio(youtubeID string, w http.ResponseWriter, r *http.Request) (error, bool) {
	ytProcessor := youtube.NewProcessor("./downloads")
	audioFile, err := ytProcessor.DownloadAudio(youtubeID)

	if err != nil {
		http.Error(w, "Failed to process Youtube video: "+err.Error(), http.StatusBadRequest)
		return nil, true
	}

	fmt.Println("Downloaded audio file:", audioFile)

	return err, false
}
