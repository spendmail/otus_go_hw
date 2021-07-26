package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrCreateEvent         = errors.New("create event error")
	ErrUpdateEvent         = errors.New("update event error")
	ErrRemoveEvent         = errors.New("removing event error")
	ErrGetDayAheadEvents   = errors.New("getting day events error")
	ErrGetWeekAheadEvents  = errors.New("getting week events error")
	ErrGetMonthAheadEvents = errors.New("getting month events error")
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error)
	RemoveEvent(ctx context.Context, event storage.Event) error
	GetDayAheadEvents(ctx context.Context) ([]storage.Event, error)
	GetWeekAheadEvents(ctx context.Context) ([]storage.Event, error)
	GetMonthAheadEvents(ctx context.Context) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger,
		storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	event, err := a.Storage.CreateEvent(ctx, event)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrCreateEvent, err.Error())
	}

	return event, err
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	event, err := a.Storage.UpdateEvent(ctx, event)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrUpdateEvent, err.Error())
	}

	return event, err
}

func (a *App) RemoveEvent(ctx context.Context, event storage.Event) error {
	err := a.Storage.RemoveEvent(ctx, event)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrRemoveEvent, err.Error())
	}

	return err
}

func (a *App) GetDayAheadEvents(ctx context.Context) ([]storage.Event, error) {
	events, err := a.Storage.GetDayAheadEvents(ctx)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrGetDayAheadEvents, err.Error())
	}

	return events, err
}

func (a *App) GetWeekAheadEvents(ctx context.Context) ([]storage.Event, error) {
	events, err := a.Storage.GetWeekAheadEvents(ctx)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrGetWeekAheadEvents, err.Error())
	}

	return events, err
}

func (a *App) GetMonthAheadEvents(ctx context.Context) ([]storage.Event, error) {
	events, err := a.Storage.GetMonthAheadEvents(ctx)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrGetMonthAheadEvents, err.Error())
	}

	return events, err
}
