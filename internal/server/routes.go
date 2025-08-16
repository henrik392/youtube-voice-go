package server

import (
	"net/http"

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
	r.Post("/process-video", web.ProcessVideoHandler)
	r.Post("/generate-voice", web.GenerateVoiceHandler)
	r.Post("/generate-voice-enhanced", web.GenerateVoiceEnhancedHandler)
	r.Post("/generate-voice-optimized", web.GenerateVoiceOptimizedHandler)
	r.Get("/serve-audio", web.ServeAudioHandler)
	
	// Component handlers for dynamic loading
	r.Get("/components/url-input", web.URLInputHandler)
	r.Get("/components/file-upload", web.FileUploadHandler)
	r.Get("/components/microphone", web.MicrophoneHandler)
	
	// Audio input handlers
	r.Post("/upload-audio", web.UploadAudioHandler)
	r.Post("/save-recording", web.SaveRecordingHandler)

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

	// For DB health check
	// jsonResp, _ := json.Marshal(s.db.Health())
	// _, _ = w.Write(jsonResp)
}
