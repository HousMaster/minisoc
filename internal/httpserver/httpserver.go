package httpserver

import (
	"app/internal/httpserver/eventserver"
	"app/internal/storage"
	"context"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/valyala/fasthttp"
)

type Server struct {
	ctx            context.Context
	log            *slog.Logger
	storage        storage.Storage
	httpServer     *echo.Echo
	validator      *validator.Validate
	addr           string
	tokenSecretKey []byte
	eventserver    *eventserver.EventServer
}

func New(ctx context.Context, log *slog.Logger, storage storage.Storage, port int, jwtKey string) *Server {

	server := &Server{
		ctx:            ctx,
		storage:        storage,
		log:            log,
		validator:      validator.New(),
		addr:           fmt.Sprintf(":%d", port),
		tokenSecretKey: []byte(jwtKey),
		eventserver:    eventserver.New(ctx, log),
	}

	return server
}

func (s *Server) Run() error {
	const op = "httpserver.Run"
	log := s.log.With(slog.String("op", op))

	log.Info("http server is running", slog.String("addr", s.addr))

	// run eventserver
	go s.eventserver.Run()

	if err := fasthttp.ListenAndServe(s.addr, s.router); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	// const op = "httpserver.Stop"
	return s.httpServer.Shutdown(ctx)
}
