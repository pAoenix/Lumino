package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

var ErrNotFound = errors.New("dataset not found")

type DatasetStore struct {
	mu        sync.RWMutex
	dataDir   string
	uploadDir string
	index     string
	items     map[string]Dataset
}

func NewDatasetStore(dataDir string) (*DatasetStore, error) {
	store := &DatasetStore{
		dataDir:   dataDir,
		uploadDir: filepath.Join(dataDir, "uploads"),
		index:     filepath.Join(dataDir, "metadata.json"),
		items:     make(map[string]Dataset),
	}
	if err := os.MkdirAll(store.uploadDir, 0o755); err != nil {
		return nil, err
	}
	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *DatasetStore) List() []Dataset {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]Dataset, 0, len(s.items))
	for _, item := range s.items {
		item.Preview = nil
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	return items
}

func (s *DatasetStore) Get(id string) (Dataset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return Dataset{}, ErrNotFound
	}
	return item, nil
}

func (s *DatasetStore) Save(dataset Dataset) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	if dataset.CreatedAt.IsZero() {
		dataset.CreatedAt = now
	}
	dataset.UpdatedAt = now
	s.items[dataset.ID] = dataset
	return s.persist()
}

func (s *DatasetStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	if !ok {
		return ErrNotFound
	}

	if err := os.Remove(s.UploadPath(item.storageName())); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	delete(s.items, id)
	if err := s.persist(); err != nil {
		s.items[id] = item
		return err
	}
	return nil
}

func (s *DatasetStore) UploadPath(fileName string) string {
	return filepath.Join(s.uploadDir, fileName)
}

func (s *DatasetStore) load() error {
	file, err := os.Open(s.index)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	var items []Dataset
	if err := json.NewDecoder(file).Decode(&items); err != nil {
		return fmt.Errorf("decode metadata: %w", err)
	}
	for _, item := range items {
		s.items[item.ID] = item
	}
	return nil
}

func (s *DatasetStore) persist() error {
	items := make([]Dataset, 0, len(s.items))
	for _, item := range s.items {
		item.Preview = nil
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})

	tmp := s.index + ".tmp"
	file, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(items); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, s.index)
}
