package internalhttp

import (
	"context"
	"net"
	"net/http"

	internalconfig "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/config"
)

type Server struct {
	server *http.Server
	logger Logger
}

type Application interface{}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type RequestHandler struct{}

func (h *RequestHandler) Hello(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Hello, World!"))
}

func NewServer(config internalconfig.HttpConf, app Application, logger Logger) *Server {
	handler := &RequestHandler{}
	mux := http.NewServeMux()
	mux.HandleFunc("/", loggingMiddleware(handler.Hello, logger))
	server := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
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
