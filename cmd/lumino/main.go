package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lumino/internal/config"
	"lumino/internal/httpserver"
	postgresrepo "lumino/internal/repository/postgres"
	"lumino/internal/service"
	"lumino/internal/storage"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg := config.Load()

	fileStorage, err := storage.NewFileStorage(cfg.DataDir)
	if err != nil {
		logger.Error("failed to create file storage", "error", err)
		os.Exit(1)
	}

	repo, err := postgresrepo.NewDatasetRepository(context.Background(), cfg.DatabaseURL, cfg.DataDir)
	if err != nil {
		logger.Error("failed to create dataset repository", "error", err)
		os.Exit(1)
	}
	defer repo.Close()

	datasetService := service.NewDatasetService(repo, fileStorage)
	server := httpserver.NewServer(datasetService, logger)

	httpServer := &http.Server{
		Addr:              cfg.Addr,
		Handler:           server.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errs := make(chan error, 1)
	go func() {
		logger.Info("lumino data platform started", "addr", cfg.Addr, "data_dir", cfg.DataDir, "database", "postgres")
		errs <- httpServer.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errs:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Error("server shutdown failed", "error", err)
			os.Exit(1)
		}
		logger.Info("server stopped")
	}
}
