//go:generate protoc --go_out ./eventpb/ --go-grpc_out ./eventpb/ ./api/EventService.proto
package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	pb "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/server/grpc/eventpb"
	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	EventDateFormat = "2006-01-02T15:04:05Z"
	ErrServerStart  = errors.New("unable to start grpc server")
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	GetZapLogger() *zap.Logger
}

type Server struct {
	server *grpc.Server
	config Config
	logger Logger
}

type Config interface {
	GetGrpcHost() string
	GetGrpcPort() string
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error)
	RemoveEvent(ctx context.Context, event storage.Event) error
	GetDayAheadEvents(ctx context.Context) ([]storage.Event, error)
	GetWeekAheadEvents(ctx context.Context) ([]storage.Event, error)
	GetMonthAheadEvents(ctx context.Context) ([]storage.Event, error)
}

type Service struct {
	pb.UnimplementedCalendarServer
	app    Application
	logger Logger
}

// CreateEvent handles creating a new event via grpc.
func (s *Service) CreateEvent(ctx context.Context, createEventRequest *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	event := storage.Event{}

	event.Title = createEventRequest.Event.Title
	event.Description = createEventRequest.Event.Description
	event.OwnerID = createEventRequest.Event.OwnerId
	event.NotificationSent = createEventRequest.Event.NotificationSent

	if beginDate, err := time.Parse(EventDateFormat, createEventRequest.Event.BeginDate); err != nil {
		event.BeginDate = time.Now()
		s.logger.Error(err.Error())
	} else {
		event.BeginDate = beginDate
	}

	if endDate, err := time.Parse(EventDateFormat, createEventRequest.Event.EndDate); err != nil {
		event.EndDate = time.Now()
		s.logger.Error(err.Error())
	} else {
		event.EndDate = endDate
	}

	event, err := s.app.CreateEvent(ctx, event)
	if err != nil {
		s.logger.Error(err.Error())
	}

	return &pb.CreateEventResponse{
		Event: &pb.Event{
			Id:               event.ID,
			Title:            event.Title,
			BeginDate:        event.BeginDate.Format(EventDateFormat),
			EndDate:          event.EndDate.Format(EventDateFormat),
			Description:      event.Description,
			OwnerId:          event.OwnerID,
			NotificationSent: event.NotificationSent,
		},
	}, nil
}

// UpdateEvent handles updating given event via grpc.
func (s *Service) UpdateEvent(ctx context.Context, updateEventRequest *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	event := storage.Event{}

	event.ID = updateEventRequest.Event.Id
	event.Title = updateEventRequest.Event.Title
	event.Description = updateEventRequest.Event.Description
	event.OwnerID = updateEventRequest.Event.OwnerId
	event.NotificationSent = updateEventRequest.Event.NotificationSent

	if beginDate, err := time.Parse(EventDateFormat, updateEventRequest.Event.BeginDate); err != nil {
		event.BeginDate = time.Now()
		s.logger.Error(err.Error())
	} else {
		event.BeginDate = beginDate
	}

	if endDate, err := time.Parse(EventDateFormat, updateEventRequest.Event.EndDate); err != nil {
		event.EndDate = time.Now()
		s.logger.Error(err.Error())
	} else {
		event.EndDate = endDate
	}

	event, err := s.app.UpdateEvent(ctx, event)
	if err != nil {
		s.logger.Error(err.Error())
	}

	return &pb.UpdateEventResponse{
		Event: &pb.Event{
			Id:               event.ID,
			Title:            event.Title,
			BeginDate:        event.BeginDate.Format(EventDateFormat),
			EndDate:          event.EndDate.Format(EventDateFormat),
			Description:      event.Description,
			OwnerId:          event.OwnerID,
			NotificationSent: event.NotificationSent,
		},
	}, nil
}

// RemoveEvent handles removing an event via grpc.
func (s *Service) RemoveEvent(ctx context.Context, removeEventRequest *pb.RemoveEventRequest) (*pb.RemoveEventResponse, error) {
	event := storage.Event{}
	event.ID = removeEventRequest.Id

	err := s.app.RemoveEvent(ctx, event)
	if err != nil {
		return &pb.RemoveEventResponse{}, err
	}

	return &pb.RemoveEventResponse{}, nil
}

// GetDayAheadEvents handles getting daily events via grpc.
func (s *Service) GetDayAheadEvents(ctx context.Context, _ *pb.GetDayAheadEventsRequest) (*pb.GetDayAheadEventsResponse, error) {
	events, err := s.app.GetDayAheadEvents(ctx)
	if err != nil {
		return &pb.GetDayAheadEventsResponse{}, err
	}

	pbEvents := &pb.GetDayAheadEventsResponse{}
	pbEvents.Items = make([]*pb.Event, len(events))

	for i, event := range events {
		pbEvents.Items[i] = &pb.Event{
			Id:               event.ID,
			Title:            event.Title,
			BeginDate:        event.BeginDate.Format(EventDateFormat),
			EndDate:          event.EndDate.Format(EventDateFormat),
			Description:      event.Description,
			OwnerId:          event.OwnerID,
			NotificationSent: event.NotificationSent,
		}
	}

	return pbEvents, nil
}

// GetWeekAheadEvents handles getting weekly events via grpc.
func (s *Service) GetWeekAheadEvents(ctx context.Context, _ *pb.GetWeekAheadEventsRequest) (*pb.GetWeekAheadEventsResponse, error) {
	events, err := s.app.GetWeekAheadEvents(ctx)
	if err != nil {
		return &pb.GetWeekAheadEventsResponse{}, err
	}

	pbEvents := &pb.GetWeekAheadEventsResponse{}
	pbEvents.Items = make([]*pb.Event, len(events))

	for i, event := range events {
		pbEvents.Items[i] = &pb.Event{
			Id:               event.ID,
			Title:            event.Title,
			BeginDate:        event.BeginDate.Format(EventDateFormat),
			EndDate:          event.EndDate.Format(EventDateFormat),
			Description:      event.Description,
			OwnerId:          event.OwnerID,
			NotificationSent: event.NotificationSent,
		}
	}

	return pbEvents, nil
}

// GetMonthAheadEvents handles getting monthly events via grpc.
func (s *Service) GetMonthAheadEvents(ctx context.Context, _ *pb.GetMonthAheadEventsRequest) (*pb.GetMonthAheadEventsResponse, error) {
	events, err := s.app.GetMonthAheadEvents(ctx)
	if err != nil {
		return &pb.GetMonthAheadEventsResponse{}, err
	}

	pbEvents := &pb.GetMonthAheadEventsResponse{}
	pbEvents.Items = make([]*pb.Event, len(events))

	for i, event := range events {
		pbEvents.Items[i] = &pb.Event{
			Id:               event.ID,
			Title:            event.Title,
			BeginDate:        event.BeginDate.Format(EventDateFormat),
			EndDate:          event.EndDate.Format(EventDateFormat),
			Description:      event.Description,
			OwnerId:          event.OwnerID,
			NotificationSent: event.NotificationSent,
		}
	}

	return pbEvents, nil
}

// NewServer returns a new grpc server instance.
func NewServer(config Config, app Application, logger Logger) *Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger.GetZapLogger()),
		)),
	)

	service := &Service{
		pb.UnimplementedCalendarServer{},
		app,
		logger,
	}

	pb.RegisterCalendarServer(server, service)

	return &Server{
		server: server,
		config: config,
		logger: logger,
	}
}

// Start launches a GRPC server.
func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", net.JoinHostPort(s.config.GetGrpcHost(), s.config.GetGrpcPort()))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrServerStart, err.Error())
	}

	if err := s.server.Serve(lsn); err != nil {
		return fmt.Errorf("%w: %s", ErrServerStart, err.Error())
	}

	return nil
}

// Stop suspends GRPC server.
func (s *Server) Stop() {
	s.server.Stop()
}
