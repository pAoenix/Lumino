package httpserver

import (
	"log/slog"
	"net/http"

	"lumino/internal/httpserver/handlers"
	"lumino/internal/service"
	"lumino/internal/web"
)

type Server struct {
	logger  *slog.Logger
	service *service.DatasetService
}

func NewServer(datasetService *service.DatasetService, logger *slog.Logger) *Server {
	return &Server{logger: logger, service: datasetService}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	datasetHandler := handlers.NewDatasetHandler(s.service)

	mux.HandleFunc("GET /api/health", handlers.HandleHealth)
	mux.HandleFunc("GET /api/datasets", datasetHandler.List)
	mux.HandleFunc("POST /api/datasets", datasetHandler.Upload)
	mux.HandleFunc("GET /api/datasets/{id}", datasetHandler.Get)
	mux.HandleFunc("DELETE /api/datasets/{id}", datasetHandler.Delete)
	mux.HandleFunc("GET /api/datasets/{id}/download", datasetHandler.Download)
	mux.HandleFunc("GET /", handleIndex)
	mux.Handle("GET /web/", http.StripPrefix("/web/", http.FileServer(http.FS(web.Files))))
	return WithLogging(s.logger, WithCORS(mux))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFileFS(w, r, web.Files, "index.html")
}
