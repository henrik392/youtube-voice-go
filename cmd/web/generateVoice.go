package web

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/henrik392/youtube-voice-go/cmd/web/components"
	"github.com/henrik392/youtube-voice-go/internal/elevenlabs"
	"github.com/henrik392/youtube-voice-go/internal/youtube"
)

func GenerateVoiceHandler(w http.ResponseWriter, r *http.Request) {
	// expects a form with a 'text' field and a 'url' field
	// 'text' is the text to be spoken
	// 'url' is the youtube video url

	videoURL := r.FormValue("url")
	text := r.FormValue("text")
	videoID := youtube.ExtractVideoID(videoURL)

	fmt.Println("URL:", videoURL)
	fmt.Println("Text:", text)
	fmt.Println("Youtube ID:", videoID)

	ytProcessor := youtube.NewProcessor("./downloads")
	audioFile, err := ytProcessor.DownloadAudio(videoURL, videoID)

	if err != nil {
		serveError(w, r, "Failed to download audio: "+err.Error())
		return
	}

	fmt.Println("Downloaded audio file:", audioFile)

	// clone the voice

	elClient := elevenlabs.NewClient(os.Getenv("ELEVENLABS_API_KEY"))

	voiceID, err := elClient.GetVoiceID(videoID)
	if err != nil {
		serveError(w, r, "Failed to clone voice: "+err.Error())
		return
	}

	audioData, err := elClient.TextToSpeech(voiceID, text)
	fmt.Println("Voice ID:", voiceID)
	if err != nil {
		serveError(w, r, "Failed to generate speech: "+err.Error())
		return
	}

	fmt.Println("Generated speech!")

	// Save the generated speech
	uuid := uuid.New()
	speechFilePath := filepath.Join("./downloads", fmt.Sprintf("%s_speech_%s.mp3", videoID, uuid.String()))
	err = elClient.SaveAudioFile(audioData, speechFilePath)
	if err != nil {
		serveError(w, r, "Failed to save speech to file: "+err.Error())
		return
	}

	fmt.Println("Saved speech to file:", speechFilePath)
	fmt.Println("Serving audio...")

	audioURL := fmt.Sprintf("/serve-audio?path=%s", url.QueryEscape(speechFilePath))

	audioPlayer := components.AudioPlayer(audioURL, "")
	err = audioPlayer.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func serveError(w http.ResponseWriter, r *http.Request, errorMessage string) {
	log.Printf("Error: %v", errorMessage)
	audioPlayer := components.AudioPlayer("", errorMessage)

	if err := audioPlayer.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
