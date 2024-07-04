package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/henrik392/youtube-voice-go/cmd/web"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", s.healthHandler)

	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)
	r.Get("/", templ.Handler(web.MainPage()).ServeHTTP)
	r.Post("/validate-url", web.ValidateURLHandler)
	r.Post("/generate-voice", web.GenerateVoiceHandler)
	r.Get("/serve-audio", web.ServeAudioHandler)

	fmt.Println("Server started on port", os.Getenv("PORT"))

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
