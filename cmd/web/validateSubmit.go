package web

import (
	"log"
	"net/http"

	"github.com/henrik392/youtube-voice-go/cmd/web/components"
)

func ValidateSubmitHandler(w http.ResponseWriter, r *http.Request) {
	component := components.SubmitButton()

	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in ValidateURLHandler: %e", err)
	}
}
