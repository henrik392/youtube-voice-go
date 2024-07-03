package web

import "net/http"

func ServeAudioHandler(w http.ResponseWriter, r *http.Request) {
	audioPath := r.URL.Query().Get("path")
	if audioPath == "" {
		http.Error(w, "Audio path is required", http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, audioPath)
}
