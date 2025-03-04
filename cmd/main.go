package main

import (
	"fmt"
	"github.com/BioMihanoid/URLShortener/internal/config"
	"github.com/BioMihanoid/URLShortener/internal/storage/sqlite"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(fmt.Sprintf("%s", "logger info"))
	log.Debug(fmt.Sprintf("%s", "logger debug"))

	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Debug("failed to init storage", err)
		os.Exit(1)
	}

	_ = storage
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
