package internalhttp

import (
	"context"
	"net"
	"net/http"

	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type Server struct {
	server *http.Server
	logger Logger
}

type Config interface {
	GetServerHost() string
	GetServerPort() string
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error)
	RemoveEvent(ctx context.Context, event storage.Event) error
	GetDayAheadEvents(ctx context.Context) ([]storage.Event, error)
	GetWeekAheadEvents(ctx context.Context) ([]storage.Event, error)
	GetMonthAheadEvents(ctx context.Context) ([]storage.Event, error)
}

type RequestHandler struct {
	App Application
}

func (h *RequestHandler) Hello(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Hello, World!"))
}

func NewServer(config Config, app Application, logger Logger) *Server {
	handler := &RequestHandler{
		App: app,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", loggingMiddleware(handler.Hello, logger))

	server := &http.Server{
		Addr:    net.JoinHostPort(config.GetServerHost(), config.GetServerPort()),
		Handler: mux,
	}

	return &Server{
		server: server,
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil {
		s.logger.Error(err.Error())
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
