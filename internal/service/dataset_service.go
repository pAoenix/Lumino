package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"lumino/internal/csvutil"
	"lumino/internal/domain"
	"lumino/internal/repository"
	"lumino/internal/storage"
)

type DatasetService struct {
	repo    repository.DatasetRepository
	storage *storage.FileStorage
}

type UploadInput struct {
	File        multipart.File
	Header      *multipart.FileHeader
	Name        string
	Description string
}

type DownloadFile struct {
	Path string
	Name string
}

func NewDatasetService(repo repository.DatasetRepository, fileStorage *storage.FileStorage) *DatasetService {
	return &DatasetService{repo: repo, storage: fileStorage}
}

func (s *DatasetService) List(ctx context.Context) ([]domain.Dataset, error) {
	return s.repo.List(ctx)
}

func (s *DatasetService) Get(ctx context.Context, id string) (domain.Dataset, error) {
	dataset, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Dataset{}, err
	}

	if strings.EqualFold(filepath.Ext(dataset.FileName), ".csv") {
		path, err := s.storage.UploadPath(dataset.StorageName())
		if err != nil {
			return domain.Dataset{}, err
		}
		_, _, preview, _, err := csvutil.Profile(path)
		if err == nil {
			dataset.Preview = preview
		}
	}
	return dataset, nil
}

func (s *DatasetService) Upload(ctx context.Context, input UploadInput) (domain.Dataset, error) {
	id := newID()
	stored, err := s.storage.Save(id, input.Header.Filename, input.File)
	if err != nil {
		return domain.Dataset{}, fmt.Errorf("save file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(input.Header.Filename))
	name := strings.TrimSpace(input.Name)
	if name == "" {
		name = strings.TrimSuffix(input.Header.Filename, ext)
	}
	now := time.Now().UTC()
	dataset := domain.Dataset{
		ID:          id,
		Name:        name,
		Description: strings.TrimSpace(input.Description),
		FileName:    input.Header.Filename,
		StoredName:  stored.StoredName,
		ContentType: contentType(input.Header.Filename, input.Header.Header.Get("Content-Type")),
		Size:        stored.Size,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if ext == ".csv" {
		columns, rows, preview, profile, err := csvutil.Profile(stored.AbsolutePath)
		if err != nil {
			_ = s.storage.Delete(stored.RelativePath)
			return domain.Dataset{}, ErrInvalidCSV
		}
		dataset.Columns = columns
		dataset.Rows = rows
		dataset.Profile = profile
		dataset.Preview = preview
	}

	if err := s.repo.Create(ctx, dataset); err != nil {
		_ = s.storage.Delete(stored.RelativePath)
		return domain.Dataset{}, fmt.Errorf("save metadata: %w", err)
	}
	return dataset, nil
}

func (s *DatasetService) Delete(ctx context.Context, id string) error {
	dataset, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	path, err := s.storage.UploadPath(dataset.StorageName())
	if err != nil {
		return err
	}
	relativePath := filepath.ToSlash(filepath.Join("uploads", filepath.Base(path)))
	if err := s.storage.Delete(relativePath); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *DatasetService) Download(ctx context.Context, id string) (DownloadFile, error) {
	dataset, err := s.repo.Get(ctx, id)
	if err != nil {
		return DownloadFile{}, err
	}
	path, err := s.storage.UploadPath(dataset.StorageName())
	if err != nil {
		return DownloadFile{}, err
	}
	return DownloadFile{Path: path, Name: dataset.FileName}, nil
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
	if _, err := io.ReadFull(rand.Reader, bytes[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes[:])
}
