package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Addr    string
	DataDir string
	DBPath  string
}

func Load() Config {
	dataDir := env("LUMINO_DATA_DIR", "data")
	return Config{
		Addr:    env("LUMINO_ADDR", ":8080"),
		DataDir: dataDir,
		DBPath:  env("LUMINO_DB_PATH", filepath.Join(dataDir, "lumino.db")),
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
