package app

import (
	"log/slog"
)

type App struct {
}

func New(log *slog.Logger, grpcPort int, storagePath string) *App {

	return &App{}
}
