package main

import (
	"log/slog"
	"os"

	"github.com/HousMaster/microsoc/internal/app"
	"github.com/HousMaster/microsoc/internal/config"
)

const (
	envLocal = "local"
	evnProd  = "prod"
)

func main() {

	//
	conf := config.MustLoad()
	//
	log := setupLogger(conf.Env)
	//
	application := app.New(log, conf.GRPC.Port, conf.StoragePath)
	_ = application

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
