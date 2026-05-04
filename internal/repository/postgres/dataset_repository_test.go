package postgres

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"lumino/internal/domain"
)

func TestDatasetRepositoryCreateListGetDelete(t *testing.T) {
	repo, cleanup := newTestRepository(t)
	defer cleanup()

	ctx := context.Background()
	dataset := testDataset("dataset-1", "Sales")
	if err := repo.Create(ctx, dataset); err != nil {
		t.Fatal(err)
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != dataset.ID {
		t.Fatalf("unexpected list: %#v", list)
	}

	got, err := repo.Get(ctx, dataset.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Columns[0] != "region" || got.Profile.Numeric[0].Column != "amount" {
		t.Fatalf("unexpected dataset payload: %#v", got)
	}

	if err := repo.Delete(ctx, dataset.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Get(ctx, dataset.ID); err != domain.ErrNotFound {
		t.Fatalf("Get after delete error = %v, want ErrNotFound", err)
	}
}

func TestDatasetRepositoryMigratesMetadataJSON(t *testing.T) {
	ctx := context.Background()
	dataDir := t.TempDir()
	legacy := testDataset("legacy-1", "Legacy")
	legacy.StoredName = ""
	payload, err := json.Marshal([]domain.Dataset{legacy})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, "metadata.json"), payload, 0o644); err != nil {
		t.Fatal(err)
	}

	resetTestDatabase(t)
	repo, err := NewDatasetRepository(ctx, testDatabaseURL(t), dataDir)
	if err != nil {
		t.Fatal(err)
	}
	cleanup := func() {
		_, _ = repo.db.ExecContext(context.Background(), `TRUNCATE datasets`)
		_ = repo.Close()
	}
	defer cleanup()

	got, err := repo.Get(ctx, "legacy-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.StoredName != legacy.FileName {
		t.Fatalf("stored name = %q, want fallback file name %q", got.StoredName, legacy.FileName)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "metadata.json")); err != nil {
		t.Fatalf("metadata.json should be kept: %v", err)
	}
}

func newTestRepository(t *testing.T) (*DatasetRepository, func()) {
	t.Helper()
	return newTestRepositoryWithDataDir(t, t.TempDir())
}

func newTestRepositoryWithDataDir(t *testing.T, dataDir string) (*DatasetRepository, func()) {
	t.Helper()
	ctx := context.Background()
	repo, err := NewDatasetRepository(ctx, testDatabaseURL(t), dataDir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.db.ExecContext(ctx, `TRUNCATE datasets`); err != nil {
		_ = repo.Close()
		t.Fatal(err)
	}
	cleanup := func() {
		_, _ = repo.db.ExecContext(context.Background(), `TRUNCATE datasets`)
		_ = repo.Close()
	}
	return repo, cleanup
}

func resetTestDatabase(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	repo, err := NewDatasetRepository(ctx, testDatabaseURL(t), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	if _, err := repo.db.ExecContext(ctx, `TRUNCATE datasets`); err != nil {
		t.Fatal(err)
	}
}

func testDatabaseURL(t *testing.T) string {
	t.Helper()
	databaseURL := os.Getenv("LUMINO_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("LUMINO_TEST_DATABASE_URL is not set")
	}
	return databaseURL
}

func testDataset(id string, name string) domain.Dataset {
	now := time.Date(2026, 5, 4, 12, 0, 0, 0, time.UTC)
	return domain.Dataset{
		ID:          id,
		Name:        name,
		Description: "desc",
		FileName:    name + ".csv",
		StoredName:  id + ".csv",
		ContentType: "text/csv; charset=utf-8",
		Size:        16,
		Rows:        2,
		Columns:     []string{"region", "amount"},
		Profile: domain.DataProfile{Numeric: []domain.NumericProfile{
			{Column: "amount", Count: 2, Min: 10, Max: 20, Avg: 15},
		}},
		CreatedAt: now,
		UpdatedAt: now,
	}
}
