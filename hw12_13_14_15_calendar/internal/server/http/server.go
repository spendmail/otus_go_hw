package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	GetHTTPHost() string
	GetHTTPPort() string
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
	App    Application
	Logger Logger
}

// Hello processes a root url.
func (h *RequestHandler) Hello(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write([]byte("Hello, World!"))
	if err != nil {
		h.Logger.Error(err.Error())
	}
}

// Create handles creating a new event.
func (h *RequestHandler) Create(writer http.ResponseWriter, request *http.Request) {
	event := storage.Event{}

	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(fmt.Sprintf("Unable to read the request: %q", err.Error()))
		return
	}

	if err = json.Unmarshal(b, &event); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(fmt.Sprintf("Unable to unmarshal the request: %q", err.Error()))
		return
	}

	event, err = h.App.CreateEvent(context.Background(), event)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(fmt.Sprintf("Unable to save the event: %q", err.Error()))
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(fmt.Sprintf("Event %q (%d) has been created successfully.", event.Title, event.ID))
}

// Update handles updating the event.
func (h *RequestHandler) Update(writer http.ResponseWriter, request *http.Request) {
	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(fmt.Sprintf("Unable to read the request: %q", err.Error()))
		return
	}

	event := storage.Event{}
	if err = json.Unmarshal(b, &event); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(fmt.Sprintf("Unable to unmarshal the request: %q", err.Error()))
		return
	}

	event, err = h.App.UpdateEvent(context.Background(), event)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(fmt.Sprintf("Unable to update the event: %q", err.Error()))
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(fmt.Sprintf("Event %q (%d) has been updated successfully.", event.Title, event.ID))
}

// Remove handles removing the event.
func (h *RequestHandler) Remove(writer http.ResponseWriter, request *http.Request) {
	event := storage.Event{}

	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(fmt.Sprintf("Unable to read the request: %q", err.Error()))
		return
	}

	if err = json.Unmarshal(b, &event); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(fmt.Sprintf("Unable to unmarshal the request: %q", err.Error()))
		return
	}

	err = h.App.RemoveEvent(context.Background(), event)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(fmt.Sprintf("Unable to remove the event: %q", err.Error()))
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(fmt.Sprintf("Event %d has been removed successfully.", event.ID))
}

// GetDayAheadEvents returns daily events.
func (h *RequestHandler) GetDayAheadEvents(writer http.ResponseWriter, request *http.Request) {
	events, err := h.App.GetDayAheadEvents(context.Background())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(fmt.Sprintf("Unable to get daily events: %q", err.Error()))
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(events)
}

// GetWeekAheadEvents returns weekly events.
func (h *RequestHandler) GetWeekAheadEvents(writer http.ResponseWriter, request *http.Request) {
	events, err := h.App.GetWeekAheadEvents(context.Background())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(fmt.Sprintf("Unable to get weekly events: %q", err.Error()))
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(events)
}

// GetMonthAheadEvents returns weekly events.
func (h *RequestHandler) GetMonthAheadEvents(writer http.ResponseWriter, request *http.Request) {
	events, err := h.App.GetMonthAheadEvents(context.Background())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(fmt.Sprintf("Unable to get monthly events: %q", err.Error()))
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(events)
}

// NewServer returns a new server instance.
func NewServer(config Config, app Application, logger Logger) *Server {
	handler := &RequestHandler{
		App:    app,
		Logger: logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", loggingMiddleware(handler.Hello, logger))
	mux.HandleFunc("/event/create", loggingMiddleware(handler.Create, logger))
	mux.HandleFunc("/event/update", loggingMiddleware(handler.Update, logger))
	mux.HandleFunc("/event/remove", loggingMiddleware(handler.Remove, logger))
	mux.HandleFunc("/event/day", loggingMiddleware(handler.GetDayAheadEvents, logger))
	mux.HandleFunc("/event/week", loggingMiddleware(handler.GetWeekAheadEvents, logger))
	mux.HandleFunc("/event/month", loggingMiddleware(handler.GetMonthAheadEvents, logger))

	server := &http.Server{
		Addr:    net.JoinHostPort(config.GetHTTPHost(), config.GetHTTPPort()),
		Handler: mux,
	}

	return &Server{
		server: server,
		logger: logger,
	}
}

// Start launches a HTTP server.
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Stop suspends HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
