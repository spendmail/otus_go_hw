package app

import (
	"context"

	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
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
	return event, err
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	event, err := a.Storage.UpdateEvent(ctx, event)
	return event, err
}

func (a *App) RemoveEvent(ctx context.Context, event storage.Event) error {
	err := a.Storage.RemoveEvent(ctx, event)
	return err
}

func (a *App) GetDayAheadEvents(ctx context.Context) ([]storage.Event, error) {
	events, err := a.Storage.GetDayAheadEvents(ctx)
	return events, err
}

func (a *App) GetWeekAheadEvents(ctx context.Context) ([]storage.Event, error) {
	events, err := a.Storage.GetWeekAheadEvents(ctx)
	return events, err
}

func (a *App) GetMonthAheadEvents(ctx context.Context) ([]storage.Event, error) {
	events, err := a.Storage.GetMonthAheadEvents(ctx)
	return events, err
}
