package web

import (
	"fmt"
	"net/http"
)

func GenerateVoiceHandler(w http.ResponseWriter, r *http.Request) {
	// expects a form with a 'text' field and a 'url' field
	// 'text' is the text to be spoken
	// 'url' is the youtube video url

	url := r.FormValue("url")
	text := r.FormValue("text")

	fmt.Println("URL:", url)
	fmt.Println("Text:", text)

	downloadButton := DownloadButton(url)

	err := downloadButton.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
