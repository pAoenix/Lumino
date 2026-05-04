package csvutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProfileCSV(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sales.csv")
	content := "name,amount,amount\nalpha,10,100\nbeta,20,200\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	columns, rows, preview, profile, err := Profile(path)
	if err != nil {
		t.Fatal(err)
	}
	if rows != 2 {
		t.Fatalf("rows = %d, want 2", rows)
	}
	if got, want := columns[2], "amount_2"; got != want {
		t.Fatalf("duplicate column = %q, want %q", got, want)
	}
	if len(preview) != 2 || preview[0]["amount"] != "10" {
		t.Fatalf("unexpected preview: %#v", preview)
	}
	if len(profile.Numeric) != 2 {
		t.Fatalf("numeric profiles = %d, want 2", len(profile.Numeric))
	}
}
