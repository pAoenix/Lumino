package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileStorage struct {
	dataDir   string
	uploadDir string
}

type StoredFile struct {
	StoredName   string
	RelativePath string
	AbsolutePath string
	Size         int64
}

func NewFileStorage(dataDir string) (*FileStorage, error) {
	store := &FileStorage{
		dataDir:   dataDir,
		uploadDir: filepath.Join(dataDir, "uploads"),
	}
	if err := os.MkdirAll(store.uploadDir, 0o755); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *FileStorage) Save(id string, originalName string, src io.Reader) (StoredFile, error) {
	ext := strings.ToLower(filepath.Ext(originalName))
	storedName := id + ext
	relativePath := filepath.ToSlash(filepath.Join("uploads", storedName))
	absolutePath, err := s.AbsolutePath(relativePath)
	if err != nil {
		return StoredFile{}, err
	}

	out, err := os.OpenFile(absolutePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return StoredFile{}, err
	}
	size, copyErr := io.Copy(out, src)
	closeErr := out.Close()
	if copyErr != nil || closeErr != nil {
		_ = os.Remove(absolutePath)
		if copyErr != nil {
			return StoredFile{}, copyErr
		}
		return StoredFile{}, closeErr
	}

	return StoredFile{
		StoredName:   storedName,
		RelativePath: relativePath,
		AbsolutePath: absolutePath,
		Size:         size,
	}, nil
}

func (s *FileStorage) AbsolutePath(relativePath string) (string, error) {
	if relativePath == "" {
		return "", errors.New("relative path is required")
	}
	clean := filepath.Clean(relativePath)
	if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("unsafe relative path: %s", relativePath)
	}
	return filepath.Join(s.dataDir, clean), nil
}

func (s *FileStorage) UploadPath(storedName string) (string, error) {
	if storedName == "" || storedName != filepath.Base(storedName) {
		return "", fmt.Errorf("unsafe stored file name: %s", storedName)
	}
	return s.AbsolutePath(filepath.ToSlash(filepath.Join("uploads", storedName)))
}

func (s *FileStorage) Delete(relativePath string) error {
	path, err := s.AbsolutePath(relativePath)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
