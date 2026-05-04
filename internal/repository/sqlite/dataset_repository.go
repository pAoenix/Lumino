package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"lumino/internal/domain"

	_ "modernc.org/sqlite"
)

type DatasetRepository struct {
	db *sql.DB
}

func NewDatasetRepository(ctx context.Context, dbPath string, dataDir string) (*DatasetRepository, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	repo := &DatasetRepository{db: db}
	if err := repo.init(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := repo.migrateMetadataJSON(ctx, filepath.Join(dataDir, "metadata.json")); err != nil {
		_ = db.Close()
		return nil, err
	}
	return repo, nil
}

func (r *DatasetRepository) Close() error {
	return r.db.Close()
}

func (r *DatasetRepository) List(ctx context.Context) ([]domain.Dataset, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description, original_file_name, stored_file_name, content_type,
		       size, rows, columns_json, profile_json, created_at, updated_at
		FROM datasets
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	datasets := make([]domain.Dataset, 0)
	for rows.Next() {
		dataset, err := scanDataset(rows)
		if err != nil {
			return nil, err
		}
		datasets = append(datasets, dataset)
	}
	return datasets, rows.Err()
}

func (r *DatasetRepository) Get(ctx context.Context, id string) (domain.Dataset, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, original_file_name, stored_file_name, content_type,
		       size, rows, columns_json, profile_json, created_at, updated_at
		FROM datasets
		WHERE id = ?
	`, id)
	dataset, err := scanDataset(row)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Dataset{}, domain.ErrNotFound
	}
	return dataset, err
}

func (r *DatasetRepository) Create(ctx context.Context, dataset domain.Dataset) error {
	return r.upsert(ctx, dataset)
}

func (r *DatasetRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM datasets WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *DatasetRepository) init(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS datasets (
		  id TEXT PRIMARY KEY,
		  name TEXT NOT NULL,
		  description TEXT NOT NULL DEFAULT '',
		  original_file_name TEXT NOT NULL,
		  stored_file_name TEXT NOT NULL,
		  relative_path TEXT NOT NULL,
		  content_type TEXT NOT NULL,
		  size INTEGER NOT NULL,
		  rows INTEGER NOT NULL DEFAULT 0,
		  columns_json TEXT NOT NULL DEFAULT '[]',
		  profile_json TEXT NOT NULL DEFAULT '{"numeric":[]}',
		  created_at TEXT NOT NULL,
		  updated_at TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_datasets_created_at ON datasets(created_at DESC);
	`)
	return err
}

func (r *DatasetRepository) migrateMetadataJSON(ctx context.Context, metadataPath string) error {
	file, err := os.Open(metadataPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	var datasets []domain.Dataset
	if err := json.NewDecoder(file).Decode(&datasets); err != nil {
		return fmt.Errorf("decode metadata json: %w", err)
	}
	for _, dataset := range datasets {
		if dataset.ID == "" {
			continue
		}
		if dataset.StoredName == "" {
			dataset.StoredName = dataset.FileName
		}
		if dataset.Profile.Numeric == nil {
			dataset.Profile.Numeric = []domain.NumericProfile{}
		}
		if dataset.CreatedAt.IsZero() {
			dataset.CreatedAt = time.Now().UTC()
		}
		if dataset.UpdatedAt.IsZero() {
			dataset.UpdatedAt = dataset.CreatedAt
		}
		if err := r.upsert(ctx, dataset); err != nil {
			return err
		}
	}
	return nil
}

func (r *DatasetRepository) upsert(ctx context.Context, dataset domain.Dataset) error {
	if dataset.Profile.Numeric == nil {
		dataset.Profile.Numeric = []domain.NumericProfile{}
	}
	columnsJSON, err := json.Marshal(dataset.Columns)
	if err != nil {
		return err
	}
	profileJSON, err := json.Marshal(dataset.Profile)
	if err != nil {
		return err
	}
	relativePath := filepath.ToSlash(filepath.Join("uploads", dataset.StorageName()))
	createdAt := dataset.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	updatedAt := dataset.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO datasets (
			id, name, description, original_file_name, stored_file_name, relative_path,
			content_type, size, rows, columns_json, profile_json, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			description = excluded.description,
			original_file_name = excluded.original_file_name,
			stored_file_name = excluded.stored_file_name,
			relative_path = excluded.relative_path,
			content_type = excluded.content_type,
			size = excluded.size,
			rows = excluded.rows,
			columns_json = excluded.columns_json,
			profile_json = excluded.profile_json,
			created_at = excluded.created_at,
			updated_at = excluded.updated_at
	`, dataset.ID, dataset.Name, dataset.Description, dataset.FileName, dataset.StorageName(), relativePath,
		dataset.ContentType, dataset.Size, dataset.Rows, string(columnsJSON), string(profileJSON),
		createdAt.Format(time.RFC3339Nano), updatedAt.Format(time.RFC3339Nano))
	return err
}

type datasetScanner interface {
	Scan(dest ...any) error
}

func scanDataset(scanner datasetScanner) (domain.Dataset, error) {
	var dataset domain.Dataset
	var columnsJSON string
	var profileJSON string
	var createdAt string
	var updatedAt string
	if err := scanner.Scan(
		&dataset.ID,
		&dataset.Name,
		&dataset.Description,
		&dataset.FileName,
		&dataset.StoredName,
		&dataset.ContentType,
		&dataset.Size,
		&dataset.Rows,
		&columnsJSON,
		&profileJSON,
		&createdAt,
		&updatedAt,
	); err != nil {
		return domain.Dataset{}, err
	}
	if err := json.Unmarshal([]byte(columnsJSON), &dataset.Columns); err != nil {
		return domain.Dataset{}, err
	}
	if err := json.Unmarshal([]byte(profileJSON), &dataset.Profile); err != nil {
		return domain.Dataset{}, err
	}
	var err error
	dataset.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return domain.Dataset{}, err
	}
	dataset.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return domain.Dataset{}, err
	}
	return dataset, nil
}
