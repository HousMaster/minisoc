package main

import (
	"app/internal/app"
	"app/internal/config"
	"app/internal/storage/sqlite"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	evnProd  = "prod"
)

func main() {

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	//
	conf := config.MustLoad()
	//
	log := setupLogger(conf.Env)
	//
	storage := sqlite.MustInit(conf.StoragePath)
	//
	application := app.New(ctx, log, storage, conf.HTTPServer.Port)
	//
	application.HttpServer.MustRun()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case evnProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
