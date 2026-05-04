package app

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed web/*
var webFiles embed.FS

type Server struct {
	cfg    Config
	logger *slog.Logger
	store  *DatasetStore
}

func NewServer(cfg Config, logger *slog.Logger) (*Server, error) {
	store, err := NewDatasetStore(cfg.DataDir)
	if err != nil {
		return nil, err
	}
	return &Server{cfg: cfg, logger: logger, store: store}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/datasets", s.handleListDatasets)
	mux.HandleFunc("POST /api/datasets", s.handleUploadDataset)
	mux.HandleFunc("GET /api/datasets/{id}", s.handleGetDataset)
	mux.HandleFunc("DELETE /api/datasets/{id}", s.handleDeleteDataset)
	mux.HandleFunc("GET /api/datasets/{id}/download", s.handleDownloadDataset)
	mux.HandleFunc("GET /", s.handleIndex)
	mux.Handle("GET /web/", http.FileServer(http.FS(webFiles)))
	return withLogging(s.logger, withCORS(mux))
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFileFS(w, r, webFiles, "web/index.html")
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleListDatasets(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.List())
}

func (s *Server) handleGetDataset(w http.ResponseWriter, r *http.Request) {
	dataset, err := s.store.Get(r.PathValue("id"))
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "dataset not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read dataset")
		return
	}

	path := s.store.UploadPath(dataset.storageName())
	if strings.EqualFold(filepath.Ext(dataset.FileName), ".csv") {
		_, _, preview, _, err := ProfileCSV(path)
		if err == nil {
			dataset.Preview = preview
		}
	}
	writeJSON(w, http.StatusOK, dataset)
}

func (s *Server) handleUploadDataset(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 100<<20)
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "upload must be multipart/form-data and no larger than 100MB")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	id := newID()
	ext := strings.ToLower(filepath.Ext(header.Filename))
	storedName := id + ext
	path := s.store.UploadPath(storedName)

	out, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create upload file")
		return
	}
	size, copyErr := io.Copy(out, file)
	closeErr := out.Close()
	if copyErr != nil || closeErr != nil {
		_ = os.Remove(path)
		writeError(w, http.StatusInternalServerError, "failed to save upload file")
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		name = strings.TrimSuffix(header.Filename, ext)
	}

	dataset := Dataset{
		ID:          id,
		Name:        name,
		Description: strings.TrimSpace(r.FormValue("description")),
		FileName:    header.Filename,
		StoredName:  storedName,
		ContentType: contentType(header.Filename, header.Header.Get("Content-Type")),
		Size:        size,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if ext == ".csv" {
		columns, rows, preview, profile, err := ProfileCSV(path)
		if err != nil {
			_ = os.Remove(path)
			writeError(w, http.StatusBadRequest, "CSV file could not be parsed")
			return
		}
		dataset.Columns = columns
		dataset.Rows = rows
		dataset.Profile = profile
		dataset.Preview = preview
	}

	if err := s.store.Save(dataset); err != nil {
		_ = os.Remove(path)
		writeError(w, http.StatusInternalServerError, "failed to save dataset metadata")
		return
	}

	writeJSON(w, http.StatusCreated, dataset)
}

func (s *Server) handleDownloadDataset(w http.ResponseWriter, r *http.Request) {
	dataset, err := s.store.Get(r.PathValue("id"))
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "dataset not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read dataset")
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", dataset.FileName))
	http.ServeFile(w, r, s.store.UploadPath(dataset.storageName()))
}

func (s *Server) handleDeleteDataset(w http.ResponseWriter, r *http.Request) {
	err := s.store.Delete(r.PathValue("id"))
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "dataset not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete dataset")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func contentType(name string, fallback string) string {
	if typ := mime.TypeByExtension(filepath.Ext(name)); typ != "" {
		return typ
	}
	if fallback != "" {
		return fallback
	}
	return "application/octet-stream"
}

func newID() string {
	var bytes [12]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes[:])
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func withLogging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}
