package app

import (
	"context"
	"log/slog"

	"app/internal/httpserver"
	"app/internal/storage"
)

type App struct {
	HttpServer *httpserver.Server
}

func New(ctx context.Context, log *slog.Logger, storage storage.Storage, port int) *App {

	return &App{
		HttpServer: httpserver.New(ctx, log, storage, port),
	}
}
