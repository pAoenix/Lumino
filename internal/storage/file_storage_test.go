package storage

import (
	"strings"
	"testing"
)

func TestFileStorageSaveDelete(t *testing.T) {
	store, err := NewFileStorage(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	stored, err := store.Save("dataset-1", "sample.CSV", strings.NewReader("a,b\n1,2\n"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := stored.StoredName, "dataset-1.csv"; got != want {
		t.Fatalf("stored name = %q, want %q", got, want)
	}
	if got, want := stored.RelativePath, "uploads/dataset-1.csv"; got != want {
		t.Fatalf("relative path = %q, want %q", got, want)
	}
	if err := store.Delete(stored.RelativePath); err != nil {
		t.Fatal(err)
	}
	if err := store.Delete(stored.RelativePath); err != nil {
		t.Fatalf("second delete should ignore missing file: %v", err)
	}
}

func TestFileStorageRejectsPathTraversal(t *testing.T) {
	store, err := NewFileStorage(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.AbsolutePath("../escape.csv"); err == nil {
		t.Fatal("expected path traversal to be rejected")
	}
	if _, err := store.UploadPath("../escape.csv"); err == nil {
		t.Fatal("expected unsafe stored file name to be rejected")
	}
}
