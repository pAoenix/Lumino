package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"lumino/internal/domain"
	"lumino/internal/service"
)

type DatasetHandler struct {
	service *service.DatasetService
}

func NewDatasetHandler(datasetService *service.DatasetService) *DatasetHandler {
	return &DatasetHandler{service: datasetService}
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *DatasetHandler) List(w http.ResponseWriter, r *http.Request) {
	datasets, err := h.service.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read datasets")
		return
	}
	writeJSON(w, http.StatusOK, datasets)
}

func (h *DatasetHandler) Get(w http.ResponseWriter, r *http.Request) {
	dataset, err := h.service.Get(r.Context(), r.PathValue("id"))
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "dataset not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read dataset")
		return
	}
	writeJSON(w, http.StatusOK, dataset)
}

func (h *DatasetHandler) Upload(w http.ResponseWriter, r *http.Request) {
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

	dataset, err := h.service.Upload(r.Context(), service.UploadInput{
		File:        file,
		Header:      header,
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	})
	if errors.Is(err, service.ErrInvalidCSV) {
		writeError(w, http.StatusBadRequest, "CSV file could not be parsed")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save dataset")
		return
	}

	writeJSON(w, http.StatusCreated, dataset)
}

func (h *DatasetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.service.Delete(r.Context(), r.PathValue("id"))
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "dataset not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete dataset")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *DatasetHandler) Download(w http.ResponseWriter, r *http.Request) {
	file, err := h.service.Download(r.Context(), r.PathValue("id"))
	if errors.Is(err, domain.ErrNotFound) {
		writeError(w, http.StatusNotFound, "dataset not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read dataset")
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", file.Name))
	http.ServeFile(w, r, file.Path)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
