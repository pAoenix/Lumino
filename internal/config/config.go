package config

import (
	"os"
)

type Config struct {
	Addr        string
	DataDir     string
	DatabaseURL string
}

func Load() Config {
	dataDir := env("LUMINO_DATA_DIR", "data")
	return Config{
		Addr:        env("LUMINO_ADDR", ":8080"),
		DataDir:     dataDir,
		DatabaseURL: env("DATABASE_URL", "postgres://postgres@localhost:5432/lumino?sslmode=disable"),
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
