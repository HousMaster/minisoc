package httpserver

import (
	"app/internal/storage"
	"context"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

type Server struct {
	ctx        context.Context
	log        *slog.Logger
	storage    storage.Storage
	httpServer *echo.Echo
	validator  *validator.Validate
	addr       string
}

func New(ctx context.Context, log *slog.Logger, storage storage.Storage, port int) *Server {

	server := &Server{
		ctx:        ctx,
		storage:    storage,
		log:        log,
		httpServer: echo.New(), // завтра поменять на fasthttp
		validator:  validator.New(),
		addr:       fmt.Sprintf(":%d", port),
	}

	return server
}

func (s *Server) SetRoutes() {

	// auth user
	// profile user
}

func (s *Server) Run() error {
	const op = "httpserver.Run"
	log := s.log.With(slog.String("op", op))

	log.Info("http server is running", slog.String("addr", s.addr))

	if err := s.httpServer.Start(s.addr); err != nil {
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
