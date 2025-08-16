package web

import (
	"log"
	"net/http"
	"strings"

	"github.com/henrik392/youtube-voice-go/cmd/web/components"
	"github.com/henrik392/youtube-voice-go/internal/youtube"
)

func ValidateURLHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ValidateURLHandler")
	url := strings.TrimSpace(r.FormValue("url"))

	component := components.URLInput(youtube.ExtractVideoID(url) != "", url)

	log.Println(youtube.ExtractVideoID(url))

	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in ValidateURLHandler: %e", err)
	}
}
