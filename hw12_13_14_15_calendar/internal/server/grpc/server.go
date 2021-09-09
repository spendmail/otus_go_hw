package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	pb "github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/server/grpc/eventpb/api"
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
func (s *Service) CreateEvent(ctx context.Context, requestEvent *pb.Event) (*pb.Event, error) {
	event := storage.Event{}

	event.Title = requestEvent.Title
	event.Description = requestEvent.Description
	event.OwnerID = requestEvent.OwnerId

	if beginDate, err := time.Parse(EventDateFormat, requestEvent.BeginDate); err != nil {
		event.BeginDate = time.Now()
		s.logger.Error(err.Error())
	} else {
		event.BeginDate = beginDate
	}

	if endDate, err := time.Parse(EventDateFormat, requestEvent.EndDate); err != nil {
		event.EndDate = time.Now()
		s.logger.Error(err.Error())
	} else {
		event.EndDate = endDate
	}

	event, err := s.app.CreateEvent(ctx, event)
	if err != nil {
		s.logger.Error(err.Error())
	}

	return &pb.Event{
		Id:          event.ID,
		Title:       event.Title,
		BeginDate:   event.BeginDate.Format(EventDateFormat),
		EndDate:     event.EndDate.Format(EventDateFormat),
		Description: event.Description,
		OwnerId:     event.OwnerID,
	}, nil
}

// UpdateEvent handles updating given event via grpc.
func (s *Service) UpdateEvent(ctx context.Context, requestEvent *pb.Event) (*pb.Event, error) {
	event := storage.Event{}

	event.ID = requestEvent.Id
	event.Title = requestEvent.Title
	event.Description = requestEvent.Description
	event.OwnerID = requestEvent.OwnerId

	if beginDate, err := time.Parse(EventDateFormat, requestEvent.BeginDate); err != nil {
		event.BeginDate = time.Now()
		s.logger.Error(err.Error())
	} else {
		event.BeginDate = beginDate
	}

	if endDate, err := time.Parse(EventDateFormat, requestEvent.EndDate); err != nil {
		event.EndDate = time.Now()
		s.logger.Error(err.Error())
	} else {
		event.EndDate = endDate
	}

	event, err := s.app.UpdateEvent(ctx, event)
	if err != nil {
		s.logger.Error(err.Error())
	}

	return &pb.Event{
		Id:          event.ID,
		Title:       event.Title,
		BeginDate:   event.BeginDate.Format(EventDateFormat),
		EndDate:     event.EndDate.Format(EventDateFormat),
		Description: event.Description,
		OwnerId:     event.OwnerID,
	}, nil
}

// RemoveEvent handles removing an event via grpc.
func (s *Service) RemoveEvent(ctx context.Context, requestEvent *pb.Event) (*pb.Message, error) {
	event := storage.Event{}
	event.ID = requestEvent.Id

	err := s.app.RemoveEvent(ctx, event)
	if err != nil {
		return &pb.Message{}, err
	}

	return &pb.Message{
		Message: fmt.Sprintf("Event %d has been successfully removed", event.ID),
	}, nil
}

// GetDayAheadEvents handles getting daily events via grpc.
func (s *Service) GetDayAheadEvents(ctx context.Context, requestEvent *pb.Empty) (*pb.Events, error) {
	events, err := s.app.GetDayAheadEvents(ctx)
	if err != nil {
		return &pb.Events{}, err
	}

	pbEvents := &pb.Events{}
	pbEvents.Items = make([]*pb.Event, len(events))

	for i, event := range events {
		pbEvents.Items[i] = &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			BeginDate:   event.BeginDate.Format(EventDateFormat),
			EndDate:     event.EndDate.Format(EventDateFormat),
			Description: event.Description,
			OwnerId:     event.OwnerID,
		}
	}

	return pbEvents, nil
}

// GetWeekAheadEvents handles getting weekly events via grpc.
func (s *Service) GetWeekAheadEvents(ctx context.Context, requestEvent *pb.Empty) (*pb.Events, error) {
	events, err := s.app.GetWeekAheadEvents(ctx)
	if err != nil {
		return &pb.Events{}, err
	}

	pbEvents := &pb.Events{}
	pbEvents.Items = make([]*pb.Event, len(events))

	for i, event := range events {
		pbEvents.Items[i] = &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			BeginDate:   event.BeginDate.Format(EventDateFormat),
			EndDate:     event.EndDate.Format(EventDateFormat),
			Description: event.Description,
			OwnerId:     event.OwnerID,
		}
	}

	return pbEvents, nil
}

// GetMonthAheadEvents handles getting monthly events via grpc.
func (s *Service) GetMonthAheadEvents(ctx context.Context, requestEvent *pb.Empty) (*pb.Events, error) {
	events, err := s.app.GetMonthAheadEvents(ctx)
	if err != nil {
		return &pb.Events{}, err
	}

	pbEvents := &pb.Events{}
	pbEvents.Items = make([]*pb.Event, len(events))

	for i, event := range events {
		pbEvents.Items[i] = &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			BeginDate:   event.BeginDate.Format(EventDateFormat),
			EndDate:     event.EndDate.Format(EventDateFormat),
			Description: event.Description,
			OwnerId:     event.OwnerID,
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
